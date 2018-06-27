package authentication

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/utils"
	"github.com/go-ldap/ldap"
)

var LDAP_PORT = "389"

func AuthenticateWithLDAP(email string, password string) bool{
	port := LDAP_PORT
	if config.GApiConfiguration.Authentication.LDAP.Port != "" {
		port = config.GApiConfiguration.Authentication.LDAP.Port
	}
	
	l, err := ldap.Dial("tcp", config.GApiConfiguration.Authentication.LDAP.Domain + ":" + port)
	if err != nil {
		utils.LogMessage("LDAP Error: " + err.Error(), utils.ErrorLogType)
	}
	err = l.Bind(email, password)
	if err != nil {
		utils.LogMessage("LDAP Credentials Error: " + err.Error(), utils.DebugLogType)
		return false
	}
	return true
}