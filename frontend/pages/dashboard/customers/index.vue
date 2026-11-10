<script setup>
import { computed, ref, watch } from 'vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
  adminOnly: true
})

useHead({ title: '顧客一覧 | GreenCart' })

const route = useRoute()
const columns = [
  { key: 'id', label: '顧客ID', sortable: true },
  { key: 'name', label: 'お名前', sortable: true },
  { key: 'email', label: 'メールアドレス' },
  { key: 'phone', label: '電話番号' },
  { key: 'created_at', label: '登録日', sortable: true },
  { key: 'detail', label: '詳細' }
]

function queryString(key) {
  const value = route.query[key]
  return typeof value === 'string' ? value.trim() : ''
}

const filters = computed(() => {
  const sortBy = queryString('sort_by').toLowerCase()

  return {
    search: queryString('search'),
    sort_by: columns.some(column => column.sortable && column.key === sortBy)
      ? sortBy
      : 'created_at',
    sort_order: queryString('sort_order').toLowerCase() === 'asc' ? 'asc' : 'desc'
  }
})

const searchInput = ref(filters.value.search)
const page = computed(() => {
  const raw = queryString('page')
  const number = /^\d+$/.test(raw) ? Number(raw) : 1
  return Number.isSafeInteger(number) && number > 0 ? number : 1
})

watch([page, filters], () => {
  searchInput.value = filters.value.search
})

function customerPage(number, options = filters.value) {
  const query = {
    page: String(number),
    sort_by: options.sort_by,
    sort_order: options.sort_order
  }

  if (options.search) query.search = options.search

  return { path: '/dashboard/customers', query }
}

function submitSearch() {
  return navigateTo(customerPage(1, { ...filters.value, search: searchInput.value.trim() }))
}

function sortLink(column) {
  const ascending = filters.value.sort_by === column && filters.value.sort_order === 'asc'

  return customerPage(1, {
    ...filters.value,
    sort_by: column,
    sort_order: ascending ? 'desc' : 'asc'
  })
}

function sortDirection(column) {
  if (column !== filters.value.sort_by) return undefined
  return filters.value.sort_order === 'asc' ? 'ascending' : 'descending'
}

function sortIcon(column) {
  const direction = sortDirection(column)
  if (!direction) return 'tabler:arrows-sort'
  return direction === 'ascending' ? 'tabler:arrow-up' : 'tabler:arrow-down'
}

const query = computed(() => ({
  ...filters.value,
  page: page.value,
  search: filters.value.search || undefined
}))

const { data: result, status, error, refresh } = await useFetch('/api/user/admin/customers', {
  key: 'admin-customers',
  query,
  lazy: true,
  transform(response) {
    if (response?.status !== 'SUCCESS' || !Array.isArray(response.data) ||
        !Number.isInteger(response.meta?.TotalPages)) {
      throw createError({ statusCode: 502 })
    }

    return {
      customers: response.data,
      meta: response.meta
    }
  }
})

const notice = useState('admin-customer-notice', () => '')
const customerNotice = ref(notice.value)
notice.value = ''

const errorCode = computed(() => error.value?.statusCode)

const errorMessage = computed(() => {
  if (errorCode.value === 401) return 'もう一度ログインしてください。'
  if (errorCode.value === 403) return '顧客一覧を表示する権限がありません。'

  const message = error.value?.data?.message
  return typeof message === 'string' && message
    ? message
    : '顧客一覧を読み込めませんでした。'
})

const pageNumbers = computed(() => {
  const total = result.value?.meta.TotalPages ?? 0
  const current = Math.min(page.value, total)
  const start = Math.max(1, Math.min(current - 2, total - 4))

  return Array.from({ length: Math.min(5, total) }, (_, index) => start + index)
})

const rowRange = computed(() => {
  const count = result.value?.customers.length ?? 0
  if (!count) return ''

  const { Page, Size } = result.value.meta
  const start = (Page - 1) * Size + 1

  return `${start}〜${start + count - 1} 件を表示`
})

const dateFormatter = new Intl.DateTimeFormat('ja-JP', {
  dateStyle: 'medium',
  timeZone: 'Asia/Tokyo'
})

function formatRegisteredDate(value) {
  if (!value) return '—'

  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '—' : dateFormatter.format(date)
}
</script>

