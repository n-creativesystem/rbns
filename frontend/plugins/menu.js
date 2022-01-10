import { IsLogin } from './login'

const menu = [
    {
        title: "permission",
        to: "/permissions",
        icon: "mdi-security",
        key: "permissions"
    },
    {
        title: "role",
        to: "/roles",
        icon: "mdi-shield-account",
        key: "roles"
    },
    {
        title: "organization",
        to: "/organizations",
        icon: "mdi-domain",
        key: "organizations"
    },
    // {
    //     "title": "resource",
    //     "href": "/resources",
    //     "icon": "mdi-semantic-web",
    //     "key": "resources"
    // }
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