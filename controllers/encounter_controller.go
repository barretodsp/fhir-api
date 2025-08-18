package controllers

import (
	"net/http"
	"strings"

	"fhir-api/models"
	"fhir-api/services"

	"github.com/gin-gonic/gin"
)

type EncounterController struct {
	service *services.EncounterService
}

func NewEncounterController(service *services.EncounterService) *EncounterController {
	return &EncounterController{service: service}
}

type GetEncounterRequest struct {
	Fields string `form:"fields" binding:"required"`
}

// GetEncounter godoc
// @Summary Get encounter by ID
// @Description Retrieves a specific encounter with optional field filtering
// @Tags Encounters
// @Accept json
// @Produce json
// @Param id path string true "Encounter ID"
// @Param fields query string false "Comma-separated list of fields to return (fhirId,fullUrl,status,class,period,practitionerId,patientId)"
// @Success 200 {object} models.Encounter
// @Failure 400 {object} map[string]string "invalid field specified"
// @Failure 404 {object} map[string]string "Encounter not found"
// @Router /encounters/{id} [get]
func (c *EncounterController) GetEncounter(ctx *gin.Context) {
	id := ctx.Param("id")

	var req GetEncounterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fields parameter is required"})
		return
	}

	fieldsParam := req.Fields
	fields := strings.Split(fieldsParam, ",")
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	validFields := []string{"fhirId", "fullUrl", "status", "class", "period", "practitionerId", "patientId"}
	for _, field := range fields {
		if !contains(validFields, field) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid field specified"})
			return
		}
	}

	encounter, err := c.service.GetEncounter(ctx.Request.Context(), id, fields)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "encounter not found"})
		return
	}

	ctx.JSON(http.StatusOK, encounter)
}

// UpdateEncounterStatus godoc
// @Summary Update encounter status
// @Description Updates the status of a specific encounter
// @Tags Encounters
// @Accept json
// @Produce json
// @Param id path string true "Encounter ID"
// @Param request body models.EncounterUpdate true "Status update payload"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 404 {object} map[string]string "Encounter not found"
func (c *EncounterController) UpdateEncounterStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.EncounterUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := c.service.UpdateEncounterStatus(ctx.Request.Context(), id, req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
