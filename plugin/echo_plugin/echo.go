package echo_plugin

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/echo_plugin/echo_middleware"

	"github.com/coreservice-io/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var match_echo = map[string]*MatchEcho{}

type MatchEcho struct {
	*echo.Echo
	Name  string
	Match func(string, string) bool // func(host string, req_uri string)
}

func GetMatchEcho(name string) *MatchEcho {
	return match_echo[name]
}

func CheckMatchedEcho(host string, req_uri string) *MatchEcho {
	for _, v := range match_echo {
		if v.Match(host, req_uri) {
			return v
		}
	}
	return nil
}

func InitMatchedEcho(name string, match func(string, string) bool) (*MatchEcho, error) {
	_, exist := match_echo[name]
	if exist {
		return nil, fmt.Errorf("MatchEcho instance <%s> has already been initialized", name)
	}
	match_echo[name] = &MatchEcho{
		echo.New(),
		name,
		match,
	}
	return match_echo[name], nil
}

type EchoServer struct {
	*echo.Echo
	Logger     log.Logger
	Http_port  int
	Keep_alive bool
	Tls        bool
	Crt_path   string
	Key_path   string
	Cert       *tls.Certificate
}

var instanceMap = map[string]*EchoServer{}

func GetInstance() *EchoServer {
	return GetInstance_("default")
}

func GetInstance_(name string) *EchoServer {
	echo_i := instanceMap[name]
	if echo_i == nil {
		basic.Logger.Errorln(name + " echo plugin null")
	}
	return echo_i
}

/*
http_port
*/
type Config struct {
	Port       int
	Keep_alive bool
	Tls        bool
	Crt_path   string
	Key_path   string
}

func Init(serverConfig Config, OnPanicHanlder func(panic_err interface{}), logger log.Logger) error {
	return Init_("default", serverConfig, OnPanicHanlder, logger)
}

// Init a new instance.
//
//	If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//	If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, serverConfig Config, OnPanicHanlder func(panic_err interface{}), logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("echo server instance <%s> has already been initialized", name)
	}

	if serverConfig.Port == 0 {
		serverConfig.Port = 8080
	}

	echoServer := &EchoServer{
		echo.New(),
		logger,
		serverConfig.Port,
		serverConfig.Keep_alive,
		serverConfig.Tls,
		serverConfig.Crt_path,
		serverConfig.Key_path,
		nil,
	}

	// cros
	echoServer.Use(middleware.CORS())

	echoServer.Echo.HideBanner = true

	// logger
	echoServer.Use(echo_middleware.LoggerWithConfig(echo_middleware.LoggerConfig{
		Logger:            logger,
		RecordFailRequest: false,
	}))
	// recover and panicHandler
	echoServer.Use(echo_middleware.RecoverWithConfig(echo_middleware.RecoverConfig{
		OnPanic: OnPanicHanlder,
	}))

	echoServer.JSONSerializer = NewJsoniter()

	instanceMap[name] = echoServer
	return nil
}

func (s *EchoServer) Start() error {

	if s.Tls {
		cert, err := tls.LoadX509KeyPair(s.Crt_path, s.Key_path)
		if err != nil {
			return err
		}
		s.Cert = &cert
		tlsconf := new(tls.Config)
		tlsconf.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return s.Cert, nil
		}

		// more safe?
		// tlsconf.MinVersion = tls.VersionTLS12

		server := http.Server{
			Addr:      ":" + strconv.Itoa(s.Http_port),
			TLSConfig: tlsconf,
		}

		server.SetKeepAlivesEnabled(s.Keep_alive)
		return s.StartServer(&server)

	} else {

		s.Echo.Server.SetKeepAlivesEnabled(s.Keep_alive)
		return s.Echo.Start(":" + strconv.Itoa(s.Http_port))
	}
}

func (s *EchoServer) ReloadCert() error {
	if s.Tls {
		cert, err := tls.LoadX509KeyPair(s.Crt_path, s.Key_path)
		if err != nil {
			return err
		}

		basic.Logger.Infoln("http server certificate reloaded")
		s.Cert = &cert
	}
	return nil
}

func (s *EchoServer) Close() {
	s.Echo.Close()
}

// check the server is indeed up
func (s *EchoServer) CheckStarted() bool {

	time_slot := time.Now().Unix()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	loop_counter := 0
	for {
		loop_counter++
		<-ticker.C
		addr_tls := s.Echo.TLSListenerAddr()
		if addr_tls != nil && strings.Contains(addr_tls.String(), ":") {
			return true
		}
		addr := s.Echo.ListenerAddr()
		if addr != nil && strings.Contains(addr.String(), ":") {
			return true
		}
		if s.Logger != nil && loop_counter%5 == 0 {
			s.Logger.Warnln("server has not started, ", time.Now().Unix()-time_slot, " second has passed")
		}
	}
}

// /////

type JsoniterHandler struct {
	json jsoniter.API
}

func NewJsoniter() *JsoniterHandler {
	return &JsoniterHandler{
		jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (j *JsoniterHandler) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := j.json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

func (j *JsoniterHandler) Deserialize(c echo.Context, i interface{}) error {
	err := j.json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
