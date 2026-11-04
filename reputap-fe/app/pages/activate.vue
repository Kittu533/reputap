<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'

definePageMeta({
  layout: false,
})

const route = useRoute()
const config = useRuntimeConfig()
const apiBase = config.public.apiBase as string

const code = ref<string>((route.query.code as string || '').trim().toLowerCase())
const checking = ref(true)
const submitting = ref(false)
const errorMessage = ref('')
const checkResult = ref<{
  can_activate: boolean
  status: string
  message: string
  store_name?: string
  review_url?: string
} | null>(null)

const successData = ref<{
  store_name: string
  review_url: string
  redirect_url: string
} | null>(null)

const form = reactive({
  business_name: '',
  category: 'coffee_shop',
  review_url: '',
  owner_name: '',
  phone: '',
  city: '',
})

const categories = [
  { label: '☕ Coffee Shop / Kafe', value: 'coffee_shop' },
  { label: '🍽️ Restoran / Rumah Makan', value: 'restaurant' },
  { label: '💈 Barbershop / Pangkas Rambut', value: 'barbershop' },
  { label: '💇‍♀️ Salon & Spa', value: 'salon' },
  { label: '🏥 Klinik & Apotek', value: 'clinic' },
  { label: '🛍️ Toko Retail / Butik', value: 'retail' },
  { label: '🏨 Hotel & Penginapan', value: 'hotel' },
  { label: '🏢 Lainnya', value: 'other' },
]

async function checkCode() {
  if (!code.value) {
    checking.value = false
    return
  }

  checking.value = true
  errorMessage.value = ''

  try {
    const res: any = await $fetch(`${apiBase}/api/v1/activate/check`, {
      params: { code: code.value },
    })
    if (res?.success && res?.data) {
      checkResult.value = res.data
    } else {
      errorMessage.value = res?.error?.message || 'Gagal memeriksa status kartu'
    }
  } catch (err: any) {
    errorMessage.value = err?.data?.error?.message || 'Koneksi ke server gagal'
  } finally {
    checking.value = false
  }
}

