<script setup lang="ts">
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
import type { Customer, Store, Product } from '~/utils/api'

const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { apiFetch } = useApi()
const toast = useToast()
const loading = ref(false)

const customers = ref<Customer[]>([])
const stores = ref<Store[]>([])
const products = ref<Product[]>([])

async function loadCustomers() {
  const r = await apiFetch<Customer[]>('/api/v1/customers', { query: { per_page: 100 } })
  if (r.success) customers.value = r.data || []
}

async function loadStores(customerId: number) {
  const r = await apiFetch<Store[]>('/api/v1/stores', { query: { per_page: 100, customer_id: customerId, status: 'active' } })
  if (r.success) stores.value = r.data || []
}

async function loadProducts() {
  const r = await apiFetch<Product[]>('/api/v1/products', { query: { per_page: 100, status: 'active' } })
  if (r.success) products.value = r.data || []
}

onMounted(() => {
  loadCustomers()
  loadProducts()
})

const state = reactive({
  customer_id: null as number | null,
  store_id: null as number | null,
  discount: 0,
  notes: '',
  items: [{ product_id: null as number | null, quantity: 1, unit_price: 0 }] as { product_id: number | null; quantity: number; unit_price: number }[],
})

watch(() => state.customer_id, (id) => {
  state.store_id = null
  stores.value = []
  if (id) loadStores(id)
})

function onProductSelect(index: number) {
  const item = state.items[index]
  if (item.product_id) {
    const product = products.value.find(p => p.id === item.product_id)
    if (product?.selling_price) {
      item.unit_price = product.selling_price
    }
  }
}

function addItem() {
  state.items.push({ product_id: null, quantity: 1, unit_price: 0 })
}

function removeItem(index: number) {
  if (state.items.length > 1) state.items.splice(index, 1)
}

const subtotal = computed(() => state.items.reduce((sum, item) => sum + item.quantity * item.unit_price, 0))
const total = computed(() => Math.max(0, subtotal.value - state.discount))

const customerOptions = computed(() => customers.value.map(c => ({ label: `${c.customer_code} — ${c.name}`, value: c.id })))
const storeOptions = computed(() => stores.value.map(s => ({ label: `${s.store_code} — ${s.name}`, value: s.id })))
const productOptions = computed(() => products.value.map(p => ({ label: `${p.sku} — ${p.name}`, value: p.id })))

function formatRp(n: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n)
}

function validate(s: typeof state): FormError[] {
  const errors: FormError[] = []
  if (!s.customer_id) errors.push({ name: 'customer_id', message: 'Customer is required' })
  if (!s.store_id) errors.push({ name: 'store_id', message: 'Store is required' })
  if (s.items.length === 0) errors.push({ name: 'items', message: 'At least one item is required' })
  for (let i = 0; i < s.items.length; i++) {
    if (!s.items[i].product_id) errors.push({ name: `items.${i}.product_id`, message: 'Product is required' })
    if (s.items[i].quantity < 1) errors.push({ name: `items.${i}.quantity`, message: 'Quantity must be at least 1' })
    if (s.items[i].unit_price <= 0) errors.push({ name: `items.${i}.unit_price`, message: 'Price must be greater than 0' })
  }
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  loading.value = true
  const body = {
    customer_id: event.data.customer_id,
    store_id: event.data.store_id,
    discount: event.data.discount,
    notes: event.data.notes || undefined,
    items: event.data.items.map(item => ({
      product_id: item.product_id,
      quantity: item.quantity,
      unit_price: item.unit_price,
    })),
  }

  const response = await apiFetch('/api/v1/orders', { method: 'POST', body })
  loading.value = false

  if (response.success) {
    toast.add({ title: 'Created', description: `Order ${response.data?.order_number} created`, color: 'success' })
    emit('saved')
  } else {
    toast.add({ title: 'Error', description: response.error?.message || 'Create failed', color: 'error' })
  }
}
</script>

<template>
  <UForm :validate="validate" :state="state" class="space-y-5 w-full" @submit="onSubmit">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <UFormField label="Customer" name="customer_id" required class="w-full">
        <USelect
          v-model="state.customer_id"
          :items="customerOptions"
          placeholder="Select customer"
          class="w-full"
        />
      </UFormField>

      <UFormField label="Destination Store" name="store_id" required class="w-full">
        <USelect
          v-model="state.store_id"
          :items="storeOptions"
          :placeholder="state.customer_id ? 'Select store' : 'Select customer first'"
          :disabled="!state.customer_id"
          class="w-full"
        />
      </UFormField>
    </div>

    <fieldset class="p-4 rounded-xl border border-slate-200/80 bg-slate-50/50">
      <legend class="text-xs font-bold uppercase tracking-wider text-slate-500 px-1">Order Items</legend>
      <div class="space-y-3 mt-2">
        <div
          v-for="(item, index) in state.items"
          :key="index"
          class="flex items-start gap-3 p-3 rounded-lg border border-slate-200 bg-white"
        >
          <div class="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-3">
            <UFormField label="Product" :name="`items.${index}.product_id`" required class="w-full">
              <USelect
                v-model="item.product_id"
                :items="productOptions"
                placeholder="Select product"
                class="w-full"
                @update:model-value="onProductSelect(index)"
              />
            </UFormField>

            <UFormField label="Quantity" :name="`items.${index}.quantity`" required class="w-full">
              <UInputNumber v-model="item.quantity" :min="1" class="w-full" />
            </UFormField>

            <UFormField label="Unit Price (Rp)" :name="`items.${index}.unit_price`" required class="w-full">
              <UInputNumber v-model="item.unit_price" :min="0" class="w-full" />
            </UFormField>
          </div>

          <UButton
            icon="i-lucide-x"
            color="error"
            variant="ghost"
            size="xs"
            class="mt-6"
            :disabled="state.items.length <= 1"
            @click="removeItem(index)"
          />
        </div>
      </div>
      <UButton
        icon="i-lucide-plus"
        label="Add Item"
        color="neutral"
        variant="outline"
        size="sm"
        class="mt-3"
        @click="addItem"
      />
    </fieldset>

    <div class="rounded-xl border border-slate-200 bg-slate-50/50 p-4 space-y-2 text-sm">
      <div class="flex justify-between text-slate-600">
        <span>Subtotal</span>
        <span class="font-semibold text-slate-800">{{ formatRp(subtotal) }}</span>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-center">
        <UFormField label="Discount (Rp)" name="discount" class="w-full">
          <UInputNumber v-model="state.discount" :min="0" class="w-full" />
        </UFormField>
        <div class="flex flex-col items-end justify-center pt-4">
          <span class="text-xs uppercase font-bold text-slate-400">Total Payable</span>
          <span class="text-xl font-bold text-blue-600">{{ formatRp(total) }}</span>
        </div>
      </div>
    </div>

    <UFormField label="Internal Order Notes" name="notes" class="w-full">
      <UTextarea v-model="state.notes" placeholder="Order notes (optional)" class="w-full" :rows="2" />
    </UFormField>

    <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100">
      <UButton color="neutral" variant="outline" label="Cancel" @click="emit('cancel')" />
      <UButton type="submit" label="Create Order" color="primary" :loading="loading" />
    </div>
  </UForm>
</template>
