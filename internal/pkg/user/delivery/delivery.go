package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/user"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
)

type Delivery struct {
	useCase user.UseCase
	utils   httpUtils.Utils
}

func NewDelivery(useCase user.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		utils:   httpUtils.NewDeliveryUtils(logger),
	}
}

func (delivery Delivery) getUserFomBody(w http.ResponseWriter, r *http.Request) (models.User, error) {
	nickname := mux.Vars(r)["nickname"]

	data, err := ioutil.ReadAll(r.Body)
	switch {
	case err == io.EOF:
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, "empty body")
		return models.User{}, err
	case err != nil:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
		return models.User{}, err
	}

	newUser := models.User{
		Nickname: nickname,
	}

	if err := json.Unmarshal(data, &newUser); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, err.Error())
		return models.User{}, err
	}

	if validationErr := newUser.Validate(); validationErr != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusBadRequest, validationErr.Error())
		return models.User{}, validationErr
	}

	return newUser, nil
}

func (delivery Delivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser, err := delivery.getUserFomBody(w, r)
	if err != nil {
		return
	}

	userCreateErr := delivery.useCase.Create(newUser)

	switch userCreateErr {
	case models.ErrAlreadyExist:
		users, err := delivery.useCase.GetWithSameNicknameAndEmail(newUser.Nickname, newUser.Email)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(users)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusConflict, data)
	case nil:
		data, err := json.Marshal(newUser)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusCreated, data)
	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, userCreateErr.Error())
	}
}

func (delivery Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	storedUser, getUserErr := delivery.useCase.GetByNickname(nickname)
	switch getUserErr {
	case models.ErrDoesNotExist:
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("user with nickname=%s does not exist", nickname))

	case nil:
		data, err := json.Marshal(storedUser)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)

	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, getUserErr.Error())
	}
}

func (delivery Delivery) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userForUpdate, err := delivery.getUserFomBody(w, r)
	if err != nil {
		return
	}

	updatedUser, userUpdateErr := delivery.useCase.UpdateByNickname(userForUpdate)
	switch userUpdateErr {
	case models.ErrDoesNotExist:
		delivery.utils.WriteResponseError(w, r, http.StatusNotFound,
			fmt.Sprintf("user with nickname=%s does not exist", userForUpdate.Nickname))

	case models.ErrConflict:
		delivery.utils.WriteResponseError(w, r, http.StatusConflict,
			fmt.Sprintf("update for nickname=%s conflicts with other user", userForUpdate.Nickname))

	case nil:
		data, err := json.Marshal(updatedUser)
		if err != nil {
			delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		delivery.utils.WriteResponse(w, r, http.StatusOK, data)

	default:
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError, userUpdateErr.Error())
	}
}
