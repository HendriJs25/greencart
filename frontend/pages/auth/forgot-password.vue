<script setup>
import { computed, reactive, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'auth',
  alias: '/forgot-password',
  key: route => route.fullPath
})

const route = useRoute()
const authStore = useAuthStore()

const isReset = computed(() =>
  /^\/forgot-password\/?$/.test(route.path) || route.query.token !== undefined
)

const token = computed(() => {
  const value = route.query.token
  return typeof value === 'string' && value.trim() ? value : ''
})

const title = computed(() =>
  isReset.value ? '新しいパスワードを設定' : 'パスワードをお忘れですか？'
)

useHead(() => ({
  title: `${title.value} | GreenCart`,
  meta: [{ name: 'referrer', content: 'no-referrer' }]
}))

const form = reactive({
  email: '',
  password_new: '',
  password_confirmation: ''
})

const passwordFields = [
  { key: 'password_new', label: '新しいパスワード' },
  { key: 'password_confirmation', label: '新しいパスワード（確認）' }
]

const shown = reactive({
  password_new: false,
  password_confirmation: false
})

const loading = ref(false)
const fieldErrors = ref({})
const error = ref('')
const success = ref('')
const tokenRejected = ref(false)
const invalidLink = computed(() => isReset.value && (!token.value || tokenRejected.value))

const inputClass = 'form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 w-full text-base'

function clearFieldError(field) {
  delete fieldErrors.value[field]

  if (field === 'password_new') {
    delete fieldErrors.value.password_confirmation
  }

  error.value = ''
}

async function handleSubmit() {
  if (loading.value || success.value || invalidLink.value) return

  error.value = ''
  const errors = {}

  if (isReset.value) {
    if (!form.password_new) {
      errors.password_new = '新しいパスワードを入力してください。'
    } else if (new TextEncoder().encode(form.password_new).length > 72) {
      errors.password_new = '新しいパスワードは72バイト以内で入力してください。'
    }

    if (!form.password_confirmation) {
      errors.password_confirmation = '確認用のパスワードを入力してください。'
    } else if (form.password_new !== form.password_confirmation) {
      errors.password_confirmation = 'パスワードが一致していません。'
    }
  } else if (!form.email.trim()) {
    errors.email = 'メールアドレスを入力してください。'
  }

  fieldErrors.value = errors
  if (Object.keys(errors).length) return

  loading.value = true

  try {
    const result = isReset.value
      ? await authStore.resetPassword(token.value, form)
      : await authStore.forgotPassword(form.email.trim())

    if (result?.status !== 'SUCCESS') {
      throw new Error('Invalid password recovery response')
    }

    form.password_new = ''
    form.password_confirmation = ''

    success.value = result.message || (isReset.value
      ? 'パスワードの再設定が完了しました。'
      : 'ご入力のメールアドレスが登録されている場合、再設定のご案内をお送りします。')
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode
    const message = (err?.data ?? err?.response?._data)?.message

    if (status === 422 && Array.isArray(message) && message.length) {
      fieldErrors.value = Object.fromEntries(
        message.map(({ field, message }) => [field, message])
      )
    } else if (isReset.value && status === 401) {
      tokenRejected.value = true
      form.password_new = ''
      form.password_confirmation = ''
      error.value = typeof message === 'string' ? message : ''
    } else {
      error.value = typeof message === 'string' && message
        ? message
        : '処理に失敗しました。しばらくしてからもう一度お試しください。'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
        <div class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2">
          <img
            src="~/assets/images/svg-graphics/fp-g.svg"
            alt=""
            class="max-w-full h-auto"
          />
        </div>

        <div class="w-full md:w-1/2 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6">
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">{{ title }}</h1>
            <p v-if="!isReset">ご登録のメールアドレスを入力してください。</p>
            <p v-else-if="!invalidLink && !success">新しいパスワードを入力してください。</p>
          </div>

          <div v-if="invalidLink" class="p-4 text-sm text-red-800 rounded-lg bg-red-50">
            <p role="alert">
              {{ error || '再設定リンクが無効か、有効期限が切れています。' }}
            </p>
            <NuxtLink to="/auth/forgot-password" class="text-green-600">
              新しい再設定リンクをリクエスト
            </NuxtLink>
          </div>

          <p
            v-else-if="success"
            class="p-4 text-sm text-green-800 rounded-lg bg-green-50"
            role="status"
          >
            {{ success }}
          </p>

          <form
            v-else
            class="needs-validation flex flex-col gap-4"
            novalidate
            :aria-busy="loading"
            @submit.prevent="handleSubmit"
          >
            <p v-if="error" class="p-4 text-sm text-red-800 rounded-lg bg-red-50" role="alert">
              {{ error }}
            </p>

            <div v-if="!isReset" class="w-full">
              <label for="recovery-email" class="sr-only">メールアドレス</label>
              <input
                id="recovery-email"
                v-model.trim="form.email"
                :class="inputClass"
                type="email"
                autocomplete="email"
                placeholder="メールアドレス"
                required
                :disabled="loading"
                :aria-invalid="!!fieldErrors.email"
                :aria-describedby="fieldErrors.email ? 'recovery-email-error' : undefined"
                @input="clearFieldError('email')"
              />
              <p v-if="fieldErrors.email" id="recovery-email-error" class="mt-1 text-sm text-red-600">
                {{ fieldErrors.email }}
              </p>
            </div>

            <template v-else>
              <div v-for="field in passwordFields" :key="field.key" class="w-full">
                <label :for="`recovery-${field.key}`" class="sr-only">{{ field.label }}</label>

                <div class="password-input-wrapper">
                  <input
                    :id="`recovery-${field.key}`"
                    v-model="form[field.key]"
                    :class="inputClass"
                    :type="shown[field.key] ? 'text' : 'password'"
                    autocomplete="new-password"
                    :placeholder="field.label"
                    required
                    :disabled="loading"
                    :aria-invalid="!!fieldErrors[field.key]"
                    :aria-describedby="fieldErrors[field.key] ? `recovery-${field.key}-error` : undefined"
                    @input="clearFieldError(field.key)"
                  />
                  <button
                    type="button"
                    class="password-toggle"
                    :aria-label="`${field.label}の表示`"
                    :aria-pressed="shown[field.key]"
                    :disabled="loading"
                    @click="shown[field.key] = !shown[field.key]"
                  >
                    <Icon :name="shown[field.key] ? 'tabler:eye' : 'tabler:eye-off'" size="20" />
                  </button>
                </div>

                <p
                  v-if="fieldErrors[field.key]"
                  :id="`recovery-${field.key}-error`"
                  class="mt-1 text-sm text-red-600"
                >
                  {{ fieldErrors[field.key] }}
                </p>
              </div>
            </template>

            <button
              type="submit"
              :disabled="loading"
              class="btn w-full inline-flex items-center justify-center bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 focus:outline-none focus:ring-4 focus:ring-green-300"
            >
              <template v-if="loading">処理中...</template>
              <template v-else-if="isReset">パスワードを再設定</template>
              <template v-else>再設定メールを送信</template>
            </button>
          </form>

          <NuxtLink
            to="/auth/signin"
            class="btn w-full inline-flex items-center justify-center bg-gray-200 text-gray-700 border-gray-200 hover:bg-gray-300"
          >
            ログイン画面に戻る
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.password-input-wrapper {
  position: relative;
}

.password-input-wrapper input {
  padding-right: 48px;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #64748b;
  cursor: pointer;
}

.password-toggle:focus-visible {
  outline: 2px solid #0aad0a;
  outline-offset: 2px;
}
</style>