//@author wuyong
//@date   2018/1/11
//@desc

package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"gateway-admin/auth"
	"gateway-admin/controller/abtest"
	"gateway-admin/lib"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BaseController struct {
	// context data
	Ctx *gin.Context

	// route controller info
	AppName        string
	ControllerName string
	ActionName     string
	MethodMapping  map[string]func() //method:routertree
}

// @todo: dispatch RESTFUL API uri path gracefully
func Dispatch(context *gin.Context) {
	app := context.Param("app")
	module := context.Param("module")
	info := context.Param("info")
	httpRequestMethod := context.Request.Method
	appModule := strings.Join([]string{httpRequestMethod, app, module}, "->")
	//fmt.Println(appModule,info)

	switch appModule {
	case "GET->abtest->policy":
		if info == "list" {
			abtestController.GetPolicyList(context)
		} else if id, err := strconv.Atoi(info); err == nil && id > 0 {
			context.Set("id", id)
			abtestController.GetPolicyById(context)
		}
	case "GET->abtest->policy_group":
		if info == "list" {
			abtestController.GetPolicyGroupList(context)
		} else if id, err := strconv.Atoi(info); err == nil && id > 0 {
			context.Set("id", id)
			abtestController.GetPolicyGroupById(context)
		}
	case "GET->abtest->runtime":
		if info == "list" {
			abtestController.GetRuntimeList(context)
		} else if id, err := strconv.Atoi(info); err == nil && id > 0 {
			context.Set("id", id)
			abtestController.GetRuntimeById(context)
		}

	case "POST->abtest->policy":
		abtestController.AddPolicy(context)
	case "POST->abtest->policy_group":
		abtestController.AddPolicyGroup(context)
	case "POST->abtest->runtime":
		abtestController.AddRuntime(context)

	case "PUT->abtest->policy":
		abtestController.UpdatePolicy(context)
	case "PUT->abtest->policy_group":
		abtestController.UpdatePolicyGroup(context)
	case "PUT->abtest->runtime":
		abtestController.UpdateRuntime(context)

	case "DELETE->abtest->policy":
		abtestController.DeletePolicy(context)
	case "DELETE->abtest->policy_group":
		abtestController.DeletePolicyGroup(context)
	case "DELETE->abtest->runtime":
		abtestController.DeleteRuntime(context)

	default:
		context.JSON(http.StatusOK, lib.GetErrorByName("UNKNOWN_METHOD"))
		context.Abort()
	}
}

// render index page -> just a navigation page
func RenderIndexPage(context *gin.Context) {
	renderName := "index"
	context.HTML(http.StatusOK, renderName, gin.H{
		"title":       "MIS",
		"description": "this is a navigation index page!",
	})
}

func Logout(appName string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var sso = &auth.SsoAuth{
			AppName: *lib.AppSsoAuthCfg.AppName,
			Version: *lib.AppSsoAuthCfg.Version,
			Host:    *lib.AppSsoAuthCfg.Host,
		}
		if _, err := context.Request.Cookie(appName + "user"); err != nil {
			logrus.Errorf("get abuser from cookie failed")
		} else {
			respCookie := &http.Cookie{
				Name:    appName + "user",
				Value:   "",
				Path:    "/",
				Expires: time.Now().Add(-1),
			}
			http.SetCookie(context.Writer, respCookie)
		}

		sso.Logout("", context.Writer, context.Request)
		context.Abort()
	} //end func
}

func RenderPage(app string) gin.HandlerFunc {
	return func(context *gin.Context) {
		module := ""
		nav := ""
		userName, ok := context.Get("userName")
		if !ok {
			logrus.Errorf("get userName failed")
		}

		moduleNav := context.Param("module_nav")
		strPaths := strings.Split(moduleNav, "/")
		if len(strPaths) == 2 {
			module = strPaths[1]
		} else if len(strPaths) == 3 {
			module = strPaths[1]
			nav = strPaths[2]
		}

		if len(module) == 0 {
			renderName := strings.Join([]string{app, "index"}, "_")
			context.HTML(http.StatusOK, renderName, gin.H{
				"userName": userName,
			})
		} else {
			renderName := ""
			if len(nav) > 0 {
				renderName = strings.Join([]string{app, module, nav}, "_")
			} else {
				renderName = strings.Join([]string{app, module, "index"}, "_")
			}

			context.HTML(http.StatusOK, renderName, gin.H{
				"userName": userName,
			})
		} //end if
	} //end func
}

func (c *BaseController) Mapping(method string, fn func()) {
	c.MethodMapping[method] = fn
}

func (c *BaseController) Register(app string, module string, method string, fn func()) {
	c.AppName = app
	c.ControllerName = module
	method = strings.Join([]string{app, module, method}, "_")
	c.Mapping(method, fn)
}
