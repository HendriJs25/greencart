<script setup>
import { computed } from 'vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
  adminOnly: true,
  validate: route => typeof route.params.id === 'string' && /^[1-9]\d{0,18}$/.test(route.params.id)
})

useHead({ title: '顧客詳細 | GreenCart' })

const route = useRoute()
const listLocation = computed(() => ({
  path: '/dashboard/customers',
  query: route.query
}))

const { data: customer, status, error, refresh } = await useFetch(
  () => `/api/user/admin/customers/${route.params.id}`,
  {
    lazy: true,
    transform(response) {
      const detail = response?.data

      if (response?.status !== 'SUCCESS' || !detail ||
          typeof detail !== 'object' || Array.isArray(detail)) {
        throw createError({ statusCode: 502 })
      }

      return detail
    }
  }
)

const errorCode = computed(() => error.value?.statusCode)

if (import.meta.server && errorCode.value === 404) {
  setResponseStatus(404)
}

const errorMessage = computed(() => {
  if (errorCode.value === 400) return '顧客IDが正しくありません。'
  if (errorCode.value === 401) return 'もう一度ログインしてください。'
  if (errorCode.value === 403) return '顧客詳細を表示する権限がありません。'
  if (errorCode.value === 404) return '顧客が見つかりませんでした。'

  const message = error.value?.data?.message
  return typeof message === 'string' && message
    ? message
    : '顧客詳細を読み込めませんでした。'
})

const fields = [
  { key: 'id', label: '顧客ID' },
  { key: 'uuid', label: 'UUID' },
  { key: 'email', label: 'メールアドレス' },
  { key: 'phone', label: '電話番号' },
  { key: 'address', label: '住所' },
  { key: 'lat', label: '緯度' },
  { key: 'lng', label: '経度' }
]

function displayValue(value) {
  return value == null || value === '' ? '未登録' : value
}
</script>

<template>
  <section class="container">
    <h1 class="mb-2 text-xl">顧客詳細</h1>

    <nav class="mb-6 text-sm text-gray-600" aria-label="パンくずリスト">
      <NuxtLink to="/dashboard" class="text-green-600">ダッシュボード</NuxtLink>
      <span aria-hidden="true"> / </span>
      <NuxtLink :to="listLocation" class="text-green-600">顧客一覧</NuxtLink>
      <span aria-hidden="true"> / </span>
      <span aria-current="page">顧客詳細</span>
    </nav>

    <div class="card card-lg">
      <div class="p-4 md:p-6">
        <p v-if="status === 'idle' || status === 'pending'" role="status">
          顧客詳細を読み込んでいます…
        </p>

        <div v-else-if="error">
          <p class="text-red-600" role="alert">{{ errorMessage }}</p>
          <NuxtLink v-if="errorCode === 401" to="/auth/signin" class="text-green-600">
            ログイン画面へ
          </NuxtLink>
          <button
            v-else-if="!errorCode || errorCode >= 500"
            type="button"
            class="btn btn-sm border-gray-300 text-gray-700"
            @click="refresh()"
          >
            もう一度読み込む
          </button>
        </div>

        <template v-else-if="customer">
          <div class="flex flex-wrap items-center gap-4 border-b border-gray-200 pb-4">
            <span class="flex w-16 h-16 items-center justify-center rounded-full bg-green-50 text-green-600">
              <Icon name="tabler:user" size="28" aria-hidden="true" />
            </span>

            <div class="min-w-0">
              <h2 class="mb-2 text-lg customer-detail-value">{{ customer.name }}</h2>
              <span
                class="inline-block rounded-full px-3 py-1 text-sm"
                :class="customer.is_verified ? 'bg-green-50 text-green-800' : 'bg-gray-100 text-gray-700'"
              >
                {{ customer.is_verified ? 'メール認証済み' : 'メール未認証' }}
              </span>
            </div>
          </div>

          <dl class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 mb-0">
            <div v-for="field in fields" :key="field.key">
              <dt class="mb-1 text-sm text-gray-600">{{ field.label }}</dt>
              <dd class="m-0 text-gray-900 customer-detail-value">
                {{ displayValue(customer[field.key]) }}
              </dd>
            </div>
          </dl>
        </template>
      </div>

      <div class="border-t border-gray-200 p-4 md:p-6">
        <NuxtLink :to="listLocation" class="btn btn-sm border-gray-300 text-gray-700">
          一覧に戻る
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.customer-detail-value {
  overflow-wrap: anywhere;
  white-space: pre-line;
}
</style>
