package main

import (
	"contact-go/config"
	"contact-go/config/db"
	"contact-go/handler"
	"contact-go/helper"
	"contact-go/helper/input"
	"contact-go/helper/logger"
	"contact-go/middleware"
	"contact-go/repository"
	"contact-go/usecase"
	"log"
	"net/http"
	"os"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(true)

	contactUC := createContactUsecase(config)

	switch config.Mode {
	case "http":
		contactHTTPHandler := handler.NewContactHTTPHandler(contactUC)
		err := NewServer(config.Port, l, contactHTTPHandler)
		if err != nil {
			l.Fatal().Err(err).Msg("server fail to start")
		}
	default:
		input := input.NewInputReader(os.Stdin)
		contactCLIHandler := handler.NewContactHandler(contactUC, input)

		menu := handler.NewMenu(contactCLIHandler, input, helper.ClearTerminal, helper.ShowMenuList)
		err := menu.ShowMenu()
		if err != nil {
			l.Fatal().Err(err).Msg("server fail to start")
		}
	}
}

func createContactUsecase(config *config.Config) usecase.ContactUsecase {
	var contactRepo repository.ContactRepository
	switch config.Storage {
	case "sql":
		switch config.Database.Driver {
		case "mysql":
			db, err := db.NewMysqlDatabase(config)
			if err != nil {
				log.Fatal(err)
			}
			contactRepo = repository.NewContactMysqlRepository(db)
		case "gorm":
			db, err := db.NewGormDatabase(config)
			if err != nil {
				log.Fatal(err)
			}
			contactRepo = repository.NewContactGormRepository(db)
		default:
			log.Fatalln("database driver not existed")
		}
	case "json":
		jsonFilePath := "data/contact.json"
		contactRepo = repository.NewContactJsonRepository(jsonFilePath)
	default:
		contactRepo = repository.NewContactRepository()
	}
	return usecase.NewContactUsecase(contactRepo)
}

func NewServer(port string, logger *logger.Logger, handler handler.ContactHTTPHandler) error {
	mux := http.NewServeMux()

	muxMiddleware := new(middleware.Middleware)
	muxMiddleware.Handler = mux

	muxMiddleware.Use(middleware.Cors)
	muxMiddleware.Use(middleware.ContentTypeJson)
	muxMiddleware.Use(
		func(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
			return middleware.Log(logger, w, r, next)
		},
	)
	muxMiddleware.Use(
		func(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
			return middleware.Error(logger, w, r, next)
		},
	)

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		switch r.Method {
		case "GET":
			handler.List(w, r)
		case "POST":
			handler.Add(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/contacts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		switch r.Method {
		case "GET":
			handler.Detail(w, r)
		case "PATCH":
			handler.Update(w, r)
		case "DELETE":
			handler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: muxMiddleware,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	logger.Info().Msgf("live on http://localhost:%s", port)
	return nil
}
