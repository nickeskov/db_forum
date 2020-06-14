package db_forum

import (
	"github.com/gorilla/mux"
	forumDelivery "github.com/nickeskov/db_forum/internal/pkg/forum/delivery"
	forumRepository "github.com/nickeskov/db_forum/internal/pkg/forum/repository"
	forumUseCase "github.com/nickeskov/db_forum/internal/pkg/forum/usecase"
	postDelivery "github.com/nickeskov/db_forum/internal/pkg/post/delivery"
	postRepository "github.com/nickeskov/db_forum/internal/pkg/post/repository"
	postUseCase "github.com/nickeskov/db_forum/internal/pkg/post/usecase"
	serviceDelivery "github.com/nickeskov/db_forum/internal/pkg/service/delivery"
	serviceRepository "github.com/nickeskov/db_forum/internal/pkg/service/repository"
	serviceUseCase "github.com/nickeskov/db_forum/internal/pkg/service/usecase"
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

const loggerKey = 1

func StartNew() {
	customLogger := logger.NewTextFormatSimpleLogger(os.Stdout, loggerKey)
	customLogger.Printf(">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now())

	// TODO(nickeskov): hardcoded database credentials
	dbConnPool, err := ConnectToDB(
		"localhost",
		"my_db_forum",
		"my_db_forum",
		"my_db_forum",
		10,
	)

	if err != nil {
		customLogger.Fatalln("cannot connect to postgres:", err)
	} else {
		customLogger.Println("successfully connected to postgres")
	}

	userRepo := userRepository.NewRepository(dbConnPool)
	forumRepo := forumRepository.NewRepository(dbConnPool)
	threadRepo := threadRepository.NewRepository(dbConnPool, forumRepo)
	postRepo := postRepository.NewRepository(dbConnPool)
	serviceRepo := serviceRepository.NewRepository(dbConnPool)

	userUC := userUseCase.NewUseCase(userRepo)
	forumUC := forumUseCase.NewUseCase(forumRepo)
	threadUC := threadUseCase.NewUseCase(threadRepo)
	postUC := postUseCase.NewUseCase(postRepo, userRepo, forumRepo, threadRepo)
	serviceUC := serviceUseCase.NewUseCase(serviceRepo)

	userHandlers := userDelivery.NewDelivery(userUC, customLogger)
	forumHandlers := forumDelivery.NewDelivery(forumUC, customLogger)
	threadHandlers := threadDelivery.NewDelivery(threadUC, customLogger)
	postHandlers := postDelivery.NewDelivery(postUC, customLogger)
	serviceHandlers := serviceDelivery.NewDelivery(serviceUC, customLogger)

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

	router.HandleFunc("/thread/{slug_or_id}/details", threadHandlers.GetThreadBySlugOrID).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/details", threadHandlers.UpdateThreadBySlugOrID).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/vote", threadHandlers.VoteThreadBySlugOrID).Methods(http.MethodPost)

	router.HandleFunc("/thread/{slug_or_id}/create", postHandlers.CreatePostsByThreadSlugOrID).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/posts", postHandlers.GetSortedPostsByThreadSlugOrID).Methods(http.MethodGet)

	router.HandleFunc("/post/{id}/details", postHandlers.GetPostInfoByID).Methods(http.MethodGet)
	router.HandleFunc("/post/{id}/details", postHandlers.UpdatePostByID).Methods(http.MethodPost)

	router.HandleFunc("/service/clear", serviceHandlers.DropAllData).Methods(http.MethodPost)
	router.HandleFunc("/service/status", serviceHandlers.GetStatus).Methods(http.MethodGet)

	// TODO(nickeskov): hardcoded server address and port
	if err := http.ListenAndServe(":5000", router); err != nil {
		customLogger.Fatalln("cannot start service:", err)
	}
}
