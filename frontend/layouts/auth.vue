
<script setup>
import { computed } from 'vue'
import Footer from '~/components/common/Footer.vue'

const route = useRoute()

const isSignin = computed(() => route.path === '/auth/signin')
</script>

<template>
  <div>
    <main>
      <div class="border-b shadow-sm">
        <nav class="relative md:flex flex-wrap w-full items-center py-3 text-center">
          <div class="container sm:flex justify-between items-center">
            <NuxtLink
              to="/"
              class="inline-block mr-4 text-lg whitespace-no-wrap"
            >
              <img
                src="~/assets/images/logo/freshcart-logo.svg"
                alt="FreshCart"
                class="inline-block align-text-top"
              />
            </NuxtLink>

            <span v-if="!isSignin" class="inline-block">
              すでにアカウントをお持ちですか？
              <NuxtLink to="/auth/signin" class="text-green-600">
                ログイン
              </NuxtLink>
            </span>
          </div>
        </nav>
      </div>

      <slot />
    </main>

    <Footer />
  </div>
</template>
