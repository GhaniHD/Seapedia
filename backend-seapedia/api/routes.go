package api

import (
	"backend-seapedia/internal/handler"
	"backend-seapedia/internal/middleware"
	"backend-seapedia/internal/model"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Review   *handler.ReviewHandler
	Store    *handler.StoreHandler
	Product  *handler.ProductHandler
	Wallet   *handler.WalletHandler
	Cart     *handler.CartHandler
	Checkout *handler.CheckoutHandler
	Order    *handler.OrderHandler
	Discount *handler.DiscountHandler
	Delivery *handler.DeliveryHandler
	Admin    *handler.AdminHandler
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRoutes(h *Handlers, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	v1 := r.Group("/api/v1")
	{
		// ============ PUBLIC (Guest) - Level 1 & 2 ============
		v1.POST("/register", h.Auth.Register)
		v1.POST("/login", h.Auth.Login)
		v1.GET("/products", h.Product.ListPublic)
		v1.GET("/products/:id", h.Product.GetPublicDetail)
		v1.GET("/stores", h.Store.ListPublic)
		v1.GET("/stores/:id", h.Store.GetPublicDetail)
		v1.POST("/reviews", h.Review.Create) // guest & logged-in user boleh submit
		v1.GET("/reviews", h.Review.List)

		// ============ SELECT ROLE (butuh temp token / token biasa) ============
		pending := v1.Group("/")
		pending.Use(middleware.AuthPending(jwtSecret))
		{
			pending.POST("/select-role", h.Auth.SelectRole)
		}

		// ============ AUTHENTICATED (butuh active_role, semua role non-guest) ============
		authed := v1.Group("/")
		authed.Use(middleware.AuthRequired(jwtSecret))
		{
			authed.POST("/logout", h.Auth.Logout)
			authed.GET("/profile", h.User.GetProfile)
			authed.POST("/roles", h.Auth.AddRole) // tambah role lain untuk akun sendiri

			// ---------- BUYER (Level 3 & 4) ----------
			buyer := authed.Group("/buyer")
			buyer.Use(middleware.RequireRole(model.RoleBuyer))
			{
				buyer.POST("/wallet/topup", h.Wallet.Topup)
				buyer.GET("/wallet", h.Wallet.GetBalance)
				buyer.GET("/wallet/transactions", h.Wallet.ListTransactions)
				buyer.POST("/addresses", h.Wallet.AddAddress)
				buyer.GET("/addresses", h.Wallet.ListAddresses)

				buyer.GET("/cart", h.Cart.GetCart)
				buyer.POST("/cart/items", h.Cart.AddItem)
				buyer.PUT("/cart/items/:productId", h.Cart.UpdateItem)
				buyer.DELETE("/cart/items/:productId", h.Cart.RemoveItem)
				buyer.DELETE("/cart", h.Cart.ClearCart)

				buyer.POST("/checkout", h.Checkout.Checkout)
				buyer.GET("/orders", h.Order.ListMine)
				buyer.GET("/orders/:id", h.Order.GetMyDetail)
				buyer.GET("/reports/spending", h.Order.SpendingReport)
			}

			// ---------- SELLER (Level 2 & 4) ----------
			seller := authed.Group("/seller")
			seller.Use(middleware.RequireRole(model.RoleSeller))
			{
				seller.POST("/store", h.Store.UpsertMyStore)
				seller.PUT("/store", h.Store.UpsertMyStore)
				seller.GET("/store", h.Store.GetMyStore)

				seller.POST("/products", h.Product.Create)
				seller.PUT("/products/:id", h.Product.Update)
				seller.DELETE("/products/:id", h.Product.Delete)
				seller.GET("/products", h.Product.ListMine)

				seller.GET("/orders", h.Order.ListIncoming)
				seller.POST("/orders/:id/process", h.Order.ProcessOrder)
				seller.GET("/reports/income", h.Order.IncomeReport)
			}

			// ---------- DRIVER (Level 5) ----------
			driver := authed.Group("/driver")
			driver.Use(middleware.RequireRole(model.RoleDriver))
			{
				driver.GET("/jobs", h.Delivery.FindJobs)
				driver.GET("/jobs/:id", h.Delivery.GetDetail)
				driver.POST("/jobs/:id/take", h.Delivery.TakeJob)
				driver.POST("/jobs/:id/complete", h.Delivery.CompleteJob)
				driver.GET("/my-jobs", h.Delivery.MyJobs)
				driver.GET("/earnings", h.Delivery.MyEarnings)
			}

			// ---------- ADMIN (Level 4 & 6) ----------
			admin := authed.Group("/admin")
			admin.Use(middleware.RequireRole(model.RoleAdmin))
			{
				admin.GET("/dashboard", h.Admin.Dashboard)
				admin.POST("/simulate-next-day", h.Admin.SimulateNextDay)

				admin.POST("/vouchers", h.Discount.CreateVoucher)
				admin.GET("/vouchers", h.Discount.ListVouchers)
				admin.POST("/promos", h.Discount.CreatePromo)
				admin.GET("/promos", h.Discount.ListPromos)
			}

			// Voucher/promo list juga bisa diakses buyer saat checkout untuk lihat kode yang tersedia
			authed.GET("/vouchers", h.Discount.ListVouchers)
			authed.GET("/promos", h.Discount.ListPromos)
		}
	}

	return r
}
