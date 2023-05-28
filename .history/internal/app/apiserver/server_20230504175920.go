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

	s.router.HandleFunc("/api/sessions", s.handleSessionsCreate()).Methods("POST")

	s.router.HandleFunc("/api/boards", s.handleGetAllBoard()).Methods("GET")
	s.router.HandleFunc("/api/boards/id/{boardId:[0-9]+}", s.handlerGetBoardCards()).Methods("GET")

	s.router.HandleFunc("/api/tasks/id/{taskId:[0-9]+}", s.handleGetTask()).Methods("GET")
	s.router.HandleFunc("/api/boards/create", s.hanleCreateBoard()).Methods("POST")

	s.router.HandleFunc("/api/boards/delete", s.hanleDeleteBoard()).Methods("POST")
	s.router.HandleFunc("/api/cards/create", s.hanleCreateCard()).Methods("POST")

	s.router.HandleFunc("/api/cards/delete", s.hanleDeleteCard()).Methods("POST")
	s.router.HandleFunc("/api/tasks/delete", s.hanleDeleteTask()).Methods("POST")

	s.router.HandleFunc("/api/tasks/create", s.hanleCreateTask()).Methods("POST")
	s.router.HandleFunc("/api/tasks/complite", s.hanleCompliteTask()).Methods("POST")

	s.router.HandleFunc("/api/workers", s.hanleGetWorkers()).Methods("GET")
	s.router.HandleFunc("/api/workers/dismiss", s.hanleWorkersDismiss()).Methods("POST")

	s.router.HandleFunc("/api/applications", s.hanleGetApplications()).Methods("GET")
	s.router.HandleFunc("/api/boards/done/{boardId:[0-9]+}", s.hanleGetBoardsDone()).Methods("GET")

	s.router.HandleFunc("/api/applications/accept", s.hanleAcceptApplications()).Methods("POST")
	s.router.HandleFunc("/api/applications/reject", s.hanleRejectApplications()).Methods("POST")

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

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
			time.Since(start),
		)
	})
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || u.Autorization.Password != req.Password {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.New(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		session.Values["user_status"] = u.Status

		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		cookie := &http.Cookie{
			Name:  "user_status",
			Value: u.Status,
			Path:  "/",
		}
		http.SetCookie(w, cookie)

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

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		status, ok := session.Values["user_status"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		boards, err := s.store.Board().GetAllBoards(id.(int), status.(string))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, boards)
	}
}

func (s *server) handlerGetBoardCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		status, ok := session.Values["user_status"]
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

		boards, err := s.store.Card().FindCardsByBoardID(id.(int), status.(string), boardIdInt)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, boards)
	}
}

func (s *server) handleGetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		status, ok := session.Values["user_status"]
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

		tasks, err := s.store.Task().FindTasksByCardID(id.(int), status.(string), taskIdInt)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for _, t := range tasks {
			t.Assigned, err = s.store.User().FindAllById(t.Assigned)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}
		s.respond(w, r, http.StatusOK, tasks)

	}
}

func (s *server) hanleCreateBoard() http.HandlerFunc {
	type request struct {
		BoardTitle string `json:"BoardTitle"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err = s.store.Board().CreateBoard(req.BoardTitle, id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type Response struct {
			Success bool `json:"success"`
		}

		s.respond(w, r, http.StatusCreated, &Response{Success: true})
	}
}

func (s *server) hanleDeleteBoard() http.HandlerFunc {
	type request struct {
		BoardId int `json:"BoardId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.Board().DeleteBoard(req.BoardId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleDeleteCard() http.HandlerFunc {
	type request struct {
		CardId int `json:"CardId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.Card().DeleteCard(req.CardId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleDeleteTask() http.HandlerFunc {
	type request struct {
		TaskId int `json:"TaskId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.Task().DeleteTask(req.TaskId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleCreateCard() http.HandlerFunc {
	type request struct {
		CardTitle string `json:"CardTitle"`
		CardDes   string `json:"CardDes"`
		BoardId   int    `json:"BoardId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			log.Print("nen [eqyz&]")
			return
		}

		err := s.store.Card().CreateCard(req.CardTitle, req.CardDes, req.BoardId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusCreated, &Response{Success: true})
	}
}

func (s *server) hanleGetWorkers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		users, err := s.store.User().FindAllByIdManager(id.(int))
		if err != nil {
			log.Print(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) hanleCreateTask() http.HandlerFunc {
	type request struct {
		TaskTitle           string        `json:"title"`
		TaskDes             string        `json:"description"`
		TaskDueDate         string        `json:"dueDate"`
		TaskSelectedWorkers []*model.User `json:"selectedWorkers"`
		CardId              int           `json:"cardId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		log.Print(req)
		err := s.store.Task().CreateTask(req.TaskTitle, req.TaskDes, req.TaskDueDate, req.TaskSelectedWorkers, req.CardId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.DBController().SendMesage(req.TaskTitle, req.TaskDes, req.TaskDueDate, req.TaskSelectedWorkers)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusCreated, &Response{Success: true})
	}
}

func (s *server) hanleCompliteTask() http.HandlerFunc {
	type request struct {
		TaskId int `json:"TaskId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		log.Print(req)

		err := s.store.Task().CompliteTask(req.TaskId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.DBController().UpdateDoneByTask(req.TaskId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleGetApplications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		applications, err := s.store.DBController().GetAppications(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, applications)

	}
}

func (s *server) hanleGetBoardsDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение идентификатора доски из URL-параметра
		boardId := mux.Vars(r)["boardId"]
		boardIdInt, err := strconv.Atoi(boardId)
		log.Print(boardIdInt)
		if err != nil {
			log.Print(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		done, err := s.store.DBController().GetBoardDone(boardIdInt)
		if err != nil {
			log.Print(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, done)
	}
}

func (s *server) hanleAcceptApplications() http.HandlerFunc {
	type request struct {
		ChatID int `json:"chatID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

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
		log.Print(req.ChatID)

		err = s.store.DBController().AcceptApplication(req.ChatID, id.(int))

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleRejectApplications() http.HandlerFunc {
	type request struct {
		ChatID int `json:"chatID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.DBController().RejectApplication(req.ChatID)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
	}
}

func (s *server) hanleWorkersDismiss() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		log.Print(req)

		err := s.store.User().DeliteUser(req.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type Response struct {
			Success bool `json:"success"`
		}
		s.respond(w, r, http.StatusOK, &Response{Success: true})
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
