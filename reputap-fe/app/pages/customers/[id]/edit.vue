<script setup lang="ts">
import type { Customer } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const router = useRouter()
const { apiFetch } = useApi()
const toast = useToast()

const customerId = route.params.id as string
const customer = ref<Customer | null>(null)
const loading = ref(true)

async function loadCustomer() {
  loading.value = true
  try {
    const res = await apiFetch<Customer>(`/api/v1/customers/${customerId}`)
    if (res.success && res.data) {
      customer.value = res.data
    }
    else {
      toast.add({ title: 'Error', description: 'Customer not found', color: 'error' })
      router.push('/customers')
    }
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCustomer()
})
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      :title="customer ? `Edit: ${customer.name}` : 'Edit Customer'"
      description="Update customer contact details and business info."
      :breadcrumbs="[
        { label: 'Customers', to: '/customers' },
        { label: customer?.customer_code || 'Edit' },
      ]"
    />

    <CommonCard>
      <div v-if="loading" class="py-12 text-center text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin mx-auto mb-2 text-blue-600" />
        <span class="text-xs font-medium">Loading customer data...</span>
      </div>
      <CustomersCustomerForm
        v-else-if="customer"
        :customer="customer"
        @saved="() => router.push('/customers')"
        @cancel="() => router.push('/customers')"
      />
    </CommonCard>
  </div>
</template>
