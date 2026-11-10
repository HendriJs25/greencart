<template>
  <header>
    <div class="border-b">
      <div class="pt-5">
        <div class="container">
          <div class="flex flex-wrap w-full items-center justify-between">
            <!-- Logo -->
            <div class="lg:w-1/6 md:w-1/2 w-2/5">
              <NuxtLink to="/" class="navbar-brand">
                <img src="/images/logo/freshcart-logo.svg" alt="FreshCart" />
              </NuxtLink>
            </div>

            <!-- Search Bar -->
            <div class="lg:w-2/5 hidden lg:block">
              <form @submit.prevent="handleSearch">
                <div class="relative">
                  <label for="searchProducts" class="invisible hidden">商品検索</label>
                  <input
                    v-model="searchQuery"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    type="search"
                    placeholder="商品を検索"
                    id="searchProducts"
                  />
                  <button class="absolute right-0 top-0 p-3" type="submit">
                    <Icon name="tabler:search" mode="svg" size="16" />
                  </button>
                </div>
              </form>
            </div>

            <!-- User Actions -->
            <div class="lg:w-1/5 text-end md:w-1/2 w-3/5">
              <div class="flex gap-7 items-center justify-end">
                <!-- User Menu -->
                <div>
                  <button @click="navigateTo(authStore.isAuthenticated ? '/account' : '/auth/signin')" class="text-gray-600">
                    <Icon name="tabler:user" mode="svg" size="22" />
                  </button>
                </div>

                <!-- Cart -->
                <div>
                  <button
                    type="button"
                    class="text-gray-600 relative"
                    @click="showCart = true"
                  >
                    <Icon name="tabler:shopping-bag" size="24" />
                    <span v-if="cartCount > 0"
                      class="absolute top-0 -mt-1 left-full rounded-full h-5 w-5 -ml-3 bg-green-600 text-white text-center font-semibold text-sm"
                    >
                      {{ cartCount }}
                      <span class="invisible">カート内の商品数</span>
                    </span>
                  </button>
                </div>

                <!-- Mobile Menu -->
                <div class="lg:hidden leading-none">
                  <button
                    @click="showMobileMenu = true"
                    class="collapsed"
                    type="button"
                    aria-label="メニューを開く"
                  >
                    <Icon name="tabler:menu-2" class="text-gray-800" size="24" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav
        class="navbar relative navbar-expand-lg lg:flex lg:flex-wrap items-center content-between text-black navbar-default"
        aria-label="Main navigation"
      >
        <div class="container max-w-7xl mx-auto w-full xl:px-4 lg:px-0">
          <div class="hidden lg:flex lg:items-center">
            <!-- Category Dropdown -->
            <div class="dropdown">
              <button
                class="mr-4 btn inline-flex items-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 active:border-green-700 focus:outline-none focus:ring-4 focus:ring-green-300"
                @click="toggleCategoryMenu"
              >
                <Icon name="tabler:layout-grid" size="16" />
                All Category
              </button>
              <ul v-if="showCategoryMenu" class="dropdown-menu">
                <li v-for="category in categories" :key="category.slug">
                  <NuxtLink :to="'/shop'" class="dropdown-item">
                    {{ category.name }}
                  </NuxtLink>
                </li>
              </ul>
            </div>

            <!-- Main Menu -->
            <ul class="navbar-nav lg:flex gap-3 lg:items-center">
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink to="/" class="nav-link">ホーム</NuxtLink>
              </li>
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink to="/shop" class="nav-link">商品一覧</NuxtLink>
              </li>
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink v-if="authStore.isAuthenticated" to="/account" class="nav-link">マイアカウント</NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </div>
  </header>

  <div v-if="showMobileMenu" class="lg:hidden border-b bg-white px-5 py-4">
    <nav class="flex flex-col gap-3">
      <NuxtLink to="/" @click="showMobileMenu = false">ホーム</NuxtLink>
      <NuxtLink to="/shop" @click="showMobileMenu = false">商品一覧</NuxtLink>
      <NuxtLink :to="authStore.isAuthenticated ? '/account' : '/auth/signin'" @click="showMobileMenu = false">マイアカウント</NuxtLink>
    </nav>
  </div>
  <!-- Login Modal -->
  

  <!-- Cart Sidebar -->
  <CartSidebar v-if="showCart" @close="showCart = false" />

  <!-- Mobile Menu -->
  <!-- <TheMobileMenu v-if="showMobileMenu" @close="showMobileMenu = false" /> -->
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import CartSidebar from '~/components/modals/CartSidebar.vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const cartItems = useState('greencart-cart-items', () => [])
const cartCount = computed(() => cartItems.value.reduce((sum, item) => sum + Number(item.quantity || 0), 0))
const searchQuery = ref('')
const showCart = ref(false)
const showMobileMenu = ref(false)
const showCategoryMenu = ref(false)
const categories = ref([])
const config = useRuntimeConfig()
const productApiBase = (config.public.apiProductBaseUrl || '/api/v1').replace(/\/$/, '')

function handleSearch() {
  if (searchQuery.value.trim()) {
    navigateTo({ path: '/shop', query: { search: searchQuery.value.trim() } })
  }
}

function toggleCategoryMenu() {
  showCategoryMenu.value = !showCategoryMenu.value
}

onMounted(async () => {
  // A failed /me request should never prevent rendering the header.
  try { await authStore.checkAuth() } catch (error) { console.warn('Session check failed', error) }
  try {
    const response = await $fetch(`${productApiBase}/categories/home`)
    categories.value = Array.isArray(response?.data) ? response.data : []
  } catch (error) {
    console.warn('Categories are not yet available', error)
  }
})
</script>