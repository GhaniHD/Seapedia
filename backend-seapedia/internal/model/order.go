package model

import (
	"time"

	"github.com/google/uuid"
)

// Status order utama. Nama status HARUS persis seperti ini karena
// harus tetap terlihat konsisten di seluruh aplikasi (lihat dokumen tugas).
const (
	StatusSedangDikemas    = "Sedang Dikemas"
	StatusMenungguPengirim = "Menunggu Pengirim"
	StatusSedangDikirim    = "Sedang Dikirim"
	StatusPesananSelesai   = "Pesanan Selesai"
	StatusDikembalikan     = "Dikembalikan"
)

const (
	DeliveryInstant  = "instant"
	DeliveryNextDay  = "next_day"
	DeliveryRegular  = "regular"
)

type Order struct {
	ID                   uuid.UUID
	OrderNo              string
	BuyerID              uuid.UUID
	StoreID              uuid.UUID
	AddressID            uuid.UUID
	DeliveryMethod       string
	Subtotal             float64
	DiscountAmount       float64
	DeliveryFee          float64
	TaxAmount            float64
	Total                float64
	Status               string
	VoucherID            *uuid.UUID
	PromoID              *uuid.UUID
	DriverID             *uuid.UUID
	Refunded             bool
	SellerIncomeReversed bool
	StockRestored        bool
	DeadlineAt           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type OrderItem struct {
	ID          uuid.UUID
	OrderID     uuid.UUID
	ProductID   uuid.UUID
	ProductName string
	Price       float64
	Quantity    int
}

type OrderStatusHistory struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Status    string
	Note      string
	CreatedAt time.Time
}
