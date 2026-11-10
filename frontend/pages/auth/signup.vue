<script setup>
import { reactive, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'auth'
})

useHead({
  title: '新規会員登録 | GreenCart'
})

const authStore = useAuthStore()

const form = reactive({
  name: '',
  email: '',
  password: '',
  password_confirmation: ''
})

const fieldErrors = ref({})
const error = ref('')
const loading = ref(false)
const showPassword = ref(false)
const showPasswordConfirmation = ref(false)
const successMessage = ref('')

const inputClass = 'form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base'

function clearFieldError(field) {
  delete fieldErrors.value[field]
  error.value = ''
}

function clearPasswordErrors() {
  clearFieldError('password')
  clearFieldError('password_confirmation')
}

async function handleSubmit() {
  if (loading.value || successMessage.value) return

  error.value = ''
  const errors = {}

  if (!form.name.trim()) errors.name = 'お名前を入力してください。'
  if (!form.email.trim()) errors.email = 'メールアドレスを入力してください。'
  if (!form.password) errors.password = 'パスワードを入力してください。'

  if (!form.password_confirmation) {
    errors.password_confirmation = '確認用のパスワードを入力してください。'
  } else if (form.password_confirmation !== form.password) {
    errors.password_confirmation = 'パスワードが一致していません。'
  }

  fieldErrors.value = errors
  if (Object.keys(errors).length) return

  loading.value = true

  try {
    const result = await authStore.signup(form)

    form.password = ''
    form.password_confirmation = ''

    successMessage.value = typeof result.message === 'string' && result.message
      ? result.message
      : 'アカウントの登録が完了しました。'
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode
    const payload = err?.data ?? err?.response?._data
    const message = payload?.message

    if (status === 422 && Array.isArray(message) && message.length > 0) {
      fieldErrors.value = Object.fromEntries(
        message.map(({ field, message }) => [field, message])
      )
    } else if (status === 409 && typeof message === 'string' && message) {
      fieldErrors.value = { email: message }
    } else {
      error.value = typeof message === 'string' && message
        ? message
        : '登録に失敗しました。しばらくしてからもう一度お試しください。'
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
            src="~/assets/images/svg-graphics/signup-g.svg"
            alt=""
            class="max-w-full h-auto"
          />
        </div>

        <div class="w-full md:w-1/2 lg:mx-1/6 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6">
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">GreenCartに新規会員登録</h1>
            <p>お名前、メールアドレス、パスワードを入力して登録してください。</p>
          </div>

          <div
            v-if="successMessage"
            class="p-4 text-sm text-green-800 rounded-lg bg-green-50"
            role="status"
          >
            <p class="font-medium">{{ successMessage }}</p>
            <p>確認メールに記載されたリンクから、メールアドレスの認証を完了してください。</p>
            <NuxtLink to="/auth/signin" class="text-green-600">
              ログイン画面へ
            </NuxtLink>
          </div>

          <form
            v-else
            class="needs-validation"
            novalidate
            :aria-busy="loading"
            @submit.prevent="handleSubmit"
          >
            <div
              v-if="error"
              class="mb-4 p-4 text-sm text-red-800 rounded-lg bg-red-50"
              role="alert"
            >
              {{ error }}
            </div>

            <div class="flex flex-wrap gap-5">
              <div class="w-full flex flex-col gap-3">
                <div class="flex flex-row gap-3">
                  <div class="w-1/2">
                    <label for="formSignupName" class="sr-only">お名前</label>
                    <input
                      id="formSignupName"
                      v-model.trim="form.name"
                      :class="inputClass"
                      type="text"
                      autocomplete="name"
                      placeholder="お名前"
                      required
                      :disabled="loading"
                      :aria-invalid="!!fieldErrors.name"
                      :aria-describedby="fieldErrors.name ? 'signup-name-error' : undefined"
                      @input="clearFieldError('name')"
                    />
                    <p v-if="fieldErrors.name" id="signup-name-error" class="mt-1 text-red-600 text-sm">
                      {{ fieldErrors.name }}
                    </p>
                  </div>

                  <div class="w-1/2">
                    <label for="formSignupEmail" class="sr-only">メールアドレス</label>
                    <input
                      id="formSignupEmail"
                      v-model.trim="form.email"
                      :class="inputClass"
                      type="email"
                      autocomplete="email"
                      placeholder="メールアドレス"
                      required
                      :disabled="loading"
                      :aria-invalid="!!fieldErrors.email"
                      :aria-describedby="fieldErrors.email ? 'signup-email-error' : undefined"
                      @input="clearFieldError('email')"
                    />
                    <p v-if="fieldErrors.email" id="signup-email-error" class="mt-1 text-red-600 text-sm">
                      {{ fieldErrors.email }}
                    </p>
                  </div>
                </div>

                <div class="w-full">
                  <label for="formSignupPassword" class="sr-only">パスワード</label>
                  <div class="password-input-wrapper">
                    <input
                      id="formSignupPassword"
                      v-model="form.password"
                      :class="inputClass"
                      :type="showPassword ? 'text' : 'password'"
                      autocomplete="new-password"
                      placeholder="パスワード"
                      required
                      :disabled="loading"
                      :aria-invalid="!!fieldErrors.password"
                      :aria-describedby="fieldErrors.password ? 'signup-password-error' : 'signup-password-hint'"
                      @input="clearPasswordErrors"
                    />
                    <button
                      type="button"
                      class="password-toggle"
                      aria-label="パスワードの表示"
                      :aria-pressed="showPassword"
                      :disabled="loading"
                      @click="showPassword = !showPassword"
                    >
                      <Icon :name="showPassword ? 'tabler:eye' : 'tabler:eye-off'" size="20" />
                    </button>
                  </div>
                  <p v-if="fieldErrors.password" id="signup-password-error" class="mt-1 text-red-600 text-sm">
                    {{ fieldErrors.password }}
                  </p>
                  <p v-else id="signup-password-hint" class="mt-1 text-sm">
                    8文字以上で、大文字・小文字・数字をそれぞれ1文字以上含めてください。
                  </p>
                </div>

                <div class="w-full">
                  <label for="formSignupPasswordConfirmation" class="sr-only">パスワード（確認）</label>
                  <div class="password-input-wrapper">
                    <input
                      id="formSignupPasswordConfirmation"
                      v-model="form.password_confirmation"
                      :class="inputClass"
                      :type="showPasswordConfirmation ? 'text' : 'password'"
                      autocomplete="new-password"
                      placeholder="パスワード（確認）"
                      required
                      :disabled="loading"
                      :aria-invalid="!!fieldErrors.password_confirmation"
                      :aria-describedby="fieldErrors.password_confirmation ? 'signup-confirmation-error' : undefined"
                      @input="clearFieldError('password_confirmation')"
                    />
                    <button
                      type="button"
                      class="password-toggle"
                      aria-label="確認用パスワードの表示"
                      :aria-pressed="showPasswordConfirmation"
                      :disabled="loading"
                      @click="showPasswordConfirmation = !showPasswordConfirmation"
                    >
                      <Icon :name="showPasswordConfirmation ? 'tabler:eye' : 'tabler:eye-off'" size="20" />
                    </button>
                  </div>
                  <p
                    v-if="fieldErrors.password_confirmation"
                    id="signup-confirmation-error"
                    class="mt-1 text-red-600 text-sm"
                  >
                    {{ fieldErrors.password_confirmation }}
                  </p>
                </div>
              </div>

              <div class="w-full grid">
                <button
                  type="submit"
                  :disabled="loading"
                  class="btn inline-flex items-center justify-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 focus:outline-none focus:ring-4 focus:ring-green-300"
                >
                  {{ loading ? '登録中...' : '新規会員登録' }}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.password-input-wrapper {
  position: relative;
  width: 100%;
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

.password-toggle:hover:enabled {
  color: #0aad0a;
}

.password-toggle:focus-visible {
  outline: 2px solid #0aad0a;
  outline-offset: 2px;
  border-radius: 4px;
}
</style>