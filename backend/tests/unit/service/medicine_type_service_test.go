package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	medicinetype "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine-type"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestMedicineTypeService() (*mocks.MockMedicineTypeRepository, medicinetype.MedicineTypeService) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := medicinetype.NewMedicineTypeService(mockRepo, cfg)
	return mockRepo, service
}

func TestMedicineTypeService_List_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	query.SortBy = "created_at"
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("List", query).Return(types, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	query.SortBy = "created_at"
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("DeletedList", query).Return(types, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	query.SortBy = "created_at"
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("ActiveList", query).Return(types, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	query.SortBy = "created_at"
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("InactiveList", query).Return(types, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	typeType := mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)

	mockRepo.On("FindByID", uint(1)).Return(typeType, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, typeType.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("medicine type not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	req := mocks.NewCreateMedicineTypeRequest("Type1", "T001", "Desc1", true)

	mockRepo.On("IsNameExists", "Type1", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.MedicineType")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Type1", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Create_NameExists(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	req := mocks.NewCreateMedicineTypeRequest("Type1", "T001", "Desc1", true)

	mockRepo.On("IsNameExists", "Type1", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Create_CodeExists(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	req := mocks.NewCreateMedicineTypeRequest("Type1", "T001", "Desc1", true)

	mockRepo.On("IsNameExists", "Type1", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	typeType := mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)
	req := mocks.NewUpdateMedicineTypeRequest("Type2", "T002", "Desc2", false)

	mockRepo.On("FindByID", uint(1)).Return(typeType, nil)
	mockRepo.On("IsNameExists", "Type2", []uint{1}).Return(false, nil)
	mockRepo.On("IsCodeExists", "T002", []uint{1}).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.MedicineType")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Type2", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()
	typeType := mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)

	mockRepo.On("FindByID", uint(1)).Return(typeType, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestMedicineTypeService()

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := service.Deactivate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
