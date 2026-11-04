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
  <div class="space-y-6">
    <CommonPageHeader
      title="Audit Logs"
      description="Track system mutation events, staff activity, and historical entity revisions."
      :breadcrumbs="[{ label: 'Audit Logs' }]"
    />

    <CommonDataTable
      v-model:page="page"
      :data="data"
      :columns="columns"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      :show-search="false"
      empty-title="No audit logs found"
      empty-description="No event records match your current filter parameters."
    >
      <template #filters>
        <div class="flex items-center gap-2 flex-wrap">
          <input
            v-model="entityTypeFilter"
            type="text"
            placeholder="Entity type..."
            class="px-3 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 w-36"
          >
          <input
            v-model="actionFilter"
            type="text"
            placeholder="Action..."
            class="px-3 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 w-32"
          >
          <input
            v-model="dateFrom"
            type="date"
            class="px-3 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-xl text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 w-36"
          >
          <input
            v-model="dateTo"
            type="date"
            class="px-3 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-xl text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 w-36"
          >
        </div>
      </template>
    </CommonDataTable>

    <!-- Changes Modal -->
    <CommonModal
      v-model:open="changesOpen"
      title="Audit Log Changes"
      :description="selectedLog ? `Event #${selectedLog.id} — ${selectedLog.action} on ${selectedLog.entity_type}` : ''"
    >
      <div v-if="selectedLog" class="space-y-4">
        <div class="grid grid-cols-2 gap-3 bg-slate-50 p-4 rounded-xl border border-slate-200/60">
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Log ID</p>
            <p class="text-sm font-semibold text-slate-900">{{ selectedLog.id }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">User ID</p>
            <p class="text-sm font-semibold text-slate-900">#{{ selectedLog.user_id }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Action</p>
            <p class="text-sm font-semibold text-slate-900">{{ selectedLog.action }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Entity</p>
            <p class="text-sm font-semibold text-slate-900">{{ selectedLog.entity_type }} {{ selectedLog.entity_id != null ? `#${selectedLog.entity_id}` : '' }}</p>
          </div>
          <div class="col-span-2">
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Created At</p>
            <p class="text-sm font-semibold text-slate-900">{{ formatDate(selectedLog.created_at) }}</p>
          </div>
        </div>

        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Old Values</p>
          <pre class="bg-slate-900 text-emerald-400 rounded-xl p-3 text-xs overflow-auto max-h-52 font-mono leading-relaxed">{{ formatJson(selectedLog.old_values) }}</pre>
        </div>

        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">New Values</p>
          <pre class="bg-slate-900 text-blue-400 rounded-xl p-3 text-xs overflow-auto max-h-52 font-mono leading-relaxed">{{ formatJson(selectedLog.new_values) }}</pre>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end w-full">
          <CommonButton variant="outline" @click="changesOpen = false">
            Close
          </CommonButton>
        </div>
      </template>
    </CommonModal>
  </div>
</template>
