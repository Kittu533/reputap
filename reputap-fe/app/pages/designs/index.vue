<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'

definePageMeta({
  middleware: 'auth',
})

interface DesignAsset {
  id: number
  order_id: number
  store_id: number
  file_path: string
  file_name: string
  file_type: string
  file_size: number
  version: number
  approval_status: 'draft' | 'pending_approval' | 'approved' | 'revision_required'
  uploaded_by: number
  created_at: string
}

const { apiFetch } = useApi()
const config = useRuntimeConfig()
const { accessToken } = useAuth()
const toast = useToast()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const statusFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<DesignAsset[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Pending Approval', value: 'pending_approval' },
  { label: 'Approved', value: 'approved' },
  { label: 'Revision Required', value: 'revision_required' },
]

const statusColor: Record<string, string> = {
  draft: 'neutral',
  pending_approval: 'warning',
  approved: 'success',
  revision_required: 'error',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const columns: TableColumn<DesignAsset>[] = [
  { accessorKey: 'id', header: 'ID' },
  { accessorKey: 'file_name', header: 'File Name' },
  {
    accessorKey: 'order_id',
    header: 'Order ID',
    cell: ({ row }) => `#${row.original.order_id}`,
  },
  {
    accessorKey: 'store_id',
    header: 'Store ID',
    cell: ({ row }) => `#${row.original.store_id}`,
  },
  { accessorKey: 'version', header: 'Version' },
  {
    accessorKey: 'approval_status',
    header: 'Approval Status',
    cell: ({ row }) => h(UBadge, {
      label: row.original.approval_status.replace(/_/g, ' '),
      color: statusColor[row.original.approval_status] || 'neutral',
      variant: 'subtle',
      size: 'xs',
      class: 'capitalize',
    }),
  },
  {
    accessorKey: 'created_at',
    header: 'Uploaded At',
    cell: ({ row }) => formatDate(row.original.created_at),
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      const asset = row.original
      const buttons = [
        h(UButton, {
          icon: 'i-lucide-download',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          onClick: () => downloadFile(asset),
        }),
      ]
      if (asset.approval_status !== 'approved') {
        buttons.push(
          h(UButton, {
            icon: 'i-lucide-check',
            color: 'success',
            variant: 'ghost',
            size: 'xs',
            onClick: () => approveAsset(asset.id),
          }),
        )
      }
      buttons.push(
        h(UButton, {
          icon: 'i-lucide-trash-2',
          color: 'error',
          variant: 'ghost',
          size: 'xs',
          onClick: () => deleteAsset(asset.id),
        }),
      )
      return h('div', { class: 'flex justify-end gap-1' }, buttons)
    },
  },
]

async function fetchData() {
  loading.value = true
  const query: Record<string, any> = {
    page: page.value,
    per_page: perPage.value,
  }
  if (statusFilter.value) query.status = statusFilter.value

  const response = await apiFetch<DesignAsset[]>('/api/v1/design-assets', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

async function downloadFile(asset: DesignAsset) {
  const apiBase = config.public.apiBase as string
  const res = await fetch(`${apiBase}/api/v1/design-assets/${asset.id}/download`, {
    headers: { Authorization: `Bearer ${accessToken.value}` },
  })
  if (!res.ok) {
    toast.add({ title: 'Download failed', color: 'error' })
    return
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = asset.file_name
  a.click()
  URL.revokeObjectURL(url)
}

async function approveAsset(id: number) {
  const response = await apiFetch(`/api/v1/design-assets/${id}/approval`, {
    method: 'PATCH',
    body: { status: 'approved' },
  })
  if (response.success) {
    toast.add({ title: 'Design approved', color: 'success' })
    fetchData()
  } else {
    toast.add({ title: 'Failed to approve', color: 'error' })
  }
}

async function deleteAsset(id: number) {
  if (!confirm('Are you sure you want to delete this design asset?')) return
  const response = await apiFetch(`/api/v1/design-assets/${id}`, { method: 'DELETE' })
  if (response.success) {
    toast.add({ title: 'Design deleted', color: 'success' })
    fetchData()
  } else {
    toast.add({ title: 'Failed to delete', color: 'error' })
  }
}

watch(statusFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())

// Upload slideover
const slideoverOpen = ref(false)
const uploading = ref(false)
const uploadForm = ref({
  order_id: undefined as number | undefined,
  store_id: undefined as number | undefined,
})
const fileInput = ref<File | null>(null)

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  fileInput.value = target.files?.[0] || null
}

async function submitUpload() {
  if (!fileInput.value || !uploadForm.value.order_id || !uploadForm.value.store_id) {
    toast.add({ title: 'Please fill all fields and select a file', color: 'warning' })
    return
  }

  uploading.value = true
  const apiBase = config.public.apiBase as string
  const formData = new FormData()
  formData.append('file', fileInput.value)
  formData.append('order_id', String(uploadForm.value.order_id))
  formData.append('store_id', String(uploadForm.value.store_id))

  try {
    const res = await fetch(`${apiBase}/api/v1/design-assets/upload`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken.value}` },
      body: formData,
    })
    if (res.ok) {
      toast.add({ title: 'Design uploaded successfully', color: 'success' })
      slideoverOpen.value = false
      uploadForm.value = { order_id: undefined, store_id: undefined }
      fileInput.value = null
      fetchData()
    } else {
      toast.add({ title: 'Upload failed', color: 'error' })
    }
  } catch {
    toast.add({ title: 'Upload failed', color: 'error' })
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Design Assets</h1>
      <UButton icon="i-lucide-upload" label="Upload Design" @click="slideoverOpen = true" />
    </div>

    <div class="flex items-center gap-4 mb-4">
      <USelect
        v-model="statusFilter"
        :items="statusOptions"
        placeholder="All Statuses"
        class="w-48"
      />
    </div>

    <UTable :data="data" :columns="columns" :loading="loading" class="w-full" />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination v-model:page="page" :items-per-page="perPage" :total="total" />
    </div>

    <USlideover v-model:open="slideoverOpen" title="Upload Design" side="right">
      <template #body>
        <div class="flex flex-col gap-4 p-4">
          <UFormField label="File">
            <input
              type="file"
              accept=".png,.jpg,.jpeg,.svg,.pdf"
              class="block w-full text-sm file:mr-4 file:rounded-md file:border-0 file:bg-primary-50 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-primary-700 hover:file:bg-primary-100"
              @change="onFileChange"
            >
          </UFormField>
          <UFormField label="Order ID">
            <UInput v-model.number="uploadForm.order_id" type="number" placeholder="Enter order ID" />
          </UFormField>
          <UFormField label="Store ID">
            <UInput v-model.number="uploadForm.store_id" type="number" placeholder="Enter store ID" />
          </UFormField>
          <div class="flex justify-end gap-2 pt-4">
            <UButton label="Cancel" color="neutral" variant="outline" @click="slideoverOpen = false" />
            <UButton label="Upload" :loading="uploading" @click="submitUpload" />
          </div>
        </div>
      </template>
    </USlideover>
  </div>
</template>
