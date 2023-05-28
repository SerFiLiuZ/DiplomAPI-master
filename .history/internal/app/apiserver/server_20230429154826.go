package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"

	"github.com/gopherschool/http-rest-api/internal/app/model"

	"github.com/gopherschool/http-rest-api/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName        = "fillin"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
	CORS = "http://localhost:3000"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.HandleFunc("/api/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/api/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/api/boards", s.handleGetAllBoard()).Methods("GET")
	s.router.HandleFunc("/api/boards/id/{boardId:[0-9]+}", s.handlerGetBoardCards()).Methods("GET")
	s.router.HandleFunc("/api/task/id/{taskId:[0-9]+}", s.handleGetTask()).Methods("GET")

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	// private := s.router.PathPrefix("/private").Subrouter()
	// private.Use(s.authenticateUser)
	// private.HandleFunc("/whoami", s.handleWhoami())

}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//setupCorsResponse(w, r)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//setupCorsResponse(w, r)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || u.Password != req.Password {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		//s.sessionStore = sessions.NewCookieStore([]byte("some-secret-key"))
		session, err := s.sessionStore.New(r, sessionName)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		session.Values["user_status"] = u.Status
		session.Options.MaxAge = 86400 // 24 часа в секундах
		session.Options.Path = "/"
		session.Options.Secure = true
		session.Options.SameSite = http.SameSiteLaxMode

		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		log.Print(session)

		type Response struct {
			Success bool `json:"success"`
		}

		s.respond(w, r, http.StatusCreated, &Response{Success: true})
	}
}

func (s *server) handleGetAllBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		log.Print("handleGetAllBoard")
		log.Print(session)

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		boards, err := s.store.Board().GetAllBoards(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, boards)
	}
}

func (s *server) handlerGetBoardCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//setupCorsResponse(w, r)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// Получение идентификатора доски из URL-параметра
		boardId := mux.Vars(r)["boardId"]
		boardIdInt, err := strconv.Atoi(boardId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		boards, err := s.store.Card().FindCardsByBoardID(id.(int), boardIdInt)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, boards)
	}
}

func (s *server) handleGetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//setupCorsResponse(w, r)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// Получение идентификатора доски из URL-параметра
		taskId := mux.Vars(r)["taskId"]
		taskIdInt, err := strconv.Atoi(taskId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		tasks, err := s.store.Task().FindTasksByCardID(id.(int), taskIdInt)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, tasks)
	}
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
