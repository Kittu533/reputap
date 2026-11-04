<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { AuditLog } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()

const UButton = resolveComponent('UButton')

const entityTypeFilter = ref('')
const actionFilter = ref('')
const dateFrom = ref('')
const dateTo = ref('')
const page = ref(1)
const perPage = ref(20)
const data = ref<AuditLog[]>([])
const total = ref(0)
const loading = ref(false)

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const columns: TableColumn<AuditLog>[] = [
  { accessorKey: 'id', header: 'ID' },
  {
    accessorKey: 'user_id',
    header: 'User ID',
    cell: ({ row }) => `#${row.original.user_id}`,
  },
  { accessorKey: 'action', header: 'Action' },
  { accessorKey: 'entity_type', header: 'Entity Type' },
  {
    accessorKey: 'entity_id',
    header: 'Entity ID',
    cell: ({ row }) => row.original.entity_id != null ? `#${row.original.entity_id}` : '—',
  },
  {
    id: 'changes',
    header: 'Changes',
    cell: ({ row }) => {
      const hasChanges = row.original.old_values || row.original.new_values
      if (!hasChanges) return '—'
      return h(UButton, {
        label: 'View',
        color: 'neutral',
        variant: 'link',
        size: 'xs',
        onClick: (e: Event) => { e.stopPropagation(); openChanges(row.original) },
      })
    },
  },
  {
    accessorKey: 'created_at',
    header: 'Created At',
    cell: ({ row }) => formatDate(row.original.created_at),
  },
]

async function fetchData() {
  loading.value = true
  const query: Record<string, any> = {
    page: page.value,
    per_page: perPage.value,
  }
  if (entityTypeFilter.value) query.entity_type = entityTypeFilter.value
  if (actionFilter.value) query.action = actionFilter.value
  if (dateFrom.value) query.date_from = dateFrom.value
  if (dateTo.value) query.date_to = dateTo.value

  const response = await apiFetch<AuditLog[]>('/api/v1/audit-logs', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

let entityTypeTimeout: ReturnType<typeof setTimeout>
watch(entityTypeFilter, () => {
  clearTimeout(entityTypeTimeout)
  entityTypeTimeout = setTimeout(() => { page.value = 1; fetchData() }, 300)
})

let actionTimeout: ReturnType<typeof setTimeout>
watch(actionFilter, () => {
  clearTimeout(actionTimeout)
  actionTimeout = setTimeout(() => { page.value = 1; fetchData() }, 300)
})

watch(dateFrom, () => { page.value = 1; fetchData() })
watch(dateTo, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())

// Changes slideover
const changesOpen = ref(false)
const selectedLog = ref<AuditLog | null>(null)

function openChanges(log: AuditLog) {
  selectedLog.value = log
  changesOpen.value = true
}

function formatJson(value: Record<string, any> | null): string {
  if (!value) return 'null'
  return JSON.stringify(value, null, 2)
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Audit Logs</h1>
    </div>

    <div class="flex items-center gap-4 mb-4 flex-wrap">
      <UInput
        v-model="entityTypeFilter"
        placeholder="Filter entity type..."
        icon="i-lucide-search"
        class="w-48"
      />
      <UInput
        v-model="actionFilter"
        placeholder="Filter action..."
        icon="i-lucide-search"
        class="w-48"
      />
      <UInput
        v-model="dateFrom"
        type="date"
        class="w-44"
      />
      <UInput
        v-model="dateTo"
        type="date"
        class="w-44"
      />
    </div>

    <UTable :data="data" :columns="columns" :loading="loading" class="w-full" />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination v-model:page="page" :items-per-page="perPage" :total="total" />
    </div>

    <!-- Changes Slideover -->
    <USlideover v-model:open="changesOpen" title="Audit Log Changes" side="right">
      <template #body>
        <div v-if="selectedLog" class="flex flex-col gap-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <p class="text-xs text-gray-500">Log ID</p>
              <p class="font-medium">{{ selectedLog.id }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">User ID</p>
              <p class="font-medium">#{{ selectedLog.user_id }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Action</p>
              <p class="font-medium">{{ selectedLog.action }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Entity</p>
              <p class="font-medium">{{ selectedLog.entity_type }} {{ selectedLog.entity_id != null ? `#${selectedLog.entity_id}` : '' }}</p>
            </div>
            <div class="col-span-2">
              <p class="text-xs text-gray-500">Created At</p>
              <p class="font-medium">{{ formatDate(selectedLog.created_at) }}</p>
            </div>
          </div>

          <div>
            <p class="text-xs text-gray-500 mb-1">Old Values</p>
            <pre class="bg-gray-100 dark:bg-gray-800 rounded p-3 text-xs overflow-auto max-h-64">{{ formatJson(selectedLog.old_values) }}</pre>
          </div>

          <div>
            <p class="text-xs text-gray-500 mb-1">New Values</p>
            <pre class="bg-gray-100 dark:bg-gray-800 rounded p-3 text-xs overflow-auto max-h-64">{{ formatJson(selectedLog.new_values) }}</pre>
          </div>
        </div>
      </template>
    </USlideover>
  </div>
</template>
