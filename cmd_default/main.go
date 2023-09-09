package cmd_default

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/color"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/cmd_default/http"
	"github.com/meson-network/bsc-data-file-utils/plugin/auto_cert_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
)

func StartDefault() {

	////////////////////////////////
	//print build mode
	basic.Logger.Infoln("build mode:", basic.GetMode())
	basic.Logger.Infoln("log level:", config.Get_config().Toml_config.Log.Level)
	//log basic info
	basic.Logger.Debugln("------------------------------------")
	basic.Logger.Debugln("working dir:", basic.WORK_DIR)
	basic.Logger.Debugln("------------------------------------")
	basic.Logger.Debugln("using user config toml file:", config.Get_config().User_config_path)
	///////////////////////////////

	//print logo first
	fmt.Println(color.Green(basic.Logo))

	//init components and run example
	InitComponent()

	//start threads jobs
	go start_jobs()

	start_components()

	//example code
	test_g_token := token_mgr.TokenMgr.GenToken()
	fmt.Println("a test token generated :", test_g_token)

	fmt.Println(color.Green("---------example code result ------------"))
	fmt.Println("color test :", color.Red("red"))
	fmt.Println("color test :", color.RedBG("red bg"))

	fmt.Println("color test :", color.Blue("blue"))
	fmt.Println("color test :", color.BlueBG("blue bg"))

	fmt.Println("color test :", color.Green("green"))
	fmt.Println("color test :", color.GreenBG("green bg"))

	fmt.Println("color test :", color.Yellow("yellow"))
	fmt.Println("color test :", color.YellowBG("yellow bg"))

	fmt.Println("color test :", color.Cyan("cyan"))
	fmt.Println("color test :", color.CyanBG("cyan bg"))

	//fmt.Println(color.Blue("---------geo ip  result ------------"))
	// if config.Get_config().Toml_config.Geo_ip.Enable {
	// 	fmt.Println(color.Blue(geo_ip_plugin.GetInstance().GetInfo("129.146.243.246")))
	// 	fmt.Println(color.Blue(geo_ip_plugin.GetInstance().GetInfo("192.168.189.125")))
	// 	fmt.Println(color.Blue(geo_ip_plugin.GetInstance().GetInfo("2600:4040:a912:a200:a438:9968:96d9:c3e4")))
	// 	fmt.Println(color.Blue(geo_ip_plugin.GetInstance().GetInfo("2600:387:1:809::3a")))
	// }

	go func() {
		fmt.Println("running and sleep...")
		for {
			time.Sleep(5 * time.Second)
			fmt.Print(".")
		}
	}()

	//waiting for exist signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	// disconnect all connected node
	basic.Logger.Infoln("cli-template exit")

}

func start_components() {

	//start the httpserver
	http.ServerStart()

}

func start_jobs() {
	//check all services already started
	if !http.ServerCheckStarted() {
		basic.Logger.Fatalln("http server not working")
	}

	// //start the auto_cert auto-updating job
	if config.Get_config().Toml_config.Auto_cert.Enable {
		auto_cert_plugin.GetInstance().AutoUpdate(func(new_crt_str, new_key_str string) {
			//reload server
			sre := http.ServerReloadCert()
			if sre != nil {
				basic.Logger.Errorln("cert change reload echo server error:" + sre.Error())
			}
		})
	}

	//start your jobs below
}
