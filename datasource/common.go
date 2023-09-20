package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"jazzine/m3u8recorder/database"
	"jazzine/m3u8recorder/model"
)

type Dao struct {
	db *sqlx.DB
}

func NewDao(db *sqlx.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (dao *Dao) FindAccountOrCreate(ctx context.Context, slug string) (*model.Account, error) {
	account, err := dao.FindAccountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("error while getting account: %w", err)
	}

	if account == nil {
		account = &model.Account{
			Slug: slug,
		}

		err := dao.SaveAccount(ctx, *account)
		if err != nil {
			return nil, fmt.Errorf("error while saving account: %w", err)
		}
	}

	// Refresh
	account, err = dao.FindAccountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("error while getting account: %w", err)
	}

	return account, nil
}

func (dao *Dao) SaveRecord(ctx context.Context, record model.Record) error {
	query := `INSERT INTO records (account_id, url, is_hidden, name, duration, thumbnail_url, res_width, res_height, source) VALUES (:account_id, :url, :is_hidden, :name, :duration, :thumbnail_url, :res_width, :res_height, :source)`
	if err := database.SaveOrUpdate(dao.db, ctx, query, record); err != nil {
		return fmt.Errorf("cannot save account: %w", err)
	}

	return nil
}

func (dao *Dao) SaveAccount(ctx context.Context, account model.Account) error {
	query := `INSERT INTO accounts (slug, thumbnail_url) VALUES (:slug, :thumbnail_url)`
	if err := database.SaveOrUpdate(dao.db, ctx, query, account); err != nil {
		return fmt.Errorf("cannot save account: %w", err)
	}

	return nil
}

func (dao *Dao) FindAccountBySlug(ctx context.Context, slug string) (*model.Account, error) {
	var entity model.Account
	query := `SELECT * FROM accounts WHERE slug=$1`
	if err := dao.db.GetContext(ctx, &entity, query, slug); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot get account: %w", err)
	}
	return &entity, nil
}
