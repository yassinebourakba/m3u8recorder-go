package model

import "time"

type Record struct {
	ID        uint64 `db:"id"`
	AccountID uint64 `db:"account_id"`
	URL       string `db:"url"`
	IsHidden  bool   `db:"is_hidden"`

	Name               *string   `db:"name"`
	DurationInSecondes int       `db:"duration"`
	ThumbnailURL       *string   `db:"thumbnail_url"`
	Views              int       `db:"views"`
	Likes              int       `db:"likes"`
	PublishedAt        time.Time `db:"published_at"`
	Width              *int      `db:"res_width"`
	Height             *int      `db:"res_height"`

	Source *string `db:"source"`
}

type Account struct {
	ID   uint64 `db:"id"`
	Slug string `db:"slug"`

	ThumbnailURL *string   `db:"thumbnail_url"`
	CreatedAt    time.Time `db:"created_at"`
}
