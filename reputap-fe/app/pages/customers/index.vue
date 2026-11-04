<script setup lang="ts">
import type { Customer } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const search = ref('')
const page = ref(1)
const perPage = ref(15)
const data = ref<Customer[]>([])
const total = ref(0)
const loading = ref(false)

// Delete Modal State
const deleteModalOpen = ref(false)
const customerToDelete = ref<Customer | null>(null)
const deleteLoading = ref(false)

async function fetchData() {
  loading.value = true
  try {
    const response = await apiFetch<Customer[]>('/api/v1/customers', {
      query: {
        page: page.value,
        per_page: perPage.value,
        search: search.value || undefined,
      },
    })
    if (response.success) {
      data.value = response.data || []
      total.value = response.meta?.total || 0
    }
  }
  finally {
    loading.value = false
  }
}

function promptDelete(customer: Customer) {
  customerToDelete.value = customer
  deleteModalOpen.value = true
}

async function confirmDelete() {
  if (!customerToDelete.value) return
  deleteLoading.value = true
  try {
    const response = await apiFetch(`/api/v1/customers/${customerToDelete.value.id}`, {
      method: 'DELETE',
    })
    if (response.success) {
      toast.add({
        title: 'Customer Deleted',
        description: `Customer "${customerToDelete.value.name}" (${customerToDelete.value.customer_code}) has been deleted.`,
        color: 'success',
      })
      deleteModalOpen.value = false
      fetchData()
    }
    else {
      toast.add({
        title: 'Failed to Delete',
        description: response.error?.message || 'Customer could not be deleted',
        color: 'error',
      })
    }
  }
  finally {
    deleteLoading.value = false
  }
}

let searchTimeout: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchData()
  }, 300)
})

watch(page, () => fetchData())

onMounted(() => fetchData())
</script>

<template>
  <div class="space-y-6">
    <!-- Header with Breadcrumbs & Action -->
    <CommonPageHeader
      title="Customers"
      description="Manage business clients, review locations, and stand ordering contacts."
      :breadcrumbs="[{ label: 'Customers' }]"
    >
      <template #actions>
        <NuxtLink
          to="/customers/new"
          class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-semibold bg-blue-600 hover:bg-blue-700 text-white shadow-xs transition-colors"
        >
          <UIcon name="i-lucide-plus" class="size-4" />
          <span>Add Customer</span>
        </NuxtLink>
      </template>
    </CommonPageHeader>

    <!-- Search & Filter Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-4 shadow-xs flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
      <div class="relative w-full max-w-sm">
        <UIcon name="i-lucide-search" class="size-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="search"
          type="text"
          placeholder="Search by name, code, email, city..."
          class="w-full pl-9 pr-4 py-2 text-sm bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
        >
      </div>

      <div class="text-xs font-semibold text-slate-500">
        Total: <span class="text-slate-900 font-bold">{{ total }}</span> clients
      </div>
    </div>

    <!-- Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div v-if="loading" class="p-12 text-center text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin mx-auto mb-2 text-blue-600" />
        <span class="text-xs font-medium">Loading customers...</span>
      </div>

      <div v-else-if="data.length === 0">
        <CommonEmptyState
          title="No customers found"
          description="There are no registered customers matching your criteria. Click below to add a new client."
          action-label="Add Customer"
          @action="navigateTo('/customers/new')"
        />
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100 bg-slate-50/70 text-[11px] font-bold uppercase tracking-wider text-slate-400">
              <th class="py-3.5 px-6">Code</th>
              <th class="py-3.5 px-6">Name</th>
              <th class="py-3.5 px-6">Business</th>
              <th class="py-3.5 px-6">Email</th>
              <th class="py-3.5 px-6">Phone</th>
              <th class="py-3.5 px-6">City</th>
              <th class="py-3.5 px-6 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm">
            <tr
              v-for="customer in data"
              :key="customer.id"
              class="hover:bg-slate-50/70 transition-colors"
            >
              <td class="py-4 px-6 font-mono text-xs font-bold text-blue-600">
                {{ customer.customer_code }}
              </td>
              <td class="py-4 px-6 font-semibold text-slate-900">
                {{ customer.name }}
              </td>
              <td class="py-4 px-6 text-slate-600">
                {{ customer.business_name || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-500">
                {{ customer.email || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-500">
                {{ customer.phone || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-600">
                {{ customer.city || '-' }}
              </td>
              <td class="py-4 px-6 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <NuxtLink
                    :to="`/customers/${customer.id}/edit`"
                    class="size-8 rounded-lg text-slate-500 hover:text-blue-600 hover:bg-blue-50 flex items-center justify-center transition-colors"
                    title="Edit Customer"
                  >
                    <UIcon name="i-lucide-pencil" class="size-4" />
                  </NuxtLink>
                  <button
                    class="size-8 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 flex items-center justify-center transition-colors"
                    title="Delete Customer"
                    @click="promptDelete(customer)"
                  >
                    <UIcon name="i-lucide-trash-2" class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div
        v-if="total > perPage"
        class="p-4 border-t border-slate-100 flex items-center justify-between"
      >
        <span class="text-xs text-slate-400">
          Showing {{ ((page - 1) * perPage) + 1 }} to {{ Math.min(page * perPage, total) }} of {{ total }}
        </span>
        <UPagination
          v-model:page="page"
          :items-per-page="perPage"
          :total="total"
        />
      </div>
    </div>

    <!-- Confirmation Delete Modal -->
    <CommonConfirmModal
      v-model:open="deleteModalOpen"
      title="Delete Customer"
      :description="`Are you sure you want to delete customer '${customerToDelete?.name}' (${customerToDelete?.customer_code})? All associated store and order history might be affected.`"
      confirm-text="Yes, Delete Customer"
      confirm-color="error"
      icon="i-lucide-alert-triangle"
      :loading="deleteLoading"
      @confirm="confirmDelete"
    />
  </div>
</template>
