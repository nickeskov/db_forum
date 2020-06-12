package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
)

type Delivery struct {
	useCase forum.UseCase
	utils   httpUtils.Utils
}

func NewDelivery(useCase forum.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		utils:   httpUtils.NewDeliveryUtils(logger),
	}
}

func (delivery Delivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	var newForum models.Forum

	if err := json.Unmarshal(data, &newForum); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if validationErr := newForum.Validate(); validationErr != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, validationErr.Error())
		return
	}

	createdForum, err := delivery.useCase.Create(newForum)
	switch err {
	case models.ErrConflict:
		existingForum, err := delivery.useCase.GetBySlug(newForum.Slug)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(existingForum)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusConflict, data)

	case models.ErrBadForeign:
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("cannot create forum, user with nickname=%s does not exits",
				newForum.User))

	case nil:
		data, err := json.Marshal(createdForum)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusCreated, data)

	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
	}

}

func (delivery Delivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	existingForum, err := delivery.useCase.GetBySlug(slug)
	switch err {
	case models.ErrDoesNotExist:
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("forum with slug=%s does not exits", slug))

	case nil:
		data, err := json.Marshal(existingForum)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)

	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (delivery Delivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	sinceNickname, desc, limit := utils.ParseSinceDescLimit(r.URL.Query())

	users, err := delivery.useCase.GetForumUsersBySlug(slug, sinceNickname, desc, limit)
	switch err {
	case models.ErrDoesNotExist:
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("forum with slug=%s does not exist", slug))

	case models.ErrInvalid:
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, models.ErrInvalid.Error())

	case nil:
		data, err := json.Marshal(users)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)

	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))
	}
}
