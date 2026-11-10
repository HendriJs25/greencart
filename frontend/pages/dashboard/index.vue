<script setup>
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: false,
  middleware: 'auth',
  adminOnly: true
})

useHead({ title: '管理ダッシュボード | GreenCart' })

if (import.meta.server) {
  useResponseHeader('Cache-Control').value = 'private, no-store'
}

const authStore = useAuthStore()
</script>

<template>
  <div>
    <header class="navbar navbar-glass">
      <div class="flex w-full h-full items-center justify-between gap-3">
        <NuxtLink to="/dashboard" class="xl:hidden font-semibold text-green-600">
          GreenCart
        </NuxtLink>
        <span class="hidden xl:block font-semibold">管理画面</span>

        <NuxtLink to="/account" class="text-sm text-green-600">
          マイアカウント
        </NuxtLink>
      </div>
    </header>

    <div class="main-wrapper">
      <nav class="navbar-vertical-nav hidden xl:block" aria-label="管理メニュー">
        <div class="navbar-vertical">
          <div class="px-4 py-5">
            <NuxtLink to="/dashboard" class="navbar-brand">
              <img src="~/assets/images/logo/freshcart-logo.svg" alt="FreshCart" />
            </NuxtLink>
          </div>

          <ul class="navbar-nav flex-col">
            <li class="nav-item">
              <NuxtLink to="/dashboard" class="nav-link active gap-2" aria-current="page">
                <Icon name="tabler:layout-dashboard" size="18" aria-hidden="true" />
                ダッシュボード
              </NuxtLink>
            </li>

            <li class="nav-item">
              <NuxtLink to="/" class="nav-link gap-2">
                <Icon name="tabler:building-store" size="18" aria-hidden="true" />
                ストアを見る
              </NuxtLink>
            </li>
          </ul>
        </div>
      </nav>

      <main class="main-content-wrapper">
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
      </main>
    </div>
  </div>
</template>

<style scoped>
.dashboard-welcome {
  background-image: url('/images/slider/slider-image-1.jpg');
  background-position: right center;
  background-repeat: no-repeat;
  background-size: cover;
}
</style>