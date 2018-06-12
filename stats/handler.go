package stats

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/cactus/gostrftime"
	"github.com/medalon/stats_from_nginx/config"
	"github.com/medalon/stats_from_nginx/db"
	"github.com/medalon/stats_from_nginx/model"
	"github.com/satyrius/gonx"
)

// ServerDB ...
type ServerDB struct {
	db db.DB
}

// NewServerDB ...
func NewServerDB(c *config.StatsNginxConfig) (*ServerDB, error) {
	s := &ServerDB{}
	conn, err := db.NewMySQL(c)
	if err != nil {
		return nil, err
	}
	s.db = conn

	return s, nil
}

// ParseLogFile ...
func ParseLogFile(logReader, nginxConfig io.Reader, format string) (model.ResMap, error) {

	resmap := make(model.ResMap)

	reader, err := gonx.NewNginxReader(logReader, nginxConfig, format)
	if err != nil {
		return nil, err
	}
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// Process the record... e.g.
		tlocal, err := rec.Field("time_local")
		if err != nil {
			return nil, err
		}
		trequest, err := rec.Field("request")
		if err != nil {
			return nil, err
		}
		if strings.Contains(trequest, "preroll") {
			treq := strings.Split(trequest, " ")
			treq = strings.Split(treq[1], "?")
			treqmap, err := url.ParseQuery(treq[1])
			if err != nil {
				return nil, err
			}

			tloc := strings.Split(tlocal, ":")
			treqmap["dtime"] = tloc

			mkey := treqmap.Get("name")
			mseckey := treqmap.Get("dtime")
			if _, ok := resmap[mkey]; !ok {
				resmap[mkey] = make(map[string]*model.ResData)
			}
			if _, ok := resmap[mkey][mseckey]; !ok {
				resmap[mkey][mseckey] = &model.ResData{
					Showcnt: 0, Clickcnt: 0,
				}
			}

			switch treqmap.Get("act") {
			case "show":
				resmap[mkey][mseckey].Showcnt++
			case "click":
				resmap[mkey][mseckey].Clickcnt++
			}

		}
	}

	return resmap, nil
}

// WriteToDb ...
func (s *ServerDB) WriteToDb(name, dtime string, shcnt, clcnt int) {
	t1, _ := time.Parse("02/Jan/2006", dtime)
	dtime = gostrftime.Format("%Y-%m-%d", t1)

	var u model.Preroll
	u.Name = name
	u.Date = dtime
	u.ShowcntDb = 0
	u.ClickcntDb = 0

	stmt, err := s.db.SelectPreroll(u)
	switch {
	case err == sql.ErrNoRows:
		u.ShowcntKg = int64(shcnt)
		u.ClickcntKg = int64(clcnt)
		_, err := s.db.CreatePreroll(u)
		if err != nil {
			fmt.Println(err)
		}
	case err != nil:
		fmt.Println(err)

	case stmt.ID > 0:
		u.ShowcntKg = stmt.ShowcntKg + int64(shcnt)
		u.ClickcntKg = stmt.ClickcntKg + int64(clcnt)
		u.ID = stmt.ID
		err := s.db.UpdatePreroll(u)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("updated")
	}
}
