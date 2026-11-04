<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const router = useRouter()
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      title="Add New Product"
      description="Define a new smart review stand product specification, dimensions, NFC chip, and pricing."
      :breadcrumbs="[
        { label: 'Products', to: '/products' },
        { label: 'New Product' },
      ]"
    />

    <CommonCard>
      <ProductsProductForm
        :product="null"
        @saved="() => router.push('/products')"
        @cancel="() => router.push('/products')"
      />
    </CommonCard>
  </div>
</template>
