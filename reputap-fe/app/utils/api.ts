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
