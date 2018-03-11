//@author wuyong
//@date   2018/1/8
//@desc

package mw

import (
	"net/http"
	"time"

	"gateway-admin/auth"

	"gateway-admin/lib"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CheckLogin(appName string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var sso = &auth.SsoAuth{
			AppName: *lib.AppSsoAuthCfg.AppName,
			Version: *lib.AppSsoAuthCfg.Version,
			Host:    *lib.AppSsoAuthCfg.Host,
		}
		userName := ""
		if cookie, error := context.Request.Cookie(appName + "user"); error != nil {
			logrus.Errorf("get " + appName + "user from cookie failed")
		} else {
			userName = cookie.Value
		}
		userInfo, err := sso.AuthUser(userName, context.Writer, context.Request)
		if err != nil {
			logrus.Error(err.Error())
			context.Abort()
		}

		if result, ok := userInfo["result"]; !ok {
			context.Abort()
		} else {
			if userName, ok := result["username"]; !ok {
				context.Abort()
			} else {
				cookie := &http.Cookie{
					Name:     appName + "user",
					Value:    userName.(string),
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Now().AddDate(0, 0, 7),
				}
				http.SetCookie(context.Writer, cookie)
				context.Set("userName", userName)
				context.Next()
			}
		}
	} //end func
}
