// @author wuyong
// @date   2017/12/29
// @desc session shared for sso when there are multi admin service

package lib

import (
	"net/http"
)

type session struct {
	Id uint64
}

type sessionMrg struct {
	sess *session
}

func (sessMrg *sessionMrg) StartSession(w http.ResponseWriter, r *http.Request) (session, error) {
	var sess session
	return sess, nil
}
