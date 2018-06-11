package db

import "github.com/medalon/stats_from_nginx/model"

// DB ...
type DB interface {
	CreatePreroll(s model.Preroll) (model.Preroll, error)
	SelectPreroll(id int64) (model.Preroll, error)
	ListPrerolls() ([]model.Preroll, error)
	UpdatePreroll(s model.Preroll) (model.Preroll, error)
	DeletePreroll(id int64) error
}
