import Vue from 'vue'
import Router from 'vue-router'
import { IsLogin, IsTenant } from './login'

Vue.use(Router)

const router = new Router({
  mode: 'history',
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
          component: () => import('@page/Permissions'),
          meta: {
            auth: true
          }
        },
        {
          path: 'roles',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'roles',
              component: () => import('@page/Roles'),
              meta: {
                auth: true
              },
            },
            {
              path: ':orgId/:id',
              name: 'roleId',
              component: () => import('@page/Role'),
              meta: {
                auth: true
              },
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
              component: () => import('@page/Organizations'),
              meta: {
                auth: true
              },
            },
            {
              path: ':id',
              name: 'organization-id',
              component: () => import('@page/Organization'),
              meta: {
                auth: true
              },
            },
            {
              path: ':id/users/:userKey',
              name: 'user-id',
              component: () => import('@page/User'),
              meta: {
                auth: true
              },
            }
          ]
        },
        {
          path: 'login',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'login',
              component: () => import('@page/Login')
            },
          ]
        },
        {
          path: 'users',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'users',
              component: () => import('@page/User')
            },
          ]
        },
        {
          path: 'tenants',
          component: () => import('@page/Parent'),
          children: [
            {
              path: '/',
              name: 'tenants',
              component: () => import('@page/Tenant')
            },
          ]
        }
        // {
        //   path: 'resources',
        //   component: () => import('@page/Parent'),
        //   children: [
        //     {
        //       path: '/',
        //       name: 'resources',
        //       component: () => import('@page/Resources')
        //     },
        //     {
        //       path: ':id',
        //       name: 'resourceId',
        //       component: () => import('@page/Resource')
        //     }
        //   ]
        // }
      ]
    },
    {
      path: '/*',
      name: 'notFound',
      redirect: '/'
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  }
})

router.beforeEach((to, from, next) => {
  if (to.meta && to.meta.auth) {
    if (IsLogin) {
      if (IsTenant) {
        next()
        return
      }
      next({
        name: 'tenants'
      })
    } else {
      next({
        name: 'login', query: {
          'to': to.path
        }
      })
    }
  } else {
    next()
  }
})

export default router