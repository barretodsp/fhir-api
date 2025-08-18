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

type PractitionerService struct {
	db          *mongo.Database
	logger      *logrus.Logger
	validFields map[string]bool
}

func NewPractitionerService(db *mongo.Database, logger *logrus.Logger) *PractitionerService {
	validFields := map[string]bool{
		"fhirId":     true,
		"givenName":  true,
		"familyName": true,
		"birthDate":  true,
		"gender":     true,
	}

	return &PractitionerService{
		db:          db,
		logger:      logger,
		validFields: validFields,
	}
}

func (s *PractitionerService) GetPractitioner(ctx context.Context, id string, fields []string) (*models.PractitionerRespose, error) {
	startTime := time.Now()
	logFields := logrus.Fields{
		"operation":       "GetPractitioner",
		"practitionerId":  id,
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

	collection := s.db.Collection("practitioners")
	var practitioner models.Practitioner

	objectID, errVal := primitive.ObjectIDFromHex(id)
	if errVal != nil {
		return nil, models.NewAppError("INVALID_INPUT", "ID Inválido", http.StatusBadRequest)
	}

	err := collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
		options.FindOne().SetProjection(projection),
	).Decode(&practitioner)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.WithFields(logFields).WithError(err).Warn("practitioner não encontrado")
			return nil, models.NewAppError("NOT_FOUND", "practitioner não encontrado", http.StatusNotFound)
		}

		s.logger.WithFields(logFields).WithError(err).Error("falha ao buscar practitioner no MongoDB")
		return nil, models.NewAppError("DATABASE_ERROR", "erro ao acessar o banco de dados", http.StatusInternalServerError)
	}

	response := s.mapToResponse(practitioner, fields)

	logFields["duration"] = time.Since(startTime).String()
	s.logger.WithFields(logFields).Info("consulta de encounter realizada com sucesso")

	return response, nil
}

func (s *PractitionerService) mapToResponse(practitioner models.Practitioner, fields []string) *models.PractitionerRespose {
	response := &models.PractitionerRespose{}

	for _, field := range fields {
		switch field {
		case "fhirId":
			response.FhirId = &practitioner.FhirId
		case "givenName":
			response.GivenName = &practitioner.GivenName
		case "familyName":
			response.FamilyName = &practitioner.FamilyName
		}
	}

	return response
}
