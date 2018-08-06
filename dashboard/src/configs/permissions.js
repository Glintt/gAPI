export const AdminPermissions = [
	'ServiceDiscovery.CreateService',
	'ServiceDiscovery.DeleteService',
	'ServiceDiscovery.UpdateService',
	'ServiceDiscovery.ManageService'
]

export function HasPermission (type, user) {
	if (!user) return false
	if (user.IsAdmin) {
		return true
	}

	if (AdminPermissions.indexOf(type) !== -1) {
		return false
	}

	return true
}
