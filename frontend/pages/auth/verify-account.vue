<script setup>
import { onMounted, ref } from 'vue'

definePageMeta({
  layout: 'auth',
  alias: '/verify-account'
})

useHead({
  title: 'メールアドレス認証 | GreenCart',
  meta: [{ name: 'referrer', content: 'no-referrer' }]
})

const route = useRoute()
const errorMessage = ref('')
const successMessage = ref('')

onMounted(async () => {
  const token = route.query.token

  if (typeof token !== 'string' || !token.trim()) {
    errorMessage.value = '認証リンクが無効です。メールに記載されたリンクをご確認ください。'
    return
  }

  try {
    const result = await $fetch('/api/user/auth/verify-account', {
      method: 'POST',
      query: { token },
      retry: 0
    })

    if (result?.status !== 'SUCCESS') {
      errorMessage.value = '認証結果を確認できませんでした。'
      return
    }

    successMessage.value = result.message || 'メールアドレスの認証が完了しました。'
  } catch (error) {
    const message = error?.data?.message
    errorMessage.value = typeof message === 'string' && message.trim()
      ? message
      : 'メールアドレスの認証に失敗しました。もう一度お試しください。'
  }
})
</script>

<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3 flex flex-col gap-4">
          <h1 class="text-xl">メールアドレス認証</h1>

          <p v-if="errorMessage" class="text-red-600" role="alert">
            {{ errorMessage }}
          </p>
          <p v-else role="status">
            {{ successMessage || 'メールアドレスを認証しています…' }}
          </p>

          <NuxtLink
            v-if="successMessage || errorMessage"
            to="/auth/signin"
            class="text-green-600"
          >
            ログインへ
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>