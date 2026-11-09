
<template>
    <section class="my-10">
        <div class="container">
            <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
                <div class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2">
                    <!-- img -->
                    <img src="~/assets/images/svg-graphics/signin-g.svg" alt="" class="max-w-full h-auto" />
                </div>
                <div class="w-full md:w-1/2 lg:mx-1/5 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6">
                    <div class="flex flex-col gap-1">
                        <h1 class="text-xl">GreenCartにログイン</h1>
                        <p>おかえりなさい！メールアドレスとパスワードを入力して<br>ログインしてください。</p>
                    </div>
                    <form @submit.prevent="handleSubmit" class="needs-validation" novalidate>
                        <div class="flex flex-col gap-5">
                            <div class="flex flex-col gap-3">
                                <div class="w-full">
                                    <!-- input -->
                                    <label for="formSigninEmail" class="invisible hidden">メールアドレス</label>
                                    <input type="email"
                                        class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                                        id="formSigninEmail" placeholder="メールアドレス" v-model="email" required />
                                    <div
                                        v-for="(error, index) in authStore.validationErrors.filter(e => e.field === 'email')"
                                        :key="index"
                                        class="mt-1 text-red-600 text-sm"
                                    >
                                        {{ error.message }}
                                    </div>
                                    <div class="invalid-feedback">有効なメールアドレスを入力してください。</div>
                                </div>
                                <div class="w-full">
                                    <!-- input -->
                                    <div class="password-field relative">
                                        <label for="formSigninPassword" class="invisible hidden">パスワード</label>
                                        <div class="password-field relative">
                                            <input :type="showPassword ? 'text' : 'password'"
                                                class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base fakePassword"
                                                id="formSigninPassword" placeholder="*****" v-model="password" required />
                                            <span @click="togglePassword">
                                                <i :class="showPassword ? 'ti ti-eye' : 'ti ti-eye-off'" class="passwordToggler"></i>
                                            </span>
                                            <div
                                                v-for="(error, index) in authStore.validationErrors.filter(e => e.field === 'password')"
                                                :key="index"
                                                class="mt-1 text-red-600 text-sm"
                                            >
                                                {{ error.message }}
                                            </div>
                                            <div class="invalid-feedback">パスワードを入力してください。</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="flex flex-col gap-4">
                                <div class="flex justify-between w-full">
                                    <div>
                                        パスワードをお忘れですか？
                                        <a href="/auth/forgot-password" class="text-green-600">再設定する</a>
                                    </div>
                                </div>
                                <!-- btn -->
                                <div class="w-full grid">
                                    <button :disabled="authStore.loading" type="submit"
                                        class="btn inline-flex items-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 hover:border-green-700 focus:outline-none focus:ring-4 focus:ring-green-300">
                                        <span v-if="authStore.loading">処理中...</span>
                                        <span v-else>ログイン</span>
                                    </button>
                                </div>
                            </div>
                            <div v-if="authStore.error" class="mt-4 text-red-600">
                                {{ authStore.error }}
                            </div>
                        </div>
                        <!-- link -->
                        <div class="mt-2">
                            アカウントをお持ちでない方はこちら
                            <a href="signup" class="text-green-600">新規会員登録</a>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </section>
</template>


<script setup>
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'auth',
})

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const handleSubmit = async () => {
    try {
        await authStore.signin(email.value, password.value)

        console.log('User:', authStore.user)
        console.log('Is Admin:', authStore.isSuperAdmin())

        if (authStore.isSuperAdmin()) {
            await navigateTo('/dashboard')
        } else {
            await navigateTo('/')
        }
    } catch (error) {
        console.error('Login error:', error)
    }
}
</script>