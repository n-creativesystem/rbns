import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'hash',
  routes: [
    {
      path: '/',
      name: 'top',
      component: () => import('@tpl/Layout.vue'),
      redirect: '/permissions',
      children: [
        {
          path: 'permissions',
          name: 'permissions',
          component: () => import('@page/Permissions')
        },
        {
          path: 'roles',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'roles',
              component: () => import('@page/Roles'),
            },
            {
              path: ':id',
              name: 'roleId',
              component: () => import('@page/Role')
            }
          ]
        },
        {
          path: 'organizations',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'organizations',
              component: () => import('@page/Organizations')
            },
            {
              path: ':id',
              name: 'organization-id',
              component: () => import('@page/Organization')
            },
            {
              path: ':id/users/:userKey',
              name: 'user-id',
              component: () => import('@page/User')
            }
          ]
        },
        {
          path: 'resources',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'resources',
              component: () => import('@page/Resources')
            },
            {
              path: ':id',
              name: 'resourceId',
              component: () => import('@page/Resource')
            }
          ]
        }
      ]
    },
    {
      path: '/*',
      name: 'notFound',
      redirect: '/'
    },
  ]
})
