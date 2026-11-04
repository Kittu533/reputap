<script setup lang="ts">
import type { AnalyticsOverview, TimeSeriesPoint, Store, StoreAnalytics } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()

const loadingOverview = ref(true)
const loadingChart = ref(true)
const loadingStoreDetail = ref(false)

const overview = ref<AnalyticsOverview | null>(null)
const chartData = ref<TimeSeriesPoint[]>([])
const stores = ref<Store[]>([])

// Filters
const selectedPeriod = ref<'daily' | 'weekly' | 'monthly'>('daily')
const selectedStoreId = ref<string>('')
const dateFrom = ref('')
const dateTo = ref('')

// Store Detail Drawer / Modal
const selectedStoreAnalytics = ref<StoreAnalytics | null>(null)
const isStoreDetailOpen = ref(false)

const maxChartValue = computed(() => {
  if (!chartData.value || chartData.value.length === 0) return 10
  const max = Math.max(...chartData.value.map(d => d.total))
  return max > 0 ? Math.ceil(max * 1.15) : 10
})

const totalChartScans = computed(() => {
  return chartData.value.reduce((acc, curr) => acc + curr.total, 0)
})

async function loadOverview() {
  loadingOverview.value = true
  try {
    const res = await apiFetch<AnalyticsOverview>('/api/v1/analytics/overview')
    if (res.success && res.data) {
      overview.value = res.data
    }
  }
  finally {
    loadingOverview.value = false
  }
}

async function loadStores() {
  try {
    const res = await apiFetch<Store[]>('/api/v1/stores', { query: { per_page: 100 } })
    if (res.success && res.data) {
      stores.value = res.data
    }
  }
  catch (e) {
    console.error('Failed to load stores for filter', e)
  }
}

async function loadTimeSeries() {
  loadingChart.value = true
  try {
    const query: Record<string, any> = {
      period: selectedPeriod.value,
    }
    if (selectedStoreId.value) {
      query.store_id = selectedStoreId.value
    }
    if (dateFrom.value) {
      query.from = dateFrom.value
    }
    if (dateTo.value) {
      query.to = dateTo.value
    }

    const res = await apiFetch<TimeSeriesPoint[]>('/api/v1/analytics/scans', { query })
    if (res.success && res.data) {
      chartData.value = res.data
    }
  }
  finally {
    loadingChart.value = false
  }
}

async function viewStoreAnalytics(storeId: number) {
  loadingStoreDetail.value = true
  isStoreDetailOpen.value = true
  try {
    const res = await apiFetch<StoreAnalytics>(`/api/v1/analytics/stores/${storeId}`)
    if (res.success && res.data) {
      selectedStoreAnalytics.value = res.data
    }
  }
  finally {
    loadingStoreDetail.value = false
  }
}

function handleFilterChange() {
  loadTimeSeries()
}

function resetDateFilters() {
  dateFrom.value = ''
  dateTo.value = ''
  selectedStoreId.value = ''
  selectedPeriod.value = 'daily'
  loadTimeSeries()
}

