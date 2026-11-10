<script setup>
if (import.meta.server) {
  useResponseHeader('Cache-Control').value = 'private, no-store'
}
</script>

<template>
  <div>
    <header class="navbar navbar-glass">
      <div class="flex w-full h-full items-center justify-between gap-3">
        <NuxtLink to="/dashboard" class="xl:hidden font-semibold text-green-600">
          GreenCart
        </NuxtLink>
        <span class="hidden xl:block font-semibold">管理画面</span>

        <div class="flex items-center gap-3">
          <NuxtLink to="/dashboard/customers" class="xl:hidden text-sm text-green-600">
            顧客一覧
          </NuxtLink>
          <NuxtLink to="/account" class="text-sm text-green-600">
            マイアカウント
          </NuxtLink>
        </div>
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
              <NuxtLink to="/dashboard" class="nav-link gap-2" exact-active-class="active">
                <Icon name="tabler:layout-dashboard" size="18" aria-hidden="true" />
                ダッシュボード
              </NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink to="/dashboard/customers" class="nav-link gap-2" exact-active-class="active">
                <Icon name="tabler:users" size="18" aria-hidden="true" />
                顧客一覧
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
        <slot />
      </main>
    </div>
  </div>
</template>