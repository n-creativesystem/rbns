package rdb

func (f *SQLStore) bus() {
	f.addPermissionBus()
	f.addRoleBus()
	f.addOrganizationBus()
	f.addUserBus()
	f.addApiKeyBus()
	f.addLoginUserBus()
	f.addTenantBas()
}
