<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui'

const { user, logout } = useAuth()

const open = ref(true)

const route = useRoute()

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Dashboard',
    icon: 'i-lucide-layout-dashboard',
    to: '/',
    active: route.path === '/',
  },
  {
    label: 'Customers',
    icon: 'i-lucide-users',
    to: '/customers',
    active: route.path.startsWith('/customers'),
  },
  {
    label: 'Stores',
    icon: 'i-lucide-store',
    to: '/stores',
    active: route.path.startsWith('/stores'),
  },
  {
    label: 'Products',
    icon: 'i-lucide-package',
    to: '/products',
    active: route.path.startsWith('/products'),
  },
  {
    label: 'Orders',
    icon: 'i-lucide-shopping-cart',
    to: '/orders',
    active: route.path.startsWith('/orders'),
  },
])

const userMenuItems = computed<DropdownMenuItem[][]>(() => [
  [{
    label: user.value?.name || 'User',
    type: 'label',
  }],
  [{
    label: 'Logout',
    icon: 'i-lucide-log-out',
    onSelect: () => logout(),
  }],
])
</script>

<template>
  <div class="flex flex-1 min-h-screen">
    <USidebar
      v-model:open="open"
      collapsible="icon"
      :ui="{
        container: 'h-screen',
      }"
    >
      <template #header>
        <div class="flex items-center gap-2 overflow-hidden">
          <div class="size-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold text-sm shrink-0">
            R
          </div>
          <span class="font-semibold text-lg truncate">REPUTAP</span>
        </div>
      </template>

      <UNavigationMenu
        :items="items"
        orientation="vertical"
        :ui="{ link: 'p-1.5 overflow-hidden' }"
      />

      <template #footer>
        <UDropdownMenu
          :items="userMenuItems"
          :content="{ align: 'center', collisionPadding: 12 }"
          :ui="{ content: 'w-(--reka-dropdown-menu-trigger-width) min-w-48' }"
        >
          <UButton
            :label="user?.name"
            :avatar="{ icon: 'i-lucide-user' }"
            trailing-icon="i-lucide-chevrons-up-down"
            color="neutral"
            variant="ghost"
            square
            class="w-full data-[state=open]:bg-elevated overflow-hidden"
            :ui="{
              trailingIcon: 'text-dimmed ms-auto',
            }"
          />
        </UDropdownMenu>
      </template>
    </USidebar>

    <div class="flex-1 flex flex-col overflow-hidden bg-default">
      <div class="h-(--ui-header-height) shrink-0 flex items-center gap-2 px-4 border-b border-default">
        <UButton
          icon="i-lucide-panel-left"
          color="neutral"
          variant="ghost"
          aria-label="Toggle sidebar"
          @click="open = !open"
        />
      </div>

      <div class="flex-1 overflow-auto p-6">
        <slot />
      </div>
    </div>
  </div>
</template>
