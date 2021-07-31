package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/broswen/taskla/pkg/storage"
)

// TODO get secret from outside of applicatoin
const SECRETSTRING = "TODOSECRETSTRING"

type Service struct {
	r storage.Repository
}

func NewService() (Service, error) {
	repo, err := storage.NewPostgres()
	if err != nil {
		return Service{}, fmt.Errorf("creating repository: %w", err)
	}
	rand.Seed(time.Now().UnixNano())
	return Service{
		r: repo,
	}, nil
}

// type to store a registration code, used for registering
type RegistrationCode struct {
	Code       string    `json:"code"`
	Expiration time.Time `json:"expiration"`
	Used       bool      `json:"used"`
}

// used to pass custom claims in jwt
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

// used to add authorization levels
type UserRole string

const (
	regular UserRole = "regular"
	admin   UserRole = "admin"
)

func (s Service) Register(username, password, code string) error {
	regCode, err := s.GetRegistrationCode(code)
	if err != nil {
		return err
	}
	if regCode.Expiration.Before(time.Now()) {
		return fmt.Errorf("registration code expired")
	}
	// validate username
	// validate password is complex enough
	existingUser, err := s.getUser(username)
	if err != nil {
		return err
	}
	if (User{}) == existingUser {
		return fmt.Errorf("username already in use")
	}
	// bcrypt password
	err = s.createUser(username, password, regular)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) createUser(username, password string, role UserRole) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	if err != nil {
		return err
	}
	_, err = s.r.Db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", username, string(hashedPassword), role)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) getUser(username string) (User, error) {
	var user User
	err := s.r.Db.QueryRow("SELECT username, password, role FROM users WHERE username = $1", username).Scan(&user.Username, &user.Password, &user.Role)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s Service) Login(username, password string) (string, error) {
	user, err := s.getUser(username)

	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("password doesn't match")
	}
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(SECRETSTRING))
	return ss, nil
}

func (s Service) CreateRegistrationCode(expiration time.Time) (RegistrationCode, error) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVQXYZ0123456789"
	code := make([]byte, 6)
	for i := 0; i < 6; i++ {
		code[i] = chars[rand.Intn(len(chars))]
	}

	var regCode RegistrationCode
	err := s.r.Db.QueryRow(
		"INSERT INTO registration_codes (code, expiration, used) VALUES ($1, $2, $3) RETURNING code, expiration, used",
		string(code), expiration, false).Scan(&regCode.Code, &regCode.Expiration, &regCode.Used)
	if err != nil {
		return RegistrationCode{}, err
	}

	return regCode, nil
}

func (s Service) GetRegistrationCode(code string) (RegistrationCode, error) {
	var regCode RegistrationCode
	err := s.r.Db.QueryRow("SELECT * FROM registration_codes WHERE code = $1 LIMIT 1", code).Scan(&regCode.Code, &regCode.Expiration, &regCode.Used)
	if err != nil {
		return RegistrationCode{}, err
	}

	return regCode, nil
}

func (s Service) UpdateRegistrationCode(code RegistrationCode) error {
	_, err := s.r.Db.Exec("UPDATE registration_codes set expiration = $2, used = $3 WHERE code = $1", code.Code, code.Expiration, code.Used)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) ValidateJWT(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETSTRING), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
