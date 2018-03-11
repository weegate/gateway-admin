// author wuyong
// date   2017/12/26
// desc

package main

import (
	"io"
	"os"
	"path/filepath"

	"gateway-admin/controller"
	"gateway-admin/lib"
	"gateway-admin/mw"
	"gateway-admin/view"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {

	InitLog()
}

func InitLog() {
	gin.DisableConsoleColor()

	accessFile, _ := os.Create(*lib.AppCfg.LogDir + "/access.log")
	gin.DefaultWriter = io.MultiWriter(accessFile)

	if *lib.AppCfg.IsDev {
		gin.DefaultWriter = io.MultiWriter(accessFile, os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
	}

	appFile, _ := os.Create(*lib.AppCfg.LogDir + "/app.log")
	logrus.SetOutput(appFile)
}

func main() {

	templateDir := *lib.AppTplCfg.FeDir + "/" + "template"
	staticDir := *lib.AppTplCfg.FeDir + "/" + "static"
	for _, checkDir := range []string{templateDir, staticDir} {
		if !lib.IsDirExist(checkDir) {
			panic(checkDir + " is not a dir; please check it!")
		}
	}

	ginEngine := gin.New()
	ginEngine.Static("/static", staticDir)
	ginEngine.Use(gin.Logger(), gin.Recovery())

	Page(ginEngine, templateDir, *lib.AppCfg.AppName)
	Api(ginEngine, *lib.AppCfg.AppName)

	//@todo use define gin engine and endless/grace service to listen
	ginEngine.Run(*lib.AppCfg.ServerPort)
}

// load template files when run the app service (default add CURD tpl files)
func InitHTMLRender(viewDir string, appName string) view.Render {
	tplDir, _ := filepath.Abs(viewDir)
	appTplDir := tplDir + "/" + appName
	if !lib.IsDirExist(appTplDir) {
		panic(appTplDir + "is not exist; please check it!")
	}

	var files []string
	var dirs []string
	moduleDirs := lib.GetFileBydir(appTplDir, files, dirs)
	numDirs := len(moduleDirs)
	if numDirs == 0 {
		panic("app view dir " + appTplDir + " have no module dirs")
	}

	render := view.New()
	for _, moduleName := range moduleDirs {
		viewTpl := view.BaseView{
			TplDir:       tplDir,
			AppName:      appName,
			ModuleName:   moduleName,
			OptPageNames: []string{"add", "list", "index", "update"},
		}
		viewTpl.RegisterRenders(&render)
	}

	return render
}

// page
func Page(ginEngine *gin.Engine, viewDir string, appName string) {
	ginEngine.Delims(*lib.AppTplCfg.LeftDelim, *lib.AppTplCfg.RightDelim)

	//interface Render is view render
	ginEngine.HTMLRender = InitHTMLRender(viewDir, appName)

	// /->index navigation page
	ginEngine.GET("/", controller.RenderIndexPage)
	// /logout
	ginEngine.GET("/logout", controller.Logout(appName))

	authRouter := ginEngine.Group("/")
	authRouter.Use(mw.CheckLogin(appName))
	{
		// for page rout. eg: /abtest/policy or  /abtest/policy/{add,info,update...}
		authRouter.GET("/"+appName+"/*module_nav", controller.RenderPage(appName))

	} //end authRouter
}

// api
func Api(ginEngine *gin.Engine, appName string) {
	//add api version for app product iterate
	apiRouter := ginEngine.Group("/v1")
	apiRouter.Use(mw.CheckLogin(appName))
	apiRouter.GET("/api/:app/:module/:info", controller.Dispatch)
	apiRouter.POST("/api/:app/:module", controller.Dispatch)
	apiRouter.PUT("/api/:app/:module", controller.Dispatch)
	apiRouter.DELETE("/api/:app/:module/:id", controller.Dispatch)
}
