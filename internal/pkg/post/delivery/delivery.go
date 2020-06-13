package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/post"
	"github.com/nickeskov/db_forum/internal/pkg/utils"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
	"strconv"
	"strings"
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

func (delivery Delivery) GetPostInfoByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	relatedQuery := r.URL.Query().Get("related")

	var related []string
	if relatedQuery != "" {
		related = strings.Split(relatedQuery, ",")
	}

	postFullInfo, err := delivery.useCase.GetPostInfoByID(id, related)

	switch {
	case errors.Is(err, models.ErrInvalid):
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())

	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("post or related does not exist in db, postID=%d, related=%+v",
				id, related),
		)

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(postFullInfo)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}

func (delivery Delivery) UpdatePostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	data, err := delivery.utils.ReadAllDataFromBody(w, r)
	if err != nil {
		return
	}

	postUpdate := models.Post{
		ID: id,
	}
	if err := json.Unmarshal(data, &postUpdate); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	updatedPost, err := delivery.useCase.UpdatePostByID(postUpdate)

	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("post does not exist in db, postID=%d", id),
		)

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(updatedPost)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}

func (delivery Delivery) GetSortedPostsByThreadSlugOrID(w http.ResponseWriter, r *http.Request) {
	threadSlugOrID := mux.Vars(r)["slug_or_id"]

	queryParams := r.URL.Query()

	sort := queryParams.Get("sort")
	sinceThreadID, desc, limit := utils.ParseSinceDescLimit(queryParams)

	posts, err := delivery.useCase.GetSortedPostsByThreadSlugOrID(threadSlugOrID, sinceThreadID,
		sort, desc, limit)

	switch {
	case errors.Is(err, models.ErrInvalid):
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())

	case errors.Is(err, models.ErrDoesNotExist):
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("thread does not exist in db, threadSlugOrID=%s", threadSlugOrID),
		)

	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))

	default:
		data, err := json.Marshal(posts)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}
