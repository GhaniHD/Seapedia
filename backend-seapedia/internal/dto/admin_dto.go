package dto

type AdminDashboardResponse struct {
	TotalUsers          int64 `json:"total_users"`
	TotalStores         int64 `json:"total_stores"`
	TotalProducts       int64 `json:"total_products"`
	TotalOrders         int64 `json:"total_orders"`
	TotalVouchers       int64 `json:"total_vouchers"`
	TotalPromos         int64 `json:"total_promos"`
	TotalDeliveryJobs   int64 `json:"total_delivery_jobs"`
	OverdueOrdersCount  int64 `json:"overdue_orders_count"`
}

type SimulateNextDayResponse struct {
	VirtualNow      string   `json:"virtual_now"`
	OverdueHandled  []string `json:"overdue_orders_handled"`
}
