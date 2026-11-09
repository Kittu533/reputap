<script setup lang="ts">
import { h, resolveComponent, ref, reactive, watch, onMounted, computed } from "vue";
import type { TableColumn } from "@nuxt/ui";
import type { QRCodeDetail, Store } from "~/utils/api";

definePageMeta({
  middleware: "auth",
});

const { apiFetch } = useApi();
const toast = useToast();
const config = useRuntimeConfig();
const { accessToken } = useAuth();

const UButton = resolveComponent("UButton");
const UBadge = resolveComponent("UBadge");

const search = ref("");
const statusFilter = ref<string | undefined>(undefined);
const page = ref(1);
const perPage = ref(20);
const data = ref<QRCodeDetail[]>([]);
const total = ref(0);
const loading = ref(false);

// Stores dropdown list
const stores = ref<Store[]>([]);
async function loadStores() {
  const res: any = await apiFetch<Store[]>("/api/v1/stores", {
    query: { per_page: 100 },
  });
  if (res.success && res.data) {
    stores.value = res.data;
  }
}

const storeOptions = computed(() => [
  { label: "Unassigned (None)", value: null },
  ...stores.value.map((s) => ({
    label: `${s.store_code} — ${s.name}`,
    value: s.id,
  })),
]);

// Batch Generation State
const batchModalOpen = ref(false);
const batchLoading = ref(false);
const batchForm = reactive({
  count: 20,
  prefix: "RP26",
});

// Edit / Assign State
const editModalOpen = ref(false);
const editLoading = ref(false);
const qrToEdit = ref<QRCodeDetail | null>(null);
const editForm = reactive({
  status: "unassigned",
  store_id: null as number | null,
});

// Delete State
const deleteModalOpen = ref(false);
const deleteLoading = ref(false);
const qrToDelete = ref<QRCodeDetail | null>(null);

const statusOptions = [
  { label: "Active", value: "active" },
  { label: "Unassigned", value: "unassigned" },
  { label: "Inactive", value: "inactive" },
  { label: "Suspended", value: "suspended" },
  { label: "Retired", value: "retired" },
];

const statusColor: Record<string, string> = {
  active: "success",
  unassigned: "info",
  inactive: "neutral",
  suspended: "warning",
  retired: "error",
};

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function truncateUrl(url: string, maxLen = 40): string {
  if (!url) return "—";
  return url.length > maxLen ? `${url.substring(0, maxLen)}…` : url;
}

const columns: TableColumn<QRCodeDetail>[] = [
  {
    accessorKey: "code",
    header: "Code",
    cell: ({ row }) =>
      h("span", { class: "font-mono text-xs font-bold" }, row.original.code),
  },
  {
    accessorKey: "store_id",
    header: "Store",
    cell: ({ row }) => {
      if (!row.original.store_id) {
        return h(
          UBadge,
          { color: "neutral", variant: "subtle", size: "xs" },
          () => "Unassigned",
        );
      }
      const matched = stores.value.find((s) => s.id === row.original.store_id);
      return h(
        "span",
        { class: "text-xs font-medium" },
        matched ? `${matched.name} (#${row.original.store_id})` : `#${row.original.store_id}`,
      );
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) =>
      h(UBadge, {
        label: row.original.status,
        color: statusColor[row.original.status] || "neutral",
        variant: "subtle",
        size: "xs",
        class: "capitalize",
      }),
  },
  {
    accessorKey: "redirect_url",
    header: "Redirect URL",
    cell: ({ row }) => {
      const url = row.original.redirect_url;
      if (!url) return "—";
      return h(
        "a",
        {
          href: url,
          target: "_blank",
          rel: "noopener noreferrer",
          class: "text-primary hover:underline text-xs",
        },
        truncateUrl(url),
      );
    },
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => formatDate(row.original.created_at),
  },
  {
    id: "actions",
    header: "",
    meta: { class: { td: "text-right" } },
    cell: ({ row }) => {
      return h("div", { class: "flex justify-end gap-1" }, [
        h(UButton, {
          icon: "i-lucide-download",
          color: "neutral",
          variant: "ghost",
          size: "xs",
          title: "Download PNG",
          onClick: () => downloadQR(row.original),
        }),
        h(UButton, {
          icon: "i-lucide-pencil",
          color: "neutral",
          variant: "ghost",
          size: "xs",
          title: "Edit / Assign Store",
          onClick: () => promptEdit(row.original),
        }),
        h(UButton, {
          icon: "i-lucide-trash-2",
          color: "error",
          variant: "ghost",
          size: "xs",
          title: "Delete QR",
          onClick: () => promptDelete(row.original),
        }),
      ]);
    },
  },
];

