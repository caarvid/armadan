package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/google/uuid"
)

type resetPassword struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	pool     *sql.DB
}

func NewResetPasswordService(reader, writer schema.Querier, pool *sql.DB) *resetPassword {
	return &resetPassword{
		dbReader: reader,
		dbWriter: writer,
		pool:     pool,
	}
}

func (rp *resetPassword) Get(ctx context.Context, token string) (*armadan.ResetPasswordToken, error) {
	resetToken, err := rp.dbReader.GetResetPasswordToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &armadan.ResetPasswordToken{
		ID:        resetToken.ID,
		UserId:    resetToken.UserID,
		ExpiresAt: armadan.ParseTime(resetToken.ExpiresAt),
	}, nil
}

func (rp *resetPassword) Create(ctx context.Context, userId string) (*armadan.ResetPasswordToken, error) {
	resetToken, err := rp.dbWriter.CreateResetPasswordToken(ctx, &schema.CreateResetPasswordTokenParams{
		UserID:    userId,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(15 * time.Minute).Format(armadan.DEFAULT_TIME_FORMAT),
	})
	if err != nil {
		return nil, err
	}

	return &armadan.ResetPasswordToken{
		ID:        resetToken.ID,
		UserId:    resetToken.UserID,
		Token:     resetToken.Token,
		ExpiresAt: armadan.ParseTime(resetToken.ExpiresAt),
	}, nil
}

func (rp *resetPassword) UpdateUserPassword(ctx context.Context, token *armadan.ResetPasswordToken, password string) error {
	hash, err := utils.GenerateHash(password, nil)
	if err != nil {
		return err
	}

	tx, err := rp.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	_, err = qtx.UpdateUserPassword(ctx, &schema.UpdateUserPasswordParams{
		ID:       token.UserId,
		Password: hash.Encode(),
	})

	if err != nil {
		return err
	}

	if err = qtx.DeleteResetPasswordToken(ctx, token.Token); err != nil {
		return err
	}

	return tx.Commit()
}
