export default {
    login: {
        provider: "/login/provider"
    },
    api: {
        v1: {
            permissions: "/api/g/v1/permissions",
            roles: "/api/g/v1/roles",
            organizations: "/api/g/v1/organizations",
            users: "/api/g/v1/organizations/{0}/users/{1}",
            // resources: "/api/g/v1/resources"
        }
    }
}