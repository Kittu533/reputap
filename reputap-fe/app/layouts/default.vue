<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

const { user, logout } = useAuth()
const route = useRoute()
const colorMode = useColorMode()
const config = useRuntimeConfig()
const apiBase = (config.public.apiBase as string) || 'https://api.kitangoding.my.id'

const sidebarOpen = ref(true)

function toggleColorMode() {
  const isCurrentlyDark = colorMode.value === 'dark' || (typeof document !== 'undefined' && document.documentElement.classList.contains('dark'))
  const nextMode = isCurrentlyDark ? 'light' : 'dark'

  colorMode.preference = nextMode
  colorMode.value = nextMode
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('nuxt-color-mode', nextMode)
  }

  if (typeof document !== 'undefined') {
    if (nextMode === 'dark') {
      document.documentElement.classList.remove('light')
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
      document.documentElement.classList.add('light')
    }
  }
}

interface NavItem {
  label: string
  icon: string
  to: string
  badge?: string | number
  badgeColor?: string
}

interface NavGroup {
  title: string
  items: NavItem[]
}

const navGroups = computed<NavGroup[]>(() => [
  {
    title: 'OVERVIEW',
    items: [
      {
        label: 'Dashboard',
        icon: 'i-lucide-layout-dashboard',
        to: '/dashboard',
      },
      {
        label: 'Analytics',
        icon: 'i-lucide-bar-chart-3',
        to: '/analytics',
      },
    ],
  },
  {
    title: 'COMMERCE',
    items: [
      {
        label: 'Customers',
        icon: 'i-lucide-users',
        to: '/customers',
      },
      {
        label: 'Stores',
        icon: 'i-lucide-store',
        to: '/stores',
      },
      {
        label: 'Products',
        icon: 'i-lucide-package',
        to: '/products',
      },
      {
        label: 'Orders',
        icon: 'i-lucide-shopping-cart',
        to: '/orders',
      },
    ],
  },
  {
    title: 'OPERATIONS',
    items: [
      {
        label: 'QR Codes',
        icon: 'i-lucide-qr-code',
        to: '/qr-codes',
      },
      {
        label: 'Designs',
        icon: 'i-lucide-palette',
        to: '/designs',
      },
      {
        label: 'Production',
        icon: 'i-lucide-factory',
        to: '/production',
      },
      {
        label: 'NFC Tools',
        icon: 'i-lucide-radio',
        to: '/nfc-tools',
      },
      {
        label: 'Quality Control',
        icon: 'i-lucide-shield-check',
        to: '/quality-control',
      },
      {
        label: 'Shipping',
        icon: 'i-lucide-truck',
        to: '/shipping',
      },
    ],
  },
  {
    title: 'SYSTEM',
    items: [
      {
        label: 'Audit Logs',
        icon: 'i-lucide-scroll-text',
        to: '/audit-logs',
      },
    ],
  },
])

function isItemActive(to: string): boolean {
  if (to === '/dashboard') return route.path === '/dashboard'
  return route.path.startsWith(to)
}

const userMenuItems = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: user.value?.name || 'Admin',
      type: 'label',
    },
    {
      label: user.value?.email || 'admin@reputap.id',
      type: 'label',
    },
  ],
  [
    {
      label: 'Sign out',
      icon: 'i-lucide-log-out',
      color: 'error',
      onSelect: () => logout(),
    },
  ],
])
</script>

