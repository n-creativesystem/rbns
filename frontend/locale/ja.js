import vuetifyJa from 'vuetify/src/locale/ja.ts'

export default {
  "$vuetify": {
    ...vuetifyJa
  },
  'translations': '翻訳',
  'no data': "データなし",
  'no record': "該当するレコードが見つかりません",
  search: "検索",
  'required fields': '必須フィールドです',
  close: '閉じる',
  save: '保存',
  permission: {
    name: "権限",
    dialog: {
      title: '権限の選択',
      add: '権限の追加'
    },
    entity: {
      name: '権限名称',
      description: '記述'
    },
    form: {
      add: {
        title: '権限の追加',
        caption: 'リソースやロールが使用する権限を定義します。'
      },
      list: {
        title: '権限の一覧',
        caption: 'リソースやロールが使用するすべての権限です。'
      }
    }
  },
  role: {
    name: 'ロール',
    dialog: {
      title: 'ロールの選択',
      add: 'ロールの追加',
      create: {
        title: 'ロールの作成'
      }
    },
    tabs: {
      settings: '設定',
      permissions: '権限',
      users: 'ユーザー'
    },
    entity: {
      name: 'ロール名称',
      description: '記述'
    },
    form: {
      title: 'ロールの一覧',
      caption: '組織のユーザーが使用するロールです。',
      add: 'ロールの作成'
    },
    info: {
      permission: {
        message: 'このロールに権限を追加します。',
        text: '権限の追加'
      }
    }
  },
  organization: {
    name: '組織',
    tabs: {
      settings: '設定',
      users: 'ユーザー'
    },
    form: {
      title: '組織の一覧',
      caption: 'ユーザー、ロールが使用するすべての組織です。',
      add: '組織の作成'
    },
    dialog: {
      create: {
        title: '組織の作成',
        user: {
          title: 'ユーザー情報',
        },
        role: {
          title: 'ロール情報'
        },
      }
    },
    entity: {
      name: '組織名称',
      description: '記述'
    },
    info: {
      user: {
        message: 'この組織にユーザーを追加します。',
        text: 'ユーザーの追加'
      }
    }
  },
  user: {
    name: 'ユーザー',
    tabs: {
      roles: 'ロール'
    },
    dialog: {
      create: {
        title: 'ユーザー情報'
      }
    },
    entity: {
      userId: 'ユーザーID',
    },
    info: {
      role: {
        message: 'このユーザーにロールを追加します',
        text: 'ロールの選択'
      }
    }
  },
  resource: {
    dialog: {
      title: 'リソースの選択',
      add: 'リソースの追加',
      create: {
        title: 'リソースの作成'
      }
    },
    form: {
      title: 'リソースの一覧',
      caption: '登録されているリソースの一覧です。',
      add: 'リソースの作成'
    },
    entity: {
      id: 'リソースID',
      description: '記述'
    },
    tabs: {
      settings: '設定',
      permissions: '権限',
    },
    info: {
      permission: {
        message: 'このリソースに権限を追加します。',
        text: '権限の追加'
      }
    }
  },
  inputs: {
    Name: '名称',
    Description: '記述',
    Add: '追加',
    download: 'ダウンロード',
    upload: 'アップロード'
  },
  menu: {
    permission: '権限の定義',
    role: 'ロールの定義',
    organization: '組織の定義',
    users: 'ユーザーの定義',
    logout: 'ログアウト'
  },
  tenant: {
    name: 'テナント',
    title: 'テナントの作成',
    caption: '組織やユーザーなどが所属するテナントの作成をします。',
    possibleCharactersValidate: `使用可能文字は半角英数字(大文字小文字) 記号は .(ピリオド) と -(ハイフン) のみです。`
  }
}