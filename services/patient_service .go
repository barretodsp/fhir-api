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

type PatientService struct {
	db          *mongo.Database
	logger      *logrus.Logger
	validFields map[string]bool
}

func NewPatientService(db *mongo.Database, logger *logrus.Logger) *PatientService {
	validFields := map[string]bool{
		"fhirId":     true,
		"givenName":  true,
		"familyName": true,
		"birthDate":  true,
		"gender":     true,
	}

	return &PatientService{
		db:          db,
		logger:      logger,
		validFields: validFields,
	}
}

func (s *PatientService) GetPatient(ctx context.Context, id string, fields []string) (*models.PatientResponse, error) {
	startTime := time.Now()
	logFields := logrus.Fields{
		"operation":       "GetPatient",
		"patientId":       id,
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

	collection := s.db.Collection("patients")
	var patient models.Patient

	objectID, errVal := primitive.ObjectIDFromHex(id)
	if errVal != nil {
		return nil, models.NewAppError("INVALID_INPUT", "ID Inválido", http.StatusBadRequest)
	}
	
	err := collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
		options.FindOne().SetProjection(projection),
	).Decode(&patient)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.WithFields(logFields).WithError(err).Warn("patient não encontrado")
			return nil, models.NewAppError("NOT_FOUND", "patient não encontrado", http.StatusNotFound)
		}

		s.logger.WithFields(logFields).WithError(err).Error("falha ao buscar patient no MongoDB")
		return nil, models.NewAppError("DATABASE_ERROR", "erro ao acessar o banco de dados", http.StatusInternalServerError)
	}

	response := s.mapToResponse(patient, fields)
	logFields["duration"] = time.Since(startTime).String()
	s.logger.WithFields(logFields).Info("consulta de encounter realizada com sucesso")

	return response, nil
}

func (s *PatientService) mapToResponse(patient models.Patient, fields []string) *models.PatientResponse {
	response := &models.PatientResponse{}

	for _, field := range fields {
		switch field {
		case "fhirId":
			response.FhirId = &patient.FhirId
		case "givenName":
			response.GivenName = &patient.GivenName
		case "familyName":
			response.FamilyName = &patient.FamilyName
		case "birthDate":
			response.BirthDate = &patient.BirthDate
		case "gender":
			response.Gender = &patient.Gender
		}
	}

	return response
}
