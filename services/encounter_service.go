package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"fhir-api/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EncounterService struct {
	db           *mongo.Database
	logger       *logrus.Logger
	validFields  map[string]bool
	validStatus  map[string]bool
	validClasses map[string]bool
}

func NewEncounterService(db *mongo.Database, logger *logrus.Logger) *EncounterService {
	validFields := map[string]bool{
		"fhirId":         true,
		"fullUrl":        true,
		"status":         true,
		"class":          true,
		"period":         true,
		"practitionerId": true,
		"patientId":      true,
	}

	validStatus := map[string]bool{
		"planned":          true,
		"in-progress":      true,
		"on-hold":          true,
		"discharged":       true,
		"completed":        true,
		"finished":         true,
		"cancelled":        true,
		"discontinued":     true,
		"entered-in-error": true,
		"unknown":          true,
	}

	validClasses := map[string]bool{
		"inpatient":   true,
		"observation": true,
		"ambulatory":  true,
		"emergency":   true,
		"virtual":     true,
		"home-health": true,
	}

	return &EncounterService{
		db:           db,
		logger:       logger,
		validFields:  validFields,
		validStatus:  validStatus,
		validClasses: validClasses,
	}
}

func (s *EncounterService) GetEncounter(ctx context.Context, id string, fields []string) (*models.EncounterResponse, error) {
	startTime := time.Now()
	logFields := logrus.Fields{
		"operation":       "GetEncounter",
		"encounterId":     id,
		"requestedFields": fields,
	}

	for _, field := range fields {
		if !s.validFields[field] {
			s.logger.WithFields(logFields).WithField("invalidField", field).Warn("campo inválido solicitado")
			return nil, models.NewAppError("INVALID_FIELD", "campo inválido solicitado: "+field, http.StatusBadRequest)
		}
	}

	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}

	collection := s.db.Collection("encounters")
	var encounter models.Encounter

	objectID, errVal := primitive.ObjectIDFromHex(id)
	if errVal != nil {
		return nil, models.NewAppError("INVALID_INPUT", "ID Inválido", http.StatusBadRequest)
	}

	err := collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
		options.FindOne().SetProjection(projection),
	).Decode(&encounter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.WithFields(logFields).WithError(err).Warn("encounter não encontrado")
			return nil, models.NewAppError("NOT_FOUND", "encounter não encontrado", http.StatusNotFound)
		}

		s.logger.WithFields(logFields).WithError(err).Error("falha ao buscar encounter no MongoDB")
		return nil, models.NewAppError("DATABASE_ERROR", "erro ao acessar o banco de dados", http.StatusInternalServerError)
	}

	response := s.mapToResponse(encounter, fields)

	logFields["duration"] = time.Since(startTime).String()
	s.logger.WithFields(logFields).Info("consulta de encounter realizada com sucesso")

	return response, nil
}

func (s *EncounterService) UpdateEncounterStatus(ctx context.Context, id, status string) error {
	startTime := time.Now()
	logFields := logrus.Fields{
		"operation":   "UpdateEncounterStatus",
		"encounterId": id,
		"newStatus":   status,
	}

	if !s.validStatus[status] {
		s.logger.WithFields(logFields).Warn("status inválido fornecido")
		return models.NewAppError("INVALID_STATUS", "status inválido: "+status, http.StatusBadRequest)
	}

	objectID, errVal := primitive.ObjectIDFromHex(id)
	if errVal != nil {
		return models.NewAppError("INVALID_INPUT", "ID Inválido", http.StatusBadRequest)
	}

	collection := s.db.Collection("encounters")
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status}},
	)

	if err != nil {
		s.logger.WithFields(logFields).WithError(err).Error("falha ao atualizar status no MongoDB")
		return models.NewAppError("DATABASE_ERROR", "erro ao atualizar o banco de dados", http.StatusInternalServerError)
	}

	if result.MatchedCount == 0 {
		s.logger.WithFields(logFields).Warn("nenhum encounter encontrado para atualização")
		return models.NewAppError("NOT_FOUND", "encounter não encontrado", http.StatusNotFound)
	}

	logFields["duration"] = time.Since(startTime).String()
	logFields["matchedCount"] = result.MatchedCount
	logFields["modifiedCount"] = result.ModifiedCount
	s.logger.WithFields(logFields).Info("status de encounter atualizado com sucesso")

	return nil
}

func (s *EncounterService) mapToResponse(encounter models.Encounter, fields []string) *models.EncounterResponse {
	response := &models.EncounterResponse{}

	for _, field := range fields {
		switch field {
		case "fhirId":
			response.FhirId = &encounter.FhirId
		case "fullUrl":
			response.FullURL = &encounter.FullUrl
		case "status":
			response.Status = &encounter.Status
		case "class":
			response.Class = &encounter.Class
		case "period":
			response.Period = &encounter.Period
		case "practitionerId":
			response.PractitionerID = &encounter.PractitionerID
		case "patientId":
			response.PatientID = &encounter.PatientID
		}
	}

	return response
}
