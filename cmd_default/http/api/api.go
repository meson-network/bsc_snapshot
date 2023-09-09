package api

import (
	"github.com/meson-network/bsc-data-file-utils/basic"

	"github.com/labstack/echo/v4"

	_ "github.com/meson-network/bsc-data-file-utils/cmd_default/http/api_docs"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/gen"
)

// for swagger
// @title           api example
// @version         1.0
// @description     api example
// @termsOfService  https://domain.com
// @contact.name    Support
// @contact.url     https://domain.com
// @contact.email   contact@domain.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes         https

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func DeclareApi(httpServer *echo.Echo) {
	// health
	configHealth(httpServer)
	// user
	configUser(httpServer)
	// kv
	configKv(httpServer)

}

func ConfigApi(httpServer *echo.Echo) {
	httpServer.GET("/*", echoSwagger.WrapHandler)
}

func GenApiDocs() {

	api_f, api_f_exist, _ := basic.PathExist("cmd_default/http/api")
	if !api_f_exist {
		basic.Logger.Errorln("api_doc_gen_search_dir folder not exist:", api_f)
		return
	}
	api_doc_f, api_doc_f_exist, _ := basic.PathExist("cmd_default/http/api_docs")
	if !api_doc_f_exist {
		basic.Logger.Errorln("api_doc_gen_output_dir folder not exist:", api_doc_f)
		return
	}

	config := &gen.Config{
		SearchDir:       api_f,
		OutputDir:       api_doc_f,
		MainAPIFile:     "api.go",
		OutputTypes:     []string{"go", "json", "yaml"},
		ParseDependency: true,
	}

	err := gen.New().Build(config)
	if err != nil {
		basic.Logger.Errorln("Gen_Api_Docs", err)
	}

}
