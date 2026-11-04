export interface ApiResponse<T = any> {
  success: boolean
  data: T
  meta?: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
  error?: {
    code: string
    message: string
    details?: { field: string; message: string }[]
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface UserInfo {
  id: number
  name: string
  email: string
  role_id: number
  role_name: string
  status: string
}

export interface Customer {
  id: number
  customer_code: string
  name: string
  business_name: string | null
  email: string | null
  phone: string | null
  address: string | null
  city: string | null
  province: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export interface Store {
  id: number
  customer_id: number
  store_code: string
  name: string
  slug: string
  category: string
  address: string | null
  google_business_name: string | null
  review_url: string | null
  review_url_verified: boolean
  status: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  sku: string
  name: string
  description: string | null
  material: string | null
  width_mm: number | null
  height_mm: number | null
  thickness_mm: number | null
  nfc_chip_type: string | null
  finishing: string | null
  production_cost: number | null
  selling_price: number | null
  status: string
  created_at: string
  updated_at: string
}

export interface Order {
  id: number
  order_number: string
  customer_id: number
  store_id: number
  order_status: string
  payment_status: string
  subtotal: number
  discount: number
  total: number
  notes: string | null
  created_by: number
  created_at: string
  updated_at: string
}

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  quantity: number
  unit_price: number
  subtotal: number
}

export interface OrderDetail extends Order {
  items: OrderItem[]
}

export interface StatusHistory {
  id: number
  order_id: number
  from_status: string | null
  to_status: string
  changed_by: number
  notes: string | null
  created_at: string
}

export interface ProductInstance {
  id: number
  order_item_id: number
  store_id: number
  serial_number: string
  nfc_status: string
  nfc_configured_at: string | null
  nfc_configured_by: number | null
  production_status: string
  qc_status: string
  created_at: string
  updated_at: string
  qr_code_id: number
  qr_code: string
  qr_status: string
}

export interface QRCodeDetail {
  id: number
  product_instance_id: number
  store_id: number
  code: string
  status: string
  created_at: string
  updated_at: string
  redirect_url: string
}

export interface DesignAsset {
  id: number
  order_id: number
  store_id: number
  file_path: string
  file_name: string
  file_type: string
  file_size: number
  version: number
  approval_status: 'draft' | 'pending_approval' | 'approved' | 'revision_required'
  uploaded_by: number
  created_at: string
}

export interface ProductionJob {
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

export interface QCRecord {
  id: number
  product_instance_id: number
  inspected_by: number
  result: 'passed' | 'failed' | 'rework_required'
  checklist_data: Record<string, boolean>
  evidence_paths: string[] | null
  notes: string | null
  inspected_at: string
  serial_number: string
  store_id: number
}

export interface Shipment {
  id: number
  order_id: number
  recipient_name: string
  recipient_address: string
  recipient_phone: string | null
  courier: string | null
  tracking_number: string | null
  shipping_status: 'waiting' | 'packed' | 'shipped' | 'delivered' | 'failed'
  shipped_at: string | null
  delivered_at: string | null
  created_at: string
  updated_at: string
}

export interface AuditLog {
  id: number
  user_id: number
  action: string
  entity_type: string
  entity_id: number | null
  old_values: Record<string, any> | null
  new_values: Record<string, any> | null
  created_at: string
}

export interface DashboardSummary {
  orders: {
    total: number
    total_revenue: number
    by_status: Record<string, number>
  }
  customers: {
    total: number
  }
  stores: {
    total: number
    active: number
    review_url_verified: number
    by_status: Record<string, number>
  }
  qr_codes: {
    total: number
    active: number
    by_status: Record<string, number>
  }
  scans: {
    total: number
    qr_count: number
    nfc_count: number
    unknown_count: number
    by_device: Record<string, number>
  }
  redirect_success_rate: number
}

export interface RecentActivity {
  id: number
  user_id: number
  user_name: string
  action: string
  entity_type: string
  entity_id: number | null
  created_at: string
}

export interface StoreScanRank {
  store_id: number
  store_name: string
  customer_name: string
  total_scans: number
  qr_scans: number
  nfc_scans: number
}

export interface AnalyticsOverview {
  total_scans: number
  qr_scans: number
  nfc_scans: number
  unknown_scans: number
  device_breakdown: Record<string, number>
  today_scans: number
  yesterday_scans: number
  this_week_scans: number
  this_month_scans: number
  top_stores: StoreScanRank[]
}

export interface TimeSeriesPoint {
  date: string
  total: number
  qr_count: number
  nfc_count: number
}

export interface StoreAnalytics {
  store_id: number
  store_name: string
  customer_name: string
  review_url: string
  status: string
  review_url_verified: boolean
  active_qr_codes: number
  total_scans: number
  qr_scans: number
  nfc_scans: number
  device_breakdown: Record<string, number>
  daily_scans: TimeSeriesPoint[]
}

