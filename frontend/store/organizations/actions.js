import urls from '@plugins/api'
import { axios } from '@plugins/axios'

export const findAll = ({ commit }) => {
  return new Promise((resolve, reject) => {
    axios.get(urls.api.v1.organizations)
      .then((result) => {
        if (result.status == 200) {
          if (result.data.organizations) {
            commit('setOrganizations', result.data.organizations)
          }
        }
        resolve(result)
      })
      .catch((err) => {
        reject(err)
      })
  })
}
