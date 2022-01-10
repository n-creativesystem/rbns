<template>
  <div v-if="finish">
    <v-card
      :tile="$vuetify.breakpoint.sm || $vuetify.breakpoint.xs"
      class="mx-auto full-width"
      flat
      max-width="640"
    >
      <div class="px-6 py-8">
        <div style="max-width:500px" class="mx-auto">
          <v-form ref="form" v-model="valid" @submit.prevent="onLogin">
            <div>
              <div class="d-flex justify-center">
                <v-img
                  lazy-src="@assets/icon.svg"
                  max-height="150"
                  max-width="150"
                  src="@assets/icon.svg"
                ></v-img>
              </div>
              <div class="d-flex justify-center">
                <h1>Role Based N Securityへようこそ</h1>
              </div>
              <div style="max-width: 350px;" class="mx-auto">
                <v-text-field
                  v-model="userId"
                  dense
                  height="48px"
                  outlined
                  placeholder="ユーザーID"
                  label="ユーザーID"
                  :hide-details="false"
                  required
                ></v-text-field>
                <v-text-field
                  v-model="password"
                  :append-icon="passwordShow ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="passwordShow ? 'text' : 'password'"
                  dense
                  height="48px"
                  name="input-password"
                  outlined
                  placeholder="パスワード"
                  @click:append="passwordShow = !passwordShow"
                  label="パスワード"
                  :hide-details="false"
                  required
                ></v-text-field>
              </div>
              <div style="max-width: 350px;" class="login-btn pb-8 mx-auto">
                <div>
                  <v-btn
                    class="full-width"
                    type="submit"
                    :disabled="!valid"
                    :loading="loading"
                    color="primary"
                  >
                    ログイン
                  </v-btn>
                </div>
                <div
                  v-if="Object.keys(providers).length > 0"
                  class="d-flex full-width justify-space-between"
                >
                  <div>
                    <div class="border"></div>
                  </div>
                  <div>
                    <span><span>or</span></span>
                  </div>
                  <div><div class="border"></div></div>
                </div>
                <template v-for="(value, key) in providers">
                  <v-btn
                    :key="key"
                    class="text-capitalize full-width"
                    :disabled="!valid"
                    :loading="loading"
                    color="light-blue lighten-5"
                    :href="`/login/${key}`"
                  >
                    Sign In With {{ value.name }}
                  </v-btn>
                </template>
              </div>
            </div>
          </v-form>
        </div>
      </div>
    </v-card>
  </div>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'login',
    mixins: [axiosMixin],
    data() {
      return {
        valid: true,
        userId: null,
        password: null,
        passwordShow: false,
        to: this.$route.query.to ? this.$route.query.to : '/',
        message: '',
        loginFailed: false,
        finish: true,
        providers: {},
      }
    },
    methods: {
      async onLogin() {
        this.$refs.form.validate()
        this.loading = true
        try {
          //
        } catch (err) {
          this.loginFailed = true
          this.message = err.message
        } finally {
          this.loading = false
        }
      },
      getProvider() {
        this.get(this.$urls.login.provider)
          .then((result) => {
            this.providers = result.data.data
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
    created() {
      if (this.IsLogin) {
        this.$router.push({ path: '/' })
      }
      try {
        this.getProvider()
      } catch (err) {
        this.finish = true
        //
      }
    },
  }
</script>

<style lang="sass" scoped>
  .border
    width: 100px
    height: 10px
    border-bottom: 1px solid black
  .bm-8
    margin-bottom: 8px
</style>
