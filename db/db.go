package db

import "github.com/medalon/stats_from_nginx/model"

// DB ...
type DB interface {
	CreatePreroll(s model.Preroll) (int64, error)
	SelectPreroll(p model.Preroll) (model.Preroll, error)
	ListPrerolls() ([]model.Preroll, error)
	UpdatePreroll(s model.Preroll) error
	DeletePreroll(id int64) error
}
