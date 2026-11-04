<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const router = useRouter()
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      title="Add New Customer"
      description="Register a new business client, contact person, and billing address."
      :breadcrumbs="[
        { label: 'Customers', to: '/customers' },
        { label: 'New Customer' },
      ]"
    />

    <CommonCard>
      <CustomersCustomerForm
        :customer="null"
        @saved="() => router.push('/customers')"
        @cancel="() => router.push('/customers')"
      />
    </CommonCard>
  </div>
</template>
