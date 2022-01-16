export default {
    login: {
        provider: "/login/provider"
    },
    api: {
        v1: {
            permissions: "/api/v1/g/permissions",
            roles: "/api/v1/g/organizations/{0}/roles",
            rolePermissions: "/api/v1/g/organizations/{0}/roles/{1}/permissions",
            organizations: "/api/v1/g/organizations",
            users: "/api/v1/g/organizations/{0}/users/{1}",
            tenants: '/api/v1/tenants'
            // resources: "/api/v1/g/resources"
        }
    }
}