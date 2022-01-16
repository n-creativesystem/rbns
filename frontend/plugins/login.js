
export default {
    data: () => ({
        IsLogin: IsLogin,
    }),
}

export const user = (window.rbnsBootstrapData && window.rbnsBootstrapData.user ? window.rbnsBootstrapData.user : {})

export const IsLogin = user.IsSignedIn || false

export const currentRole = user.Role || 'Viewer'

export const IsTenant = user.IsTenant || false