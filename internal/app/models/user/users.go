package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
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

type Dto struct {
	ID       uint   `validate:"omitempty"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=10,lte=32"`
}

type User struct {
	ctx context.Context
	Dto
}

func New(ctx context.Context) *User {
	r := router.Request(ctx)

	return &User{
		ctx: ctx,
		Dto: Dto{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		},
	}
}

func (u *User) Create() error {
	op := "model.user.Create"

	db := database.CtxDatabase(u.ctx)
	valid := validate.Validation(u.ctx)

	err := valid.Struct(u.Dto)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	u.Email = strings.ToLower(u.Email)
	u.Password = string(hashedPassword)

	row := db.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, u.Email, u.Password,
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

	modelSession := session.New(u.ctx)
	err = modelSession.Create(u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Auth() error {
	op := "model.user.Auth"

	db := database.CtxDatabase(u.ctx)
	valid := validate.Validation(u.ctx)

	err := valid.Struct(u.Dto)
	if err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	password := u.Password

	row := db.QueryRow("SELECT * FROM users WHERE email = $1", u.Email)
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

	modelSession := session.New(u.ctx)
	err = modelSession.Create(u.ID)
	if err != nil {
		return err
	}

	return nil
}
