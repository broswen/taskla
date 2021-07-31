package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

type RegisterUserRequest struct {
	*User
	RegistrationCode string `json:"code"`
	Password         string `json:"password"`
}

func (ru *RegisterUserRequest) Bind(r *http.Request) error {
	if ru.User == nil {
		return errors.New("missing User fields")
	}
	return nil
}

type RegisterUserResponse struct {
	Username string `json:"username"`
}

func (ru *RegisterUserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Register(s Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Info().Msg("register 1")
		data := &RegisterUserRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Error().Err(err).Msgf("Couldn't bind RegisterUserRequest")
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		regCode, err := s.GetRegistrationCode(data.RegistrationCode)
		if err != nil {
			log.Warn().Err(err).Msgf("Registration code expired")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		if regCode.Expiration.Before(time.Now()) || regCode.Used {
			log.Warn().Err(err).Msgf("Registration code expired")
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		err = s.createUser(data.Username, data.Password, regular)
		if err != nil {
			oplog.Error().Err(err).Msgf("Couldn't create user")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		regCode.Used = true
		if err := s.UpdateRegistrationCode(regCode); err != nil {
			oplog.Error().Err(err).Msgf("Couldn't update registration code: %v", regCode)
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		render.Render(w, r, &RegisterUserResponse{data.Username})
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr *LoginRequest) Bind(r *http.Request) error {
	if lr.Username == "" {
		return errors.New("missing username")
	}

	if lr.Password == "" {
		return errors.New("missing password")
	}
	return nil
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (lr *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Login(s Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		data := &LoginRequest{}
		if err := render.Bind(r, data); err != nil {
			oplog.Error().Err(err).Msg("Couldn't bind LoginRequest")
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		jwt, err := s.Login(data.Username, data.Password)
		if err != nil {
			oplog.Error().Err(err).Msg("Couldn't login")
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, &LoginResponse{jwt})
	}
}

// jwt validation middleware

func JWT(s Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			oplog := httplog.LogEntry(r.Context())
			jwt := strings.Split(r.Header.Get("Authorization"), " ")[1]
			if jwt == "" {
				oplog.Error().Msg("Missing JWT")
				render.Render(w, r, ErrInvalidRequest(errors.New("missing jwt")))
				return
			}

			claims, err := s.ValidateJWT(jwt)
			if err != nil {
				oplog.Error().Err(err).Msg("Invalid or expired JWT")
				render.Render(w, r, ErrInternalServer(err))
				return
			}

			ctx := context.WithValue(r.Context(), "subject", claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
