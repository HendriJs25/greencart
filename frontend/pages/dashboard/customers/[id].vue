<script setup>
import { computed, ref, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
  adminOnly: true,
  validate: route => typeof route.params.id === 'string' && /^[1-9]\d{0,18}$/.test(route.params.id)
})

useHead({ title: '顧客詳細 | GreenCart' })

const route = useRoute()
const authStore = useAuthStore()
const listLocation = computed(() => ({
  path: '/dashboard/customers',
  query: route.query
}))

const { data: customer, status, error, refresh, clear } = await useFetch(
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

const deleting = ref(false)
const deleteError = ref('')
const deleteNeedsSignin = ref(false)
const deleteNotice = useState('admin-customer-notice', () => '')
const deleteDialogOpen = ref(false)
const deleteTrigger = ref(null)
const deleteDialog = ref(null)
const cancelDeleteButton = ref(null)

function openDeleteDialog() {
  deleteError.value = ''
  deleteNeedsSignin.value = false
  deleteDialogOpen.value = true
}

function closeDeleteDialog() {
  if (!deleting.value) deleteDialogOpen.value = false
}

watch(deleteDialogOpen, (open, _, onCleanup) => {
  if (!open) {
    deleteTrigger.value?.focus()
    return
  }

  const appRoot = deleteTrigger.value?.closest('#__nuxt')
  const wasInert = appRoot?.inert ?? false
  const previousOverflow = document.body.style.overflow

  if (appRoot) appRoot.inert = true
  document.body.style.overflow = 'hidden'

  onCleanup(() => {
    if (appRoot) appRoot.inert = wasInert
    document.body.style.overflow = previousOverflow
  })

  cancelDeleteButton.value?.focus()
}, { flush: 'post' })

watch(deleting, (pending) => {
  if (deleteDialogOpen.value) {
    const target = pending ? deleteDialog.value : cancelDeleteButton.value
    target?.focus()
  }
}, { flush: 'post' })

function keepFocusInside(event) {
  const panel = deleteDialog.value
  const buttons = panel.querySelectorAll('button:not([disabled]), a[href]')
  const first = buttons[0]
  const last = buttons[buttons.length - 1]
  const focused = document.activeElement

  if (!first) {
    event.preventDefault()
  } else if (focused === panel || focused === (event.shiftKey ? first : last)) {
    event.preventDefault()
    const target = event.shiftKey ? last : first
    target.focus()
  }
}

async function handleDelete() {
  if (!deleteDialogOpen.value || !customer.value || deleting.value) return

  const id = route.params.id

  const query = { ...listLocation.value.query }
  delete query.page

  deleting.value = true
  deleteError.value = ''
  deleteNeedsSignin.value = false
  let message = '顧客を削除しました。'

  try {
    await authStore.deleteCustomer(id)
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode

    if (status === 404) {
      message = '顧客が見つかりませんでした。'
    } else {
      const backendMessage = (err?.data ?? err?.response?._data)?.message
      deleteNeedsSignin.value = status === 401 || status === 403
      deleteError.value = deleteNeedsSignin.value
        ? '顧客を削除するには、もう一度ログインしてください。'
        : typeof backendMessage === 'string' && backendMessage
          ? backendMessage
          : '顧客を削除できませんでした。もう一度お試しください。'
      return
    }
  } finally {
    deleting.value = false
  }

  deleteNotice.value = message
  closeDeleteDialog()
  clearNuxtData('admin-customers')
  clear()

  return navigateTo({ path: '/dashboard/customers', query }, { replace: true })
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
        <div class="flex flex-wrap items-center gap-2">
          <NuxtLink :to="listLocation" class="btn btn-sm border-gray-300 text-gray-700">
            一覧に戻る
          </NuxtLink>

          <button
            v-if="status === 'success' && customer"
            ref="deleteTrigger"
            type="button"
            class="btn btn-sm inline-flex items-center gap-2 text-red-600 border-red-600 hover:bg-red-600 hover:text-white disabled:opacity-50"
            :disabled="deleting"
            aria-haspopup="dialog"
            @click="openDeleteDialog"
          >
            <Icon name="tabler:trash" size="16" aria-hidden="true" />
            {{ deleting ? '削除中…' : '顧客を削除' }}
          </button>
        </div>

      </div>
    </div>

    <Teleport v-if="deleteDialogOpen" to="body">
      <div class="delete-overlay" @click.self="closeDeleteDialog">
        <section
          ref="deleteDialog"
          class="delete-dialog bg-white rounded-lg shadow-lg p-6"
          role="alertdialog"
          aria-modal="true"
          aria-labelledby="customer-delete-title"
          aria-describedby="customer-delete-description"
          tabindex="-1"
          @keydown.esc.stop.prevent="closeDeleteDialog"
          @keydown.tab="keepFocusInside"
        >
          <div class="flex items-center gap-2 mb-3">
            <Icon name="tabler:alert-triangle" size="24" class="text-red-600" aria-hidden="true" />
            <h2 id="customer-delete-title" class="mb-0 text-lg">顧客を削除しますか？</h2>
          </div>

          <p id="customer-delete-description" class="mb-4 customer-detail-value">
            {{ customer?.name }}（ID {{ route.params.id }}）を顧客一覧から削除します。
          </p>

          <p v-if="deleteError" class="mb-2 text-sm text-red-600" role="alert">
            {{ deleteError }}
          </p>
          <NuxtLink v-if="deleteNeedsSignin" to="/auth/signin" class="text-sm text-green-600">
            ログイン画面へ
          </NuxtLink>
          <p v-if="deleting" class="sr-only" role="status">削除中…</p>

          <div class="mt-4 flex flex-wrap justify-end gap-2">
            <button
              ref="cancelDeleteButton"
              type="button"
              class="btn border-gray-300 text-gray-700 hover:bg-gray-100 disabled:opacity-50"
              :disabled="deleting"
              @click="closeDeleteDialog"
            >
              キャンセル
            </button>
            <button
              type="button"
              class="btn bg-red-600 border-red-600 text-white hover:bg-red-700 hover:text-white disabled:opacity-50"
              :disabled="deleting"
              @click="handleDelete"
            >
              {{ deleting ? '削除中…' : '削除する' }}
            </button>
          </div>
        </section>
      </div>
    </Teleport>
  </section>
</template>

<style scoped>
.customer-detail-value {
  overflow-wrap: anywhere;
  white-space: pre-line;
}

.delete-overlay {
  position: fixed;
  inset: 0;
  z-index: 1100;
  display: grid;
  place-items: center;
  padding: 1rem;
  background: rgb(15 23 42 / 50%);
}

.delete-dialog {
  width: min(100%, 28rem);
  max-height: calc(100dvh - 2rem);
  overflow: auto;
}
</style>
