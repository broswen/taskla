package server

import (
	"fmt"
	"net/http"

	"github.com/broswen/taskla/pkg/auth"
	"github.com/broswen/taskla/pkg/group"
	"github.com/broswen/taskla/pkg/task"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router       chi.Router
	authService  auth.Service
	groupService group.Service
	taskService  task.Service
	logger       zerolog.Logger
}

func New() (Server, error) {
	log.Info().Msg("Initializing server...")
	authService, err := auth.NewService()
	if err != nil {
		return Server{}, fmt.Errorf("creating auth service: %w", err)
	}

	groupService, err := group.NewService()
	if err != nil {
		return Server{}, fmt.Errorf("creating group service: %w", err)
	}

	taskService, err := task.NewService()
	if err != nil {
		return Server{}, fmt.Errorf("creating group service: %w", err)
	}

	logger := httplog.NewLogger("taskla", httplog.Options{
		JSON: true,
	})
	return Server{
		router:       chi.NewRouter(),
		authService:  authService,
		groupService: groupService,
		taskService:  taskService,
		logger:       logger,
	}, nil
}

func (s Server) Start() error {
	log.Info().Msg("Starting server")
	return http.ListenAndServe(":8080", s.router)
}

func (s Server) Routes() {

	s.router.Use(httplog.RequestLogger(s.logger))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// GET /ping for health check
	s.router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	// routes for login and registering
	s.router.HandleFunc("/register", auth.Register(s.authService))
	s.router.HandleFunc("/login", auth.Login(s.authService))

	// behind jwt middleware

	s.router.Route("/group", func(r chi.Router) {
		r.Use(auth.JWT(s.authService))

		r.Get("/", group.Get(s.groupService))
		r.Get("/{id}/task", task.GetTasksByGroup(s.taskService, s.groupService))
		r.Post("/", group.Create(s.groupService))
		r.Get("/{id}", group.GetById(s.groupService))
		r.Put("/{id}", group.Update(s.groupService))
		r.Delete("/{id}", group.Delete(s.groupService))
	})

	s.router.Route("/task", func(r chi.Router) {
		r.Use(auth.JWT(s.authService))

		r.Get("/", task.Get(s.taskService))
		r.Post("/", task.Create(s.taskService))
		r.Get("/{id}", task.GetById(s.taskService))
		r.Put("/{id}", task.Update(s.taskService))
		r.Delete("/{id}", task.Delete(s.taskService))
	})

}
