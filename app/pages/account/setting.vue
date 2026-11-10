<template>
    <main>
        <!-- section -->
        <section>
          <div class="container">
            <div class="grid grid-cols-12">
                <div class="col-span-12">
                    <div class="flex justify-between items-center md:hidden py-4">
                        <!-- heading -->
                        <h3 class="text-md mb-0">Account Setting</h3>
                        <!-- button -->
                        <button
                            class="btn inline-flex items-center gap-x-2 bg-transparent text-gray-600 border-gray-300 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300 md:hidden btn-sm ms-3"
                            type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasAccount"
                            aria-controls="offcanvasAccount">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                stroke-linejoin="round"
                                class="icon icon-tabler icons-tabler-outline icon-tabler-indent-increase">
                                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                                <path d="M20 6l-11 0" />
                                <path d="M20 12l-7 0" />
                                <path d="M20 18l-11 0" />
                                <path d="M4 8l4 4l-4 4" />
                            </svg>
                        </button>
                    </div>
                </div>
                <div class="lg:col-span-2 md:col-span-4 col-span-12 border-r hidden md:block">
                    <div class="pt-10 md:pr-10">
                        <!-- nav -->
                        <ul class="nav flex-col nav-pills nav-pills-dark">
                            <!-- nav item -->
                            <li class="nav-item">
                                <a class="nav-link active gap-2" href="/account">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18"
                                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                        stroke-linecap="round" stroke-linejoin="round"
                                        class="icon icon-tabler icons-tabler-outline icon-tabler-shopping-bag">
                                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                                        <path
                                            d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z" />
                                        <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
                                    </svg>
                                    Your Orders
                                </a>
                            </li>
                            <!-- nav item -->
                            <li class="nav-item">
                                <a class="nav-link gap-2" aria-current="page" href="/account/setting">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18"
                                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                        stroke-linecap="round" stroke-linejoin="round"
                                        class="icon icon-tabler icons-tabler-outline icon-tabler-settings">
                                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                                        <path
                                            d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z" />
                                        <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
                                    </svg>
                                    Settings
                                </a>
                            </li>
                            <!-- nav item -->
                            <li class="nav-item">
                                <hr class="pt-3 mt-5" />
                            </li>
                            <!-- nav item -->
                            <li  class="nav-item">
                                <button type="button" class="nav-link gap-2" @click="handleLogout">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18"
                                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                        stroke-linecap="round" stroke-linejoin="round"
                                        class="icon icon-tabler icons-tabler-outline icon-tabler-logout">
                                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                                        <path
                                            d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2" />
                                        <path d="M9 12h12l-3 -3" />
                                        <path d="M18 15l3 -3" />
                                    </svg>
                                    Log out
                                </button>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="lg:col-span-10 md:col-span-8 col-span-12">
					<div class="py-6 md:p-6 lg:p-10">
						<div class="mb-6">
							<!-- heading -->
							<h2 class="mb-0 text-lg">Account Setting</h2>
						</div>
						<div>
							<!-- heading -->
							<h5 class="mb-4">Account details</h5>
							<div class="grid grid-cols-12">
								<div class="col-span-5">
									<!-- form -->
									<form @submit.prevent="handleSubmit" enctype="multipart/form-data">
										<!-- input -->
										<div class="mb-3">
											<label class="inline-block text-gray-800 font-medium mb-2">Name</label>
											<input v-model="formData.name" type="text"
												class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
												placeholder="jitu chauhan" />
										</div>
										<!-- input -->
										<div class="mb-3">
											<label class="inline-block text-gray-800 font-medium mb-2">Email</label>
											<input
                                                v-model="formData.email"
                                                type="email"
                                                readonly
                                                class="form-control bg-gray-100"
                                            />
										</div>
										<!-- input -->
										<div class="mb-5">
											<label class="inline-block text-gray-800 font-medium mb-2">Phone</label>
											<input v-model="formData.phone" type="text"
												class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
												placeholder="Phone number" />
										</div>
										<div class="mb-5">
											<label class="inline-block text-gray-800 font-medium mb-2">Address</label>
											<textarea
                                                v-model="formData.address"
												class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
												placeholder="Phone number">
                                            </textarea>
										</div>
										<div class="grid grid-cols-2 gap-3 mb-4">
                                            <div>
                                                <label class="block mb-2">Latitude</label>
                                                <input
                                                    v-model.number="formData.lat"
                                                    type="number"
                                                    step="any"
                                                    min="-90"
                                                    max="90"
                                                    class="form-control"
                                                />
                                            </div>

                                            <div>
                                                <label class="block mb-2">Longitude</label>
                                                <input
                                                    v-model.number="formData.lng"
                                                    type="number"
                                                    step="any"
                                                    min="-180"
                                                    max="180"
                                                    class="form-control"
                                                />
                                            </div>
                                        </div>

                                        <button
                                            type="button"
                                            :disabled="isGettingLocation"
                                            @click="useCurrentLocation"
                                            class="text-green-600 mb-4"
                                        >
                                            {{ isGettingLocation
                                                ? 'Mengambil lokasi...'
                                                : 'Gunakan lokasi saya' }}
                                        </button>

                                        <p v-if="isProfileLoading" class="text-gray-500 mb-3">
                                            Memuat profil...
                                        </p>

                                        <p v-if="errorMessage" class="text-red-600 mb-3">
                                            {{ errorMessage }}
                                        </p>

                                        <p v-if="successMessage" class="text-green-600 mb-3">
                                            {{ successMessage }}
                                        </p>
                                        <!-- <div class="mb-5">
											<label class="inline-block text-gray-800 font-medium mb-2">Photo</label>
											<input 
                                                type="file" 
                                                @change="handleFileUpload"
                                                accept="image/*"
                                                class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                                                placeholder="Pilih foto profil" 
                                            />
										</div> -->
										<!-- button -->
										<div class="mb-3">
											<button :disabled="authStore.loading || isProfileLoading"  type="submit"
												class="btn inline-flex items-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 active:border-green-700 focus:outline-none focus:ring-4 focus:ring-green-100">
												<span v-if="authStore.loading">Loading...</span>
                                                <span v-else>Save Details</span>
											</button>
										</div>
									</form>
								</div>
							</div>
						</div>
						<hr class="my-10" />
						<div class="lg:pr-14">
							<!-- heading -->
							<h5 class="mb-4">Password</h5>
							<form class="grid grid-cols-1 lg:grid-cols-2 gap-x-4">
								<!-- input -->
								<div class="mb-3">
									<label class="inline-block text-gray-800 font-medium mb-2">New Password</label>
									<input type="password"
										class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
										placeholder="**********" />
								</div>
								<!-- input -->
								<div class="mb-3">
									<label class="inline-block text-gray-800 font-medium mb-2">Current
										Password</label>
									<input type="password"
										class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(10,173,10,.25)] focus:ring-green-600 focus:ring-0 focus:border-green-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
										placeholder="**********" />
								</div>
								<!-- input -->
								<div class="">
									<p class="mb-4">
										Can’t remember your current password?
										<a href="#" class="text-green-600">Reset your password.</a>
									</p>
									<a href="#"
										class="btn inline-flex items-center gap-x-2 bg-green-600 text-white border-green-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-green-700 hover:border-green-700 active:bg-green-700 active:border-green-700 focus:outline-none focus:ring-4 focus:ring-green-100">
										Save Password
									</a>
								</div>
							</form>
						</div>
						<hr class="my-10" />
						<div>
							<!-- heading -->
							<h5 class="mb-4">Delete Account</h5>
							<p class="mb-2">Would you like to delete your account?</p>
							<p class="mb-5">This account contain 12 orders, Deleting your account will remove all
								the order details associated with it.</p>
							<!-- btn -->
							<a href="#"
								class="btn inline-flex items-center gap-x-2 bg-transparent text-red-600 border-red-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-red-600 hover:border-red-600 active:bg-red-600 active:border-red-600 focus:outline-none focus:ring-4 focus:ring-red-100">
								I want to delete my account
							</a>
						</div>
					</div>
				</div>
            </div>
          </div>
        </section>
    </main>

