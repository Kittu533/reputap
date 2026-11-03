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
  <UForm :validate="validate" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormField label="Customer" name="customer_id" required>
      <USelect
        v-model="state.customer_id"
        :items="customerOptions"
        placeholder="Select customer"
        class="w-full"
      />
    </UFormField>

    <UFormField label="Store" name="store_id" required>
      <USelect
        v-model="state.store_id"
        :items="storeOptions"
        :placeholder="state.customer_id ? 'Select store' : 'Select customer first'"
        :disabled="!state.customer_id"
        class="w-full"
      />
    </UFormField>

    <fieldset>
      <legend class="text-sm font-medium text-(--ui-text) mb-1.5">Order Items</legend>
      <div class="space-y-3">
        <div
          v-for="(item, index) in state.items"
          :key="index"
          class="flex items-start gap-2 p-3 rounded-lg border border-(--ui-border)"
        >
          <div class="flex-1 grid grid-cols-3 gap-2">
            <UFormField label="Product" :name="`items.${index}.product_id`" required>
              <USelect
                v-model="item.product_id"
                :items="productOptions"
                placeholder="Select product"
                class="w-full"
                @update:model-value="onProductSelect(index)"
              />
            </UFormField>

            <UFormField label="Qty" :name="`items.${index}.quantity`" required>
              <UInputNumber v-model="item.quantity" :min="1" class="w-full" />
            </UFormField>

            <UFormField label="Unit Price" :name="`items.${index}.unit_price`" required>
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
        class="mt-2"
        @click="addItem"
      />
    </fieldset>

    <div class="rounded-lg border border-(--ui-border) p-3 space-y-1 text-sm">
      <div class="flex justify-between">
        <span class="text-(--ui-text-muted)">Subtotal</span>
        <span>{{ formatRp(subtotal) }}</span>
      </div>
      <UFormField label="Discount (Rp)" name="discount">
        <UInputNumber v-model="state.discount" :min="0" class="w-full" />
      </UFormField>
      <div class="flex justify-between font-semibold text-base pt-1 border-t border-(--ui-border)">
        <span>Total</span>
        <span>{{ formatRp(total) }}</span>
      </div>
    </div>

    <UFormField label="Notes" name="notes">
      <UTextarea v-model="state.notes" placeholder="Order notes (optional)" class="w-full" :rows="2" />
    </UFormField>

    <div class="flex justify-end gap-2 pt-4">
      <UButton color="neutral" variant="outline" label="Cancel" @click="emit('cancel')" />
      <UButton type="submit" label="Create Order" :loading="loading" />
    </div>
  </UForm>
</template>
