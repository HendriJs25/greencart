<script setup>
import { reactive, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

const props = defineProps({
  profile: { type: Object, required: true }
})

const emit = defineEmits(['saved'])
const authStore = useAuthStore()

const form = reactive({
  name: props.profile.name ?? '',
  phone: props.profile.phone ?? '',
  address: props.profile.address ?? '',
  lat: props.profile.lat ?? '',
  lng: props.profile.lng ?? ''
})

const fields = [
  { key: 'name', label: 'お名前', type: 'text', autocomplete: 'name' },
  { key: 'phone', label: '電話番号', type: 'tel', autocomplete: 'tel', placeholder: '090-1234-5678' },
  { key: 'address', label: '住所', autocomplete: 'street-address' },
  { key: 'lat', label: '緯度', type: 'number', min: -90, max: 90 },
  { key: 'lng', label: '経度', type: 'number', min: -180, max: 180 }
]

const inputClass = 'form-control border border-gray-300 text-gray-900 rounded-lg focus:ring-green-600 focus:border-green-600 block p-2 px-3 w-full text-base'
const saving = ref(false)
const locating = ref(false)
const fieldErrors = ref({})
const saveError = ref('')
const locationError = ref('')
const success = ref(false)
const needsSignin = ref(false)

function clearFieldError(field) {
  delete fieldErrors.value[field]
  success.value = false

  if (field === 'lat' || field === 'lng') {
    locationError.value = ''
  }
}

function normalizePhone(value) {
  const digits = value.normalize('NFKC')
    .replace(/[\s()-]/g, '')
    .replace(/^\+81/, '')
    .replace(/^0/, '')

  const phone = Number(digits)
  return /^[1-9]\d*$/.test(digits) && Number.isSafeInteger(phone) ? phone : null
}

function readCoordinate(value, min, max) {
  if (String(value ?? '').trim() === '') return null

  const number = Number(value)
  return Number.isFinite(number) && number >= min && number <= max ? number : null
}

function useCurrentLocation() {
  if (locating.value || saving.value) return

  locationError.value = ''

  if (!navigator.geolocation) {
    locationError.value = '現在地を取得できません。緯度と経度を入力してください。'
    return
  }

  locating.value = true

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      form.lat = coords.latitude
      form.lng = coords.longitude
      clearFieldError('lat')
      clearFieldError('lng')
      locating.value = false
    },
    () => {
      locationError.value = '現在地を取得できませんでした。緯度と経度を入力してください。'
      locating.value = false
    },
    { timeout: 10000 }
  )
}

async function saveProfile() {
  if (saving.value || locating.value) return

  saveError.value = ''
  needsSignin.value = false
  success.value = false

  const payload = {
    name: form.name.trim(),
    phone: normalizePhone(form.phone),
    address: form.address.trim(),
    lat: readCoordinate(form.lat, -90, 90),
    lng: readCoordinate(form.lng, -180, 180)
  }

  const errors = {}

  if (!payload.name) errors.name = 'お名前を入力してください。'
  if (payload.phone === null) errors.phone = '日本の電話番号を入力してください。'
  if (!payload.address) errors.address = '住所を入力してください。'
  if (payload.lat === null) errors.lat = '緯度を -90 から 90 の範囲で入力してください。'
  if (payload.lng === null) errors.lng = '経度を -180 から 180 の範囲で入力してください。'

  fieldErrors.value = errors
  if (Object.keys(errors).length) return

  saving.value = true

  try {
    await authStore.updateProfile(payload)

    const saved = { ...payload, phone: `+81${payload.phone}` }
    Object.assign(form, saved)
    emit('saved', saved)
    success.value = true
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode
    const message = (err?.data ?? err?.response?._data)?.message

    if (status === 422 && Array.isArray(message) && message.length) {
      fieldErrors.value = Object.fromEntries(
        message.map(({ field, message }) => [field, message])
      )
    } else if (status === 401 || status === 403) {
      needsSignin.value = true
      saveError.value = 'プロフィールを更新するには、もう一度ログインしてください。'
    } else {
      saveError.value = typeof message === 'string' && message
        ? message
        : 'プロフィールを更新できませんでした。もう一度お試しください。'
    }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <form
    class="mt-6 border-t border-gray-200 pt-4"
    novalidate
    :aria-busy="saving"
    @submit.prevent="saveProfile"
  >
    <h2 class="mb-4 text-lg">プロフィールを編集</h2>

    <fieldset :disabled="saving" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="field in fields"
        :key="field.key"
        :class="{ 'profile-address': field.key === 'address' }"
      >
        <label :for="`profile-${field.key}`" class="block mb-1 text-sm font-medium text-gray-800">
          {{ field.label }}
        </label>

        <textarea
          v-if="field.key === 'address'"
          :id="`profile-${field.key}`"
          v-model="form.address"
          :class="inputClass"
          :autocomplete="field.autocomplete"
          rows="2"
          required
          :aria-invalid="!!fieldErrors[field.key]"
          :aria-describedby="fieldErrors[field.key] ? `profile-${field.key}-error` : undefined"
          @input="clearFieldError(field.key)"
        />

        <input
          v-else
          :id="`profile-${field.key}`"
          v-model="form[field.key]"
          :class="inputClass"
          :type="field.type"
          :autocomplete="field.autocomplete"
          :placeholder="field.placeholder"
          :min="field.min"
          :max="field.max"
          :step="field.type === 'number' ? 'any' : undefined"
          :disabled="locating && field.type === 'number'"
          required
          :aria-invalid="!!fieldErrors[field.key]"
          :aria-describedby="fieldErrors[field.key] ? `profile-${field.key}-error` : undefined"
          @input="clearFieldError(field.key)"
        />

        <p
          v-if="fieldErrors[field.key]"
          :id="`profile-${field.key}-error`"
          class="mt-1 mb-0 text-sm text-red-600"
        >
          {{ fieldErrors[field.key] }}
        </p>
      </div>
    </fieldset>

    <p class="mt-2 mb-1 text-sm text-gray-600">
      住所に対応する緯度と経度を入力してください。
    </p>

    <button
      type="button"
      class="btn btn-sm bg-white text-gray-700 border-gray-300"
      :disabled="saving || locating"
      @click="useCurrentLocation"
    >
      {{ locating ? '現在地を取得中…' : '現在地を使う' }}
    </button>

    <p v-if="locationError" class="mt-2 mb-0 text-sm text-red-600" role="alert">
      {{ locationError }}
    </p>

    <p v-if="saveError" class="mt-2 mb-0 text-sm text-red-600" role="alert">
      {{ saveError }}
    </p>

    <NuxtLink v-if="needsSignin" to="/auth/signin" class="text-sm text-green-600">
      ログイン画面へ
    </NuxtLink>

    <p v-if="success" class="mt-2 mb-0 text-sm text-green-600" role="status">
      プロフィールを更新しました。
    </p>

    <div class="mt-4">
      <button
        type="submit"
        class="btn bg-green-600 text-white border-green-600 hover:bg-green-700 hover:border-green-700 disabled:opacity-50"
        :disabled="saving || locating"
      >
        {{ saving ? '保存中…' : '変更を保存' }}
      </button>
    </div>
  </form>
</template>

<style scoped>
.profile-address {
  grid-column: 1 / -1;
}
</style>