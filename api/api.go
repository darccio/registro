package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/mozillazg/request"
)

// Bind authenticates an user using given credentials.
func Bind(username, password string) (ok bool, err error) {
	c := new(http.Client)
	rq := request.NewRequest(c)
	rq.Data = map[string]string{
		"username": username,
		"password": password,
	}
	rq.BasicAuth = request.BasicAuth{
		Username: os.Getenv("REGISTRO_BINDUSERNAME"),
		Password: os.Getenv("REGISTRO_BINDPASSWORD"),
	}
	rs, err := rq.Post(os.Getenv("REGISTRO_ENDPOINT"))
	if err != nil {
		return
	}
	if !rs.OK() {
		err = errors.New(rs.Reason())
	}
	if rs.StatusCode == http.StatusNoContent {
		ok = true
	}
	return
}
