<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'

const props = withDefaults(
  defineProps<{
    data: any[]
    columns?: TableColumn<any>[]
    loading?: boolean
    total?: number
    perPage?: number
    searchPlaceholder?: string
    showSearch?: boolean
    emptyTitle?: string
    emptyDescription?: string
    emptyActionLabel?: string
  }>(),
  {
    columns: undefined,
    loading: false,
    total: 0,
    perPage: 20,
    searchPlaceholder: 'Search records...',
    showSearch: true,
    emptyTitle: 'No records found',
    emptyDescription: 'No records match your search or filter criteria.',
    emptyActionLabel: '',
  }
)

const emit = defineEmits<{
  (e: 'emptyAction'): void
}>()

const search = defineModel<string>('search', { default: '' })
const page = defineModel<number>('page', { default: 1 })
</script>

<template>
  <div class="space-y-4 w-full">
    <!-- Search, Filters & Actions Toolbar -->
    <div
      v-if="showSearch || $slots.filters || $slots.actions"
      class="bg-white rounded-2xl border border-slate-200/80 p-4 shadow-xs flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 w-full"
    >
      <div class="flex flex-1 flex-wrap items-center gap-3">
        <!-- Search Input -->
        <div v-if="showSearch" class="relative w-full sm:max-w-sm">
          <UIcon name="i-lucide-search" class="size-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          <input
            v-model="search"
            type="text"
            :placeholder="searchPlaceholder"
            class="w-full pl-9.5 pr-8 py-2 text-sm bg-slate-50 border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
          >
          <button
            v-if="search"
            type="button"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
            @click="search = ''"
          >
            <UIcon name="i-lucide-x" class="size-3.5" />
          </button>
        </div>

        <!-- Custom Filters Slot -->
        <slot name="filters" />
      </div>

      <div class="flex items-center gap-3 justify-between sm:justify-end">
        <div v-if="total !== undefined" class="text-xs font-semibold text-slate-500 shrink-0">
          Total: <span class="text-slate-900 font-bold">{{ total }}</span>
        </div>
        <slot name="actions" />
      </div>
    </div>

    <!-- Data Table Container -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden w-full">
      <!-- Loading State -->
      <div v-if="loading" class="py-16 text-center text-slate-400">
        <UIcon name="i-lucide-loader-2" class="size-7 animate-spin mx-auto mb-2 text-blue-600" />
        <span class="text-xs font-medium">Loading data...</span>
      </div>

      <!-- Empty State -->
      <div v-else-if="!data || data.length === 0">
        <slot name="empty">
          <CommonEmptyState
            :title="emptyTitle"
            :description="emptyDescription"
            :action-label="emptyActionLabel"
            @action="emit('emptyAction')"
          />
        </slot>
      </div>

      <!-- Table Slot or Nuxt UI Table -->
      <div v-else class="overflow-x-auto w-full">
        <slot name="table">
          <UTable
            v-if="columns"
            :data="data"
            :columns="columns"
            class="w-full"
          >
            <!-- Forward dynamic cell slots -->
            <template v-for="(_, slotName) in $slots" #[slotName]="slotProps">
              <slot :name="slotName" v-bind="slotProps" />
            </template>
          </UTable>
          <slot v-else />
        </slot>
      </div>

      <!-- Pagination -->
      <CommonPagination
        v-if="total > 0"
        v-model:page="page"
        :total="total"
        :per-page="perPage"
      />
    </div>
  </div>
</template>