async function fetchData() {
  loading.value = true;
  const query: Record<string, any> = {
    page: page.value,
    per_page: perPage.value,
  };
  if (search.value) query.search = search.value;
  if (statusFilter.value) query.status = statusFilter.value;

  const response = await apiFetch<QRCodeDetail[]>("/api/v1/qr-codes", {
    query,
  });
  if (response.success) {
    data.value = response.data || [];
    total.value = response.meta?.total || 0;
  }
  loading.value = false;
}

// Edit QR prompt & save
function promptEdit(qr: QRCodeDetail) {
  qrToEdit.value = qr;
  editForm.status = qr.status;
  editForm.store_id = qr.store_id || null;
  editModalOpen.value = true;
}

async function saveEditQR() {
  if (!qrToEdit.value) return;
  editLoading.value = true;
  try {
    const res: any = await apiFetch(`/api/v1/qr-codes/${qrToEdit.value.id}`, {
      method: "PUT",
      body: {
        status: editForm.status,
        store_id: editForm.store_id || null,
      },
    });

    if (res.success) {
      toast.add({
        title: "QR Code Updated",
        description: `QR Code ${qrToEdit.value.code} updated successfully.`,
        color: "success",
      });
      editModalOpen.value = false;
      fetchData();
    } else {
      toast.add({
        title: "Update Failed",
        description: res.error?.message || "Failed to update QR code.",
        color: "error",
      });
    }
  } catch (err: any) {
    toast.add({
      title: "Update Error",
      description: err?.message,
      color: "error",
    });
  } finally {
    editLoading.value = false;
  }
}

// Delete QR prompt & confirm
function promptDelete(qr: QRCodeDetail) {
  qrToDelete.value = qr;
  deleteModalOpen.value = true;
}

async function confirmDeleteQR() {
  if (!qrToDelete.value) return;
  deleteLoading.value = true;
  try {
    const res: any = await apiFetch(`/api/v1/qr-codes/${qrToDelete.value.id}`, {
      method: "DELETE",
    });

    if (res.success) {
      toast.add({
        title: "QR Code Deleted",
        description: `QR Code ${qrToDelete.value.code} deleted successfully.`,
        color: "success",
      });
      deleteModalOpen.value = false;
      fetchData();
    } else {
      toast.add({
        title: "Delete Failed",
        description: res.error?.message || "Failed to delete QR code.",
        color: "error",
      });
    }
  } catch (err: any) {
    toast.add({
      title: "Delete Error",
      description: err?.message,
      color: "error",
    });
  } finally {
    deleteLoading.value = false;
  }
}

// Batch generation
async function handleBatchGenerate() {
  batchLoading.value = true;
  try {
    const res: any = await apiFetch("/api/v1/qr-codes/batch", {
      method: "POST",
      body: {
        count: Number(batchForm.count),
        prefix: batchForm.prefix.trim().toUpperCase(),
      },
    });

    if (res.success) {
      toast.add({
        title: "Batch Generated",
        description: `Successfully generated ${batchForm.count} new QR codes.`,
        color: "success",
      });
      batchModalOpen.value = false;
      fetchData();
    } else {
      toast.add({
        title: "Generation Failed",
        description: res.error?.message || "Failed to generate batch.",
        color: "error",
      });
    }
  } catch (err: any) {
    toast.add({
      title: "Generation Error",
      description: err?.message,
      color: "error",
    });
  } finally {
    batchLoading.value = false;
  }
}

async function downloadQR(qr: QRCodeDetail) {
  const apiBase = config.public.apiBase as string;
  const response = await fetch(
    `${apiBase}/api/v1/qr-codes/${qr.id}/download?size=512`,
    {
      headers: { Authorization: `Bearer ${accessToken.value}` },
    },
  );
  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `${qr.code}.png`;
  a.click();
  URL.revokeObjectURL(url);
}

