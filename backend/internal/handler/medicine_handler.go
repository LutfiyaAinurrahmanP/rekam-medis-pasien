package handler

import (
	"net/http"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)


type MedicineHandler struct {
  service medicine.MedicineService
}

func NewMedicineHandler(s medicine.MedicineService) *MedicineHandler{
	return &MedicineHandler{
		service: s,
	}
}

func (h *MedicineHandler) ListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicines retrieved successfully", res)
}