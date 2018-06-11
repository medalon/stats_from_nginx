package db

import (

	// This line is must for working MySQL database

	"github.com/jmoiron/sqlx"

	"github.com/medalon/stats_from_nginx/config"
	"github.com/medalon/stats_from_nginx/model"
)

// MySQL provides api for work with mysql database
type MySQL struct {
	conn *sqlx.DB
}

// NewMySQL creates a new instance of database API
func NewMySQL(c *config.StatsNginxConfig) (*MySQL, error) {
	conn, err := sqlx.Open("mysql", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	m := &MySQL{}
	m.conn = conn
	return m, nil
}

// CreatePreroll creates preroll entry in database
func (m *MySQL) CreatePreroll(s model.Preroll) (model.Preroll, error) {
	res, err := m.conn.Exec(
		"INSERT INTO `preroll` (`name`, `date`, `showcnt_kg`, `showcnt_db`, `clickcnt_kg`, `clickcnt_db`) VALUES (?, ?, ?)", s.Name, s.Date, s.ShowcntKg, s.ShowcntDb, s.ClickcntKg, s.ClickcntDb,
	)
	if err != nil {
		return s, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return s, err
	}
	s.ID = id
	return s, nil
}

// SelectPreroll selects preroll entry from database
func (m *MySQL) SelectPreroll(id int64) (model.Preroll, error) {
	var s model.Preroll
	err := m.conn.Get(&s, "SELECT * FROM `preroll` WHERE id=?", id)
	return s, err
}

// ListPrerolls returns array of prerolls entries from database
func (m *MySQL) ListPrerolls() ([]model.Preroll, error) {
	prerolls := []model.Preroll{}
	err := m.conn.Select(&prerolls, "SELECT * FROM `preroll`")
	return prerolls, err
}

// UpdatePreroll updates preroll entry in database
func (m *MySQL) UpdatePreroll(s model.Preroll) (model.Preroll, error) {
	tx := m.conn.MustBegin()
	tx.MustExec(
		"UPDATE `preroll` SET `name` = ?, `date` = ?, `showcnt_kg` = ?, `showcnt_db` = ?, `clickcnt_kg` = ?, `clickcnt_db` = ? WHERE `id` = ?",
		s.Name, s.Date, s.ShowcntKg, s.ShowcntDb, s.ClickcntKg, s.ClickcntDb, s.ID,
	)
	err := tx.Commit()

	if err != nil {
		return s, err
	}
	var i model.Preroll
	err = m.conn.Get(&i, "SELECT * FROM `preroll` WHERE id=?", s.ID)
	return i, err
}

// DeletePreroll deletes preroll entry from database
func (m *MySQL) DeletePreroll(id int64) error {
	_, err := m.conn.Exec("DELETE FROM `preroll` WHERE id=?", id)
	return err
}
