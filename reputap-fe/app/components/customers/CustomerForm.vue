<script setup lang="ts">
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
import type { Customer } from '~/utils/api'

const props = defineProps<{
  customer: Customer | null
}>()

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { apiFetch } = useApi()
const toast = useToast()
const loading = ref(false)

const state = reactive({
  name: props.customer?.name || '',
  business_name: props.customer?.business_name || '',
  email: props.customer?.email || '',
  phone: props.customer?.phone || '',
  address: props.customer?.address || '',
  city: props.customer?.city || '',
  province: props.customer?.province || '',
  notes: props.customer?.notes || '',
})

watch(() => props.customer, (c) => {
  state.name = c?.name || ''
  state.business_name = c?.business_name || ''
  state.email = c?.email || ''
  state.phone = c?.phone || ''
  state.address = c?.address || ''
  state.city = c?.city || ''
  state.province = c?.province || ''
  state.notes = c?.notes || ''
})

function validate(s: typeof state): FormError[] {
  const errors: FormError[] = []
  if (!s.name) errors.push({ name: 'name', message: 'Name is required' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  loading.value = true

  const body: Record<string, any> = { name: event.data.name }
  if (event.data.business_name) body.business_name = event.data.business_name
  if (event.data.email) body.email = event.data.email
  if (event.data.phone) body.phone = event.data.phone
  if (event.data.address) body.address = event.data.address
  if (event.data.city) body.city = event.data.city
  if (event.data.province) body.province = event.data.province
  if (event.data.notes) body.notes = event.data.notes

  const isEdit = !!props.customer
  const path = isEdit ? `/api/v1/customers/${props.customer!.id}` : '/api/v1/customers'
  const method = isEdit ? 'PUT' : 'POST'

  const response = await apiFetch<Customer>(path, { method, body })
  loading.value = false

  if (response.success) {
    toast.add({
      title: isEdit ? 'Updated' : 'Created',
      description: `Customer ${response.data?.customer_code} ${isEdit ? 'updated' : 'created'}`,
      color: 'success',
    })
    emit('saved')
  }
  else {
    toast.add({
      title: 'Error',
      description: response.error?.message || 'Save failed',
      color: 'error',
    })
  }
}
</script>

<template>
  <UForm :validate="validate" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormField label="Name" name="name" required>
      <UInput v-model="state.name" placeholder="Customer name" class="w-full" />
    </UFormField>

    <UFormField label="Business Name" name="business_name">
      <UInput v-model="state.business_name" placeholder="Business/company name" class="w-full" />
    </UFormField>

    <div class="grid grid-cols-2 gap-4">
      <UFormField label="Email" name="email">
        <UInput v-model="state.email" type="email" placeholder="email@example.com" class="w-full" />
      </UFormField>

      <UFormField label="Phone" name="phone">
        <UInput v-model="state.phone" placeholder="08xxxxxxxxxx" class="w-full" />
      </UFormField>
    </div>

    <UFormField label="Address" name="address">
      <UTextarea v-model="state.address" placeholder="Full address" class="w-full" :rows="2" />
    </UFormField>

    <div class="grid grid-cols-2 gap-4">
      <UFormField label="City" name="city">
        <UInput v-model="state.city" placeholder="City" class="w-full" />
      </UFormField>

      <UFormField label="Province" name="province">
        <UInput v-model="state.province" placeholder="Province" class="w-full" />
      </UFormField>
    </div>

    <UFormField label="Notes" name="notes">
      <UTextarea v-model="state.notes" placeholder="Internal notes" class="w-full" :rows="2" />
    </UFormField>

    <div class="flex justify-end gap-2 pt-4">
      <UButton color="neutral" variant="outline" label="Cancel" @click="emit('cancel')" />
      <UButton type="submit" :label="customer ? 'Update' : 'Create'" :loading="loading" />
    </div>
  </UForm>
</template>
