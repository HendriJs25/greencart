<script setup>
import { computed, reactive, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
  adminOnly: true
})

useHead({ title: '顧客を追加 | GreenCart' })

const route = useRoute()
const authStore = useAuthStore()
const listLocation = computed(() => ({
  path: '/dashboard/customers',
  query: route.query
}))
const notice = useState('admin-customer-notice', () => '')

const form = reactive({
  name: '',
  email: '',
  password: '',
  password_confirmation: '',
  phone: '',
  address: '',
  lat: '',
  lng: ''
})

const coordinateFields = [
  { key: 'lat', label: '緯度', type: 'number', min: -90, max: 90 },
  { key: 'lng', label: '経度', type: 'number', min: -180, max: 180 }
]

const fields = [
  { key: 'name', label: 'お名前', type: 'text', required: true, autocomplete: 'name' },
  { key: 'email', label: 'メールアドレス', type: 'email', required: true, autocomplete: 'email' },
  { key: 'password', label: '初期パスワード', type: 'password', required: true, autocomplete: 'new-password' },
  { key: 'password_confirmation', label: '初期パスワード（確認）', type: 'password', required: true, autocomplete: 'new-password' },
  { key: 'phone', label: '電話番号', type: 'tel', autocomplete: 'tel', maxLength: 17, placeholder: '09012345678' },
  { key: 'address', label: '住所', type: 'text', autocomplete: 'street-address' },
  ...coordinateFields
]

const inputClass = 'form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 w-full text-base'
const visiblePasswords = reactive({ password: false, password_confirmation: false })
const saving = ref(false)
const created = ref(false)
const fieldErrors = ref({})
const saveError = ref('')
const needsSignin = ref(false)

function clearFieldError(field) {
  delete fieldErrors.value[field]

  if (field === 'password') {
    delete fieldErrors.value.password_confirmation
  }
}

async function handleSubmit() {
  if (saving.value || created.value) return

  saveError.value = ''
  needsSignin.value = false
  const errors = {}

  const payload = {
    name: form.name.trim(),
    email: form.email.trim(),
    password: form.password,
    password_confirmation: form.password_confirmation
  }

  if (!payload.name) errors.name = 'お名前を入力してください。'
  if (!payload.email) errors.email = 'メールアドレスを入力してください。'

  if (!payload.password) {
    errors.password = '初期パスワードを入力してください。'
  } else if (new TextEncoder().encode(payload.password).length > 72) {
    errors.password = 'パスワードは72バイト以内で入力してください。'
  }

  if (!payload.password_confirmation) {
    errors.password_confirmation = '確認用のパスワードを入力してください。'
  } else if (payload.password !== payload.password_confirmation) {
    errors.password_confirmation = 'パスワードが一致していません。'
  }

  const phone = form.phone.trim()
  if (phone) {
    if (/^[0-9]{1,17}$/.test(phone)) {
      payload.phone = phone
    } else {
      errors.phone = '電話番号は半角数字17桁以内で入力してください。'
    }
  }

  const address = form.address.trim()
  if (address) payload.address = address

  for (const field of coordinateFields) {
    const raw = String(form[field.key] ?? '').trim()
    if (!raw) continue

    const value = Number(raw)
    if (!Number.isFinite(value) || value < field.min || value > field.max) {
      errors[field.key] = `${field.label}は ${field.min} から ${field.max} の範囲で入力してください。`
    } else {
      payload[field.key] = value
    }
  }

  fieldErrors.value = errors
  if (Object.keys(errors).length) return

  saving.value = true
  let result

  try {
    result = await authStore.createCustomer(payload)
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode
    const message = (err?.data ?? err?.response?._data)?.message

    if (status === 422 && Array.isArray(message) && message.length) {
      fieldErrors.value = Object.fromEntries(message.map(({ field, message }) => [field, message]))
    } else if (status === 409) {
      fieldErrors.value = {
        email: typeof message === 'string' && message
          ? message
          : 'このメールアドレスは既に使用されています。'
      }
    } else if (status === 401 || status === 403) {
      needsSignin.value = true
      saveError.value = '顧客を登録するには、もう一度ログインしてください。'
    } else {
      saveError.value = typeof message === 'string' && message
        ? message
        : '顧客を登録できませんでした。もう一度お試しください。'
    }
    return
  } finally {
    saving.value = false
  }

  created.value = true
  form.password = ''
  form.password_confirmation = ''
  notice.value = typeof result.message === 'string' && result.message
    ? result.message
    : '顧客を登録しました。'
  clearNuxtData('admin-customers')

  const query = { ...listLocation.value.query }
  delete query.page

  return navigateTo({ path: '/dashboard/customers', query }, { replace: true })
}
</script>

