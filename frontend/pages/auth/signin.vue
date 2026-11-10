
<script setup>
import { ref, reactive, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'auth'
})

useHead({
  title: 'ログイン | GreenCart'
})

const authStore = useAuthStore()

const {
  error: loginError,
  loading,
  validationErrors
} = storeToRefs(authStore)

// Form state
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const submitError = ref('')

// Frontend validation
const localErrors = reactive({
  email: '',
  password: ''
})

// Backend + frontend validation
const emailError = computed(() => {
  return localErrors.email ||
    (validationErrors.value ?? []).find(
      item => item.field === 'email'
    )?.message || ''
})

const passwordError = computed(() => {
  return localErrors.password ||
    (validationErrors.value ?? []).find(
      item => item.field === 'password'
    )?.message || ''
})

const generalError = computed(() => {
  return loginError.value || submitError.value
})

// Clear validation when user types
function clearFieldError(field) {
  localErrors[field] = ''

  authStore.validationErrors =
    (authStore.validationErrors ?? []).filter(
      item => item.field !== field
    )

  authStore.error = null
  submitError.value = ''
}

// Validate before sending request
function validateForm() {
  localErrors.email = ''
  localErrors.password = ''

  if (!email.value.trim()) {
    localErrors.email = 'メールアドレスを入力してください。'
  } else if (
    !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())
  ) {
    localErrors.email = '有効なメールアドレスを入力してください。'
  }

  if (!password.value) {
    localErrors.password = 'パスワードを入力してください。'
  }

  return !localErrors.email && !localErrors.password
}

// Login
async function handleSubmit() {
  if (loading.value) return

  submitError.value = ''
  authStore.error = null
  authStore.validationErrors = []

  if (!validateForm()) return

  try {
    await authStore.signin(
      email.value.trim(),
      password.value
    )

    // Login successful
    const destination = authStore.isAdmin ? '/dashboard' : '/'
    await navigateTo(destination, { replace: true })

  } catch (err) {
    console.error('[GreenCart] Login failed:', err)

    // Backend validation errors are handled by Pinia.
    const hasFieldErrors =
      (authStore.validationErrors ?? []).length > 0

    if (!authStore.error && !hasFieldErrors) {
      const backendMessage = err?.data?.message

      submitError.value =
        typeof backendMessage === 'string'
          ? backendMessage
          : 'ログインに失敗しました。もう一度お試しください。'
    }
  }
}
</script>


<template>
  <section class="my-10">
    <div class="container">
      <div
        class="flex flex-wrap justify-center items-center gap-8 lg:gap-16"
      >
        <!-- FreshCart original illustration -->
        <div
          class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2"
        >
          <img
            src="~/assets/images/svg-graphics/signin-g.svg"
            alt="お買い物のイラスト"
            class="max-w-full h-auto"
          />
        </div>

        <!-- Login form -->
        <div
          class="w-full md:w-1/2 lg:mx-1/5 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6"
        >
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">
              GreenCartにログイン
            </h1>

            <p>
              おかえりなさい！メールアドレスとパスワードを入力して
              ログインしてください。
            </p>
          </div>

          <form
            class="needs-validation"
            novalidate
            @submit.prevent="handleSubmit"
          >
            <div class="flex flex-col gap-5">
              <div class="flex flex-col gap-3">

                <!-- Email -->
                <div class="w-full">
                  <label
                    for="formSigninEmail"
                    class="invisible hidden"
                  >
                    メールアドレス
                  </label>

                 <input
                   id="formSigninEmail"
                   v-model.trim="email"
                   type="email"
                   autocomplete="email"
                   placeholder="メールアドレス"
                   :aria-invalid="!!emailError"
                   :aria-describedby="emailError ? 'email-error' : undefined"
                   @input="clearFieldError('email')"
                   class="form-control border border-gray-300 text-gray-900 rounded-lg block p-2 px-3 w-full text-base"
                 />

                 <p
                   v-if="emailError"
                   id="email-error"
                   class="mt-1 text-red-600 text-sm"
                 >
                   {{ emailError }}
                 </p>
                </div>

                <!-- Password -->

                <div class="w-full">
                  <label for="formSigninPassword" class="invisible hidden">
                    パスワード
                  </label>

                  <div class="password-input-wrapper">
                    <input
                      id="formSigninPassword"
                      v-model="password"
                      :type="showPassword ? 'text' : 'password'"
                      autocomplete="current-password"
                      placeholder="パスワード"
                      :aria-invalid="!!passwordError"
                      :aria-describedby="passwordError ? 'password-error' : undefined"
                      @input="clearFieldError('password')"
                      class="form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 w-full text-base"
                    />

                    <button
                      type="button"
                      class="password-toggle"
                      :aria-label="showPassword ? 'パスワードを隠す' : 'パスワードを表示'"
                      :aria-pressed="showPassword"
                      @click="showPassword = !showPassword"
                    >
                      <Icon
                        :name="showPassword ? 'tabler:eye' : 'tabler:eye-off'"
                        size="20"
                      />
                    </button>
                  </div>

                  <p
                    v-if="passwordError"
                    id="password-error"
                    class="mt-1 text-red-600 text-sm"
                  >
                    {{ passwordError }}
                  </p>
                </div>

              </div>

              <!-- Forgot password -->
              <div class="flex flex-col gap-4">
                <div class="flex justify-between w-full">
                  <div>
                    パスワードをお忘れですか？
                    <NuxtLink
                      to="/auth/forgot-password"
                      class="text-green-600"
                    >
                      再設定する
                    </NuxtLink>
                  </div>
                </div>

                <div
                  v-if="generalError"
                  role="alert"
                  class="text-red-600 text-sm"
                >
                  {{ generalError }}
                </div>

                <!-- Submit -->
                <div class="w-full grid">
                  <button
                    type="submit"
                    :disabled="loading"
                    class="btn inline-flex items-center justify-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 focus:outline-none focus:ring-4 focus:ring-green-300"
                  >
                    <span v-if="loading">
                      ログイン中...
                    </span>

                    <span v-else>
                      ログイン
                    </span>
                  </button>
                </div>
              </div>

            </div>

            <!-- Signup -->
            <div class="mt-2">
              アカウントをお持ちでない方はこちら
              <NuxtLink
                to="/auth/signup"
                class="text-green-600"
              >
                新規会員登録
              </NuxtLink>
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
  width: 100%;
  padding-right: 48px;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  z-index: 2;
}

.password-toggle:hover {
  color: #0aad0a;
}

.password-toggle:focus-visible {
  outline: 2px solid #0aad0a;
  outline-offset: 2px;
  border-radius: 4px;
}
</style>