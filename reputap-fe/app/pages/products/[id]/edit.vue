<script setup lang="ts">
import type { Product } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const router = useRouter()
const { apiFetch } = useApi()
const toast = useToast()

const productId = route.params.id as string
const product = ref<Product | null>(null)
const loading = ref(true)

async function loadProduct() {
  loading.value = true
  try {
    const res = await apiFetch<Product>(`/api/v1/products/${productId}`)
    if (res.success && res.data) {
      product.value = res.data
    }
    else {
      toast.add({ title: 'Error', description: 'Product not found', color: 'error' })
      router.push('/products')
    }
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProduct()
})
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      :title="product ? `Edit: ${product.name}` : 'Edit Product'"
      description="Update stand specifications, materials, and selling price."
      :breadcrumbs="[
        { label: 'Products', to: '/products' },
        { label: product?.sku || 'Edit' },
      ]"
    />

    <CommonCard>
      <div v-if="loading" class="py-12 text-center text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin mx-auto mb-2 text-blue-600" />
        <span class="text-xs font-medium">Loading product data...</span>
      </div>
      <ProductsProductForm
        v-else-if="product"
        :product="product"
        @saved="() => router.push('/products')"
        @cancel="() => router.push('/products')"
      />
    </CommonCard>
  </div>
</template>
