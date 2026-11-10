<script setup>
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
  adminOnly: true
})

useHead({ title: '管理ダッシュボード | GreenCart' })
const authStore = useAuthStore()
</script>

<template>
  <section class="container">
    <div class="card bg-gray-100 border-0 rounded-2xl dashboard-welcome">
      <div class="card-body lg:p-12">
        <h1 class="mb-2 text-xl">管理ダッシュボード</h1>
        <p class="mb-6">
          {{ authStore.user?.name || '管理者' }}さん、おかえりなさい。
        </p>
        <NuxtLink
          to="/"
          class="btn inline-flex items-center bg-green-600 text-white border-green-600 hover:text-white hover:bg-green-700 hover:border-green-700"
        >
          ストアを見る
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.dashboard-welcome {
  background-image: url('/images/slider/slider-image-1.jpg');
  background-position: right center;
  background-repeat: no-repeat;
  background-size: cover;
}
</style>