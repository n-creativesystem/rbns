import vuetifyEn from 'vuetify/src/locale/en.ts'

export default {
  "$vuetify": {
    ...vuetifyEn
  },
  'translations': 'Translations',
  'no data': "no data",
  'no record': "no record",
  search: "search",
  'required fields': 'required fields',
  close: 'Close',
  save: 'Save',
  permission: {
    dialog: {
      title: 'Permission selection',
      add: 'Add permission'
    },
    entity: {
      name: 'permission name',
      description: 'permission description'
    },
    form: {
      add: {
        title: 'Add permission',
        caption: 'Defines the permissions used by the resource.'
      },
      list: {
        title: 'List of permission',
        caption: 'All permissions used by the resource.'
      }
    }
  },
  role: {
    dialog: {
      title: 'Role selection',
      add: 'Add role',
      create: {
        title: 'Creating a role'
      }
    },
    tabs: {
      settings: 'settings',
      permissions: 'permission',
      users: 'users'
    },
    entity: {
      name: 'Role name',
      description: 'description'
    },
    form: {
      title: 'List of roles',
      caption: 'The role used by the resource.',
      add: 'Creating a role'
    },
    info: {
      permission: {
        message: 'Add permissions to this role.',
        text: 'Add permission'
      }
    }
  },
  organization: {
    tabs: {
      settings: 'settings',
      users: 'users'
    },
    form: {
      title: 'List of organizations',
      caption: 'All organizations used by users and roles.',
      add: 'Creating an organization'
    },
    dialog: {
      create: {
        title: 'Creating an organization',
        user: {
          title: 'User information',
        },
        role: {
          title: 'Role information'
        },
      }
    },
    entity: {
      name: 'Organization name',
      description: 'description'
    },
    info: {
      user: {
        message: 'Add users to this organization.',
        text: 'Add user'
      }
    }
  },
  user: {
    tabs: {
      roles: 'roles'
    },
    dialog: {
      create: {
        title: 'User information'
      }
    },
    entity: {
      userId: 'User ID',
    },
    info: {
      role: {
        message: 'Add a role to this user',
        text: 'Role selection'
      }
    }
  },
  inputs: {
    Name: 'name',
    Description: 'description',
    Add: 'add'
  },
  menu: {
    permission: 'Definition of permissions',
    role: 'Role definition',
    organization: 'Organization definition'
  }
}