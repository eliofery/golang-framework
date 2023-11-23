package session

import (
	"context"
	"fmt"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/validate"
)

type Dto struct {
	ID        uint   `validate:"omitempty"`
	UserID    uint   `validate:"required,min=1"`
	TokenHash string `validate:"required"`
}

type Session struct {
	ctx context.Context
	Dto
}

func New(ctx context.Context) *Session {
	return &Session{
		ctx: ctx,
	}
}

func (s *Session) Create(userID uint) error {
	op := "model.session.Create"

	w := router.ResponseWriter(s.ctx)
	db := database.CtxDatabase(s.ctx)
	valid := validate.Validation(s.ctx)

	token, err := rand.SessionToken()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.UserID = userID
	s.TokenHash = rand.HashToken(token)

	err = valid.Struct(s.Dto)
	if err != nil {
		return err
	}

	row := db.QueryRow(`
        INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)
        ON CONFLICT (user_id) DO
        UPDATE SET token_hash = $2
        RETURNING id;
    `, s.UserID, s.TokenHash)
	err = row.Scan(&s.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cookie.Set(w, cookie.Session, token)

	return nil
}
