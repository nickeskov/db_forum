package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	httpUtils "github.com/nickeskov/db_forum/internal/pkg/utils/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
)

type Delivery struct {
	useCase forum.UseCase
	logger  logger.Logger
}

func NewDelivery(useCase forum.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		logger:  logger,
	}
}

func (delivery Delivery) writeResponseError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	err := httpUtils.WriteResponseError(w, code, msg)
	if err != nil {
		delivery.logger.HttpLogCallerError(r.Context(), err, err)
	}
}

func (delivery Delivery) writeResponse(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	err := httpUtils.WriteResponse(w, code, data)
	if err != nil {
		delivery.logger.HttpLogCallerError(r.Context(), err, err)
	}
}

func (delivery Delivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	switch {
	case err == io.EOF:
		delivery.writeResponseError(w, r, http.StatusBadRequest, "empty body")
	case err != nil:
		delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
	}

	var newForum models.Forum

	if err := json.Unmarshal(data, &newForum); err != nil {
		delivery.writeResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if validationErr := newForum.Validate(); validationErr != nil {
		delivery.writeResponseError(w, r, http.StatusBadRequest, validationErr.Error())
		return
	}

	createdForum, err := delivery.useCase.Create(newForum)
	switch err {
	case models.ErrConflict:
		existingForum, err := delivery.useCase.GetBySlug(newForum.Slug)
		if err != nil {
			delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(existingForum)
		if err != nil {
			delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.writeResponse(w, r, http.StatusConflict, data)

	case models.ErrBadForeign:
		delivery.writeResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("cannot create forum, user with nickname=%s does not exits", newForum.User))

	case nil:
		data, err := json.Marshal(createdForum)
		if err != nil {
			delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.writeResponse(w, r, http.StatusCreated, data)

	default:
		delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
	}

}

func (delivery Delivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	existingForum, err := delivery.useCase.GetBySlug(slug)
	switch err {
	case models.ErrDoesNotExist:
		delivery.writeResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("forum with slug=%s does not exits", slug))
	case nil:
		data, err := json.Marshal(existingForum)
		if err != nil {
			delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.writeResponse(w, r, http.StatusOK, data)
	default:
		delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (delivery Delivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	urlQuery := r.URL.Query()

	sinceNickname := urlQuery.Get("since")

	var limit string
	if limit = urlQuery.Get("limit"); limit == "" {
		limit = "1"
	}

	var desc string
	if desc = urlQuery.Get(""); desc == "" {
		desc = "false"
	}

	users, err := delivery.useCase.GetForumUsersBySlug(slug, sinceNickname, desc, limit)
	switch err {
	case models.ErrDoesNotExist:
		delivery.writeResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("forum with slug=%s does not exist", slug))

	case models.ErrInvalid:
		delivery.writeResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("forum with slug=%s does not exist", slug))

	case nil:
		data, err := json.Marshal(users)
		if err != nil {
			delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.writeResponse(w, r, http.StatusOK, data)

	default:
		delivery.writeResponseError(w, r, http.StatusInternalServerError, err.Error())
	}
}
