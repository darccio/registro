package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/mozillazg/request"
	"github.com/vjeantet/goldap/message"
)

// Bind authenticates an user using given credentials.
func Bind(dn message.LDAPDN, password message.OCTETSTRING) (ok bool, err error) {
	c := new(http.Client)
	rq := request.NewRequest(c)
	username, err := resolveUsername(dn, os.Getenv("REGISTRO_BASEDN"))
	if err != nil {
		return
	}
	data := map[string]string{
		"username": username,
		"password": string(password),
	}
	js, err := json.Marshal(data)
	if err != nil {
		return
	}
	rq.Json, err = simplejson.NewJson(js)
	if err != nil {
		return
	}
	rq.BasicAuth = request.BasicAuth{
		Username: os.Getenv("REGISTRO_BINDUSERNAME"),
		Password: os.Getenv("REGISTRO_BINDPASSWORD"),
	}
	rs, err := rq.Post(os.Getenv("REGISTRO_ENDPOINT"))
	if !rs.OK() {
		err = errors.New(rs.Reason())
	}
	if err != nil {
		return
	}
	if rs.StatusCode == http.StatusNoContent {
		ok = true
	}
	return
}

func resolveUsername(dn message.LDAPDN, baseDN string) (username string, err error) {
	sDN := string(dn)
	eos := strings.LastIndex(sDN, baseDN)
	if eos == -1 || (len(sDN) != eos+len(baseDN)) {
		err = errors.New("invalid base distinguished name")
		return
	}
	fields := strings.Count(sDN[:eos], "=")
	if fields != 1 {
		err = errors.New("invalid distinguished name")
		return
	}
	bos := strings.Index(sDN, "=")
	username = sDN[bos+1 : eos-1]
	return
}
