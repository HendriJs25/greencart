<template>
    <footer class="bg-gray-200 py-8">
      <div class="container">
        <!-- Categories and Links Section -->
        <div class="flex flex-wrap md:gap-4 lg:gap-0 py-4 mb-6">
          <!-- Categories Column -->
          <div class="w-full md:w-full lg:w-1/3 flex flex-col gap-4 mb-6">
            <h6>カテゴリー</h6>
            <div class="flex flex-wrap">
              <div class="w-1/2">
                <ul class="flex flex-col gap-2">
                  <li v-for="(category, index) in categoriesLeft" :key="index">
                    <NuxtLink :to="'/shop'" class="inline-block hover:text-green-600">
                      {{ category }}
                    </NuxtLink>
                  </li>
                </ul>
              </div>
              <div class="w-1/2">
                <ul class="flex flex-col gap-2">
                  <li v-for="(category, index) in categoriesRight" :key="index">
                    <NuxtLink :to="'/shop'" class="inline-block hover:text-green-600">
                      {{ category }}
                    </NuxtLink>
                  </li>
                </ul>
              </div>
            </div>
          </div>
  
          <!-- Other Links Columns -->
          <div class="w-full md:w-full lg:w-2/3">
            <div class="flex flex-wrap">
              <div v-for="(section, index) in footerSections" :key="index" 
                   class="w-1/2 sm:w-1/2 md:w-1/4 flex flex-col gap-4 mb-6">
                <h6>{{ section.title }}</h6>
                <ul class="flex flex-col gap-2">
                  <li v-for="(link, linkIndex) in section.links" :key="linkIndex">
                    <NuxtLink :to="link.url || '#'" class="inline-block hover:text-green-600">
                      {{ link.text }}
                    </NuxtLink>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
  
        <!-- Payment Partners Section -->
        <div class="border-t py-4 border-gray-300">
          <div class="gap-y-4 flex flex-wrap items-center justify-center lg:justify-start">
            <div class="lg:w-2/5 lg:text-left text-center">
              <div class="flex md:flex-row flex-col gap-3 md:gap-6 items-center">
                <div class="text-gray-900">お支払い方法</div>
                <ul class="flex items-center flex-row gap-4">
                  <li v-for="(payment, index) in paymentMethods" :key="index">
                    <img :src="`/images/payment/${payment}.svg`" :alt="payment" class="h-6" />
                  </li>
                </ul>
              </div>
            </div>
            <div class="lg:w-3/5 flex justify-end">
              <div class="flex flex-col md:flex-row items-center gap-3 md:gap-6">
                <div class="text-gray-900">GreenCartのアプリで便利にお買い物</div>
                <ul class="flex flex-row gap-2">
                  <li>
                    <img src="~/assets/images/appbutton/appstore-btn.svg" alt="App Store" class="h-10" />
                  </li>
                  <li>
                    <img src="~/assets/images/appbutton/googleplay-btn.svg" alt="Google Play" class="h-10" />
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
  
        <!-- Copyright Section -->
        <div class="border-t py-4 border-gray-300">
          <div class="flex flex-col md:flex-row items-center gap-3">
            <div class="md:w-1/2">
              <span class="text-sm text-gray-500">
                © {{ new Date().getFullYear() }} GreenCart ストア。
              </span>
            </div>
            <div class="md:w-1/2 flex md:justify-end items-center">
              <div class="flex flex-row gap-5 items-center">
                <div class="text-gray-500">公式SNS</div>
                <ul class="flex items-center justify-end text-sm gap-1">
                  <li v-for="(social, index) in socialMedia" :key="index">
                    <a :href="social.url" class="inline-flex justify-center items-center align-middle text-center select-none border font-normal whitespace-no-wrap rounded leading-normal no-underline h-8 w-8 border-gray-300 hover:border-green-600 hover:text-green-600 transition ease-in-out">
                      <Icon :name="social.icon" size="16" />
                    </a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  </template>
  
  <script setup>
  const categoriesLeft = [
    '野菜・果物',
    '朝食・インスタント食品',
    'パン・ビスケット',
    '小麦粉・お米・豆類',
    'ソース・スプレッド',
    'オーガニック・グルメ',
    'ベビー用品',
    '掃除用品',
    'パーソナルケア'
  ]
  
  const categoriesRight = [
    '乳製品・パン・卵',
    '冷たい飲み物・ジュース',
    'お茶・コーヒー・飲料',
    'スパイス・食用油',
    '鶏肉・肉類・魚介',
    'その他の食品',
    'ヘルスケア',
    '生活・オフィス用品',
    'ペット用品'
  ]
  
  const footerSections = [
    {
      title: 'GreenCartについて',
      links: [
        { text: '会社概要', url: '/about' },
        { text: '私たちについて', url: '/about' },
        { text: 'ブログ', url: '/blog' },
        { text: 'ヘルプセンター', url: '/help' },
        { text: '大切にしていること', url: '/values' }
      ]
    },
    {
      title: 'お客様向け情報',
      links: [
        { text: 'お支払い', url: '/payments' },
        { text: '配送について', url: '/shipping' },
        { text: '返品について', url: '/returns' },
        { text: 'よくあるご質問', url: '/faq' },
        { text: 'ご購入手続き', url: '/checkout' }
      ]
    },
    {
      title: 'パートナー向け情報',
      links: [
        { text: 'パートナー募集', url: '/shopper' },
        { text: 'パートナー向け情報', url: '/shopper/register' },
        { text: '収益について', url: '/shopper/earnings' },
        { text: 'ご利用ガイド', url: '/guides' },
        { text: '新しい販売店', url: '/retailers' }
      ]
    },
    {
      title: 'GreenCartプログラム',
      links: [
        { text: 'GreenCartプログラム', url: '/programs' },
        { text: 'ギフトカード', url: '/gift-cards' },
        { text: 'キャンペーン・クーポン', url: '/coupons' },
        { text: '広告について', url: '/ads' },
        { text: '採用情報', url: '/careers' }
      ]
    }
  ]
  
  const paymentMethods = [
    'amazonpay',
    'american-express',
    'mastercard',
    'paypal',
    'visa'
  ]
  
  const socialMedia = [
    { icon: 'tabler:brand-facebook', url: '#' },
    { icon: 'tabler:brand-x', url: '#' },
    { icon: 'tabler:brand-instagram', url: '#' }
  ]
  </script>