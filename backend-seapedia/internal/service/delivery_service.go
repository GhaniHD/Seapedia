package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type DeliveryService interface {
	FindAvailableJobs(ctx context.Context) ([]dto.DeliveryJobResponse, error)
	GetJobDetail(ctx context.Context, jobID uuid.UUID) (*dto.DeliveryJobResponse, error)
	TakeJob(ctx context.Context, driverID, jobID uuid.UUID) error
	CompleteJob(ctx context.Context, driverID, jobID uuid.UUID) error
	MyEarnings(ctx context.Context, driverID uuid.UUID) (*dto.DriverEarningResponse, error)
	MyJobs(ctx context.Context, driverID uuid.UUID) ([]dto.DeliveryJobResponse, error)
}

type deliveryService struct {
	deliveryRepo repository.DeliveryRepository
	orderRepo    repository.OrderRepository
	storeRepo    repository.StoreRepository
	addressRepo  repository.AddressRepository
}

func NewDeliveryService(deliveryRepo repository.DeliveryRepository, orderRepo repository.OrderRepository,
	storeRepo repository.StoreRepository, addressRepo repository.AddressRepository) DeliveryService {
	return &deliveryService{deliveryRepo: deliveryRepo, orderRepo: orderRepo, storeRepo: storeRepo, addressRepo: addressRepo}
}

func (s *deliveryService) toJobResp(ctx context.Context, d *model.Delivery) *dto.DeliveryJobResponse {
	resp := &dto.DeliveryJobResponse{ID: d.ID, OrderID: d.OrderID, Fee: d.Fee, Status: d.Status, TakenAt: d.TakenAt, CompletedAt: d.CompletedAt}
	if o, err := s.orderRepo.FindByID(ctx, d.OrderID); err == nil {
		resp.OrderNo = o.OrderNo
		if st, err := s.storeRepo.FindByID(ctx, o.StoreID); err == nil {
			resp.StoreName = st.Name
		}
		if a, err := s.addressRepo.FindByID(ctx, o.AddressID); err == nil {
			resp.Address = a.Detail
		}
	}
	return resp
}

// FindAvailableJobs: hanya order yang statusnya "Menunggu Pengirim" DAN sudah punya delivery
// job berstatus "available" yang muncul. Order "Sedang Dikemas" tidak akan pernah muncul karena
// delivery job baru dibuat oleh seller saat ProcessOrder (Level 4/5).
func (s *deliveryService) FindAvailableJobs(ctx context.Context) ([]dto.DeliveryJobResponse, error) {
	jobs, err := s.deliveryRepo.ListAvailable(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryJobResponse, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, *s.toJobResp(ctx, &j))
	}
	return out, nil
}

func (s *deliveryService) GetJobDetail(ctx context.Context, jobID uuid.UUID) (*dto.DeliveryJobResponse, error) {
	d, err := s.deliveryRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	return s.toJobResp(ctx, d), nil
}

// TakeJob: atomic di level SQL (WHERE status='available'), jadi 2 driver tidak bisa
// ambil order yang sama secara bersamaan (business rule Level 5).
func (s *deliveryService) TakeJob(ctx context.Context, driverID, jobID uuid.UUID) error {
	d, err := s.deliveryRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}
	if err := s.deliveryRepo.TakeJob(ctx, jobID, driverID); err != nil {
		return err
	}
	if err := s.orderRepo.SetDriver(ctx, d.OrderID, driverID); err != nil {
		return err
	}
	if err := s.orderRepo.UpdateStatus(ctx, d.OrderID, model.StatusSedangDikirim); err != nil {
		return err
	}
	return s.orderRepo.AddStatusHistory(ctx, &model.OrderStatusHistory{ID: uuid.New(), OrderID: d.OrderID, Status: model.StatusSedangDikirim, Note: "Diambil oleh driver"})
}

func (s *deliveryService) CompleteJob(ctx context.Context, driverID, jobID uuid.UUID) error {
	d, err := s.deliveryRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}
	if d.DriverID == nil || *d.DriverID != driverID {
		return errors.New("job ini bukan milik Anda")
	}
	if err := s.deliveryRepo.CompleteJob(ctx, jobID); err != nil {
		return err
	}
	if err := s.orderRepo.UpdateStatus(ctx, d.OrderID, model.StatusPesananSelesai); err != nil {
		return err
	}
	return s.orderRepo.AddStatusHistory(ctx, &model.OrderStatusHistory{ID: uuid.New(), OrderID: d.OrderID, Status: model.StatusPesananSelesai, Note: "Dikonfirmasi selesai oleh driver"})
}

func (s *deliveryService) MyEarnings(ctx context.Context, driverID uuid.UUID) (*dto.DriverEarningResponse, error) {
	jobs, err := s.deliveryRepo.ListByDriver(ctx, driverID)
	if err != nil {
		return nil, err
	}
	resp := &dto.DriverEarningResponse{}
	for _, j := range jobs {
		if j.Status == model.DeliveryJobCompleted {
			resp.CompletedJobs++
			resp.TotalEarning += j.DriverEarning
		}
	}
	return resp, nil
}

func (s *deliveryService) MyJobs(ctx context.Context, driverID uuid.UUID) ([]dto.DeliveryJobResponse, error) {
	jobs, err := s.deliveryRepo.ListByDriver(ctx, driverID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryJobResponse, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, *s.toJobResp(ctx, &j))
	}
	return out, nil
}
