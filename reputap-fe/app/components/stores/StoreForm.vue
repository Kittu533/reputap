<script setup lang="ts">
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
import type { Store, Customer } from '~/utils/api'

const props = defineProps<{
  store: Store | null
}>()

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { apiFetch } = useApi()
const toast = useToast()
const loading = ref(false)

// Fetch customers for dropdown
const customers = ref<Customer[]>([])
async function loadCustomers() {
  const response = await apiFetch<Customer[]>('/api/v1/customers', { query: { per_page: 100 } })
  if (response.success) {
    customers.value = response.data || []
  }
}
onMounted(() => loadCustomers())

const customerOptions = computed(() =>
  customers.value.map(c => ({ label: `${c.customer_code} — ${c.name}`, value: c.id }))
)

const categoryOptions = [
  { label: 'Coffee Shop', value: 'coffee_shop' },
  { label: 'Restaurant', value: 'restaurant' },
  { label: 'Barbershop', value: 'barbershop' },
  { label: 'Salon', value: 'salon' },
  { label: 'Clinic', value: 'clinic' },
  { label: 'Retail', value: 'retail' },
  { label: 'Hotel', value: 'hotel' },
  { label: 'Other', value: 'other' },
]

const state = reactive({
  customer_id: props.store?.customer_id || null as number | null,
  name: props.store?.name || '',
  category: props.store?.category || 'other',
  address: props.store?.address || '',
  google_business_name: props.store?.google_business_name || '',
  review_url: props.store?.review_url || '',
})

watch(() => props.store, (s) => {
  state.customer_id = s?.customer_id || null
  state.name = s?.name || ''
  state.category = s?.category || 'other'
  state.address = s?.address || ''
  state.google_business_name = s?.google_business_name || ''
  state.review_url = s?.review_url || ''
})

function validate(s: typeof state): FormError[] {
  const errors: FormError[] = []
  if (!props.store && !s.customer_id) errors.push({ name: 'customer_id', message: 'Customer is required' })
  if (!s.name) errors.push({ name: 'name', message: 'Store name is required' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  loading.value = true

  const isEdit = !!props.store
  let path: string
  let method: string
  let body: Record<string, any>

  if (isEdit) {
    path = `/api/v1/stores/${props.store!.id}`
    method = 'PUT'
    body = {
      name: event.data.name,
      category: event.data.category,
      address: event.data.address || undefined,
      google_business_name: event.data.google_business_name || undefined,
    }
  } else {
    path = '/api/v1/stores'
    method = 'POST'
    body = {
      customer_id: event.data.customer_id,
      name: event.data.name,
      category: event.data.category,
      address: event.data.address || undefined,
      google_business_name: event.data.google_business_name || undefined,
      review_url: event.data.review_url || undefined,
    }
  }

  const response = await apiFetch<Store>(path, { method, body })
  loading.value = false

  if (response.success) {
    toast.add({
      title: isEdit ? 'Updated' : 'Created',
      description: `Store ${response.data?.store_code} ${isEdit ? 'updated' : 'created'}`,
      color: 'success',
    })
    emit('saved')
  } else {
    toast.add({ title: 'Error', description: response.error?.message || 'Save failed', color: 'error' })
  }
}
</script>

<template>
  <UForm :validate="validate" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormField v-if="!store" label="Customer" name="customer_id" required>
      <USelect
        v-model="state.customer_id"
        :items="customerOptions"
        placeholder="Select customer"
        class="w-full"
      />
    </UFormField>

    <UFormField label="Store Name" name="name" required>
      <UInput v-model="state.name" placeholder="e.g. Kopi Senja Bandung" class="w-full" />
    </UFormField>

    <UFormField label="Category" name="category">
      <USelect v-model="state.category" :items="categoryOptions" class="w-full" />
    </UFormField>

    <UFormField label="Address" name="address">
      <UTextarea v-model="state.address" placeholder="Full address" class="w-full" :rows="2" />
    </UFormField>

    <UFormField label="Google Business Name" name="google_business_name">
      <UInput v-model="state.google_business_name" placeholder="As it appears on Google" class="w-full" />
    </UFormField>

    <UFormField v-if="!store" label="Google Review URL" name="review_url">
      <UInput v-model="state.review_url" placeholder="https://search.google.com/local/writereview?placeid=..." class="w-full" />
    </UFormField>

    <div class="flex justify-end gap-2 pt-4">
      <UButton color="neutral" variant="outline" label="Cancel" @click="emit('cancel')" />
      <UButton type="submit" :label="store ? 'Update' : 'Create'" :loading="loading" />
    </div>
  </UForm>
</template>
