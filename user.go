package registro

import (
	"fmt"
	"os"

	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

// User is a retrieved user from API
type User struct {
	Username  string
	Email     string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ToLDAPEntry converts to LDAP SearchResultEntry
func (u *User) ToLDAPEntry() (ul message.SearchResultEntry) {
	ul = ldap.NewSearchResultEntry(fmt.Sprintf("uid=%s,%s", u.Username, os.Getenv("REGISTRO_USERSDN")))
	ul.AddAttribute("uid", message.AttributeValue(u.Username))
	ul.AddAttribute("mail", message.AttributeValue(u.Email))
	ul.AddAttribute("cn", message.AttributeValue(u.FirstName))
	ul.AddAttribute("sn", message.AttributeValue(u.LastName))
	ul.AddAttribute("objectClass", message.AttributeValue("inetOrgPerson"))
	return
}
