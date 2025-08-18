package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type AppError struct {
	Code       string `json:"code"`              
	Message    string `json:"message"`           // Mensagem do erro (humano)
	StatusCode int    `json:"statusCode"`        // Código HTTP
	Details    string `json:"details,omitempty"` // Detalhes adicionais
}

// Error implementa a interface error padrão
func (e *AppError) Error() string {
	return e.Message
}

// LogError registra o erro no logger com contexto adicional
func (e *AppError) LogError(logger *logrus.Logger, context map[string]interface{}) {
	fields := logrus.Fields{
		"errorCode":  e.Code,
		"statusCode": e.StatusCode,
	}

	// Adiciona contexto adicional se fornecido
	for k, v := range context {
		fields[k] = v
	}

	// Determina o nível de log baseado no status code
	switch {
	case e.StatusCode >= 500:
		logger.WithFields(fields).Error(e.Message)
	case e.StatusCode >= 400:
		logger.WithFields(fields).Warn(e.Message)
	default:
		logger.WithFields(fields).Info(e.Message)
	}
}

// NewAppError cria um novo erro da aplicação
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAppErrorWithDetails cria um erro com detalhes adicionais
func NewAppErrorWithDetails(code, message string, statusCode int, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// Erros pré-definidos
var (
	ErrInternalServer = NewAppError("INTERNAL_ERROR", "Erro interno do servidor", http.StatusInternalServerError)
	ErrNotFound       = NewAppError("NOT_FOUND", "Recurso não encontrado", http.StatusNotFound)
	ErrBadRequest     = NewAppError("BAD_REQUEST", "Requisição inválida", http.StatusBadRequest)
	ErrUnauthorized   = NewAppError("UNAUTHORIZED", "Não autorizado", http.StatusUnauthorized)
	ErrForbidden      = NewAppError("FORBIDDEN", "Acesso proibido", http.StatusForbidden)
)

// LogInternalError registra erros internos com stack trace
func LogInternalError(logger *logrus.Logger, err error, context map[string]interface{}) {
	fields := logrus.Fields{
		"error": err.Error(),
	}

	for k, v := range context {
		fields[k] = v
	}

	logger.WithFields(fields).Error("Erro interno")
}
