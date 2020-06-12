package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/thread"
	"github.com/nickeskov/db_forum/internal/pkg/utils"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
)

type Delivery struct {
	useCase thread.UseCase
	utils   httpUtils.Utils
}

func NewDelivery(useCase thread.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		utils:   httpUtils.NewDeliveryUtils(logger),
	}
}

func (delivery Delivery) CreateThread(w http.ResponseWriter, r *http.Request) {
	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	var newThread models.Thread

	if err := json.Unmarshal(data, &newThread); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if validationErr := newThread.Validate(); validationErr != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, validationErr.Error())
		return
	}

	newThread.Forum = mux.Vars(r)["slug"]

	createdThread, err := delivery.useCase.Create(newThread)
	switch {
	case errors.Is(err, models.ErrConflict):
		existingThread, err := delivery.useCase.GetBySlugOrID(newThread.Slug)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
				fmt.Sprintf("%+v", err))
			return
		}

		data, err := json.Marshal(existingThread)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusConflict, data)

	case errors.Is(err, models.ErrBadForeign):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("cannot create thread (author=%s or forum=%s does not exist)",
				newThread.Author, newThread.Forum))

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(createdThread)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusCreated, data)
	}
}

func (delivery Delivery) GetThreadsByForumSlug(w http.ResponseWriter, r *http.Request) {
	forumSlug := mux.Vars(r)["slug"]

	since, desc, limit := utils.ParseSinceDescLimit(r.URL.Query())

	threads, err := delivery.useCase.GetThreadsByForumSlug(forumSlug, since, desc, limit)
	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("forum with slug=%s does not exits", forumSlug))

	case errors.Is(err, models.ErrInvalid):
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, "bad request")

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(threads)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}

func (delivery Delivery) GetThreadBySlugOrID(w http.ResponseWriter, r *http.Request) {
	slugOrID := mux.Vars(r)["slug_or_id"]

	threads, err := delivery.useCase.GetBySlugOrID(slugOrID)
	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("thread with slug_or_id=%s does not exits", slugOrID))

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(threads)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}

func (delivery Delivery) UpdateThreadBySlugOrID(w http.ResponseWriter, r *http.Request) {
	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	var threadUpdate models.Thread

	if err := json.Unmarshal(data, &threadUpdate); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	slugOrID := mux.Vars(r)["slug_or_id"]

	updatedThread, err := delivery.useCase.UpdateBySlugOrID(slugOrID, threadUpdate)
	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("thread with slug_or_id=%s does not exits", slugOrID))

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(updatedThread)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}

func (delivery Delivery) VoteThreadBySlugOrID(w http.ResponseWriter, r *http.Request) {
	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	var vote models.Vote

	if err := json.Unmarshal(data, &vote); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	slugOrID := mux.Vars(r)["slug_or_id"]

	updatedThread, err := delivery.useCase.VoteBySlugOrID(slugOrID, vote)
	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("thread with slug_or_id=%s or author=%s does not exits",
				slugOrID, vote.Nickname))

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(updatedThread)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}
