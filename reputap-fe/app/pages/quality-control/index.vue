<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { QCRecord } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const resultFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<QCRecord[]>([])
const total = ref(0)
const loading = ref(false)

const resultOptions = [
  { label: 'Passed', value: 'passed' },
  { label: 'Failed', value: 'failed' },
  { label: 'Rework Required', value: 'rework_required' },
]

const resultColor: Record<string, string> = {
  passed: 'success',
  failed: 'error',
  rework_required: 'warning',
}

const checklistItems: { key: string; label: string }[] = [
  { key: 'visual_check_ok', label: 'Visual check OK' },
  { key: 'nfc_chip_responsive', label: 'NFC chip responsive' },
  { key: 'qr_code_scannable', label: 'QR code scannable' },
  { key: 'dimensions_correct', label: 'Dimensions correct' },
  { key: 'finishing_quality_ok', label: 'Finishing quality OK' },
  { key: 'no_defects_found', label: 'No defects found' },
]

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function truncate(text: string | null, maxLen = 30): string {
  if (!text) return '—'
  return text.length > maxLen ? `${text.substring(0, maxLen)}…` : text
}

function formatResult(result: string): string {
  return result.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

const columns: TableColumn<QCRecord>[] = [
  { accessorKey: 'id', header: 'ID' },
  { accessorKey: 'serial_number', header: 'Serial Number' },
  {
    accessorKey: 'result',
    header: 'Result',
    cell: ({ row }) => {
      const result = row.original.result
      return h(UBadge, { variant: 'subtle', color: resultColor[result] || 'neutral', class: 'capitalize' }, () => formatResult(result))
    },
  },
  {
    accessorKey: 'inspected_at',
    header: 'Inspected At',
    cell: ({ row }) => formatDate(row.original.inspected_at),
  },
  {
    accessorKey: 'notes',
    header: 'Notes',
    cell: ({ row }) => truncate(row.original.notes),
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      return h('div', { class: 'flex justify-end gap-1' }, [
        h(UButton, {
          icon: 'i-lucide-eye',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          onClick: (e: Event) => { e.stopPropagation(); openView(row.original) },
        }),
      ])
    },
  },
]

