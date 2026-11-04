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
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Quality Control</h1>
      <UButton icon="i-lucide-plus" label="New Inspection" @click="openCreate" />
    </div>

    <div class="flex items-center gap-4 mb-4">
      <USelect
        v-model="resultFilter"
        :items="resultOptions"
        placeholder="All Results"
        class="w-48"
      />
    </div>

    <UTable :data="data" :columns="columns" :loading="loading" class="w-full" />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination v-model:page="page" :items-per-page="perPage" :total="total" />
    </div>

    <!-- Create Slideover -->
    <USlideover v-model:open="createOpen" title="New Inspection" side="right">
      <template #body>
        <div class="flex flex-col gap-4">
          <UFormField label="Product Instance ID" required>
            <UInput v-model.number="form.product_instance_id" type="number" placeholder="Enter product instance ID" />
          </UFormField>

          <UFormField label="Result" required>
            <USelect
              v-model="form.result"
              :items="resultOptions"
              placeholder="Select result"
              class="w-full"
            />
          </UFormField>

          <UFormField label="QC Checklist">
            <div class="flex flex-col gap-2 mt-1">
              <label
                v-for="item in checklistItems"
                :key="item.key"
                class="flex items-center gap-2 cursor-pointer"
              >
                <input
                  v-model="form.checklist[item.key]"
                  type="checkbox"
                  class="rounded border-gray-300"
                >
                <span class="text-sm">{{ item.label }}</span>
              </label>
            </div>
          </UFormField>

          <UFormField label="Notes">
            <UTextarea v-model="form.notes" placeholder="Additional notes..." :rows="3" />
          </UFormField>

          <div class="flex justify-end gap-2 pt-4">
            <UButton label="Cancel" color="neutral" variant="ghost" @click="createOpen = false" />
            <UButton label="Submit" :loading="submitting" @click="submitCreate" />
          </div>
        </div>
      </template>
    </USlideover>

    <!-- View Slideover -->
    <USlideover v-model:open="viewOpen" title="QC Inspection Details" side="right">
      <template #body>
        <div v-if="viewLoading" class="flex items-center justify-center py-8">
          <span class="text-sm text-gray-500">Loading...</span>
        </div>
        <div v-else-if="viewRecord" class="flex flex-col gap-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <p class="text-xs text-gray-500">ID</p>
              <p class="font-medium">{{ viewRecord.id }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Serial Number</p>
              <p class="font-medium">{{ viewRecord.serial_number }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Product Instance ID</p>
              <p class="font-medium">{{ viewRecord.product_instance_id }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Store ID</p>
              <p class="font-medium">{{ viewRecord.store_id }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Inspected By</p>
              <p class="font-medium">#{{ viewRecord.inspected_by }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">Inspected At</p>
              <p class="font-medium">{{ formatDate(viewRecord.inspected_at) }}</p>
            </div>
          </div>

          <div>
            <p class="text-xs text-gray-500 mb-1">Result</p>
            <UBadge :color="resultColor[viewRecord.result] || 'neutral'" variant="subtle" class="capitalize">
              {{ formatResult(viewRecord.result) }}
            </UBadge>
          </div>

          <div>
            <p class="text-xs text-gray-500 mb-2">Checklist</p>
            <div class="flex flex-col gap-1">
              <div
                v-for="item in checklistItems"
                :key="item.key"
                class="flex items-center gap-2 text-sm"
              >
                <span v-if="viewRecord.checklist_data?.[item.key]" class="text-green-600">✓</span>
                <span v-else class="text-red-500">✗</span>
                <span>{{ item.label }}</span>
              </div>
            </div>
          </div>

          <div v-if="viewRecord.evidence_paths?.length">
            <p class="text-xs text-gray-500 mb-1">Evidence</p>
            <ul class="list-disc list-inside text-sm">
              <li v-for="path in viewRecord.evidence_paths" :key="path">{{ path }}</li>
            </ul>
          </div>

          <div>
            <p class="text-xs text-gray-500 mb-1">Notes</p>
            <p class="text-sm">{{ viewRecord.notes || '—' }}</p>
          </div>
        </div>
      </template>
    </USlideover>
  </div>
</template>
