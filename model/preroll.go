package model

// Preroll ...
type Preroll struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Date       string `db:"date"`
	ShowcntKg  int64  `db:"showcnt_kg"`
	ShowcntDb  int64  `db:"showcnt_db"`
	ClickcntKg int64  `db:"clickcnt_kg"`
	ClickcntDb int64  `db:"clickcnt_db"`
}