</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

definePageMeta({
    middleware: ['auth'],
    requiresAuth: true
})
const router = useRouter()
const authStore = useAuthStore()

const latitude = ref(null)
const longitude = ref(null)
const error = ref(null)

const formData = reactive({
    name: '',
    email: '',
    phone: '',
    address: '',
    lat: null,
    lng: null
})

const errorMessage = ref('')
const successMessage = ref('')
const isProfileLoading = ref(true)
const isGettingLocation = ref(false)

// Backend otomatis menambahkan prefix +81.
// Kirim hanya angka nasional tanpa angka 0 awal.
const normalizePhone = (phone) => {
    let digits = String(phone)
        .trim()
        .replace(/[\s()-]/g, '')

    if (digits.startsWith('+81')) {
        digits = digits.slice(3)
    } else if (/^81\d{9,10}$/.test(digits)) {
        digits = digits.slice(2)
    }

    if (digits.startsWith('0')) {
        digits = digits.slice(1)
    }

    if (!/^\d{9,10}$/.test(digits)) {
        throw new Error('Nomor telepon Jepang tidak valid')
    }

    return Number(digits)
}

const parseCoordinate = (value, label, min, max) => {
    if (value === null || value === '') {
        throw new Error(`${label} wajib diisi`)
    }

    const number = Number(value)

    if (!Number.isFinite(number) ||
        number < min || number > max) {
        throw new Error(`${label} tidak valid`)
    }

    return number
}

