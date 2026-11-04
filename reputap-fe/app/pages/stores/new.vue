<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const router = useRouter()
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      title="Register New Store"
      description="Add a physical store location, assign it to a client, and configure its Google review destination URL."
      :breadcrumbs="[
        { label: 'Stores', to: '/stores' },
        { label: 'New Store' },
      ]"
    />

    <CommonCard>
      <StoresStoreForm
        :store="null"
        @saved="() => router.push('/stores')"
        @cancel="() => router.push('/stores')"
      />
    </CommonCard>
  </div>
</template>