onMounted(() => {
  loadOverview()
  loadStores()
  loadTimeSeries()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-purple-600 bg-purple-50 px-2.5 py-0.5 rounded-full border border-purple-100">
            Real-time Scan Telemetry
          </span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl mt-1.5">
          Scan & NFC Analytics
        </h1>
        <p class="mt-1 text-sm text-slate-500 max-w-2xl">
          Track customer taps, QR camera scans, device distributions, and individual store conversion rates.
        </p>
      </div>

      <div class="flex items-center gap-2.5 shrink-0">
        <UButton
          icon="i-lucide-refresh-cw"
          color="neutral"
          variant="outline"
          label="Refresh Data"
          :loading="loadingOverview || loadingChart"
          @click="() => { loadOverview(); loadTimeSeries() }"
        />
      </div>
    </div>

    <!-- Stat Cards Overview -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      <!-- Total Scans -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Total Scans & Taps</span>
          <div class="size-10 rounded-xl bg-purple-50 text-purple-600 border border-purple-100 flex items-center justify-center">
            <UIcon name="i-lucide-activity" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ overview?.total_scans || 0 }}
          </p>
          <div class="mt-2 flex items-center gap-2 text-xs text-slate-500">
            <span class="text-blue-600 font-semibold">{{ overview?.qr_scans || 0 }} QR</span>
            <span>•</span>
            <span class="text-emerald-600 font-semibold">{{ overview?.nfc_scans || 0 }} NFC</span>
          </div>
        </div>
      </div>

      <!-- Today Scans -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Today vs Yesterday</span>
          <div class="size-10 rounded-xl bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center">
            <UIcon name="i-lucide-calendar" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ overview?.today_scans || 0 }}
          </p>
          <div class="mt-2 text-xs text-slate-500">
            <span>Yesterday: </span>
            <span class="font-semibold text-slate-700">{{ overview?.yesterday_scans || 0 }} scans</span>
          </div>
        </div>
      </div>

      <!-- This Week -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">This Week</span>
          <div class="size-10 rounded-xl bg-emerald-50 text-emerald-600 border border-emerald-100 flex items-center justify-center">
            <UIcon name="i-lucide-trending-up" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ overview?.this_week_scans || 0 }}
          </p>
          <div class="mt-2 text-xs text-slate-500">
            Current calendar week
          </div>
        </div>
      </div>

      <!-- This Month -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">This Month</span>
          <div class="size-10 rounded-xl bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center">
            <UIcon name="i-lucide-bar-chart-2" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ overview?.this_month_scans || 0 }}
          </p>
          <div class="mt-2 text-xs text-slate-500">
            Current calendar month
          </div>
        </div>
      </div>
    </div>

    <!-- Time Series Chart Section with Filters -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between border-b border-slate-100 pb-5">
        <div>
          <h2 class="text-base font-bold text-slate-900">Scan Activity Trends</h2>
          <p class="text-xs text-slate-500 mt-0.5">
            Aggregated traffic from smart review stands (Total: {{ totalChartScans }} scans in current range)
          </p>
        </div>

        <!-- Filter Controls -->
        <div class="flex flex-wrap items-center gap-3">
          <!-- Period Selector Tabs -->
          <div class="inline-flex rounded-lg bg-slate-100 p-1 text-xs font-medium text-slate-600">
            <button
              type="button"
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedPeriod === 'daily' ? 'bg-white text-slate-900 shadow-xs font-semibold' : 'hover:text-slate-900'"
              @click="() => { selectedPeriod = 'daily'; handleFilterChange() }"
            >
              Daily
            </button>
            <button
              type="button"
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedPeriod === 'weekly' ? 'bg-white text-slate-900 shadow-xs font-semibold' : 'hover:text-slate-900'"
              @click="() => { selectedPeriod = 'weekly'; handleFilterChange() }"
            >
              Weekly
            </button>
            <button
              type="button"
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedPeriod === 'monthly' ? 'bg-white text-slate-900 shadow-xs font-semibold' : 'hover:text-slate-900'"
              @click="() => { selectedPeriod = 'monthly'; handleFilterChange() }"
            >
              Monthly
            </button>
          </div>

          <!-- Store Selector -->
          <select
            v-model="selectedStoreId"
            class="text-xs rounded-lg border border-slate-200 px-3 py-1.5 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            @change="handleFilterChange"
          >
            <option value="">All Stores</option>
            <option v-for="s in stores" :key="s.id" :value="String(s.id)">
              {{ s.name }}
            </option>
          </select>

          <!-- Date Filters -->
          <input
            v-model="dateFrom"
            type="date"
            class="text-xs rounded-lg border border-slate-200 px-2.5 py-1.5 bg-white text-slate-700"
            @change="handleFilterChange"
          />
          <span class="text-xs text-slate-400">to</span>
          <input
            v-model="dateTo"
            type="date"
            class="text-xs rounded-lg border border-slate-200 px-2.5 py-1.5 bg-white text-slate-700"
            @change="handleFilterChange"
          />

          <button
            v-if="dateFrom || dateTo || selectedStoreId"
            type="button"
            class="text-xs text-slate-500 hover:text-slate-800 underline px-1"
            @click="resetDateFilters"
          >
            Reset
          </button>
        </div>
      </div>

      <!-- Legend -->
      <div class="flex items-center gap-5 text-xs text-slate-600">
        <div class="flex items-center gap-1.5">
          <span class="size-3 rounded-sm bg-blue-500" />
          <span>QR Scans</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="size-3 rounded-sm bg-emerald-500" />
          <span>NFC Taps</span>
        </div>
      </div>

      <!-- Chart Container -->
      <div class="min-h-64 flex items-end justify-center pt-4">
        <div v-if="loadingChart" class="h-64 flex items-center justify-center w-full text-xs text-slate-400">
          Loading scan data...
        </div>

        <div v-else-if="chartData.length === 0" class="h-64 flex flex-col items-center justify-center w-full text-slate-400">
          <UIcon name="i-lucide-bar-chart-3" class="size-8 mb-2 opacity-50" />
          <p class="text-xs">No scan events recorded for this selected time window.</p>
        </div>

        <!-- CSS/SVG Bar Chart -->
        <div v-else class="w-full flex items-end gap-2 overflow-x-auto pb-4 pt-6 h-64 border-b border-slate-200">
          <div
            v-for="point in chartData"
            :key="point.date"
            class="flex-1 min-w-8 flex flex-col items-center h-full justify-end group relative"
          >
            <!-- Hover Tooltip -->
            <div class="absolute -top-10 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-900 text-white text-[10px] rounded px-2 py-1 pointer-events-none whitespace-nowrap z-20 shadow-md">
              <span class="font-bold">{{ point.date }}</span>: {{ point.total }} total ({{ point.qr_count }} QR, {{ point.nfc_count }} NFC)
            </div>

            <!-- Stacked Bar -->
            <div
              class="w-full max-w-6 rounded-t overflow-hidden flex flex-col-reverse transition-all group-hover:brightness-95"
              :style="{ height: `${Math.max(8, (point.total / maxChartValue) * 100)}%` }"
            >
              <!-- QR segment -->
              <div
                class="bg-blue-500 w-full"
                :style="{ height: point.total > 0 ? `${(point.qr_count / point.total) * 100}%` : '0%' }"
              />
              <!-- NFC segment -->
              <div
                class="bg-emerald-500 w-full"
                :style="{ height: point.total > 0 ? `${(point.nfc_count / point.total) * 100}%` : '0%' }"
              />
            </div>

            <!-- Date Label -->
            <span class="text-[10px] text-slate-400 mt-2 truncate w-full text-center">
              {{ point.date.length > 5 ? point.date.slice(5) : point.date }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 2-Column Section: Device Breakdown & Top Stores -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Device Breakdown (1 Col) -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-4">
        <div>
          <h2 class="text-base font-bold text-slate-900">Device Breakdown</h2>
          <p class="text-xs text-slate-400 mt-0.5">Visitor hardware operating systems</p>
        </div>

        <div class="space-y-4 pt-2">
          <!-- Android -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-smartphone" class="size-4 text-emerald-600" />
                <span class="font-medium text-slate-700">Android</span>
              </div>
              <span class="font-bold text-slate-900">{{ overview?.device_breakdown?.android || 0 }}</span>
            </div>
            <div class="w-full h-2 rounded-full bg-slate-100 overflow-hidden">
              <div
                class="h-full bg-emerald-500 rounded-full"
                :style="{ width: overview?.total_scans ? `${((overview.device_breakdown?.android || 0) / overview.total_scans) * 100}%` : '0%' }"
              />
            </div>
          </div>

          <!-- iOS -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-apple" class="size-4 text-slate-800" />
                <span class="font-medium text-slate-700">iOS (iPhone)</span>
              </div>
              <span class="font-bold text-slate-900">{{ overview?.device_breakdown?.ios || 0 }}</span>
            </div>
            <div class="w-full h-2 rounded-full bg-slate-100 overflow-hidden">
              <div
                class="h-full bg-slate-800 rounded-full"
                :style="{ width: overview?.total_scans ? `${((overview.device_breakdown?.ios || 0) / overview.total_scans) * 100}%` : '0%' }"
              />
            </div>
          </div>

          <!-- Desktop -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-monitor" class="size-4 text-blue-600" />
                <span class="font-medium text-slate-700">Desktop / PC</span>
              </div>
              <span class="font-bold text-slate-900">{{ overview?.device_breakdown?.desktop || 0 }}</span>
            </div>
            <div class="w-full h-2 rounded-full bg-slate-100 overflow-hidden">
              <div
                class="h-full bg-blue-500 rounded-full"
                :style="{ width: overview?.total_scans ? `${((overview.device_breakdown?.desktop || 0) / overview.total_scans) * 100}%` : '0%' }"
              />
            </div>
          </div>

          <!-- Other -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-help-circle" class="size-4 text-slate-400" />
                <span class="font-medium text-slate-700">Other / Unknown</span>
              </div>
              <span class="font-bold text-slate-900">{{ overview?.device_breakdown?.other || 0 }}</span>
            </div>
            <div class="w-full h-2 rounded-full bg-slate-100 overflow-hidden">
              <div
                class="h-full bg-slate-400 rounded-full"
                :style="{ width: overview?.total_scans ? `${((overview.device_breakdown?.other || 0) / overview.total_scans) * 100}%` : '0%' }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Top Stores Ranking (2 Cols) -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
        <div class="px-6 py-5 border-b border-slate-200/80 flex items-center justify-between">
          <div>
            <h2 class="text-base font-bold text-slate-900">Top Performing Stores</h2>
            <p class="text-xs text-slate-400 mt-0.5">Highest review stand scan counts</p>
          </div>
        </div>

        <div v-if="overview?.top_stores?.length === 0" class="p-8 text-center text-xs text-slate-400">
          No store scan activity recorded yet.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-100 bg-slate-50/60 text-[11px] font-bold uppercase tracking-wider text-slate-400">
                <th class="py-3 px-6">Store</th>
                <th class="py-3 px-6">Customer</th>
                <th class="py-3 px-6">Total Scans</th>
                <th class="py-3 px-6">QR / NFC</th>
                <th class="py-3 px-6 text-right">Details</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-sm">
              <tr
                v-for="store in overview?.top_stores"
                :key="store.store_id"
                class="hover:bg-slate-50/70 transition-colors"
              >
                <td class="py-4 px-6 font-semibold text-slate-900">
                  {{ store.store_name }}
                </td>
                <td class="py-4 px-6 text-xs text-slate-500">
                  {{ store.customer_name }}
                </td>
                <td class="py-4 px-6 font-bold text-purple-700">
                  {{ store.total_scans }}
                </td>
                <td class="py-4 px-6 text-xs text-slate-600">
                  <span class="text-blue-600 font-semibold">{{ store.qr_scans }} QR</span> /
                  <span class="text-emerald-600 font-semibold">{{ store.nfc_scans }} NFC</span>
                </td>
                <td class="py-4 px-6 text-right">
                  <button
                    type="button"
                    class="text-xs font-semibold text-blue-600 hover:text-blue-800"
                    @click="viewStoreAnalytics(store.store_id)"
                  >
                    View store &rarr;
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Store Specific Analytics Slideover / Drawer -->
    <UModal v-model:open="isStoreDetailOpen">
      <template #content>
        <div class="p-6 space-y-5">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <div>
              <h3 class="text-lg font-bold text-slate-900">
                {{ selectedStoreAnalytics?.store_name || 'Store Analytics' }}
              </h3>
              <p class="text-xs text-slate-400">
                {{ selectedStoreAnalytics?.customer_name }} • {{ selectedStoreAnalytics?.status }}
              </p>
            </div>
            <UButton
              icon="i-lucide-x"
              color="neutral"
              variant="ghost"
              size="xs"
              @click="isStoreDetailOpen = false"
            />
          </div>

          <div v-if="loadingStoreDetail" class="py-8 text-center text-xs text-slate-400">
            Loading store performance data...
          </div>

          <div v-else-if="selectedStoreAnalytics" class="space-y-4">
            <!-- Review URL -->
            <div class="bg-slate-50 p-3 rounded-xl border border-slate-200/70 text-xs">
              <span class="font-semibold text-slate-700">Review Destination:</span>
              <a
                :href="selectedStoreAnalytics.review_url"
                target="_blank"
                class="block text-blue-600 truncate underline mt-0.5"
              >
                {{ selectedStoreAnalytics.review_url || 'No URL configured' }}
              </a>
            </div>

            <!-- Stats Grid -->
            <div class="grid grid-cols-3 gap-3">
              <div class="bg-purple-50/60 p-3 rounded-xl border border-purple-100">
                <span class="text-[10px] font-bold text-slate-400 uppercase">Total Scans</span>
                <p class="text-lg font-bold text-purple-700">{{ selectedStoreAnalytics.total_scans }}</p>
              </div>
              <div class="bg-blue-50/60 p-3 rounded-xl border border-blue-100">
                <span class="text-[10px] font-bold text-slate-400 uppercase">QR Scans</span>
                <p class="text-lg font-bold text-blue-700">{{ selectedStoreAnalytics.qr_scans }}</p>
              </div>
              <div class="bg-emerald-50/60 p-3 rounded-xl border border-emerald-100">
                <span class="text-[10px] font-bold text-slate-400 uppercase">NFC Taps</span>
                <p class="text-lg font-bold text-emerald-700">{{ selectedStoreAnalytics.nfc_scans }}</p>
              </div>
            </div>

            <!-- Active QRs -->
            <div class="text-xs text-slate-500">
              <span>Active Physical Stands: </span>
              <span class="font-bold text-slate-800">{{ selectedStoreAnalytics.active_qr_codes }} deployed</span>
            </div>

            <!-- Daily Activity List (Last 30 Days) -->
            <div>
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider mb-2">Daily Activity</h4>
              <div v-if="selectedStoreAnalytics.daily_scans.length === 0" class="text-xs text-slate-400 py-3 text-center">
                No recent daily scan events.
              </div>
              <div v-else class="max-h-48 overflow-y-auto divide-y divide-slate-100 border border-slate-200 rounded-xl">
                <div
                  v-for="day in selectedStoreAnalytics.daily_scans"
                  :key="day.date"
                  class="flex items-center justify-between px-3 py-2 text-xs"
                >
                  <span class="font-medium text-slate-700">{{ day.date }}</span>
                  <div class="flex items-center gap-2">
                    <span class="text-blue-600">{{ day.qr_count }} QR</span>
                    <span>•</span>
                    <span class="text-emerald-600">{{ day.nfc_count }} NFC</span>
                    <span class="font-bold text-slate-900 ml-2">({{ day.total }} total)</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>
