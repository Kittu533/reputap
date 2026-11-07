<script setup lang="ts">
interface BreadcrumbItem {
  label: string
  to?: string
}

defineProps<{
  title: string
  description?: string
  breadcrumbs?: BreadcrumbItem[]
}>()
</script>

<template>
  <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
    <div>
      <!-- Breadcrumbs -->
      <nav v-if="breadcrumbs && breadcrumbs.length > 0" class="flex items-center gap-1.5 text-xs font-medium text-slate-400 mb-2">
        <NuxtLink to="/dashboard" class="hover:text-slate-700 transition-colors flex items-center gap-1">
          <UIcon name="i-lucide-home" class="size-3.5" />
          <span>Home</span>
        </NuxtLink>
        <template v-for="(item, idx) in breadcrumbs" :key="idx">
          <UIcon name="i-lucide-chevron-right" class="size-3 text-slate-300" />
          <NuxtLink
            v-if="item.to"
            :to="item.to"
            class="hover:text-slate-700 transition-colors"
          >
            {{ item.label }}
          </NuxtLink>
          <span v-else class="text-slate-700 font-semibold">{{ item.label }}</span>
        </template>
      </nav>

      <!-- Title & Description -->
      <h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
        {{ title }}
      </h1>
      <p v-if="description" class="mt-1 text-sm text-slate-500 leading-relaxed">
        {{ description }}
      </p>
    </div>

    <!-- Actions slot -->
    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2.5 shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>
