package stats

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/medalon/stats_from_nginx/model"
	"github.com/satyrius/gonx"
)

// ParseLogFile ...
func ParseLogFile(logReader, nginxConfig io.Reader, format string) model.ResMap {

	resmap := make(model.ResMap)

	reader, err := gonx.NewNginxReader(logReader, nginxConfig, format)
	if err != nil {
		panic(err)
	}
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// Process the record... e.g.
		tlocal, err := rec.Field("time_local")
		if err != nil {
			fmt.Println(err)
		}
		trequest, err := rec.Field("request")
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(trequest, "preroll") {
			treq := strings.Split(trequest, " ")
			treq = strings.Split(treq[1], "?")
			treqmap, err := url.ParseQuery(treq[1])
			if err != nil {
				fmt.Println(err)
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

	return resmap
}
