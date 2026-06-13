package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	medicine "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestMedicineService() (*mocks.MockMedicineRepository, medicine.MedicineService) {
	mockRepo := new(mocks.MockMedicineRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := medicine.NewMedicineService(mockRepo, cfg)
	return mockRepo, service
}

func TestMedicineService_List_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "created_at"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("List", query).Return(medicines, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "created_at"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("DeletedList", query).Return(medicines, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_AvailableList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "created_at"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("AvailableList", query).Return(medicines, int64(2), nil)

	res, err := service.AvailableList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_LowStockList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "stock"
	query.SortDir = "asc"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("LowStockList", query).Return(medicines, int64(2), nil)

	res, err := service.LowStockList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_OutStockList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "stock"
	query.SortDir = "asc"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("OutStockList", query).Return(medicines, int64(2), nil)

	res, err := service.OutStockList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "created_at"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("ActiveList", query).Return(medicines, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	query := mocks.NewMedicinePaginationQuery(1, 10)
	query.SortBy = "created_at"
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("InactiveList", query).Return(medicines, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)

	mockRepo.On("FindByID", uint(1)).Return(med, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, med.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestMedicineService()

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("medicine not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(med, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, med.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	req := mocks.NewCreateMedicineRequest("Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)

	mockRepo.On("Create", mock.AnythingOfType("*models.Medicine")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Med1", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)
	req := mocks.NewUpdateMedicineRequest("Med2", "Gen2", "Brand2", 2, "250mg", "Manuf2", "Capsule", 20000)

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Medicine")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Med2", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_AddStock_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)
	req := &dto.AddStockRequest{Quantity: 50}

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("AddStock", uint(1), 50).Return(nil)

	err := service.AddStock(1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_ReduceStock_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)
	req := &dto.ReduceStockRequest{Quantity: 50}

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("ReduceStock", uint(1), 50).Return(nil)

	err := service.ReduceStock(1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_ReduceStock_Insufficient(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 20, 15000, false)
	req := &dto.ReduceStockRequest{Quantity: 50}

	mockRepo.On("FindByID", uint(1)).Return(med, nil)

	err := service.ReduceStock(1, req)

	assert.Error(t, err)
	assert.Equal(t, "insufficient stock", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)
	med.IsActive = false

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("Activate", uint(1)).Return(nil)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)
	req := &dto.DeactivateMedicineRequest{}

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := service.Deactivate(1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)

	mockRepo.On("FindByID", uint(1)).Return(med, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(med, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineService()
	med := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(med, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
