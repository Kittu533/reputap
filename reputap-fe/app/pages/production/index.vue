<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'

definePageMeta({
  middleware: 'auth',
})

interface ProductionJob {
  id: number
  product_instance_id: number
  assigned_to: number | null
  status: 'pending' | 'in_progress' | 'completed' | 'failed'
  checklist_data: Record<string, boolean> | null
  started_at: string | null
  completed_at: string | null
  notes: string | null
  created_at: string
  updated_at: string
  serial_number: string
  store_id: number
}

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const statusFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<ProductionJob[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Completed', value: 'completed' },
  { label: 'Failed', value: 'failed' },
]

const statusColor: Record<string, string> = {
  pending: 'neutral',
  in_progress: 'info',
  completed: 'success',
  failed: 'error',
}

const checklistItems = [
  { key: 'akrilik_tersedia', label: 'Akrilik tersedia' },
  { key: 'cutting_selesai', label: 'Cutting selesai' },
  { key: 'finishing_selesai', label: 'Finishing selesai' },
  { key: 'vinyl_dicetak', label: 'Vinyl dicetak' },
  { key: 'nfc_tersedia', label: 'NFC tersedia' },
  { key: 'qr_tersedia', label: 'QR tersedia' },
  { key: 'assembly_selesai', label: 'Assembly selesai' },
  { key: 'packaging_selesai', label: 'Packaging selesai' },
]

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function getNextStatuses(current: string): { label: string, value: string }[] {
  if (current === 'pending') return [{ label: 'In Progress', value: 'in_progress' }]
  if (current === 'in_progress') return [{ label: 'Completed', value: 'completed' }, { label: 'Failed', value: 'failed' }]
  return []
}

