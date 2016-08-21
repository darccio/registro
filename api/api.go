package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/mozillazg/request"
	"github.com/vjeantet/goldap/message"
)

// Bind authenticates an user using given credentials.
func Bind(dn message.LDAPDN, password message.OCTETSTRING) (ok bool, err error) {
	rq := newRequest()
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
	rs, err := rq.Post(endpoint("users/bind"))
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

func newRequest() (rq *request.Request) {
	c := new(http.Client)
	rq = request.NewRequest(c)
	rq.BasicAuth = request.BasicAuth{
		Username: os.Getenv("REGISTRO_BINDUSERNAME"),
		Password: os.Getenv("REGISTRO_BINDPASSWORD"),
	}
	return
}

func endpoint(op string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("REGISTRO_ENDPOINT"), op)
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

// GetUsers retrieves users
func GetUsers(filter message.Filter) (users []interface{}, err error) {
	rq := newRequest()
	rq.Params, err = resolveUserFilter(filter)
	if err != nil {
		return
	}
	rs, err := rq.Get(endpoint("users"))
	if !rs.OK() {
		err = errors.New(rs.Reason())
	}
	if err != nil {
		return
	}
	return
}

func resolveUserFilter(filter message.Filter) (params map[string]string, err error) {
	params = make(map[string]string)
	switch filter.(type) {
	case message.FilterEqualityMatch:
		feq := filter.(message.FilterEqualityMatch)
		key := string(feq.AttributeDesc())
		switch key {
		case "uid":
			key = "username"
		case "cn":
			key = "first_name"
		case "sn":
			key = "last_name"
		}
		params[key] = string(feq.AssertionValue())
	default:
		err = errors.New("filter not supported")
	}
	return
}
