
<script setup>
import { ref, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

const { data: profile, error, pending, refresh } =
  await useAsyncData('customer-profile', () => {
    return authStore.getProfile()
  })

const name = ref('')
const saving = ref(false)
const saveError = ref('')
const success = ref(false)

watch(profile, (value) => {
  name.value = value?.name ?? ''
}, { immediate: true })

async function saveProfile() {
  saving.value = true
  saveError.value = ''
  success.value = false

  try {
    await authStore.updateProfile({
      name: name.value
    })

    await refresh()
    success.value = true
  } catch (err) {
    saveError.value =
      err.data?.message || err.message || 'Gagal menyimpan profil'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="container">
    <div class="grid grid-cols-12">
      <AccountSidebar />
      <div class="lg:col-span-10 md:col-span-8 col-span-12">
  <main class="account-page">
    <h1>Profil Customer</h1>

    <p v-if="pending">Memuat profil...</p>

    <div v-else-if="error">
      <p role="alert">Gagal mengambil profil customer.</p>
      <button @click="refresh()">Coba Lagi</button>
    </div>

    <form v-else-if="profile" @submit.prevent="saveProfile">
      <label for="name">Nama</label>
      <input id="name" v-model.trim="name" required />

      <label for="email">Email</label>
      <input
        id="email"
        :value="profile.email || authStore.user?.email"
        readonly
      />

      <p v-if="saveError" role="alert">{{ saveError }}</p>
      <p v-if="success">Profil berhasil diperbarui.</p>

      <button :disabled="saving || !name">
        {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
      </button>
    </form>

    <p v-else>Data profil tidak tersedia.</p>
  </main>
      </div>
    </div>
  </section>
</template>

<style scoped>
.account-page {
  width: min(720px, 92%);
  margin: 48px auto;
  min-height: 400px;
}

form {
  display: grid;
  gap: 14px;
  margin-top: 24px;
  padding: 24px;
  border: 1px solid #ddd;
  border-radius: 12px;
}

input { padding: 12px; border: 1px solid #ccc; border-radius: 6px; }
button { padding: 12px; background: #00a000; color: white; border: 0; border-radius: 6px; cursor: pointer; }
button:disabled { opacity: 0.6; }
</style>
