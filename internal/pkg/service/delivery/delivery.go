package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/nickeskov/db_forum/internal/pkg/service"
	httpUtils "github.com/nickeskov/db_forum/pkg/http"
	"github.com/nickeskov/db_forum/pkg/logger"
	"net/http"
)

type Delivery struct {
	useCase service.UseCase
	utils   httpUtils.Utils
}

func NewDelivery(useCase service.UseCase, logger logger.Logger) Delivery {
	return Delivery{
		useCase: useCase,
		utils:   httpUtils.NewDeliveryUtils(logger),
	}
}

func (delivery Delivery) DropAllData(w http.ResponseWriter, r *http.Request) {
	if err := delivery.useCase.DropAllData(); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))
	} else {
		delivery.utils.WriteResponse(w, r, http.StatusOK, nil)
	}
}

func (delivery Delivery) GetStatus(w http.ResponseWriter, r *http.Request) {
	if status, err := delivery.useCase.GetStatus(); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))
	} else if data, err := json.Marshal(status); err != nil {
		delivery.utils.WriteResponseError(w, r, http.StatusInternalServerError,
			fmt.Sprintf("%+v", err))
	} else {
		delivery.utils.WriteResponse(w, r, http.StatusOK, data)
	}
}
