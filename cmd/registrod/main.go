package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/registro/api"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

func main() {
	server := ldap.NewServer()
	routes := ldap.NewRouteMux()
	routes.Bind(handleBind)
	routes.Search(handleSearch)
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

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()
	baseObject := string(r.BaseObject())
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultUnwillingToPerform)
	if baseObject == "" {
		w.Write(res)
		m.Abandon()
		return
	}
	var (
		err error
	)
	switch baseObject {
	case os.Getenv("REGISTRO_USERSDN"):
		res, err = handleUserSearch(w, m)
	case os.Getenv("REGISTRO_GROUPSDN"):
		res, err = handleGroupSearch(w, m)
	}
	if err != nil {
		log.Printf("%s", err)
	}
	w.Write(res)
	m.Abandon()
}

func handleUserSearch(w ldap.ResponseWriter, m *ldap.Message) (res message.SearchResultDone, err error) {
	r := m.GetSearchRequest()
	// TODO register all searches (log.DEBUG)
	res = ldap.NewSearchResultDoneResponse(ldap.LDAPResultUnwillingToPerform)
	users, err := api.GetUsers(r.Filter())
	if err != nil {
		return
	}
	if len(users) == 0 {
		res = ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
		return
	}
	for _, u := range users {
		w.Write(u.ToLDAPEntry())
	}
	res = ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	return
}

func handleGroupSearch(w ldap.ResponseWriter, m *ldap.Message) (res message.SearchResultDone, err error) {
	res = ldap.NewSearchResultDoneResponse(ldap.LDAPResultUnwillingToPerform)
	return
}
