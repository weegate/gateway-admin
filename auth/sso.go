// @author wuyong
// @date   2017/12/26
// @desc
// @todo use gorequest

package auth

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SsoAuth struct {
	AppName string
	Version string
	Host    string
}

func init() {

}

// firstly check ticket is from cookie, if not, check ticket from request params,
// if not in cookie and params, redirect sso login,return empty info(nothing), so u must check this, and abort this request process;
// if username is checked successfully, return user info, if not, get user info from sso auth api and return user info.
func (sso *SsoAuth) AuthUser(username string, w http.ResponseWriter, r *http.Request) (map[string]map[string]interface{}, error) {
	ticket := ""
	fromIp := false

	redirectUrl := "http://" + r.Host + r.RequestURI

	successData := map[string]map[string]interface{}{
		"result": {},
		"status": {},
	}

	if cookie, error := r.Cookie("sso-ticket"); error == nil {
		ticket = cookie.Value
	} else if paramTicket, ok := r.URL.Query()["sso_ticket"]; ok && len(paramTicket) > 0 {
		fromIp = true
		ticket = paramTicket[0]
	} else {
		sso.Redirect(redirectUrl, w, r)
		return successData, nil
	}

	if len(ticket) > 0 {
		tickets := strings.Split(ticket, "-")
		if len(username) > 0 && len(tickets) == 2 {
			md5Handle := md5.New()
			md5Handle.Write([]byte(tickets[0] + username))
			md5Sum := fmt.Sprintf("%x", string(md5Handle.Sum(nil)))

			if md5Sum == tickets[1] {
				successData["result"] = map[string]interface{}{"username": username}
				successData["status"] = map[string]interface{}{"status_code": 0, "status_reason": ""}
			}
		} else {
			url := sso.Host + "/sso/auth"
			params := "sso_ticket=" + ticket + "&sso_version=" + sso.Version + "&redirect=" + redirectUrl + "&app_name=" + sso.AppName
			authUrl := url + "?" + params
			client := &http.Client{}
			req, error := http.NewRequest("GET", authUrl, strings.NewReader(""))
			if error != nil {
				return nil, errors.New(authUrl + " don't access")
			}
			req.Header.Add("User-Agent", sso.AppName)
			req.Header.Add("Referer", redirectUrl)

			respon, error := client.Do(req)
			if error != nil {
				return nil, errors.New("auth response is failed")
			}

			defer respon.Body.Close()
			body, error := ioutil.ReadAll(respon.Body)
			if error != nil {
				return nil, errors.New("get response body failed")
			}

			if json.Unmarshal(body, &successData) != nil {
				return nil, errors.New("json unmarshal failed")
			}

		}

		if successData != nil && fromIp {
			cookie := &http.Cookie{
				Name:     "sso-ticket",
				Value:    ticket,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().AddDate(0, 0, 7),
			}
			http.SetCookie(w, cookie)
		}

		return successData, nil
	}

	if _, error := r.Cookie("sso-ticket"); error == nil {
		respCookie := &http.Cookie{
			Name:    "sso-ticket",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1),
		}
		http.SetCookie(w, respCookie)

	}

	return successData, nil
}

func (sso *SsoAuth) getUserData(url string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	userData := map[string]interface{}{
		"result": map[string]string{},
		"status": map[string]string{},
	}

	return userData, nil
}

func (sso *SsoAuth) Logout(redirectUrl string, w http.ResponseWriter, r *http.Request) {
	respCookie := &http.Cookie{
		Name:    "sso-ticket",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1),
	}
	http.SetCookie(w, respCookie)

	sso.Redirect(redirectUrl, w, r)
}

func (sso *SsoAuth) Redirect(redirectUrl string, w http.ResponseWriter, r *http.Request) {
	var redirect string
	if len(redirectUrl) > 0 {
		redirect = redirectUrl
	} else {
		redirect = "http://" + r.Host + strings.Split(r.RequestURI, "/")[0]
	}
	logoutUrl := sso.Host + "/sso/login?" + "redirect=" + redirect + "&version=" + sso.Version + "&app_name=" + sso.AppName
	http.Redirect(w, r, logoutUrl, http.StatusFound)
}
