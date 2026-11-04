<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    total: number
    perPage?: number
  }>(),
  {
    perPage: 20,
  }
)

const page = defineModel<number>('page', { default: 1 })

const from = computed(() => (props.total === 0 ? 0 : (page.value - 1) * props.perPage + 1))
const to = computed(() => Math.min(page.value * props.perPage, props.total))
</script>

<template>
  <div class="px-6 py-3.5 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-slate-500 bg-white">
    <div>
      Showing <span class="font-semibold text-slate-900">{{ from }}</span> to <span class="font-semibold text-slate-900">{{ to }}</span> of <span class="font-semibold text-slate-900">{{ total }}</span> results
    </div>
    <div v-if="total > perPage">
      <UPagination
        v-model:page="page"
        :items-per-page="perPage"
        :total="total"
        size="sm"
      />
    </div>
  </div>
</template>
