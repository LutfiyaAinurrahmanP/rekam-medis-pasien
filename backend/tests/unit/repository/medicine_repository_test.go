package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMedicineRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("List", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("DeletedList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_AvailableList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("AvailableList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.AvailableList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_LowStockList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("LowStockList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.LowStockList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_OutStockList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("OutStockList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.OutStockList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_ActiveList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("ActiveList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_InactiveList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	query := mocks.NewMedicinePaginationQuery(1, 10)
	medicines := mocks.NewTestMedicineList(2)

	mockRepo.On("InactiveList", query).Return(medicines, int64(2), nil)

	res, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	expectedMed := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)

	mockRepo.On("FindByID", uint(1)).Return(expectedMed, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedMed.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)
	expectedMed := mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedMed, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedMed.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Medicine")).Return(nil)

	err := mockRepo.Create(&models.Medicine{Name: "Med1"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.Medicine")).Return(nil)

	err := mockRepo.Update(&models.Medicine{ID: 1, Name: "Med1"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_AddStock_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("AddStock", uint(1), 50).Return(nil)

	err := mockRepo.AddStock(1, 50)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_ReduceStock_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("ReduceStock", uint(1), 20).Return(nil)

	err := mockRepo.ReduceStock(1, 20)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
