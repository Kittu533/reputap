<script setup lang="ts">
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
import type { Product } from '~/utils/api'

const props = defineProps<{
  product: Product | null
}>()

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { apiFetch } = useApi()
const toast = useToast()
const loading = ref(false)

const state = reactive({
  sku: props.product?.sku || '',
  name: props.product?.name || '',
  description: props.product?.description || '',
  material: props.product?.material || '',
  width_mm: props.product?.width_mm ?? null as number | null,
  height_mm: props.product?.height_mm ?? null as number | null,
  thickness_mm: props.product?.thickness_mm ?? null as number | null,
  nfc_chip_type: props.product?.nfc_chip_type || 'NTAG213',
  finishing: props.product?.finishing || '',
  production_cost: props.product?.production_cost ?? null as number | null,
  selling_price: props.product?.selling_price ?? null as number | null,
})

watch(() => props.product, (p) => {
  state.sku = p?.sku || ''
  state.name = p?.name || ''
  state.description = p?.description || ''
  state.material = p?.material || ''
  state.width_mm = p?.width_mm ?? null
  state.height_mm = p?.height_mm ?? null
  state.thickness_mm = p?.thickness_mm ?? null
  state.nfc_chip_type = p?.nfc_chip_type || 'NTAG213'
  state.finishing = p?.finishing || ''
  state.production_cost = p?.production_cost ?? null
  state.selling_price = p?.selling_price ?? null
})

function validate(s: typeof state): FormError[] {
  const errors: FormError[] = []
  if (!props.product && !s.sku) errors.push({ name: 'sku', message: 'SKU is required' })
  if (!s.name) errors.push({ name: 'name', message: 'Name is required' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  loading.value = true

  const isEdit = !!props.product
  const path = isEdit ? `/api/v1/products/${props.product!.id}` : '/api/v1/products'
  const method = isEdit ? 'PUT' : 'POST'

  const body: Record<string, any> = {
    name: event.data.name,
  }

  if (!isEdit) body.sku = event.data.sku
  if (event.data.description) body.description = event.data.description
  if (event.data.material) body.material = event.data.material
  if (event.data.width_mm != null) body.width_mm = event.data.width_mm
  if (event.data.height_mm != null) body.height_mm = event.data.height_mm
  if (event.data.thickness_mm != null) body.thickness_mm = event.data.thickness_mm
  if (event.data.nfc_chip_type) body.nfc_chip_type = event.data.nfc_chip_type
  if (event.data.finishing) body.finishing = event.data.finishing
  if (event.data.production_cost != null) body.production_cost = event.data.production_cost
  if (event.data.selling_price != null) body.selling_price = event.data.selling_price

  const response = await apiFetch<Product>(path, { method, body })
  loading.value = false

  if (response.success) {
    toast.add({
      title: isEdit ? 'Updated' : 'Created',
      description: `Product ${response.data?.sku} ${isEdit ? 'updated' : 'created'}`,
      color: 'success',
    })
    emit('saved')
  } else {
    toast.add({
      title: 'Error',
      description: response.error?.message || 'Save failed',
      color: 'error',
    })
  }
}
</script>

<template>
  <UForm :validate="validate" :state="state" class="space-y-5 w-full" @submit="onSubmit">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <UFormField label="SKU / Item Code" name="sku" :required="!product" class="w-full">
        <UInput
          v-model="state.sku"
          placeholder="e.g. STAND-ACRYLIC-01"
          class="w-full"
          :readonly="!!product"
          :disabled="!!product"
        />
      </UFormField>

      <UFormField label="Product Name" name="name" required class="w-full">
        <UInput v-model="state.name" placeholder="e.g. Premium Acrylic NFC Stand" class="w-full" />
      </UFormField>
    </div>

    <UFormField label="Description" name="description" class="w-full">
      <UTextarea v-model="state.description" placeholder="Product description, stand features and materials" class="w-full" :rows="3" />
    </UFormField>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <UFormField label="Material" name="material" class="w-full">
        <UInput v-model="state.material" placeholder="e.g. Acrylic, Wood, Metal" class="w-full" />
      </UFormField>

      <UFormField label="Finishing" name="finishing" class="w-full">
        <UInput v-model="state.finishing" placeholder="e.g. Glossy, Matte, Laser Engraved" class="w-full" />
      </UFormField>

      <UFormField label="NFC Chip Type" name="nfc_chip_type" class="w-full">
        <UInput v-model="state.nfc_chip_type" placeholder="NTAG213 / NTAG215" class="w-full" />
      </UFormField>
    </div>

    <fieldset class="p-4 rounded-xl border border-slate-200/80 bg-slate-50/50">
      <legend class="text-xs font-bold uppercase tracking-wider text-slate-500 px-1">Dimensions (mm)</legend>
      <div class="grid grid-cols-3 gap-4 mt-2">
        <UFormField label="Width (mm)" name="width_mm" class="w-full">
          <UInputNumber v-model="state.width_mm" placeholder="e.g. 100" :min="0" class="w-full" />
        </UFormField>
        <UFormField label="Height (mm)" name="height_mm" class="w-full">
          <UInputNumber v-model="state.height_mm" placeholder="e.g. 150" :min="0" class="w-full" />
        </UFormField>
        <UFormField label="Thickness (mm)" name="thickness_mm" class="w-full">
          <UInputNumber v-model="state.thickness_mm" placeholder="e.g. 5" :min="0" class="w-full" />
        </UFormField>
      </div>
    </fieldset>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <UFormField label="Production Cost (IDR)" name="production_cost" class="w-full">
        <UInputNumber v-model="state.production_cost" placeholder="0" :min="0" class="w-full" />
      </UFormField>

      <UFormField label="Selling Price (IDR)" name="selling_price" class="w-full">
        <UInputNumber v-model="state.selling_price" placeholder="0" :min="0" class="w-full" />
      </UFormField>
    </div>

    <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100">
      <UButton color="neutral" variant="outline" label="Cancel" @click="emit('cancel')" />
      <UButton type="submit" :label="product ? 'Save Changes' : 'Create Product'" color="primary" :loading="loading" />
    </div>
  </UForm>
</template>
