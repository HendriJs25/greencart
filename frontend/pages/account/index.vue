<script setup>
import { computed, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth'
})

useHead({ title: 'プロフィール | GreenCart' })

if (import.meta.server) {
  useResponseHeader('Cache-Control').value = 'private, no-store'
}

const authStore = useAuthStore()
const profileKey = `customer-profile:${authStore.user?.uuid ?? ''}`

const { data: profile, status, error, refresh } = await useLazyAsyncData(
  profileKey,
  async () => {
    const data = await authStore.getProfile()

    if (!data) {
      throw createError({ statusCode: 502 })
    }

    return data
  },
  { pick: ['name', 'email', 'phone', 'address', 'photo', 'lat', 'lng'] }
)

const failedPhoto = ref('')
const sessionExpired = computed(() => error.value?.statusCode === 401)

const loadError = computed(() => {
  if (sessionExpired.value) {
    return 'ログインの有効期限が切れました。もう一度ログインしてください。'
  }

  const message = error.value?.data?.message
  return typeof message === 'string' && message
    ? message
    : 'プロフィールを読み込めませんでした。'
})

const photoInput = ref(null)
const uploading = ref(false)
const uploadError = ref('')
const uploadSuccess = ref(false)
const uploadNeedsSignin = ref(false)

async function handlePhotoChange(event) {
  const input = event.currentTarget
  const file = input.files?.[0]
  input.value = ''

  if (!file || uploading.value || !profile.value) return

  uploadError.value = ''
  uploadSuccess.value = false
  uploadNeedsSignin.value = false

  if (!/\.(jpe?g|png)$/i.test(file.name) || file.size === 0) {
    uploadError.value = 'JPGまたはPNGの画像ファイルを選択してください。'
    return
  }

  if (file.size > 2 * 1024 * 1024) {
    uploadError.value = '2MB以下の画像を選択してください。'
    return
  }

  uploading.value = true

  try {
    const photo = await authStore.uploadProfileImage(file)

    failedPhoto.value = ''
    profile.value = { ...profile.value, photo }
    uploadSuccess.value = true
  } catch (err) {
    const status = err?.response?.status ?? err?.statusCode
    const message = (err?.data ?? err?.response?._data)?.message

    if (status === 401 || status === 403) {
      uploadNeedsSignin.value = true
      uploadError.value = '写真を変更するには、もう一度ログインしてください。'
    } else {
      uploadError.value = typeof message === 'string' && message
        ? message
        : '写真をアップロードできませんでした。もう一度お試しください。'
    }
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <main class="container">
    <div class="grid grid-cols-12">
      <AccountSidebar />

      <section class="lg:col-span-10 md:col-span-8 col-span-12 py-6 md:p-6 lg:p-10">
        <h1 class="mb-6 text-xl">プロフィール</h1>

        <p v-if="status === 'pending' || status === 'idle'" role="status">
          プロフィールを読み込んでいます…
        </p>

        <div v-else-if="error" class="rounded-lg bg-red-50 p-4">
          <p class="mb-3 text-sm text-red-800" role="alert">
            {{ loadError }}
          </p>

          <NuxtLink v-if="sessionExpired" to="/auth/signin" class="text-green-600">
            ログイン画面へ
          </NuxtLink>

          <button
            v-else
            type="button"
            class="btn btn-sm bg-white text-gray-700 border-gray-300"
            @click="refresh()"
          >
            もう一度読み込む
          </button>
        </div>

        <div v-else-if="profile" class="card rounded-lg border border-gray-200 bg-white">
          <div class="p-4 md:p-6">
            <div class="flex flex-wrap items-center gap-4 border-b border-gray-200 pb-4">
              <div class="flex w-16 h-16 items-center justify-center overflow-hidden rounded-full bg-green-50 text-green-600">
                <img
                  v-if="profile.photo && failedPhoto !== profile.photo"
                  :src="profile.photo"
                  alt=""
                  referrerpolicy="no-referrer"
                  class="w-full h-full object-cover"
                  @error="failedPhoto = profile.photo"
                />
                <Icon v-else name="tabler:user" size="28" aria-hidden="true" />
              </div>

              <div class="profile-text min-w-0">
                <h2 class="mb-1 text-lg">{{ profile.name || 'お客様' }}</h2>
                <p class="mb-0 text-gray-600">
                  {{ profile.email || 'メールアドレス未登録' }}
                </p>
              </div>
            </div>

            <div class="pt-4" :aria-busy="uploading">
              <input
                ref="photoInput"
                type="file"
                accept=".jpg,.jpeg,.png,image/jpeg,image/png"
                class="hidden"
                :disabled="uploading"
                @change="handlePhotoChange"
              />

              <button
                type="button"
                class="btn btn-sm inline-flex items-center gap-2 bg-green-600 text-white border-green-600 hover:bg-green-700 hover:border-green-700 disabled:opacity-50 disabled:pointer-events-none"
                :disabled="uploading"
                aria-describedby="profile-photo-help"
                @click="photoInput?.click()"
              >
                <Icon name="tabler:camera" size="16" aria-hidden="true" />
                {{ uploading ? 'アップロード中…' : 'プロフィール写真を変更' }}
              </button>

              <p id="profile-photo-help" class="mt-2 mb-0 text-sm text-gray-600">
                JPG・PNG（最大2MB）
              </p>

              <p v-if="uploadError" class="mt-2 mb-0 text-sm text-red-600" role="alert">
                {{ uploadError }}
              </p>

              <NuxtLink v-if="uploadNeedsSignin" to="/auth/signin" class="text-sm text-green-600">
                ログイン画面へ
              </NuxtLink>

              <p v-if="uploadSuccess" class="mt-2 mb-0 text-sm text-green-600" role="status">
                プロフィール写真を更新しました。
              </p>
            </div>

            <AccountProfileForm
              :profile="profile"
              @saved="profile = { ...profile, ...$event }"
            />
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.profile-text {
  overflow-wrap: anywhere;
}
</style>