<script setup lang="ts">
import type { Order, Customer, Store, ProductionJob, QCRecord, Shipment } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const { user } = useAuth()

const loading = ref(true)

// Stats
const orders = ref<Order[]>([])
const customers = ref<Customer[]>([])
const stores = ref<Store[]>([])
const productionJobs = ref<ProductionJob[]>([])
const qcRecords = ref<QCRecord[]>([])
const shipments = ref<Shipment[]>([])

const totalRevenue = computed(() => {
  return orders.value.reduce((acc, curr) => acc + (Number(curr.total) || 0), 0)
})

const activeStoresCount = computed(() => {
  return stores.value.filter(s => s.status === 'active').length
})

const verifiedStoresCount = computed(() => {
  return stores.value.filter(s => s.review_url_verified).length
})

const pendingProductionCount = computed(() => {
  return productionJobs.value.filter(j => j.status === 'pending' || j.status === 'in_progress').length
})

const qcPassRate = computed(() => {
  if (qcRecords.value.length === 0) return 100
  const passed = qcRecords.value.filter(q => q.result === 'passed').length
  return Math.round((passed / qcRecords.value.length) * 100)
})

function formatRupiah(amount: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(amount)
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

async function loadDashboardData() {
  loading.value = true
  try {
    const [ordersRes, custRes, storesRes, prodRes, qcRes, shipRes] = await Promise.all([
      apiFetch<Order[]>('/api/v1/orders', { query: { per_page: 5 } }),
      apiFetch<Customer[]>('/api/v1/customers', { query: { per_page: 50 } }),
      apiFetch<Store[]>('/api/v1/stores', { query: { per_page: 50 } }),
      apiFetch<ProductionJob[]>('/api/v1/production/jobs', { query: { per_page: 50 } }),
      apiFetch<QCRecord[]>('/api/v1/quality-controls', { query: { per_page: 50 } }),
      apiFetch<Shipment[]>('/api/v1/shipments', { query: { per_page: 50 } }),
    ])

    if (ordersRes.success) orders.value = ordersRes.data || []
    if (custRes.success) customers.value = custRes.data || []
    if (storesRes.success) stores.value = storesRes.data || []
    if (prodRes.success) productionJobs.value = prodRes.data || []
    if (qcRes.success) qcRecords.value = qcRes.data || []
    if (shipRes.success) shipments.value = shipRes.data || []
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Welcome Header Banner -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-blue-600 bg-blue-50 px-2.5 py-0.5 rounded-full border border-blue-100">
            System Operational
          </span>
          <span class="text-xs text-slate-400">• {{ new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl mt-1.5">
          Welcome back, {{ user?.name || 'Admin' }} 👋
        </h1>
        <p class="mt-1 text-sm text-slate-500 leading-relaxed max-w-2xl">
          Here is your Done-for-You smart review stand manufacturing and operations overview today.
        </p>
      </div>

      <div class="flex items-center gap-2.5 shrink-0">
        <UButton
          icon="i-lucide-refresh-cw"
          color="neutral"
          variant="outline"
          label="Refresh"
          :loading="loading"
          @click="loadDashboardData"
        />
        <NuxtLink
          to="/orders/new"
          class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-semibold bg-blue-600 hover:bg-blue-700 text-white shadow-xs transition-colors"
        >
          <UIcon name="i-lucide-plus" class="size-4" />
          <span>New Order</span>
        </NuxtLink>
      </div>
    </div>

    <!-- 4 Key Stat Cards -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      <!-- Stat 1: Total Revenue -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs hover:border-slate-300 transition-all">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Total Revenue</span>
          <div class="size-10 rounded-xl bg-emerald-50 text-emerald-600 border border-emerald-100 flex items-center justify-center">
            <UIcon name="i-lucide-dollar-sign" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ formatRupiah(totalRevenue) }}
          </p>
          <div class="mt-2 flex items-center gap-1.5 text-xs text-slate-500">
            <span class="font-semibold text-emerald-600 inline-flex items-center">
              <UIcon name="i-lucide-trending-up" class="size-3.5 mr-0.5" />
              Active
            </span>
            <span>across {{ orders.length }} recorded orders</span>
          </div>
        </div>
      </div>

      <!-- Stat 2: Orders Count -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs hover:border-slate-300 transition-all">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Orders Placed</span>
          <div class="size-10 rounded-xl bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center">
            <UIcon name="i-lucide-shopping-cart" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ orders.length }}
          </p>
          <div class="mt-2 flex items-center gap-1.5 text-xs text-slate-500">
            <NuxtLink to="/orders" class="font-semibold text-blue-600 hover:underline">
              View order pipeline &rarr;
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Stat 3: Stores & Customers -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs hover:border-slate-300 transition-all">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Clients & Stores</span>
          <div class="size-10 rounded-xl bg-purple-50 text-purple-600 border border-purple-100 flex items-center justify-center">
            <UIcon name="i-lucide-store" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ activeStoresCount }} <span class="text-base font-normal text-slate-400">/ {{ stores.length }}</span>
          </p>
          <div class="mt-2 flex items-center gap-1.5 text-xs text-slate-500">
            <span class="font-medium text-emerald-600">{{ verifiedStoresCount }} verified URLs</span>
            <span>from {{ customers.length }} clients</span>
          </div>
        </div>
      </div>

      <!-- Stat 4: Production & QC -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs hover:border-slate-300 transition-all">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Production & QC</span>
          <div class="size-10 rounded-xl bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center">
            <UIcon name="i-lucide-factory" class="size-5" />
          </div>
        </div>
        <div class="mt-3">
          <p class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
            {{ pendingProductionCount }} <span class="text-base font-normal text-slate-400">active jobs</span>
          </p>
          <div class="mt-2 flex items-center gap-1.5 text-xs text-slate-500">
            <span class="font-semibold text-teal-600">{{ qcPassRate }}% QC pass rate</span>
            <span>across {{ qcRecords.length }} checks</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Done-For-You Workflow Stepper Pipeline -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="text-base font-bold text-slate-900">Done-for-You Manufacturing & Fulfilment Flow</h2>
          <p class="text-xs text-slate-500">End-to-end status flow from customer order to live smart stand delivery</p>
        </div>
        <span class="text-xs font-medium text-slate-400 hidden sm:inline">5-Stage System</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-5 gap-3 pt-2">
        <div class="flex flex-col p-3 rounded-xl border border-slate-200/80 bg-slate-50/60">
          <div class="flex items-center gap-2 mb-1">
            <span class="size-6 rounded-full bg-blue-600 text-white text-xs font-bold flex items-center justify-center">1</span>
            <span class="text-xs font-semibold text-slate-900">Order & Store</span>
          </div>
          <p class="text-[11px] text-slate-500">Review URL verified, order confirmed</p>
        </div>

        <div class="flex flex-col p-3 rounded-xl border border-slate-200/80 bg-slate-50/60">
          <div class="flex items-center gap-2 mb-1">
            <span class="size-6 rounded-full bg-blue-600 text-white text-xs font-bold flex items-center justify-center">2</span>
            <span class="text-xs font-semibold text-slate-900">Design Asset</span>
          </div>
          <p class="text-[11px] text-slate-500">Stand artwork uploaded & approved</p>
        </div>

        <div class="flex flex-col p-3 rounded-xl border border-slate-200/80 bg-slate-50/60">
          <div class="flex items-center gap-2 mb-1">
            <span class="size-6 rounded-full bg-blue-600 text-white text-xs font-bold flex items-center justify-center">3</span>
            <span class="text-xs font-semibold text-slate-900">Production</span>
          </div>
          <p class="text-[11px] text-slate-500">Acrylic cut, printed, NFC & QR encoded</p>
        </div>

        <div class="flex flex-col p-3 rounded-xl border border-slate-200/80 bg-slate-50/60">
          <div class="flex items-center gap-2 mb-1">
            <span class="size-6 rounded-full bg-blue-600 text-white text-xs font-bold flex items-center justify-center">4</span>
            <span class="text-xs font-semibold text-slate-900">Quality Control</span>
          </div>
          <p class="text-[11px] text-slate-500">Physical & scan hardware inspection</p>
        </div>

        <div class="flex flex-col p-3 rounded-xl border border-slate-200/80 bg-slate-50/60">
          <div class="flex items-center gap-2 mb-1">
            <span class="size-6 rounded-full bg-emerald-600 text-white text-xs font-bold flex items-center justify-center">5</span>
            <span class="text-xs font-semibold text-slate-900">Shipping</span>
          </div>
          <p class="text-[11px] text-slate-500">Packed, courier tracking & delivered</p>
        </div>
      </div>
    </div>

    <!-- Main 2-Column Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Recent Orders (2 cols) -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
          <div class="px-6 py-5 border-b border-slate-200/80 flex items-center justify-between">
            <div>
              <h2 class="text-base font-bold text-slate-900">Recent Orders</h2>
              <p class="text-xs text-slate-500 mt-0.5">Latest smart stand purchase orders</p>
            </div>
            <NuxtLink
              to="/orders"
              class="text-xs font-semibold text-blue-600 hover:text-blue-700 flex items-center gap-1"
            >
              <span>View all</span>
              <UIcon name="i-lucide-arrow-right" class="size-3.5" />
            </NuxtLink>
          </div>

          <div v-if="orders.length === 0" class="p-8 text-center">
            <CommonEmptyState
              title="No orders yet"
              description="Create your first client order to start the Done-For-You manufacturing pipeline."
              action-label="Create Order"
              @action="navigateTo('/orders/new')"
            />
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="border-b border-slate-100 bg-slate-50/60 text-[11px] font-bold uppercase tracking-wider text-slate-400">
                  <th class="py-3 px-6">Order #</th>
                  <th class="py-3 px-6">Date</th>
                  <th class="py-3 px-6">Total</th>
                  <th class="py-3 px-6">Order Status</th>
                  <th class="py-3 px-6">Payment</th>
                  <th class="py-3 px-6 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-sm">
                <tr
                  v-for="ord in orders"
                  :key="ord.id"
                  class="hover:bg-slate-50/70 transition-colors"
                >
                  <td class="py-4 px-6 font-semibold text-slate-900">
                    <NuxtLink :to="`/orders/${ord.id}`" class="text-blue-600 hover:underline">
                      {{ ord.order_number }}
                    </NuxtLink>
                  </td>
                  <td class="py-4 px-6 text-xs text-slate-500">
                    {{ formatDate(ord.created_at) }}
                  </td>
                  <td class="py-4 px-6 font-semibold text-slate-900">
                    {{ formatRupiah(ord.total) }}
                  </td>
                  <td class="py-4 px-6">
                    <CommonStatusBadge :status="ord.order_status" size="xs" />
                  </td>
                  <td class="py-4 px-6">
                    <CommonStatusBadge :status="ord.payment_status" size="xs" />
                  </td>
                  <td class="py-4 px-6 text-right">
                    <NuxtLink
                      :to="`/orders/${ord.id}`"
                      class="inline-flex items-center gap-1 text-xs font-semibold text-slate-600 hover:text-slate-900 px-2.5 py-1 rounded-md hover:bg-slate-100 transition-colors"
                    >
                      <span>Manage</span>
                      <UIcon name="i-lucide-chevron-right" class="size-3.5" />
                    </NuxtLink>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Operations Snapshot (Production & QC) -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-base font-bold text-slate-900">Factory Production Queue</h2>
              <p class="text-xs text-slate-500 mt-0.5">Physical stands currently in manufacturing</p>
            </div>
            <NuxtLink to="/production" class="text-xs font-semibold text-blue-600 hover:underline">
              Go to queue &rarr;
            </NuxtLink>
          </div>

          <div v-if="productionJobs.length === 0" class="text-sm text-slate-400 py-6 text-center">
            No active production jobs at this moment.
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="job in productionJobs.slice(0, 4)"
              :key="job.id"
              class="flex items-center justify-between p-3.5 rounded-xl border border-slate-200/70 bg-slate-50/50 hover:bg-slate-50 transition-colors"
            >
              <div class="flex items-center gap-3">
                <div class="size-9 rounded-lg bg-white border border-slate-200 text-slate-700 flex items-center justify-center font-mono font-bold text-xs">
                  #{{ job.id }}
                </div>
                <div>
                  <p class="text-xs font-bold text-slate-900 font-mono">
                    {{ job.serial_number }}
                  </p>
                  <p class="text-[11px] text-slate-400">
                    Started: {{ job.started_at ? formatDate(job.started_at) : 'Not started' }}
                  </p>
                </div>
              </div>
              <CommonStatusBadge :status="job.status" size="xs" />
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Shortcuts & Quick Info (1 col) -->
      <div class="space-y-6">
        <!-- Quick Action Shortcuts -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-5">
          <h2 class="text-base font-bold text-slate-900 mb-1">Quick Actions</h2>
          <p class="text-xs text-slate-400 mb-4">Frequently used operational tasks</p>

          <div class="grid grid-cols-1 gap-2.5">
            <NuxtLink
              to="/orders/new"
              class="flex items-center gap-3 p-3 rounded-xl border border-slate-200 hover:border-blue-300 hover:bg-blue-50/40 text-slate-700 hover:text-blue-700 transition-all group"
            >
              <div class="size-8 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center group-hover:scale-105 transition-transform">
                <UIcon name="i-lucide-plus" class="size-4" />
              </div>
              <div class="text-left">
                <p class="text-xs font-bold text-slate-900 group-hover:text-blue-700">New Order</p>
                <p class="text-[11px] text-slate-400">Create new client purchase order</p>
              </div>
            </NuxtLink>

            <NuxtLink
              to="/customers/new"
              class="flex items-center gap-3 p-3 rounded-xl border border-slate-200 hover:border-blue-300 hover:bg-blue-50/40 text-slate-700 hover:text-blue-700 transition-all group"
            >
              <div class="size-8 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center group-hover:scale-105 transition-transform">
                <UIcon name="i-lucide-user-plus" class="size-4" />
              </div>
              <div class="text-left">
                <p class="text-xs font-bold text-slate-900 group-hover:text-blue-700">Add Customer</p>
                <p class="text-[11px] text-slate-400">Register new business customer</p>
              </div>
            </NuxtLink>

            <NuxtLink
              to="/stores/new"
              class="flex items-center gap-3 p-3 rounded-xl border border-slate-200 hover:border-blue-300 hover:bg-blue-50/40 text-slate-700 hover:text-blue-700 transition-all group"
            >
              <div class="size-8 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center group-hover:scale-105 transition-transform">
                <UIcon name="i-lucide-store" class="size-4" />
              </div>
              <div class="text-left">
                <p class="text-xs font-bold text-slate-900 group-hover:text-blue-700">Register Store</p>
                <p class="text-[11px] text-slate-400">Add store location & Google URL</p>
              </div>
            </NuxtLink>

            <NuxtLink
              to="/designs/new"
              class="flex items-center gap-3 p-3 rounded-xl border border-slate-200 hover:border-blue-300 hover:bg-blue-50/40 text-slate-700 hover:text-blue-700 transition-all group"
            >
              <div class="size-8 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center group-hover:scale-105 transition-transform">
                <UIcon name="i-lucide-upload" class="size-4" />
              </div>
              <div class="text-left">
                <p class="text-xs font-bold text-slate-900 group-hover:text-blue-700">Upload Artwork</p>
                <p class="text-[11px] text-slate-400">Attach stand design for order</p>
              </div>
            </NuxtLink>
          </div>
        </div>

        <!-- Smart Stand Hardware & Redirect Status -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-5">
          <h2 class="text-base font-bold text-slate-900 mb-1">NFC & QR Hardware</h2>
          <p class="text-xs text-slate-400 mb-3">Live redirect engine status</p>

          <div class="space-y-3">
            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 border border-slate-200/70">
              <div class="flex items-center gap-2.5">
                <div class="size-2 rounded-full bg-emerald-500 animate-pulse" />
                <span class="text-xs font-semibold text-slate-700">302 Redirect Engine</span>
              </div>
              <span class="text-[11px] font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200">
                Active (Async Log)
              </span>
            </div>

            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 border border-slate-200/70">
              <span class="text-xs text-slate-600">Total Active Shipments</span>
              <span class="text-xs font-bold text-slate-900">{{ shipments.length }} shipments</span>
            </div>

            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 border border-slate-200/70">
              <span class="text-xs text-slate-600">Audit Logs Recorded</span>
              <NuxtLink to="/audit-logs" class="text-xs font-semibold text-blue-600 hover:underline">
                View logs &rarr;
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
