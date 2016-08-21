package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/registro/api"
	ldap "github.com/vjeantet/ldapserver"
)

func main() {
	server := ldap.NewServer()
	routes := ldap.NewRouteMux()
	routes.Bind(handleBind)
	server.Handle(routes)
	go server.ListenAndServe("0.0.0.0:389")
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)
	server.Stop()
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
	ok, err := api.Bind(r.Name(), r.AuthenticationSimple())
	if ok {
		res.SetResultCode(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	if err == nil {
		err = errors.New("invalid credentials")
	}
	res.SetDiagnosticMessage(err.Error())
	w.Write(res)
}
