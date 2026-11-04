<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()
const config = useRuntimeConfig()
const redirectBase = (config.public.redirectBase as string) || 'http://localhost:8081/r'

const isSupported = ref(false)
const activeTab = ref<'write' | 'read'>('write')
const instances = ref<any[]>([])
const selectedInstanceId = ref<number | null>(null)
const autoAdvance = ref(true)

// Writing state
const isWriting = ref(false)
const writeStatus = ref<'idle' | 'waiting_tap' | 'success' | 'error'>('idle')
const writeMessage = ref('')
let writeAbortController: AbortController | null = null

// Reading state
const isReading = ref(false)
const readResult = ref<{
  serialNumber?: string
  records: Array<{ type: string; data: string }>
  matchedCode?: string
  dbInfo?: any
} | null>(null)
let readAbortController: AbortController | null = null

const selectedInstance = computed(() => {
  return instances.value.find(i => i.id === selectedInstanceId.value) || null
})

const nfcPayload = computed(() => {
  if (!selectedInstance.value) return ''
  const code = selectedInstance.value.qr_code || selectedInstance.value.serial_number
  return `${redirectBase}/${code}?src=nfc`
})

async function fetchInstances() {
  try {
    const res: any = await apiFetch('/api/v1/product-instances', {
      query: { per_page: 50 },
    })
    if (res?.data) {
      instances.value = res.data
      if (!selectedInstanceId.value && instances.value.length > 0) {
        // Default to first unconfigured card or first card
        const unconfigured = instances.value.find(i => i.nfc_status !== 'configured')
        selectedInstanceId.value = unconfigured ? unconfigured.id : instances.value[0].id
      }
    }
  } catch (err: any) {
    toast.add({
      title: 'Gagal memuat data kartu',
      description: err?.message,
      color: 'error',
    })
  }
}

async function startWriteNFC() {
  if (!isSupported.value) {
    toast.add({
      title: 'Web NFC tidak didukung',
      description: 'Gunakan browser Chrome di Android untuk write NFC langsung dari web.',
      color: 'warning',
    })
    return
  }

  if (!nfcPayload.value) return

  isWriting.value = true
  writeStatus.value = 'waiting_tap'
  writeMessage.value = 'Dekatkan kartu / akrilik ke bagian belakang HP...'

  try {
    writeAbortController = new AbortController()
    const ndef = new (window as any).NDEFReader()

    await ndef.write({
      records: [
        {
          recordType: 'url',
          data: nfcPayload.value,
        },
      ],
    }, { signal: writeAbortController.signal })

    // Success writing to physical chip
    writeStatus.value = 'success'
    writeMessage.value = `Berhasil! NFC terisi: ${nfcPayload.value}`

    // Haptic feedback if available
    if (navigator.vibrate) {
      navigator.vibrate([100, 50, 100])
    }

    // Update backend DB
    if (selectedInstance.value) {
      await apiFetch(`/api/v1/product-instances/${selectedInstance.value.id}/nfc-status`, {
        method: 'PATCH',
        body: { status: 'configured' },
      })
      selectedInstance.value.nfc_status = 'configured'
    }

    toast.add({
      title: 'NFC Berhasil Ditulis',
      description: `${selectedInstance.value?.serial_number} sudah aktif.`,
      color: 'success',
    })

    // Auto advance to next card
    if (autoAdvance.value) {
      const currentIndex = instances.value.findIndex(i => i.id === selectedInstanceId.value)
      if (currentIndex !== -1 && currentIndex < instances.value.length - 1) {
        selectedInstanceId.value = instances.value[currentIndex + 1].id
      }
    }
  } catch (err: any) {
    if (err.name !== 'AbortError') {
      writeStatus.value = 'error'
      writeMessage.value = `Gagal menulis: ${err.message || 'Coba dekatkan kembali kartu.'}`
    }
  } finally {
    isWriting.value = false
  }
}

function cancelWrite() {
  if (writeAbortController) {
    writeAbortController.abort()
    writeAbortController = null
  }
  isWriting.value = false
  writeStatus.value = 'idle'
  writeMessage.value = ''
}