async function onSubmit() {
  errorMessage.value = ''

  if (!form.business_name.trim()) {
    errorMessage.value = 'Nama toko wajib diisi'
    return
  }
  if (!form.review_url.trim().startsWith('http')) {
    errorMessage.value = 'Link Google Review harus diawali https://'
    return
  }
  if (!form.owner_name.trim()) {
    errorMessage.value = 'Nama pemilik/PIC wajib diisi'
    return
  }
  if (!form.phone.trim()) {
    errorMessage.value = 'Nomor WhatsApp wajib diisi'
    return
  }

  submitting.value = true

  try {
    const res: any = await $fetch(`${apiBase}/api/v1/activate`, {
      method: 'POST',
      body: {
        code: code.value,
        business_name: form.business_name.trim(),
        category: form.category,
        review_url: form.review_url.trim(),
        owner_name: form.owner_name.trim(),
        phone: form.phone.trim(),
        city: form.city.trim(),
      },
    })

    if (res?.success && res?.data) {
      successData.value = res.data
    } else {
      errorMessage.value = res?.error?.message || 'Aktivasi gagal, silakan coba lagi'
    }
  } catch (err: any) {
    errorMessage.value = err?.data?.error?.message || 'Terjadi kesalahan sistem saat aktivasi'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (code.value) {
    checkCode()
  } else {
    checking.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 flex flex-col justify-between py-6 px-4 sm:px-6">
    <div class="max-w-md w-full mx-auto">
      <!-- Brand Header -->
      <div class="text-center mb-6">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-blue-600 text-white shadow-lg shadow-blue-500/30 mb-3">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">REPUTAP</h1>
        <p class="text-xs font-medium text-slate-500 uppercase tracking-widest mt-0.5">Smart Review & Tap System</p>
      </div>

      <!-- State: No Code Provided -->
      <div v-if="!code" class="bg-white rounded-2xl p-6 shadow-sm border border-slate-200 text-center">
        <div class="w-12 h-12 bg-amber-50 text-amber-600 rounded-full flex items-center justify-center mx-auto mb-3 text-xl">
          🔍
        </div>
        <h2 class="text-lg font-bold text-slate-900 mb-1">Masukkan Kode Kartu</h2>
        <p class="text-sm text-slate-500 mb-4">Ketik kode yang tertera di bagian belakang kartu atau stiker akrilik Anda.</p>
        <div class="flex gap-2">
          <input
            v-model="code"
            type="text"
            placeholder="Contoh: rpt-0001"
            class="flex-1 px-4 py-2.5 border border-slate-300 rounded-xl text-sm uppercase tracking-wider font-semibold focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            @click="checkCode"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl text-sm transition"
          >
            Cek
          </button>
        </div>
      </div>

      <!-- State: Loading / Checking -->
      <div v-else-if="checking" class="bg-white rounded-2xl p-8 shadow-sm border border-slate-200 text-center">
        <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto mb-3"></div>
        <p class="text-sm font-medium text-slate-600">Memeriksa kartu {{ code.toUpperCase() }}...</p>
      </div>

      <!-- State: Success Screen -->
      <div v-else-if="successData" class="bg-white rounded-2xl p-6 sm:p-8 shadow-md border border-slate-200 text-center">
        <div class="w-16 h-16 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl shadow-sm">
          ✓
        </div>
        <span class="inline-block px-3 py-1 bg-emerald-50 text-emerald-700 text-xs font-semibold rounded-full uppercase tracking-wider mb-2">
          Kartu Berhasil Aktif
        </span>
        <h2 class="text-xl font-bold text-slate-900 mb-1">{{ successData.store_name }}</h2>
        <p class="text-xs text-slate-400 mb-5 font-mono">Kode Unit: {{ successData.code.toUpperCase() }}</p>

        <div class="bg-slate-50 border border-slate-200 rounded-xl p-4 text-left mb-6 text-sm text-slate-600 leading-relaxed space-y-2">
          <p class="font-medium text-slate-800">🎉 Selamat! Kartu ini sudah terhubung.</p>
          <p>Mulai sekarang, siapa pun yang tap NFC atau scan QR di akrilik ini akan langsung membuka ulasan Google toko Anda.</p>
        </div>

        <div class="space-y-3">
          <a
            :href="successData.redirect_url"
            target="_blank"
            class="block w-full py-3 px-4 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-xl text-sm transition shadow-sm"
          >
            ⚡ Tes Tap Google Review Sekarang
          </a>
          <NuxtLink
            to="/login"
            class="block w-full py-3 px-4 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-xl text-sm transition"
          >
            Masuk ke Panel Admin
          </NuxtLink>
        </div>
      </div>

      <!-- State: Already Active -->
      <div v-else-if="checkResult && !checkResult.can_activate && checkResult.status === 'active'" class="bg-white rounded-2xl p-6 shadow-sm border border-slate-200 text-center">
        <div class="w-14 h-14 bg-blue-50 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-3 text-2xl">
          🏪
        </div>
        <h2 class="text-lg font-bold text-slate-900 mb-1">Kartu Sudah Aktif</h2>
        <p class="text-sm text-slate-600 mb-2 font-medium">{{ checkResult.store_name }}</p>
        <p class="text-xs text-slate-400 mb-5 font-mono">Kode: {{ code.toUpperCase() }}</p>
        <p class="text-xs text-slate-500 mb-6 bg-slate-50 p-3 rounded-lg border border-slate-100">
          Kartu ini sudah terdaftar dan siap digunakan pelanggan untuk memberikan bintang ulasan Google.
        </p>

        <a
          v-if="checkResult.review_url"
          :href="checkResult.review_url"
          target="_blank"
          class="block w-full py-2.5 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl text-sm transition"
        >
          Buka Link Google Review
        </a>
      </div>

      <!-- State: Not Found -->
      <div v-else-if="checkResult && checkResult.status === 'not_found'" class="bg-white rounded-2xl p-6 shadow-sm border border-slate-200 text-center">
        <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-full flex items-center justify-center mx-auto mb-3 text-xl">
          ⚠️
        </div>
        <h2 class="text-lg font-bold text-slate-900 mb-1">Kartu Tidak Ditemukan</h2>
        <p class="text-sm text-slate-500 mb-4">Kode <span class="font-mono font-semibold text-slate-700">{{ code.toUpperCase() }}</span> tidak terdaftar di sistem Reputap.</p>
        <button
          @click="code = ''; checkResult = null"
          class="px-4 py-2 text-sm text-blue-600 hover:underline font-medium"
        >
          Coba kode lain
        </button>
      </div>

      <!-- State: Ready to Activate Form -->
      <div v-else class="bg-white rounded-2xl p-6 sm:p-7 shadow-sm border border-slate-200">
        <div class="flex items-center justify-between pb-4 mb-5 border-b border-slate-100">
          <div>
            <span class="inline-block px-2.5 py-0.5 bg-blue-50 text-blue-700 text-[11px] font-semibold rounded-md uppercase tracking-wider">
              Ready Stock Unit
            </span>
            <h2 class="text-lg font-bold text-slate-900 mt-1">Aktivasi Toko Anda</h2>
          </div>
          <span class="text-xs font-mono font-bold bg-slate-100 text-slate-700 px-2.5 py-1.5 rounded-lg border border-slate-200">
            {{ code.toUpperCase() }}
          </span>
        </div>

        <form @submit.prevent="onSubmit" class="space-y-4">
          <!-- Store Name -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1.5">
              Nama Toko / Bisnis <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="form.business_name"
              type="text"
              required
              placeholder="Contoh: Kopi Kenangan Senopati"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Category -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1.5">
              Kategori Usaha <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="form.category"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option v-for="cat in categories" :key="cat.value" :value="cat.value">
                {{ cat.label }}
              </option>
            </select>
          </div>

          <!-- Google Review Link -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider">
                Link Google Review <span class="text-rose-500">*</span>
              </label>
            </div>
            <input
              v-model="form.review_url"
              type="url"
              required
              placeholder="https://g.page/r/.../review"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <p class="text-[11px] text-slate-500 mt-1.5 leading-relaxed bg-blue-50/60 text-blue-900/80 p-2.5 rounded-lg border border-blue-100">
              💡 <strong>Cara dapat link:</strong> Buka <em>Google Maps</em> di HP &rarr; Profil Bisnis Anda &rarr; Pilih <em>"Minta Ulasan"</em> &rarr; Salin Tautan Ulasan.
            </p>
          </div>

          <!-- Owner / PIC Name -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1.5">
              Nama Pemilik / PIC <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="form.owner_name"
              type="text"
              required
              placeholder="Contoh: Budi Santoso"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- WhatsApp Number -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1.5">
              Nomor WhatsApp <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="form.phone"
              type="tel"
              required
              placeholder="Contoh: 081234567890"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- City -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1.5">
              Kota / Daerah <span class="text-slate-400 font-normal">(Opsional)</span>
            </label>
            <input
              v-model="form.city"
              type="text"
              placeholder="Contoh: Jakarta Selatan"
              class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 text-rose-700 text-xs rounded-xl font-medium">
            {{ errorMessage }}
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="submitting"
            class="w-full py-3.5 px-4 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-bold rounded-xl text-sm transition shadow-md shadow-blue-500/20 mt-2 flex items-center justify-center gap-2"
          >
            <span v-if="submitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span>{{ submitting ? 'Mengaktifkan Kartu...' : '🚀 Aktifkan Kartu Sekarang' }}</span>
          </button>
        </form>
      </div>
    </div>

    <!-- Footer -->
    <footer class="mt-8 text-center text-xs text-slate-400">
      &copy; {{ new Date().getFullYear() }} REPUTAP &bull; Smart Review Automation
    </footer>
  </div>
</template>
