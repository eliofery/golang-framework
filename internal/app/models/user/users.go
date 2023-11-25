package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
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

type service struct {
	ctx context.Context
	User
}

func New(ctx context.Context, user User) *service {
	return &service{
		ctx:  ctx,
		User: user,
	}
}

func (u *service) SignUp() error {
	op := "model.user.SignUp"

	d, v := database.CtxDatabase(u.ctx), validate.Validation(u.ctx)

	err := v.Struct(u.User)
	if err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := d.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, u.Email, string(hashedPassword),
	)
	err = row.Scan(&u.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return errors.Public(err, ErrEmailAlreadyExists.Error())
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	err = session.New(u.ctx, session.Session{UserID: u.ID}).Create()
	if err != nil {
		return err
	}

	return nil
}

func (u *service) SignIn() error {
	op := "model.user.SignIn"

	d, v := database.CtxDatabase(u.ctx), validate.Validation(u.ctx)

	err := v.Struct(u.User)
	if err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	password := u.Password

	row := d.QueryRow("SELECT * FROM users WHERE email = $1", u.Email)
	err = row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, ErrLoginOrPassword.Error())
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.Public(err, ErrLoginOrPassword.Error())
	}

	err = session.New(u.ctx, session.Session{UserID: u.ID}).Create()
	if err != nil {
		return err
	}

	return nil
}
