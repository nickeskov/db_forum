package db_forum

import (
	"github.com/gorilla/mux"
	forumDelivery "github.com/nickeskov/db_forum/internal/pkg/forum/delivery"
	forumRepository "github.com/nickeskov/db_forum/internal/pkg/forum/repository"
	forumUseCase "github.com/nickeskov/db_forum/internal/pkg/forum/usecase"
	threadDelivery "github.com/nickeskov/db_forum/internal/pkg/thread/delivery"
	threadRepository "github.com/nickeskov/db_forum/internal/pkg/thread/repository"
	threadUseCase "github.com/nickeskov/db_forum/internal/pkg/thread/usecase"
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

	// TODO(nickeskov): hardcoded database credentials
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
	forumRepo := forumRepository.NewRepository(dbConnPool)
	threadRepo := threadRepository.NewRepository(dbConnPool, forumRepo)

	userUC := userUseCase.NewUseCase(userRepo)
	forumUC := forumUseCase.NewUseCase(forumRepo)
	threadUC := threadUseCase.NewUseCase(threadRepo)

	userHandlers := userDelivery.NewDelivery(userUC, customLogger)
	forumHandlers := forumDelivery.NewDelivery(forumUC, customLogger)
	threadHandlers := threadDelivery.NewDelivery(threadUC, customLogger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(middleware.JsonContentTypeMiddleware)

	router.HandleFunc("/user/{nickname}/profile", userHandlers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/create", userHandlers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", userHandlers.UpdateUser).Methods(http.MethodPost)

	router.HandleFunc("/forum/create", forumHandlers.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/details", forumHandlers.GetForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/users", forumHandlers.GetForumUsers).Methods(http.MethodGet)

	router.HandleFunc("/forum/{slug}/create", threadHandlers.CreateThread).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/threads", threadHandlers.GetThreadsByForumSlug).Methods(http.MethodGet)

	// TODO(nickeskov): hardcoded server address and port
	if err := http.ListenAndServe(":5000", router); err != nil {
		customLogger.Fatalln("cannot start service:", err)
	}
}