const columns: TableColumn<ProductionJob>[] = [
  { accessorKey: 'id', header: 'ID' },
  { accessorKey: 'serial_number', header: 'Serial Number' },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => h(UBadge, {
      label: row.original.status.replace(/_/g, ' '),
      color: statusColor[row.original.status] || 'neutral',
      variant: 'subtle',
      size: 'xs',
      class: 'capitalize',
    }),
  },
  {
    accessorKey: 'assigned_to',
    header: 'Assigned To',
    cell: ({ row }) => row.original.assigned_to ? `#${row.original.assigned_to}` : '—',
  },
  {
    accessorKey: 'started_at',
    header: 'Started At',
    cell: ({ row }) => formatDate(row.original.started_at),
  },
  {
    accessorKey: 'completed_at',
    header: 'Completed At',
    cell: ({ row }) => formatDate(row.original.completed_at),
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      const job = row.original
      const nextStatuses = getNextStatuses(job.status)
      const buttons: any[] = []
      if (nextStatuses.length > 0) {
        buttons.push(
          h(UButton, {
            icon: 'i-lucide-arrow-right-circle',
            color: 'neutral',
            variant: 'ghost',
            size: 'xs',
            label: 'Update',
            onClick: () => openUpdateStatus(job),
          }),
        )
      }
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

  const response = await apiFetch<ProductionJob[]>('/api/v1/production/jobs', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

watch(statusFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())

// Create job slideover
const createOpen = ref(false)
const creating = ref(false)
const createForm = ref({
  product_instance_id: undefined as number | undefined,
  notes: '',
})

async function submitCreate() {
  if (!createForm.value.product_instance_id) {
    toast.add({ title: 'Please enter a product instance ID', color: 'warning' })
    return
  }

  creating.value = true
  const body: Record<string, any> = {
    product_instance_id: createForm.value.product_instance_id,
  }
  if (createForm.value.notes) body.notes = createForm.value.notes

  const response = await apiFetch('/api/v1/production/jobs', { method: 'POST', body })
  if (response.success) {
    toast.add({ title: 'Job created', color: 'success' })
    createOpen.value = false
    createForm.value = { product_instance_id: undefined, notes: '' }
    fetchData()
  } else {
    toast.add({ title: 'Failed to create job', color: 'error' })
  }
  creating.value = false
}

// Update status slideover
const updateOpen = ref(false)
const updating = ref(false)
const selectedJob = ref<ProductionJob | null>(null)
const updateForm = ref({
  status: undefined as string | undefined,
  notes: '',
  checklist_data: {} as Record<string, boolean>,
})

function openUpdateStatus(job: ProductionJob) {
  selectedJob.value = job
  const existingChecklist = job.checklist_data || {}
  const checklist: Record<string, boolean> = {}
  for (const item of checklistItems) {
    checklist[item.key] = existingChecklist[item.key] || false
  }
  updateForm.value = {
    status: undefined,
    notes: job.notes || '',
    checklist_data: checklist,
  }
  updateOpen.value = true
}

async function submitUpdate() {
  if (!selectedJob.value || !updateForm.value.status) {
    toast.add({ title: 'Please select a status', color: 'warning' })
    return
  }

  updating.value = true
  const body: Record<string, any> = {
    status: updateForm.value.status,
    checklist_data: updateForm.value.checklist_data,
  }
  if (updateForm.value.notes) body.notes = updateForm.value.notes

  const response = await apiFetch(`/api/v1/production/jobs/${selectedJob.value.id}/status`, {
    method: 'PATCH',
    body,
  })
  if (response.success) {
    toast.add({ title: 'Status updated', color: 'success' })
    updateOpen.value = false
    fetchData()
  } else {
    toast.add({ title: 'Failed to update status', color: 'error' })
  }
  updating.value = false
}
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="Production Queue"
      description="Track manufacturing assembly, NFC encoding, and hardware inspection jobs."
      :breadcrumbs="[{ label: 'Production' }]"
    >
      <template #actions>
        <CommonButton icon="i-lucide-plus" @click="createOpen = true">
          Create Job
        </CommonButton>
      </template>
    </CommonPageHeader>

    <CommonDataTable
      v-model:page="page"
      :data="data"
      :columns="columns"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      :show-search="false"
      empty-title="No production jobs"
      empty-description="There are currently no production jobs matching your filter."
      empty-action-label="Create Job"
      @empty-action="createOpen = true"
    >
      <template #filters>
        <USelect
          v-model="statusFilter"
          :items="statusOptions"
          placeholder="All Statuses"
          class="w-44"
        />
      </template>
    </CommonDataTable>

    <!-- Create Job Modal -->
    <CommonModal
      v-model:open="createOpen"
      title="Create Production Job"
      description="Initiate an assembly and encoding job for a product instance."
    >
      <div class="space-y-4">
        <CommonFormField label="Product Instance ID" required>
          <CommonInput
            v-model.number="createForm.product_instance_id"
            type="number"
            placeholder="Enter product instance ID"
          />
        </CommonFormField>
        <CommonFormField label="Notes">
          <textarea
            v-model="createForm.notes"
            rows="3"
            placeholder="Optional production notes"
            class="w-full p-3 text-sm bg-white border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
          />
        </CommonFormField>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2 w-full">
          <CommonButton variant="outline" @click="createOpen = false">
            Cancel
          </CommonButton>
          <CommonButton :loading="creating" @click="submitCreate">
            Create Job
          </CommonButton>
        </div>
      </template>
    </CommonModal>

    <!-- Update Status Modal -->
    <CommonModal
      v-model:open="updateOpen"
      title="Update Job Status"
      :description="selectedJob ? `Job #${selectedJob.id} — ${selectedJob.serial_number}` : ''"
    >
      <div v-if="selectedJob" class="space-y-4">
        <CommonFormField label="New Status" required>
          <USelect
            v-model="updateForm.status"
            :items="getNextStatuses(selectedJob.status)"
            placeholder="Select status"
            class="w-full"
          />
        </CommonFormField>

        <CommonFormField label="Notes">
          <textarea
            v-model="updateForm.notes"
            rows="3"
            placeholder="Optional status notes"
            class="w-full p-3 text-sm bg-white border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
          />
        </CommonFormField>

        <div>
          <label class="block text-xs font-bold uppercase tracking-wider text-slate-700 mb-2">Checklist</label>
          <div class="space-y-2 bg-slate-50 p-3 rounded-xl border border-slate-200/60">
            <label
              v-for="item in checklistItems"
              :key="item.key"
              class="flex items-center gap-2.5 text-xs font-medium text-slate-700 cursor-pointer"
            >
              <input
                v-model="updateForm.checklist_data[item.key]"
                type="checkbox"
                class="size-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500/20"
              >
              <span>{{ item.label }}</span>
            </label>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2 w-full">
          <CommonButton variant="outline" @click="updateOpen = false">
            Cancel
          </CommonButton>
          <CommonButton :loading="updating" @click="submitUpdate">
            Save Changes
          </CommonButton>
        </div>
      </template>
    </CommonModal>
  </div>
</template>
