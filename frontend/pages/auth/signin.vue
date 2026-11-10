
<script setup>
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()

const email = ref('')
const password = ref('')

async function handleSignin() {
  try {
    await authStore.signin(email.value, password.value)

    await navigateTo('/account', {
      replace: true
    })
  } catch (err) {
    // Pesan kesalahan ditampilkan dari authStore.error
  }
}
</script>

<template>
  <main class="auth-page">
    <form class="auth-card" @submit.prevent="handleSignin">
      <h1>GreenCart</h1>
      <h2>Login Customer</h2>

      <label for="email">Email</label>
      <input
        id="email"
        v-model.trim="email"
        type="email"
        autocomplete="email"
        required
      />

      <label for="password">Password</label>
      <input
        id="password"
        v-model="password"
        type="password"
        autocomplete="current-password"
        required
      />

      <p v-if="authStore.error" role="alert" class="error">
        {{ authStore.error }}
      </p>

      <button :disabled="authStore.loading">
        {{ authStore.loading ? 'Memproses...' : 'Login' }}
      </button>
    </form>
  </main>
</template>

<style scoped>
.auth-page {
  min-height: 80vh;
  display: grid;
  place-items: center;
  background: #f8faf8;
}

.auth-card {
  width: min(420px, 92%);
  padding: 32px;
  border: 1px solid #ddd;
  border-radius: 12px;
  background: white;
  display: grid;
  gap: 14px;
}

.auth-card h1 { color: #00a000; }
.auth-card input { padding: 12px; border: 1px solid #ccc; border-radius: 6px; }
.auth-card button { padding: 12px; border: 0; border-radius: 6px; background: #00a000; color: white; cursor: pointer; }
.auth-card button:disabled { opacity: 0.6; }
.error { color: #c62828; }
</style>