let searchTimeout: ReturnType<typeof setTimeout>;
watch(search, () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    page.value = 1;
    fetchData();
  }, 300);
});
watch(statusFilter, () => {
  page.value = 1;
  fetchData();
});
watch(page, () => fetchData());
onMounted(() => {
  fetchData();
  loadStores();
});
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="QR Codes"
      description="Monitor generated QR tags, tracking codes, and redirection targets."
      :breadcrumbs="[{ label: 'QR Codes' }]"
    >
      <template #actions>
        <CommonButton
          icon="i-lucide-plus"
          @click="batchModalOpen = true"
        >
          Generate Batch QR
        </CommonButton>
      </template>
    </CommonPageHeader>

    <CommonDataTable
      v-model:search="search"
      v-model:page="page"
      :data="data"
      :columns="columns"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      search-placeholder="Search by code..."
      empty-title="No QR codes found"
      empty-description="No generated QR codes match your criteria."
      empty-action-label="Generate Batch QR"
      @empty-action="batchModalOpen = true"
    >
      <template #filters>
        <USelect
          v-model="statusFilter"
          :items="statusOptions"
          placeholder="All Status"
          class="w-44"
        />
      </template>
    </CommonDataTable>

    <!-- Modal: Batch Generate QR -->
    <UModal v-model:open="batchModalOpen" title="Generate Batch QR Codes">
      <template #body>
        <div class="space-y-4 py-2">
          <p class="text-xs text-slate-500">
            Generate unit QR Code fisik baru langsung ke database (status awal: unassigned / siap pakai).
          </p>

          <UFormField label="Jumlah Unit" name="count">
            <UInput
              v-model="batchForm.count"
              type="number"
              min="1"
              max="100"
              placeholder="Contoh: 20"
              class="w-full"
            />
          </UFormField>

          <UFormField label="Prefix Nomor Seri" name="prefix">
            <UInput
              v-model="batchForm.prefix"
              placeholder="Contoh: RP26"
              class="w-full uppercase font-mono"
            />
          </UFormField>
        </div>
      </template>

      <template #footer>
        <div class="flex items-center justify-end gap-2.5 w-full">
          <UButton
            label="Batal"
            color="neutral"
            variant="outline"
            :disabled="batchLoading"
            @click="batchModalOpen = false"
          />
          <UButton
            label="Generate Sekarang"
            color="primary"
            icon="i-lucide-sparkles"
            :loading="batchLoading"
            @click="handleBatchGenerate"
          />
        </div>
      </template>
    </UModal>

    <!-- Modal: Edit QR Code & Assign Store -->
    <UModal v-model:open="editModalOpen" :title="`Edit QR Code: ${qrToEdit?.code || ''}`">
      <template #body>
        <div class="space-y-4 py-2">
          <UFormField label="Hubungkan ke Toko / Store" name="store_id">
            <USelect
              v-model="editForm.store_id"
              :items="storeOptions"
              placeholder="Pilih Toko..."
              class="w-full"
            />
          </UFormField>

          <UFormField label="Status QR Code" name="status">
            <USelect
              v-model="editForm.status"
              :items="statusOptions"
              class="w-full capitalize"
            />
          </UFormField>

          <div v-if="qrToEdit?.redirect_url" class="p-3 bg-slate-50 rounded-xl border border-slate-200 text-xs">
            <span class="text-slate-400 block mb-0.5">Redirect URL saat ini:</span>
            <a :href="qrToEdit.redirect_url" target="_blank" class="text-primary font-mono hover:underline break-all">
              {{ qrToEdit.redirect_url }}
            </a>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="flex items-center justify-end gap-2.5 w-full">
          <UButton
            label="Batal"
            color="neutral"
            variant="outline"
            :disabled="editLoading"
            @click="editModalOpen = false"
          />
          <UButton
            label="Simpan Perubahan"
            color="primary"
            :loading="editLoading"
            @click="saveEditQR"
          />
        </div>
      </template>
    </UModal>

    <!-- Modal: Delete Confirmation -->
    <CommonConfirmModal
      v-model:open="deleteModalOpen"
      title="Delete QR Code"
      :description="`Apakah kamu yakin ingin menghapus QR Code '${qrToDelete?.code}'? Unit fisik instance dan log scan terkait juga akan dihapus.`"
      confirm-text="Ya, Hapus QR"
      confirm-color="error"
      icon="i-lucide-alert-triangle"
      :loading="deleteLoading"
      @confirm="confirmDeleteQR"
    />
  </div>
</template>
