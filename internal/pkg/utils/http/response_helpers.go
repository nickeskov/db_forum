package http

import (
	"encoding/json"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"net/http"
)

func WriteResponseError(w http.ResponseWriter, code int, msg string) error {
	w.WriteHeader(code)

	data, err := json.Marshal(models.NewError(msg))
	if err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return err
	}

	if _, err := w.Write(data); err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return err
	}

	return nil
}

func WriteResponse(w http.ResponseWriter, code int, data []byte) error {
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return err
	}

	return nil
}
