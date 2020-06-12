package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/post"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
)

type Delivery struct {
	useCase post.UseCase
	utils   httpUtils.Utils
}

func NewDelivery(useCase post.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		utils:   httpUtils.NewDeliveryUtils(logger),
	}
}

func (delivery Delivery) CreatePostsByThreadSlugOrID(w http.ResponseWriter, r *http.Request) {
	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	var newPosts models.Posts

	if err := json.Unmarshal(data, &newPosts); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	threadSlugOrID := mux.Vars(r)["slug_or_id"]

	createdPosts, err := delivery.useCase.CreatePostsByThreadSlugOrID(threadSlugOrID, newPosts)

	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("thread, forum or user does not exits with threadSlugOrID=%s",
				threadSlugOrID))

	case errors.Is(err, models.ErrConflict):
		delivery.utils.WriteResponseError(w, r, http.StatusConflict,
			fmt.Sprintf("one or many parent posts not exists in thread with threadSlugOrID=%s",
				threadSlugOrID))

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(createdPosts)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusCreated, data)
	}
}
