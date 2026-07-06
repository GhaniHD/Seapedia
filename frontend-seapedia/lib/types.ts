// Mirrors backend-seapedia/internal/dto/*.go

export type Role = "admin" | "seller" | "buyer" | "driver";

export const ORDER_STATUS = {
  PACKING: "Sedang Dikemas",
  WAITING_COURIER: "Menunggu Pengirim",
  SHIPPING: "Sedang Dikirim",
  DONE: "Pesanan Selesai",
  RETURNED: "Dikembalikan",
} as const;

export type DeliveryMethod = "instant" | "next_day" | "regular";

// ---------- User / Auth ----------
export interface UserResponse {
  id: string;
  name: string;
  email: string;
}

export interface ProfileResponse {
  id: string;
  name: string;
  email: string;
  roles: Role[];
  active_role: Role | "";
  wallet_balance?: number;
  store_income?: number;
  driver_earning?: number;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  need_role_selection: boolean;
  roles: Role[];
  active_role: Role | "";
  user: UserResponse;
}

export interface SelectRoleRequest {
  role: Role;
}

export interface AddRoleRequest {
  role: Role;
}

// ---------- Store ----------
export interface UpsertStoreRequest {
  name: string;
  description: string;
}

export interface StoreResponse {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

// ---------- Product ----------
export interface UpsertProductRequest {
  name: string;
  description: string;
  price: number;
  stock: number;
}

export interface ProductResponse {
  id: string;
  store_id: string;
  store_name?: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  created_at: string;
}

// ---------- Review ----------
export interface CreateReviewRequest {
  reviewer_name: string;
  rating: number;
  comment: string;
}

export interface ReviewResponse {
  id: string;
  reviewer_name: string;
  rating: number;
  comment: string;
  created_at: string;
}

// ---------- Wallet ----------
export interface TopupRequest {
  amount: number;
}

export interface WalletResponse {
  balance: number;
}

export interface WalletTransactionResponse {
  id: string;
  type: string;
  amount: number;
  description: string;
  created_at: string;
}

export interface UpsertAddressRequest {
  label: string;
  detail: string;
  is_default: boolean;
}

export interface AddressResponse {
  id: string;
  label: string;
  detail: string;
  is_default: boolean;
}

// ---------- Cart ----------
export interface AddCartItemRequest {
  product_id: string;
  quantity: number;
}

export interface UpdateCartItemRequest {
  quantity: number;
}

export interface CartItemResponse {
  id: string;
  product_id: string;
  name: string;
  price: number;
  quantity: number;
  subtotal: number;
}

export interface CartResponse {
  store_id: string | null;
  store_name?: string;
  items: CartItemResponse[];
  subtotal: number;
}

// ---------- Checkout ----------
export interface CheckoutRequest {
  address_id: string;
  delivery_method: DeliveryMethod;
  discount_code?: string;
}

export interface CheckoutSummaryResponse {
  order_id: string;
  order_no: string;
  subtotal: number;
  discount_amount: number;
  discount_kind?: "voucher" | "promo" | "";
  delivery_fee: number;
  tax_amount: number;
  tax_rate: number;
  total: number;
  status: string;
  created_at: string;
}

// ---------- Order ----------
export interface OrderItemResponse {
  product_name: string;
  price: number;
  quantity: number;
}

export interface StatusHistoryResponse {
  status: string;
  note: string;
  created_at: string;
}

export interface OrderResponse {
  id: string;
  order_no: string;
  store_name?: string;
  buyer_name?: string;
  delivery_method: DeliveryMethod;
  subtotal: number;
  discount_amount: number;
  delivery_fee: number;
  tax_amount: number;
  total: number;
  status: string;
  deadline_at?: string;
  items?: OrderItemResponse[];
  status_history?: StatusHistoryResponse[];
  created_at: string;
}

export interface SpendingReportResponse {
  total_orders: number;
  total_spending: number;
}

export interface IncomeReportResponse {
  total_orders: number;
  total_income: number;
  total_reversed: number;
}

// ---------- Discount ----------
export interface VoucherResponse {
  id: string;
  code: string;
  discount_type: "percent" | "fixed";
  discount_value: number;
  expiry_date: string;
  usage_limit: number;
  usage_count: number;
}

export interface PromoResponse {
  id: string;
  code: string;
  discount_type: "percent" | "fixed";
  discount_value: number;
  expiry_date: string;
}

// ---------- API error shape ----------
export interface ApiErrorBody {
  error?: string;
  message?: string;
}
