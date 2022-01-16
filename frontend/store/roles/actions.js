import urls from '@plugins/api'
import { axios } from '@plugins/axios'

export const findAll = ({ commit }, { organizationId }) => {
  return new Promise((resolve, reject) => {
    axios.get(urls.api.v1.roles.format(organizationId))
      .then((result) => {
        if (result.status == 200) {
          if (result.data.roles) {
            commit('setRoles', result.data.roles)
          }
        }
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export const add = (_, { organizationId, roles }) => {
  return new Promise((resolve, reject) => {
    axios.post(urls.api.v1.roles.format(organizationId), {
      roles: roles
    })
      .then((result) => {
        resolve(result)
      }).catch((err) => {
        reject(err)
      });
  })
}

export const put = (_, { organizationId, roleId, data }) => {
  return new Promise((resolve, reject) => {
    axios.put(`${urls.api.v1.roles.format(organizationId)}/${roleId}`, {
      name: data.name,
      description: data.description,
    })
      .then(result => {
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export const remove = (_, { organizationId, roleId }) => {
  return new Promise((resolve, reject) => {
    axios.delete(`${urls.api.v1.roles.format(organizationId)}/${roleId}`)
      .then((result) => {
        resolve(result)
      }).catch((err) => {
        reject(err)
      });
  })
}

export const findById = (_, { organizationId, roleId }) => {
  return new Promise((resolve, reject) => {
    let role = {
      name: '',
      description: '',
      permissions: [],
      error: undefined
    }
    axios.get(`${urls.api.v1.roles.format(organizationId)}/${roleId}`)
      .then((result) => {
        if (result.status == 200) {
          const data = result.data
          if (data) {
            role = {
              name: data.name || '',
              description: data.description || '',
              permissions: data.permissions || [],
            }
          }
        }
        resolve(role)
      }).catch((err) => {
        role.error = err
        reject(role)
      });
  })
}

export const findPermissions = ({ commit }, { organizationId, roleId }) => {
  return new Promise((resolve, reject) => {
    axios.get(`${urls.api.v1.rolePermissions.format(organizationId, roleId)}`)
      .then((result) => {
        commit('setRolePermissions', result.data.permissions)
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export const putRolePermissions = (_, { organizationId, roleId, permissions }) => {
  return new Promise((resolve, reject) => {
    axios.put(`${urls.api.v1.rolePermissions.format(organizationId, roleId)}`, {
      permissions: permissions
    })
      .then((result) => {
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export const removeRolePermissions = (_, { organizationId, roleId, permissionId }) => {
  return new Promise((resolve, reject) => {
    axios.delete(`${urls.api.v1.rolePermissions.format(organizationId, roleId)}/${permissionId}`)
      .then((result) => {
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}
