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

type Session struct {
	ID        uint   `validate:"omitempty"`
	UserID    uint   `validate:"required,min=1"`
	TokenHash string `validate:"required"`
}

type service struct {
	ctx context.Context
	Session
}

func New(ctx context.Context, session Session) *service {
	return &service{
		ctx:     ctx,
		Session: session,
	}
}

func (s *service) Create() error {
	op := "model.session.SignUp"

	w, d, v := router.ResponseWriter(s.ctx), database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	token, err := rand.SessionToken()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.TokenHash = rand.HashToken(token)

	err = v.Struct(s.Session)
	if err != nil {
		return err
	}

	row := d.QueryRow(`
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
