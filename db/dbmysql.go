package db

import (

	// This line is must for working MySQL database
	_ "github.com/go-sql-driver/mysql"
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
func (m *MySQL) CreatePreroll(s model.Preroll) (int64, error) {
	res, err := m.conn.Exec(
		"INSERT INTO `preroll` (`name`, `date`, `showcnt_kg`, `showcnt_db`, `clickcnt_kg`, `clickcnt_db`) VALUES (?, ?, ?, ?, ?, ?)", s.Name, s.Date, s.ShowcntKg, s.ShowcntDb, s.ClickcntKg, s.ClickcntDb,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SelectPreroll selects preroll entry from database
func (m *MySQL) SelectPreroll(p model.Preroll) (model.Preroll, error) {
	var s model.Preroll
	err := m.conn.Get(&s, "SELECT `id`, `showcnt_kg`, `showcnt_db`, `clickcnt_kg`, `clickcnt_db` FROM `preroll` WHERE `name`=? AND `date`=?", p.Name, p.Date)
	return s, err
}

// ListPrerolls returns array of prerolls entries from database
func (m *MySQL) ListPrerolls() ([]model.Preroll, error) {
	prerolls := []model.Preroll{}
	err := m.conn.Select(&prerolls, "SELECT * FROM `preroll`")
	return prerolls, err
}

// UpdatePreroll updates preroll entry in database
func (m *MySQL) UpdatePreroll(s model.Preroll) error {
	tx := m.conn.MustBegin()
	tx.MustExec(
		"UPDATE `preroll` SET `name` = ?, `date` = ?, `showcnt_kg` = ?, `showcnt_db` = ?, `clickcnt_kg` = ?, `clickcnt_db` = ? WHERE `id` = ?",
		s.Name, s.Date, s.ShowcntKg, s.ShowcntDb, s.ClickcntKg, s.ClickcntDb, s.ID,
	)
	err := tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

// DeletePreroll deletes preroll entry from database
func (m *MySQL) DeletePreroll(id int64) error {
	_, err := m.conn.Exec("DELETE FROM `preroll` WHERE id=?", id)
	return err
}