async function startReadNFC() {
  if (!isSupported.value) {
    toast.add({
      title: 'Web NFC tidak didukung',
      description: 'Gunakan browser Chrome di Android untuk scan NFC.',
      color: 'warning',
    })
    return
  }

  isReading.value = true
  readResult.value = null

  try {
    readAbortController = new AbortController()
    const ndef = new (window as any).NDEFReader()
    await ndef.scan({ signal: readAbortController.signal })

    ndef.onreading = async (event: any) => {
      if (navigator.vibrate) {
        navigator.vibrate(100)
      }

      const records: Array<{ type: string; data: string }> = []
      let matchedCode = ''

      for (const record of event.message.records) {
        let textData = ''
        if (record.recordType === 'text') {
          const textDecoder = new TextDecoder(record.encoding)
          textData = textDecoder.decode(record.data)
        } else if (record.recordType === 'url') {
          const textDecoder = new TextDecoder()
          textData = textDecoder.decode(record.data)
          // extract code
          const m = textData.match(/\/r\/([^?]+)/)
          if (m) matchedCode = m[1]
        } else {
          textData = `[${record.recordType}]`
        }

        records.push({
          type: record.recordType,
          data: textData,
        })
      }

      let dbInfo = null
      if (matchedCode) {
        try {
          const res: any = await apiFetch(`/api/v1/activate/check?code=${matchedCode}`)
          if (res?.data) dbInfo = res.data
        } catch (_) {}
      }

      readResult.value = {
        serialNumber: event.serialNumber,
        records,
        matchedCode,
        dbInfo,
      }
    }
  } catch (err: any) {
    if (err.name !== 'AbortError') {
      toast.add({
        title: 'Gagal scan NFC',
        description: err.message,
        color: 'error',
      })
    }
    isReading.value = false
  }
}

function stopRead() {
  if (readAbortController) {
    readAbortController.abort()
    readAbortController = null
  }
  isReading.value = false
}

function copyPayload() {
  if (!nfcPayload.value) return
  navigator.clipboard.writeText(nfcPayload.value)
  toast.add({
    title: 'URL Disalin',
    description: nfcPayload.value,
    color: 'success',
  })
}

onMounted(() => {
  isSupported.value = typeof window !== 'undefined' && 'NDEFReader' in window
  fetchInstances()
})