<template>
  <section class="container">
    <div class="mb-2 flex flex-wrap items-center justify-between gap-3">
      <h1 class="mb-0 text-xl">顧客一覧</h1>
      <NuxtLink
        :to="{ path: '/dashboard/customers/create', query: customerPage(page).query }"
        class="btn inline-flex items-center gap-2 bg-green-600 text-white border-green-600 hover:text-white hover:bg-green-700 hover:border-green-700"
      >
        <Icon name="tabler:user-plus" size="18" aria-hidden="true" />
        顧客を追加
      </NuxtLink>
    </div>

    <nav class="mb-6 text-sm text-gray-600" aria-label="パンくずリスト">
      <NuxtLink to="/dashboard" class="text-green-600">ダッシュボード</NuxtLink>
      <span aria-hidden="true"> / </span>
      <span aria-current="page">顧客一覧</span>
    </nav>

    <p
      v-if="customerNotice"
      class="mb-4 rounded-lg bg-gray-100 p-4 text-sm text-gray-700"
      role="status"
    >
      {{ customerNotice }}
    </p>

    <div class="card card-lg">
      <form class="flex flex-wrap gap-2 p-4 md:p-6" role="search" @submit.prevent="submitSearch">
        <label for="customer-search" class="sr-only">名前、メールアドレス、住所で検索</label>
        <input
          id="customer-search"
          v-model="searchInput"
          type="search"
          placeholder="名前、メールアドレス、住所"
          class="form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 w-full md:w-1/2 text-base"
        />
        <button
          type="submit"
          class="btn bg-green-600 text-white border-green-600 hover:text-white hover:bg-green-700 hover:border-green-700"
        >
          検索
        </button>
      </form>

      <div class="overflow-x-auto">
        <table
          class="w-full whitespace-nowrap text-left table-hover"
          :aria-busy="status === 'pending'"
        >
          <thead class="bg-gray-200 text-gray-700">
            <tr>
              <th
                v-for="column in columns"
                :key="column.key"
                scope="col"
                class="px-6 py-3"
                :aria-sort="sortDirection(column.key)"
              >
                <NuxtLink
                  v-if="column.sortable"
                  :to="sortLink(column.key)"
                  aria-current-value="false"
                  class="inline-flex items-center gap-2 text-inherit hover:text-green-600"
                >
                  {{ column.label }}
                  <Icon :name="sortIcon(column.key)" size="16" aria-hidden="true" />
                </NuxtLink>
                <template v-else>{{ column.label }}</template>
              </th>
            </tr>
          </thead>

          <tbody>
            <tr v-if="status === 'idle' || status === 'pending'">
              <td :colspan="columns.length" class="px-6 py-4">
                <p class="mb-0" role="status">顧客一覧を読み込んでいます…</p>
              </td>
            </tr>

            <tr v-else-if="error">
              <td :colspan="columns.length" class="px-6 py-4">
                <p class="text-red-600" role="alert">{{ errorMessage }}</p>
                <NuxtLink v-if="errorCode === 401" to="/auth/signin" class="text-green-600">
                  ログイン画面へ
                </NuxtLink>
                <button
                  v-else-if="errorCode !== 403"
                  type="button"
                  class="btn btn-sm border-gray-300"
                  @click="refresh()"
                >
                  もう一度読み込む
                </button>
              </td>
            </tr>

            <tr v-else-if="!result?.customers.length">
              <td :colspan="columns.length" class="px-6 py-4">
                <p class="mb-0" role="status">顧客が見つかりませんでした。</p>
              </td>
            </tr>

            <template v-else>
              <tr v-for="customer in result.customers" :key="customer.id" class="border-b border-gray-200">
                <td class="px-6 py-3">{{ customer.id }}</td>
                <th scope="row" class="px-6 py-3 font-medium">
                  <div class="flex items-center gap-3">
                    <span class="flex w-8 h-8 items-center justify-center rounded-full bg-green-50 text-green-600">
                      <Icon name="tabler:user" size="16" aria-hidden="true" />
                    </span>
                    {{ customer.name }}
                  </div>
                </th>
                <td class="px-6 py-3">{{ customer.email }}</td>
                <td class="px-6 py-3">{{ customer.phone || '未登録' }}</td>
                <td class="px-6 py-3">{{ formatRegisteredDate(customer.created_at) }}</td>
                <td class="px-6 py-3">
                  <NuxtLink
                    :to="{ path: `/dashboard/customers/${customer.id}`, query: customerPage(page).query }"
                    :aria-label="`${customer.name}の詳細を見る`"
                    class="btn btn-sm border-gray-300 text-gray-700 hover:text-green-600"
                  >
                    詳細を見る
                  </NuxtLink>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <div
        v-if="status === 'success' && result"
        class="flex flex-wrap items-center justify-between gap-3 border-t border-gray-200 px-6 py-4"
      >
        <p class="mb-0 text-sm text-gray-600">
          <span v-if="rowRange">{{ rowRange }} / </span>全 {{ result.meta.Total }} 件
        </p>

        <nav
          v-if="result.meta.TotalPages > 1 || page > 1"
          class="flex flex-wrap items-center gap-2"
          aria-label="ページ切り替え"
        >
          <NuxtLink
            v-if="page > 1"
            :to="customerPage(Math.min(page - 1, Math.max(result.meta.TotalPages, 1)))"
            aria-current-value="false"
            class="btn btn-sm border-gray-300 text-gray-700"
          >
            前へ
          </NuxtLink>

          <NuxtLink
            v-for="number in pageNumbers"
            :key="number"
            :to="customerPage(number)"
            :aria-label="`${number}ページ`"
            :aria-current-value="number === page ? 'page' : 'false'"
            class="btn btn-sm"
            :class="number === page ? 'bg-green-600 text-white border-green-600' : 'border-gray-300 text-gray-700'"
          >
            {{ number }}
          </NuxtLink>

          <NuxtLink
            v-if="result.meta.HasNext"
            :to="customerPage(page + 1)"
            aria-current-value="false"
            class="btn btn-sm border-gray-300 text-gray-700"
          >
            次へ
          </NuxtLink>
        </nav>
      </div>
    </div>
  </section>
</template>
