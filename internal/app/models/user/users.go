package user

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var (
	ErrEmailAlreadyExists = errors.New("email адрес уже существует")
)

type Dto struct {
	ID       uint
	Email    string
	Password string
}

type User struct {
	db *sql.DB
	Dto
}

func New(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) SignUp() error {
	op := "model.user.SignUp"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	u.Email = strings.ToLower(u.Email)
	u.Password = string(hashedPassword)

	row := u.db.QueryRow(
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

	return nil
}
