<script setup lang="ts">
import type { FormSubmitEvent, FormError } from '@nuxt/ui'

definePageMeta({
  layout: false,
})

const { login, isAuthenticated } = useAuth()

if (isAuthenticated.value) {
  navigateTo('/dashboard')
}

const state = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const errorMessage = ref('')

function validate(state: any): FormError[] {
  const errors: FormError[] = []
  if (!state.email) errors.push({ name: 'email', message: 'Email is required' })
  if (!state.password) errors.push({ name: 'password', message: 'Password is required' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  loading.value = true
  errorMessage.value = ''

  const result = await login({
    email: event.data.email,
    password: event.data.password,
  })

  loading.value = false

  if (result.success) {
    navigateTo('/dashboard')
  }
  else {
    errorMessage.value = result.error || 'Login failed'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-default">
    <UCard class="w-full max-w-sm">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold text-highlighted">REPUTAP</h1>
          <p class="text-sm text-muted mt-1">Smart Review Stand Management</p>
        </div>
      </template>

      <UForm :validate="validate" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormField label="Email" name="email">
          <UInput
            v-model="state.email"
            type="email"
            placeholder="admin@reputap.id"
            icon="i-lucide-mail"
            class="w-full"
          />
        </UFormField>

        <UFormField label="Password" name="password">
          <UInput
            v-model="state.password"
            type="password"
            placeholder="Enter your password"
            icon="i-lucide-lock"
            class="w-full"
          />
        </UFormField>

        <div v-if="errorMessage" class="p-3 rounded-md bg-error-50 dark:bg-error-950 text-error text-sm">
          {{ errorMessage }}
        </div>

        <UButton type="submit" block :loading="loading">
          Sign In
        </UButton>
      </UForm>
    </UCard>
  </div>
</template>
