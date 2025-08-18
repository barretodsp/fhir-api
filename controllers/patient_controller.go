package controllers

import (
	"net/http"
	"strings"

	"fhir-api/services"

	"github.com/gin-gonic/gin"
)

type PatientController struct {
	service *services.PatientService
}

func NewPatientController(service *services.PatientService) *PatientController {
	return &PatientController{service: service}
}

type GetPatientRequest struct {
	Fields string `form:"fields" binding:"required"`
}

// GetPatientByID godoc
// @Summary Retorna um paciente
// @Description Busca paciente pelo ID
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID do paciente"
// @Success 200 {object} models.Patient
// @Failure 400 {object} map[string]string
// @Router /patients/{id} [get]
func (c *PatientController) GetPatient(ctx *gin.Context) {
	id := ctx.Param("id")

	var req GetPatientRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fields parameter is required"})
		return
	}

	fieldsParam := req.Fields
	fields := strings.Split(fieldsParam, ",")
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	validFields := []string{"id", "fhirId", "givenName", "familyName", "birthDate", "gender"}
	for _, field := range fields {
		if !containsPatientFields(validFields, field) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid field specified"})
			return
		}
	}

	patient, err := c.service.GetPatient(ctx.Request.Context(), id, fields)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

func containsPatientFields(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
