package auth

import (
	"context"
	"errors"
	"fmt"

	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/macemacemace/seraphim/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrWeakPassword       = errors.New("password must be 8 to 72 characters")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type userStore interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

func NewService(store userStore, secret []byte) *Service {
	return &Service{store: store, secret: secret}
}

func (s *Service) Register(ctx context.Context, email, password, name string) (db.User, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if _, err := mail.ParseAddress(email); err != nil {
		return db.User{}, ErrInvalidEmail
	}
	if len(password) < 8 || len(password) > 72 {
		return db.User{}, ErrWeakPassword

	}
	name = strings.TrimSpace(name)

	if name == "" {
		return db.User{}, ErrNameRequired
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)

	if err != nil {
		return db.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
	})
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.User{}, ErrEmailTaken
		}
		return db.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil

}

func (s *Service) Login(ctx context.Context, email, password string) (db.User, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	user, err := s.store.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrInvalidCredentials
		}

		return db.User{}, fmt.Errorf("get user: %w", err)

	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return db.User{}, ErrInvalidCredentials
	}

	return user, nil
}

type Service struct {
	store  userStore
	secret []byte
}
