package user_permission
import (
	"github.com/Glintt/gAPI/api/users"
	"strings"
)

const (
	USER_PERMISSION_CHECK = "id in (select service_id from gapi_user_services_permissions c where c.user_id = '##USER_ID##') or 1 = ##IS_USER_ADMIN##"
)

func AppendPermissionFilterToQuery(query string, table string, user users.User) string {
	query = query + " and "
	permissionQuery := USER_PERMISSION_CHECK

	permissionQuery = strings.Replace(permissionQuery, "##USER_ID##", user.Id.Hex(), -1)
	isAdminValue := "0"
	if user.IsAdmin {
		isAdminValue = "1"
	} 
	permissionQuery = strings.Replace(permissionQuery, "##IS_USER_ADMIN##", isAdminValue, -1)

	query = query + "(" + table + "." + permissionQuery + ")"
	return query
}