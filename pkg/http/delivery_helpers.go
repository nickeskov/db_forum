package http

import (
	"github.com/nickeskov/db_forum/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
)

type Utils interface {
	GetLogger() logger.Logger
	WriteResponseError(w http.ResponseWriter, r *http.Request, code int, msg string)
	WriteResponse(w http.ResponseWriter, r *http.Request, code int, data []byte)
	ReadAllDataFromBody(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type deliveryUtils struct {
	logger logger.Logger
}

func NewDeliveryUtils(logger logger.Logger) Utils {
	return deliveryUtils{
		logger: logger,
	}
}

func (utils deliveryUtils) GetLogger() logger.Logger {
	return utils.logger
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

func (utils deliveryUtils) ReadAllDataFromBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	switch data, err := ioutil.ReadAll(r.Body); err {
	case nil:
		return data, nil
	case io.EOF:
		utils.WriteResponseError(w, r, http.StatusBadRequest, "empty body")
		return nil, err
	default:
		utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
		return nil, err
	}
}
