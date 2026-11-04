<script setup lang="ts">
import type { Store } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const router = useRouter()
const { apiFetch } = useApi()
const toast = useToast()

const storeId = route.params.id as string
const store = ref<Store | null>(null)
const loading = ref(true)

async function loadStore() {
  loading.value = true
  try {
    const res = await apiFetch<Store>(`/api/v1/stores/${storeId}`)
    if (res.success && res.data) {
      store.value = res.data
    }
    else {
      toast.add({ title: 'Error', description: 'Store not found', color: 'error' })
      router.push('/stores')
    }
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStore()
})
</script>

<template>
  <div class="space-y-6 w-full">
    <CommonPageHeader
      :title="store ? `Edit: ${store.name}` : 'Edit Store'"
      description="Update store address, category, and review destination URL."
      :breadcrumbs="[
        { label: 'Stores', to: '/stores' },
        { label: store?.store_code || 'Edit' },
      ]"
    />

    <CommonCard>
      <div v-if="loading" class="py-12 text-center text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin mx-auto mb-2 text-blue-600" />
        <span class="text-xs font-medium">Loading store data...</span>
      </div>
      <StoresStoreForm
        v-else-if="store"
        :store="store"
        @saved="() => router.push('/stores')"
        @cancel="() => router.push('/stores')"
      />
    </CommonCard>
  </div>
</template>
