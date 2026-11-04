<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const router = useRouter()
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      title="Create New Order"
      description="Place a client order for smart review stands, specify quantities, pricing, and destination store."
      :breadcrumbs="[
        { label: 'Orders', to: '/orders' },
        { label: 'New Order' },
      ]"
    />

    <CommonCard>
      <OrdersOrderForm
        @saved="() => router.push('/orders')"
        @cancel="() => router.push('/orders')"
      />
    </CommonCard>
  </div>
</template>
