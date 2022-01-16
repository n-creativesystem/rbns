import { IsLogin } from './login'

const menu = [
    {
        title: "permission",
        to: "/permissions",
        icon: "mdi-security",
        key: "permissions"
    },
    {
        title: "organization",
        to: "/organizations",
        icon: "mdi-domain",
        key: "organizations"
    },
    {
        title: "role",
        to: "/roles",
        icon: "mdi-shield-account",
        key: "roles"
    },
    {
        title: "users",
        to: "/users",
        icon: "mdi-account",
        key: "users"
    },
]

if (IsLogin) {
    menu.push({
        title: "logout",
        href: "/logout",
        icon: "mdi-logout",
        key: "logout"
    })
}

export default menu