package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

var (
	ErrEmailAlreadyExists = errors.New("email адрес уже существует")
	ErrLoginOrPassword    = errors.New("неверный логин или пароль")
)

type User struct {
	ID       uint   `validate:"omitempty"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=10,lte=32"`
}

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) SignUp(user *User) error {
	op := "model.user.SignUp"

	d, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	err := v.Struct(user)
	if err != nil {
		return err
	}

	user.Email = strings.ToLower(user.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := d.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email, string(hashedPassword),
	)
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return errors.Public(err, ErrEmailAlreadyExists.Error())
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	err = session.NewService(s.ctx).Create(&session.Session{UserID: user.ID})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(user *User) error {
	op := "model.user.SignIn"

	d, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	err := v.Struct(user)
	if err != nil {
		return err
	}

	user.Email = strings.ToLower(user.Email)
	password := user.Password

	row := d.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, ErrLoginOrPassword.Error())
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.Public(err, ErrLoginOrPassword.Error())
	}

	err = session.NewService(s.ctx).Create(&session.Session{UserID: user.ID})
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentUser(r *http.Request) (*User, error) {
	op := "model.user.CurrentUser"

	userData := &User{}

	token, err := cookie.Get(r, cookie.Session)
	if err != nil {
		return userData, fmt.Errorf("%s: %w", op, err)
	}
	tokenHash := rand.HashToken(token)

	d := database.CtxDatabase(r.Context())
	row := d.QueryRow(`
       SELECT users.id, users.email, users.password
       FROM users
       INNER JOIN sessions ON users.id = sessions.user_id
       WHERE sessions.token_hash = $1;
   `, tokenHash)
	err = row.Scan(&userData.ID, &userData.Email, &userData.Password)
	if err != nil {
		return userData, fmt.Errorf("%s: %w", op, err)
	}

	return userData, nil
}