const useCurrentLocation = () => {
    errorMessage.value = ''

    if (!navigator.geolocation) {
        errorMessage.value = 'Geolocation tidak didukung'
        return
    }

    isGettingLocation.value = true

    navigator.geolocation.getCurrentPosition(
        (position) => {
            formData.lat = position.coords.latitude
            formData.lng = position.coords.longitude
            isGettingLocation.value = false
        },
        (error) => {
            errorMessage.value = error.message
            isGettingLocation.value = false
        },
        {
            enableHighAccuracy: false,
            timeout: 10000
        }
    )
}

const handleFileUpload = (event) => {
    const file = event.target.files[0]
    if (file) {
        // Validasi tipe file
        if (!file.type.includes('image/')) {
            alert('File harus berupa gambar')
            return
        }
        
        // Validasi ukuran file (maksimal 2MB)
        if (file.size > 2 * 1024 * 1024) {
            alert('Ukuran file maksimal 2MB')
            return
        }

        // Buat URL untuk preview (opsional)
        formData.value.photo = file
    }
}

const handleLogout = () => {
    try {
        authStore.logout()
        router.push('/auth/signin')
    } catch (error) {
        console.error('Error during logout:', error)
    }
}

const handleSubmit = async () => {
    errorMessage.value = ''
    successMessage.value = ''

    try {
        if (!formData.name.trim()) {
            throw new Error('Nama wajib diisi')
        }

        if (!formData.address.trim()) {
            throw new Error('Alamat wajib diisi')
        }

        const payload = {
            name: formData.name.trim(),
            phone: normalizePhone(formData.phone),
            address: formData.address.trim(),
            lat: parseCoordinate(
                formData.lat, 'Latitude', -90, 90
            ),
            lng: parseCoordinate(
                formData.lng, 'Longitude', -180, 180
            )
        }

        await authStore.updateProfile(payload)

        // Sinkronkan kembali nilai form
        const profile = authStore.profile

        if (profile) {
            formData.phone = profile.phone ?? ''
            formData.name = profile.name ?? ''
            formData.address = profile.address ?? ''
        }

        successMessage.value =
            'Profil berhasil diperbarui'

    } catch (error) {
        errorMessage.value = error.message
    }
}


onMounted(async () => {
    isProfileLoading.value = true
    errorMessage.value = ''

    try {
        await authStore.getProfile()

        const profile = authStore.profile

        if (!profile) {
            throw new Error('Data profil tidak ditemukan')
        }

        Object.assign(formData, {
            name: profile.name ?? '',
            email: profile.email ?? '',
            phone: profile.phone ?? '',
            address: profile.address ?? '',
            lat: profile.lat ?? null,
            lng: profile.lng ?? null
        })

    } catch (error) {
        errorMessage.value = error.message

        if (!authStore.isAuthenticated) {
            await navigateTo('/auth/signin')
        }
    } finally {
        isProfileLoading.value = false
    }
})

</script>