async function fetchData() {
  loading.value = true
  const query: Record<string, any> = {
    page: page.value,
    per_page: perPage.value,
  }
  if (resultFilter.value) query.result = resultFilter.value

  const response = await apiFetch<QCRecord[]>('/api/v1/quality-controls', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

watch(resultFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())

// Create slideover
const createOpen = ref(false)
const submitting = ref(false)
const form = reactive({
  product_instance_id: null as number | null,
  result: undefined as string | undefined,
  checklist: {
    visual_check_ok: false,
    nfc_chip_responsive: false,
    qr_code_scannable: false,
    dimensions_correct: false,
    finishing_quality_ok: false,
    no_defects_found: false,
  } as Record<string, boolean>,
  notes: '',
})

function resetForm() {
  form.product_instance_id = null
  form.result = undefined
  form.notes = ''
  for (const key of Object.keys(form.checklist)) {
    form.checklist[key] = false
  }
}

function openCreate() {
  resetForm()
  createOpen.value = true
}

async function submitCreate() {
  if (!form.product_instance_id || !form.result) {
    toast.add({ title: 'Please fill required fields', color: 'warning' })
    return
  }

  submitting.value = true
  const response = await apiFetch('/api/v1/quality-controls', {
    method: 'POST',
    body: {
      product_instance_id: form.product_instance_id,
      result: form.result,
      checklist_data: { ...form.checklist },
      notes: form.notes || undefined,
    },
  })

  if (response.success) {
    toast.add({ title: 'Inspection created', color: 'success' })
    createOpen.value = false
    fetchData()
  }
  else {
    toast.add({ title: response.error?.message || 'Failed to create inspection', color: 'error' })
  }
  submitting.value = false
}

// View slideover
const viewOpen = ref(false)
const viewRecord = ref<QCRecord | null>(null)
const viewLoading = ref(false)

async function openView(record: QCRecord) {
  viewLoading.value = true
  viewOpen.value = true
  const response = await apiFetch<QCRecord>(`/api/v1/quality-controls/${record.id}`)
  if (response.success) {
    viewRecord.value = response.data
  }
  else {
    viewRecord.value = record
  }
  viewLoading.value = false
}
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="Quality Control"
      description="Inspect finished hardware, verify NFC signals, and record defect evaluations."
      :breadcrumbs="[{ label: 'Quality Control' }]"
    >
      <template #actions>
        <CommonButton icon="i-lucide-plus" @click="openCreate">
          New Inspection
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
      empty-title="No inspections found"
      empty-description="There are currently no QC inspection records matching your filter."
      empty-action-label="New Inspection"
      @empty-action="openCreate"
    >
      <template #filters>
        <USelect
          v-model="resultFilter"
          :items="resultOptions"
          placeholder="All Results"
          class="w-48"
        />
      </template>
    </CommonDataTable>

    <!-- Create Inspection Modal -->
    <CommonModal
      v-model:open="createOpen"
      title="New QC Inspection"
      description="Record a quality verification test and checklist outcome."
    >
      <div class="space-y-4">
        <CommonFormField label="Product Instance ID" required>
          <CommonInput
            v-model.number="form.product_instance_id"
            type="number"
            placeholder="Enter product instance ID"
          />
        </CommonFormField>

        <CommonFormField label="Result" required>
          <USelect
            v-model="form.result"
            :items="resultOptions"
            placeholder="Select result"
            class="w-full"
          />
        </CommonFormField>

        <div>
          <label class="block text-xs font-bold uppercase tracking-wider text-slate-700 mb-2">QC Checklist</label>
          <div class="space-y-2 bg-slate-50 p-3 rounded-xl border border-slate-200/60">
            <label
              v-for="item in checklistItems"
              :key="item.key"
              class="flex items-center gap-2.5 text-xs font-medium text-slate-700 cursor-pointer"
            >
              <input
                v-model="form.checklist[item.key]"
                type="checkbox"
                class="size-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500/20"
              >
              <span>{{ item.label }}</span>
            </label>
          </div>
        </div>

        <CommonFormField label="Notes">
          <textarea
            v-model="form.notes"
            rows="3"
            placeholder="Additional inspection notes..."
            class="w-full p-3 text-sm bg-white border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
          />
        </CommonFormField>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-2 w-full">
          <CommonButton variant="outline" @click="createOpen = false">
            Cancel
          </CommonButton>
          <CommonButton :loading="submitting" @click="submitCreate">
            Submit Inspection
          </CommonButton>
        </div>
      </template>
    </CommonModal>

    <!-- View Inspection Modal -->
    <CommonModal
      v-model:open="viewOpen"
      title="QC Inspection Details"
      :description="viewRecord ? `Inspection Record #${viewRecord.id}` : ''"
    >
      <div v-if="viewLoading" class="flex items-center justify-center py-8 text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin mr-2 text-blue-600" />
        <span class="text-xs font-medium">Loading details...</span>
      </div>
      <div v-else-if="viewRecord" class="space-y-4">
        <div class="grid grid-cols-2 gap-3 bg-slate-50 p-4 rounded-xl border border-slate-200/60">
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">ID</p>
            <p class="text-sm font-semibold text-slate-900">{{ viewRecord.id }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Serial Number</p>
            <p class="text-sm font-semibold text-slate-900">{{ viewRecord.serial_number }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Product Instance ID</p>
            <p class="text-sm font-semibold text-slate-900">{{ viewRecord.product_instance_id }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Store ID</p>
            <p class="text-sm font-semibold text-slate-900">{{ viewRecord.store_id }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Inspected By</p>
            <p class="text-sm font-semibold text-slate-900">#{{ viewRecord.inspected_by }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Inspected At</p>
            <p class="text-sm font-semibold text-slate-900">{{ formatDate(viewRecord.inspected_at) }}</p>
          </div>
        </div>

        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Result</p>
          <UBadge :color="resultColor[viewRecord.result] || 'neutral'" variant="subtle" class="capitalize">
            {{ formatResult(viewRecord.result) }}
          </UBadge>
        </div>

        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-2">Checklist Items</p>
          <div class="space-y-1.5 bg-slate-50 p-3 rounded-xl border border-slate-200/60">
            <div
              v-for="item in checklistItems"
              :key="item.key"
              class="flex items-center gap-2 text-xs font-medium"
            >
              <UIcon
                v-if="viewRecord.checklist_data?.[item.key]"
                name="i-lucide-check-circle-2"
                class="size-4 text-emerald-600"
              />
              <UIcon
                v-else
                name="i-lucide-x-circle"
                class="size-4 text-rose-500"
              />
              <span :class="viewRecord.checklist_data?.[item.key] ? 'text-slate-800' : 'text-slate-500'">{{ item.label }}</span>
            </div>
          </div>
        </div>

        <div v-if="viewRecord.evidence_paths?.length">
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Evidence</p>
          <ul class="list-disc list-inside text-xs text-slate-600 space-y-0.5">
            <li v-for="path in viewRecord.evidence_paths" :key="path">{{ path }}</li>
          </ul>
        </div>

        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Notes</p>
          <p class="text-xs text-slate-600 bg-slate-50 p-3 rounded-xl border border-slate-200/60">{{ viewRecord.notes || 'No inspection notes.' }}</p>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end w-full">
          <CommonButton variant="outline" @click="viewOpen = false">
            Close
          </CommonButton>
        </div>
      </template>
    </CommonModal>
  </div>
</template>
