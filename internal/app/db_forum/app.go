package db_forum

import (
	"github.com/gorilla/mux"
	userDelivery "github.com/nickeskov/db_forum/internal/pkg/user/delivery"
	userRepository "github.com/nickeskov/db_forum/internal/pkg/user/repository"
	userUseCase "github.com/nickeskov/db_forum/internal/pkg/user/usecase"
	"github.com/nickeskov/db_forum/pkg/logger"
	"github.com/nickeskov/db_forum/pkg/middleware"
	"net/http"
	"os"
	"time"
)

func StartNew() {
	customLogger := logger.NewTextFormatSimpleLogger(os.Stdout)
	customLogger.Printf(">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now())

	// TODO(nickeskov): hardcode
	dbConnPool, err := ConnectToDB(
		"localhost",
		"my_db_forum",
		"my_db_forum",
		"my_db_forum",
	)
	if err != nil {
		customLogger.Fatalln("cannot connect to postgres:", err)
	} else {
		customLogger.Println("successfully connected to postgres")
	}

	userRepo := userRepository.NewRepository(dbConnPool)

	userUC := userUseCase.NewUseCase(userRepo)

	userHandlers := userDelivery.NewDelivery(userUC, customLogger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(middleware.JsonContentTypeMiddleware)

	router.HandleFunc("/user/{nickname}/profile", userHandlers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/create", userHandlers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", userHandlers.UpdateUser).Methods(http.MethodPost)

	// TODO(nickeskov): hardcode
	if err := http.ListenAndServe(":5000", router); err != nil {
		customLogger.Fatalln("cannot start service:", err)
	}
}