onBeforeUnmount(() => {
  cancelWrite()
  stopRead()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white flex items-center gap-2">
          <span>📡</span> NFC Workshop Tools
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Flashing payload dan inspeksi chip NFC NTAG213 langsung dari browser.
        </p>
      </div>

      <!-- Browser Support Badge -->
      <UBadge
        :color="isSupported ? 'success' : 'warning'"
        variant="subtle"
        size="md"
        class="gap-1.5"
      >
        <span class="w-2 h-2 rounded-full" :class="isSupported ? 'bg-emerald-500 animate-pulse' : 'bg-amber-500'"></span>
        {{ isSupported ? 'Web NFC Ready (Android Chrome)' : 'Browser Desktop / iOS (Manual Mode)' }}
      </UBadge>
    </div>

    <!-- Notice if Web NFC not supported -->
    <div v-if="!isSupported" class="bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 rounded-2xl p-4 text-sm text-amber-800 dark:text-amber-200 flex items-start gap-3">
      <span class="text-xl">💡</span>
      <div>
        <p class="font-semibold">Akses hardware NFC browser dibatasi di perangkat ini.</p>
        <p class="text-xs mt-1 leading-relaxed text-amber-700 dark:text-amber-300">
          Untuk flashing 1-klik dengan menempelkan kartu ke belakang layar, buka halaman ini di <strong>HP Android (Chrome)</strong>.<br />
          Jika menggunakan <strong>iPhone</strong> atau <strong>USB Reader ACR122U</strong> di komputer, gunakan URL payload yang disediakan di bawah via app <em>NFC Tools</em>.
        </p>
      </div>
    </div>

    <!-- Mode Tabs -->
    <div class="flex border-b border-slate-200 dark:border-slate-800 gap-4">
      <button
        @click="activeTab = 'write'"
        class="pb-3 text-sm font-semibold border-b-2 transition flex items-center gap-2"
        :class="activeTab === 'write' ? 'border-blue-600 text-blue-600 dark:text-blue-400' : 'border-transparent text-slate-500 hover:text-slate-700'"
      >
        <span>⚡</span> Flash / Tulis Kartu
      </button>
      <button
        @click="activeTab = 'read'"
        class="pb-3 text-sm font-semibold border-b-2 transition flex items-center gap-2"
        :class="activeTab === 'read' ? 'border-blue-600 text-blue-600 dark:text-blue-400' : 'border-transparent text-slate-500 hover:text-slate-700'"
      >
        <span>🔍</span> Inspeksi / Baca Kartu
      </button>
    </div>

    <!-- TAB 1: WRITE MODE -->
    <div v-if="activeTab === 'write'" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Unit Selector -->
      <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-5 shadow-sm space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-bold text-slate-900 dark:text-white text-sm">Pilih Unit Akrilik</h2>
          <span class="text-xs text-slate-400">{{ instances.length }} unit</span>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-1.5">
            Unit Serial / QR Code
          </label>
          <select
            v-model="selectedInstanceId"
            class="w-full px-3 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-300 dark:border-slate-700 rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option v-for="inst in instances" :key="inst.id" :value="inst.id">
              {{ inst.serial_number }} &bull; {{ inst.qr_code || 'No QR' }} [{{ inst.nfc_status }}]
            </option>
          </select>
        </div>

        <div v-if="selectedInstance" class="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-xl border border-slate-200 dark:border-slate-700 space-y-2 text-xs">
          <div class="flex justify-between">
            <span class="text-slate-500">Serial Number:</span>
            <span class="font-mono font-bold">{{ selectedInstance.serial_number }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-500">Kode QR:</span>
            <span class="font-mono font-bold">{{ selectedInstance.qr_code }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-slate-500">NFC Status:</span>
            <UBadge
              :color="selectedInstance.nfc_status === 'configured' ? 'success' : 'neutral'"
              variant="subtle"
              size="xs"
              class="capitalize font-mono"
            >
              {{ selectedInstance.nfc_status }}
            </UBadge>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <input
            id="autoAdvance"
            v-model="autoAdvance"
            type="checkbox"
            class="rounded text-blue-600 focus:ring-blue-500"
          />
          <label for="autoAdvance" class="text-xs text-slate-600 dark:text-slate-400">
            Auto-pindah ke unit selanjutnya setelah flash
          </label>
        </div>
      </div>

      <!-- Right Column: Writer Action -->
      <div class="lg:col-span-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
        <div class="space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
            <div>
              <span class="text-xs font-semibold uppercase tracking-wider text-blue-600 dark:text-blue-400">
                Payload Writer
              </span>
              <h2 class="text-base font-bold text-slate-900 dark:text-white mt-0.5">
                URL Payload untuk Chip NTAG213
              </h2>
            </div>
            <span class="text-xs font-mono bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-300 px-2.5 py-1 rounded-md">
              NFC Tag Record: URI
            </span>
          </div>

          <!-- Payload URL Display -->
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1">
              Data yang akan ditulis ke chip:
            </label>
            <div class="flex gap-2">
              <input
                type="text"
                readonly
                :value="nfcPayload"
                class="flex-1 px-3.5 py-2.5 bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl font-mono text-sm text-slate-700 dark:text-slate-300"
              />
              <button
                @click="copyPayload"
                class="px-4 py-2.5 bg-slate-200 hover:bg-slate-300 dark:bg-slate-700 dark:hover:bg-slate-600 rounded-xl text-xs font-semibold transition"
              >
                Copy
              </button>
            </div>
          </div>

          <!-- Status Box during write -->
          <div
            v-if="writeStatus !== 'idle'"
            class="p-4 rounded-xl border text-sm flex items-center gap-3 transition"
            :class="{
              'bg-blue-50 border-blue-200 text-blue-900 dark:bg-blue-950/40 dark:border-blue-800 dark:text-blue-200': writeStatus === 'waiting_tap',
              'bg-emerald-50 border-emerald-200 text-emerald-900 dark:bg-emerald-950/40 dark:border-emerald-800 dark:text-emerald-200': writeStatus === 'success',
              'bg-rose-50 border-rose-200 text-rose-900 dark:bg-rose-950/40 dark:border-rose-800 dark:text-rose-200': writeStatus === 'error',
            }"
          >
            <span v-if="writeStatus === 'waiting_tap'" class="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></span>
            <span v-else-if="writeStatus === 'success'" class="text-lg">✅</span>
            <span v-else class="text-lg">⚠️</span>
            <span class="font-medium text-xs sm:text-sm">{{ writeMessage }}</span>
          </div>
        </div>

        <!-- Buttons -->
        <div class="mt-8 pt-4 border-t border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row gap-3">
          <button
            v-if="!isWriting"
            @click="startWriteNFC"
            :disabled="!selectedInstance"
            class="flex-1 py-3.5 px-4 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-bold rounded-xl text-sm transition shadow-md shadow-blue-500/20 flex items-center justify-center gap-2"
          >
            <span>📡</span> Tulis NFC Sekarang
          </button>
          <button
            v-else
            @click="cancelWrite"
            class="flex-1 py-3.5 px-4 bg-rose-600 hover:bg-rose-700 text-white font-bold rounded-xl text-sm transition"
          >
            Batal Menulis
          </button>
        </div>
      </div>
    </div>

    <!-- TAB 2: READ / INSPECT MODE -->
    <div v-else class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 shadow-sm space-y-6">
      <div class="max-w-xl">
        <h2 class="text-base font-bold text-slate-900 dark:text-white">Inspeksi Chip Fisik</h2>
        <p class="text-xs text-slate-500 mt-1">
          Tempelkan kartu apapun ke belakang HP untuk membaca isi payload dan memverifikasi statusnya di database Reputap.
        </p>

        <div class="mt-4 flex gap-3">
          <button
            v-if="!isReading"
            @click="startReadNFC"
            class="py-2.5 px-5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold rounded-xl text-sm transition flex items-center gap-2"
          >
            <span>🔍</span> Mulai Scan Kartu
          </button>
          <button
            v-else
            @click="stopRead"
            class="py-2.5 px-5 bg-rose-600 hover:bg-rose-700 text-white font-bold rounded-xl text-sm transition"
          >
            Stop Scan
          </button>
        </div>
      </div>

      <!-- Scan result -->
      <div v-if="readResult" class="p-5 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="font-bold text-slate-900 dark:text-white text-sm">Hasil Pembacaan Chip</h3>
          <span class="text-xs font-mono text-slate-400">UID: {{ readResult.serialNumber || 'Unknown' }}</span>
        </div>

        <div class="space-y-2">
          <div v-for="(rec, idx) in readResult.records" :key="idx" class="p-3 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-700 text-xs">
            <span class="font-bold uppercase tracking-wider text-slate-400 mr-2">[{{ rec.type }}]</span>
            <span class="font-mono text-blue-600 dark:text-blue-400 break-all">{{ rec.data }}</span>
          </div>
        </div>

        <!-- DB Match info -->
        <div v-if="readResult.dbInfo" class="mt-4 p-4 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 rounded-xl text-xs space-y-1.5 text-emerald-900 dark:text-emerald-200">
          <p class="font-bold">✓ Kartu Terdaftar di Database Reputap</p>
          <p>Kode: <span class="font-mono font-bold">{{ readResult.dbInfo.code }}</span> &bull; Status: <span class="capitalize font-semibold">{{ readResult.dbInfo.status }}</span></p>
          <p v-if="readResult.dbInfo.store_name">Toko: <span class="font-semibold">{{ readResult.dbInfo.store_name }}</span></p>
        </div>
      </div>
    </div>
  </div>
</template>
