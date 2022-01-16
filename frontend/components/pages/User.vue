<template>
  <page-layout>
    <form-parts class="input-and-list">
      <input-form :title="title" :caption="caption">
        <template v-slot:content>
          <n-form @submit="onSubmit" ref="form">
            <v-row>
              <div class="col-12 col-md-4">
                <required-text
                  name="name"
                  :label="$t('inputs.Name')"
                  id="name"
                  v-model="name"
                ></required-text>
              </div>
              <div class="col-12 col-md-4">
                <btn-tfn type="submit">
                  <v-icon>mdi-send</v-icon>
                  {{ $t('inputs.Add') }}
                </btn-tfn>
              </div>
            </v-row>
          </n-form>
          <v-row>
            <div class="col-12 col-md-4">
              <v-file-input
                label="csvファイル"
                hint="ユーザー情報が記載されたCSVファイルのアップロード"
                accept="text/csv"
                show-size
                v-model="file"
              ></v-file-input>
            </div>
            <div class="col-12 col-md-4">
              <btn-tfn @click="onUpload" class="mr-1">
                <v-icon>mdi-cloud-upload</v-icon>
                {{ $t('inputs.upload') }}
              </btn-tfn>
              <btn-tfn @click="onDownload">
                <v-icon>mdi-download</v-icon>
                {{ $t('inputs.download') }}
              </btn-tfn>
            </div>
          </v-row>
        </template>
      </input-form>
    </form-parts>
  </page-layout>
</template>

<script>
  const toBase64 = (file) =>
    new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.readAsDataURL(file)
      reader.onload = () => resolve(reader.result)
      reader.onerror = (error) => reject(error)
    })
  export default {
    name: 'User',
    data: () => ({
      name: '',
      file: null,
    }),
    methods: {
      onSubmit() {},
      onDownload() {
        const blob = new Blob(['ID,名称'], { type: 'text/csv' })
        const link = document.createElement('a')
        link.href = window.URL.createObjectURL(blob)
        link.download = 'user.csv'
        link.click()
      },
      async onUpload() {
        const file = await toBase64(this.file)
        // const formData = new FormData()
        // formData.append('file', this.file)
        // formData.append('fileType.', 'lazyQuote')
        this.$axios
          .post('/api/v1/g/users/files', {
            data: file,
            fileType: {
              name: 'user.csv',
            },
          })
          .then((result) => {
            console.log(result)
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
    computed: {
      title() {
        return ''
      },
      caption() {
        return ''
      },
    },
  }
</script>

<style></style>
