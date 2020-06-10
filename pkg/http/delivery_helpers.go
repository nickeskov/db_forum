package http

import (
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
)

type Utils interface {
	WriteResponseError(w http.ResponseWriter, r *http.Request, code int, msg string)
	WriteResponse(w http.ResponseWriter, r *http.Request, code int, data []byte)
	GetLogger() logger.Logger
}

type deliveryUtils struct {
	logger logger.Logger
}

func NewDeliveryUtils(logger logger.Logger) Utils {
	return deliveryUtils{
		logger: logger,
	}
}

func (utils deliveryUtils) WriteResponseError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	if err := WriteResponseError(w, code, msg); err != nil {
		utils.logger.HttpLogCallerError(r.Context(), err, err)
	}
}

func (utils deliveryUtils) WriteResponse(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if err := WriteResponse(w, code, data); err != nil {
		utils.logger.HttpLogCallerError(r.Context(), err, err)
	}
}

func (utils deliveryUtils) GetLogger() logger.Logger {
	return utils.logger
}
