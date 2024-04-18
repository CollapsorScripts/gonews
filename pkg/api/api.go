package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"newsaggr/pkg/database/model"
	"newsaggr/pkg/logger"
	"strconv"
	"time"
)

type Service struct {
	r      *mux.Router
	Server *http.Server
}

// New - создание экземпляра сервиса апи
func New() *Service {
	srv := &Service{r: mux.NewRouter()}
	srv.Server = srv.endpoints()

	return srv
}

func (h Service) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//------------------Подзагрузка маршрутов-----------------------------------

func (h *Service) endpoints() *http.Server {
	h.r.HandleFunc("/news/{n}", h.posts).Methods(http.MethodGet, http.MethodOptions)
	h.r.Use(cors.Default().Handler, h.middleware)
	// веб-приложение
	h.r.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir("C:/Users/Alex/GolandProjects/newsaggr/cmd/entrypoint/webapp"))))

	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "application/json"},
	})
	handler := crs.Handler(h.r)

	srv := &http.Server{
		Addr:         ":80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.AllowAll().Handler(handler),
	}

	return srv
}

func (h *Service) posts(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	news, err := model.FindLimit(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
