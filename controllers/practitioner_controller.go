package controllers

import (
	"net/http"
	"strings"

	_ "fhir-api/models"
	"fhir-api/services"

	"github.com/gin-gonic/gin"
)

type PractitionerController struct {
	service *services.PractitionerService
}

func NewPractitionerController(service *services.PractitionerService) *PractitionerController {
	return &PractitionerController{service: service}
}

type GetPractitionerRequest struct {
	Fields string `form:"fields" binding:"required"`
}

// GetPractitioner godoc
// @Summary      Busca um Practitioner por ID
// @Description  Retorna um Practitioner específico, filtrando os campos desejados.
// @Tags         practitioners
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "ID do Practitioner"
// @Param        fields query     string  true  "Lista de campos separados por vírgula (ex: fhirId,fullUrl,status)"
// @Failure      400    {object}  map[string]string "Erro de validação nos parâmetros"
// @Failure      404    {object}  map[string]string "Practitioner não encontrado"
// @Router       /practitioner/{id} [get]
func (c *PractitionerController) GetPractitioner(ctx *gin.Context) {
	id := ctx.Param("id")

	var req GetPractitionerRequest
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
		if !containsPractitionerFields(validFields, field) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid field specified"})
			return
		}
	}

	practitioner, err := c.service.GetPractitioner(ctx.Request.Context(), id, fields)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "practitioner not found"})
		return
	}

	ctx.JSON(http.StatusOK, practitioner)
}

func containsPractitionerFields(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
