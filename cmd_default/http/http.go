package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/cmd_default/http/api"
	"github.com/meson-network/bsc-data-file-utils/plugin/echo_plugin"
	common_api "github.com/meson-network/bsc-data-file-utils/src/common/http/api"
)

// httpServer example
func ServerStart() {

	toml_conf := config.Get_config().Toml_config

	//init matched echo

	///////////////////
	api_echo, err := echo_plugin.InitMatchedEcho("api", func(host, req_uri string) bool {
		return strings.HasPrefix(host, "api")
	})
	if err != nil {
		basic.Logger.Fatalln(err)
	} else {
		api.DeclareApi(api_echo.Echo)
		api.ConfigApi(api_echo.Echo)
		api_echo.HTTPErrorHandler = func(err error, ctx echo.Context) {
			ctx.JSON(http.StatusOK, common_api.API_META_STATUS{Meta_status: -1, Meta_message: err.Error()})
		}
	}

	///////html echo /////////////
	html_echo, err := echo_plugin.InitMatchedEcho("www", func(host, req_uri string) bool {
		return strings.HasPrefix(host, "www")
	})

	if err == nil {

		https_html_abs_dir, https_html_abs_dir_exist, _ := basic.PathExist(toml_conf.Https.Html_dir)
		if !https_html_abs_dir_exist {
			basic.Logger.Fatalln("https.html_dir:" + toml_conf.Https.Html_dir + " not exist on disk")
		}

		html_file := filepath.Join(https_html_abs_dir, "index.html")
		_, err := os.Stat(html_file)
		if err != nil {
			basic.Logger.Fatalln("index.html file not found inside " + https_html_abs_dir)
		}

		html_echo.Static("/*", https_html_abs_dir)
		//router all traffic to html.page as for html mode
		html_echo.HTTPErrorHandler = func(err error, ctx echo.Context) {
			ctx.File(html_file)
		}
	}

	/////////////////////

	//config https server
	https_srv := echo_plugin.GetInstance_("https")
	if https_srv != nil {

		https_srv.Any("/*", func(ctx echo.Context) error {
			res := ctx.Response()
			req := ctx.Request()

			matched_echo := echo_plugin.CheckMatchedEcho(req.Host, req.RequestURI)
			if matched_echo == nil {
				html_echo.ServeHTTP(res, req)
			} else {
				matched_echo.ServeHTTP(res, req)
			}

			return nil
		})

		go func() {
			err := https_srv.Start()
			if err != nil {
				basic.Logger.Fatalln("default https echo server start err:", err)
			}
		}()
	}

	//config http server
	//http just redirect to https
	http_srv := echo_plugin.GetInstance_("http")
	if http_srv != nil {
		if https_srv != nil {
			http_srv.Any("/*", func(ctx echo.Context) error {
				domain := strings.Split(ctx.Request().Host, ":")[0]
				return ctx.Redirect(301, "https://"+domain+":"+strconv.Itoa(https_srv.Http_port)+ctx.Request().URL.String())
			})
		}

		go func() {
			err := http_srv.Start()
			if err != nil {
				basic.Logger.Fatalln("default http echo server start err:", err)
			}
		}()
	}

}

func ServerReloadCert() error {
	return echo_plugin.GetInstance_("https").ReloadCert()
}

func ServerCheckStarted() bool {

	http_srv := echo_plugin.GetInstance_("http")
	if http_srv != nil {
		if !http_srv.CheckStarted() {
			return false
		}
	}

	https_srv := echo_plugin.GetInstance_("https")
	if https_srv != nil {
		if !https_srv.CheckStarted() {
			return false
		}
	}

	return true
}