<template>
  <section class="container">
    <h1 class="mb-2 text-xl">顧客を追加</h1>

    <nav class="mb-6 text-sm text-gray-600" aria-label="パンくずリスト">
      <NuxtLink to="/dashboard" class="text-green-600">ダッシュボード</NuxtLink>
      <span aria-hidden="true"> / </span>
      <NuxtLink :to="listLocation" class="text-green-600">顧客一覧</NuxtLink>
      <span aria-hidden="true"> / </span>
      <span aria-current="page">顧客を追加</span>
    </nav>

    <div class="card card-lg">
      <form class="p-4 md:p-6" novalidate :aria-busy="saving" @submit.prevent="handleSubmit">
        <h2 class="mb-4 text-lg">顧客情報</h2>

        <p v-if="saveError" class="mb-2 rounded-lg bg-red-50 p-4 text-sm text-red-800" role="alert">
          {{ saveError }}
        </p>
        <NuxtLink v-if="needsSignin" to="/auth/signin" class="inline-block mb-4 text-green-600">
          ログイン画面へ
        </NuxtLink>

        <fieldset :disabled="saving || created" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="field in fields" :key="field.key">
            <label :for="`customer-${field.key}`" class="block mb-1 text-sm font-medium text-gray-800">
              {{ field.label }}
              <span v-if="field.required" class="text-red-600" aria-hidden="true">*</span>
              <span v-else class="text-gray-500">（任意）</span>
            </label>

            <div :class="{ 'password-input-wrapper': field.type === 'password' }">
              <input
                :id="`customer-${field.key}`"
                v-model="form[field.key]"
                :type="visiblePasswords[field.key] ? 'text' : field.type"
                :autocomplete="field.autocomplete"
                :maxlength="field.maxLength"
                :placeholder="field.placeholder"
                :min="field.min"
                :max="field.max"
                :step="field.type === 'number' ? 'any' : undefined"
                :required="field.required"
                :class="inputClass"
                :aria-invalid="!!fieldErrors[field.key]"
                :aria-describedby="fieldErrors[field.key]
                  ? `customer-${field.key}-error`
                  : field.key === 'password' ? 'customer-password-hint' : undefined"
                @input="clearFieldError(field.key)"
              />
              <button
                v-if="field.type === 'password'"
                type="button"
                class="password-toggle"
                :aria-label="`${field.label}の表示`"
                :aria-pressed="visiblePasswords[field.key]"
                @click="visiblePasswords[field.key] = !visiblePasswords[field.key]"
              >
                <Icon :name="visiblePasswords[field.key] ? 'tabler:eye' : 'tabler:eye-off'" size="20" />
              </button>
            </div>

            <p
              v-if="fieldErrors[field.key]"
              :id="`customer-${field.key}-error`"
              class="mt-1 mb-0 text-sm text-red-600"
            >
              {{ fieldErrors[field.key] }}
            </p>
            <p v-else-if="field.key === 'password'" id="customer-password-hint" class="mt-1 mb-0 text-sm text-gray-600">
              8文字以上で、大文字・小文字・数字をそれぞれ1文字以上含めてください。
            </p>
          </div>
        </fieldset>

        <div class="mt-6 flex flex-wrap gap-2">
          <button
            type="submit"
            class="btn bg-green-600 text-white border-green-600 hover:text-white hover:bg-green-700 hover:border-green-700 disabled:opacity-50"
            :disabled="saving || created"
          >
            {{ saving ? '登録中…' : '顧客を登録' }}
          </button>
          <NuxtLink :to="listLocation" class="btn border-gray-300 text-gray-700">
            一覧に戻る
          </NuxtLink>
        </div>
      </form>
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
}

.password-toggle:focus-visible {
  outline: 2px solid #0aad0a;
  outline-offset: 2px;
}
</style>