<template>
  <div class="min-h-screen flex bg-slate-50 text-slate-800 antialiased font-sans">
    <!-- Desktop & Tablet Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-30 flex flex-col bg-white border-r border-slate-200/90 transition-all duration-300 ease-in-out"
      :class="sidebarOpen ? 'w-64' : 'w-20'"
    >
      <!-- Brand Header -->
      <div class="h-16 flex items-center justify-between px-4 border-b border-slate-200/80 bg-white">
        <NuxtLink to="/dashboard" class="flex items-center gap-3 overflow-hidden">
          <div class="size-9 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white font-black text-base shadow-sm shrink-0">
            R
          </div>
          <div v-if="sidebarOpen" class="flex flex-col truncate">
            <span class="font-bold text-base tracking-tight text-slate-900 leading-none">REPUTAP</span>
            <span class="text-[10px] font-semibold uppercase tracking-wider text-slate-400 mt-1">Smart Stand DFY</span>
          </div>
        </NuxtLink>
        <button
          v-if="sidebarOpen"
          class="size-8 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 flex items-center justify-center transition-colors"
          title="Collapse sidebar"
          @click="sidebarOpen = false"
        >
          <UIcon name="i-lucide-chevron-left" class="size-4" />
        </button>
      </div>

      <!-- Navigation Links -->
      <div class="flex-1 overflow-y-auto px-3 py-4 space-y-6">
        <div v-for="(group, gIdx) in navGroups" :key="gIdx" class="space-y-1">
          <p
            v-if="sidebarOpen"
            class="px-3 text-[11px] font-bold uppercase tracking-wider text-slate-400 mb-2"
          >
            {{ group.title }}
          </p>
          <div v-else class="h-px bg-slate-100 my-2" />

          <NuxtLink
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            class="group flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all"
            :class="[
              isItemActive(item.to)
                ? 'bg-blue-50 text-blue-700 font-semibold shadow-xs'
                : 'text-slate-600 hover:bg-slate-100/80 hover:text-slate-900',
              !sidebarOpen ? 'justify-center px-0' : '',
            ]"
            :title="!sidebarOpen ? item.label : undefined"
          >
            <UIcon
              :name="item.icon"
              class="size-5 shrink-0 transition-colors"
              :class="isItemActive(item.to) ? 'text-blue-600' : 'text-slate-400 group-hover:text-slate-600'"
            />
            <span v-if="sidebarOpen" class="truncate">{{ item.label }}</span>
            <span
              v-if="sidebarOpen && item.badge"
              class="ml-auto px-1.5 py-0.5 text-[11px] font-semibold rounded-full"
              :class="item.badgeColor || 'bg-slate-100 text-slate-600'"
            >
              {{ item.badge }}
            </span>
          </NuxtLink>
        </div>
      </div>

      <!-- User Profile Footer -->
      <div class="p-3 border-t border-slate-200/80 bg-white">
        <UDropdownMenu
          :items="userMenuItems"
          :content="{ align: 'center', collisionPadding: 12 }"
        >
          <button
            class="w-full flex items-center gap-3 p-2 rounded-xl text-left hover:bg-slate-50 transition-colors border border-transparent hover:border-slate-200/70"
            :class="!sidebarOpen ? 'justify-center' : ''"
          >
            <div class="size-9 rounded-full bg-slate-100 border border-slate-200 text-slate-700 font-bold flex items-center justify-center shrink-0">
              {{ user?.name ? user.name.charAt(0).toUpperCase() : 'A' }}
            </div>
            <div v-if="sidebarOpen" class="flex-1 min-w-0">
              <p class="text-xs font-semibold text-slate-900 truncate">
                {{ user?.name || 'Administrator' }}
              </p>
              <p class="text-[11px] text-slate-400 truncate">
                {{ user?.role_name || 'Admin' }}
              </p>
            </div>
            <UIcon
              v-if="sidebarOpen"
              name="i-lucide-chevrons-up-down"
              class="size-4 text-slate-400 ml-auto shrink-0"
            />
          </button>
        </UDropdownMenu>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div
      class="flex-1 flex flex-col min-w-0 transition-all duration-300 ease-in-out"
      :class="sidebarOpen ? 'pl-64' : 'pl-20'"
    >
      <!-- Top Sticky Header -->
      <header class="sticky top-0 z-20 h-16 bg-white/95 backdrop-blur-md border-b border-slate-200/80 px-6 flex items-center justify-between shadow-xs">
        <div class="flex items-center gap-3">
          <button
            v-if="!sidebarOpen"
            class="size-8 rounded-lg text-slate-500 hover:text-slate-800 hover:bg-slate-100 flex items-center justify-center transition-colors"
            title="Expand sidebar"
            @click="sidebarOpen = true"
          >
            <UIcon name="i-lucide-menu" class="size-5" />
          </button>
          
          <div class="hidden sm:flex items-center gap-2 text-xs font-medium text-slate-400">
            <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 font-semibold">
              <span class="size-1.5 rounded-full bg-emerald-500 animate-pulse" />
              API Connected
            </span>
            <span class="text-slate-300">|</span>
            <span class="text-slate-500">Jakarta, ID (WIB)</span>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <a
            :href="`${apiBase}/swagger/index.html`"
            target="_blank"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 border border-slate-200 transition-colors"
          >
            <UIcon name="i-lucide-file-code-2" class="size-3.5 text-slate-500" />
            <span>Swagger API</span>
          </a>

          <button
            type="button"
            class="size-8 rounded-lg text-slate-500 hover:text-slate-800 hover:bg-slate-100 border border-slate-200 flex items-center justify-center transition-colors"
            :title="colorMode.value === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
            @click="toggleColorMode"
          >
            <UIcon :name="colorMode.value === 'dark' ? 'i-lucide-sun' : 'i-lucide-moon'" class="size-4" />
          </button>

          <div class="h-4 w-px bg-slate-200" />

          <NuxtLink
            to="/orders/new"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white shadow-xs transition-colors"
          >
            <UIcon name="i-lucide-plus" class="size-4" />
            <span>Create Order</span>
          </NuxtLink>
        </div>
      </header>

      <!-- Main Body Canvas -->
      <main class="flex-1 p-6 md:p-8 bg-[#f8fafc]">
        <slot />
      </main>
    </div>
  </div>
</template>
