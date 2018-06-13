package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/medalon/stats_from_nginx/config"
	"github.com/medalon/stats_from_nginx/stats"
)

var format = "main"
var logFile string
var mesto string
var geo int

func init() {
	flag.StringVar(&logFile, "log", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&mesto, "mesto", "kg", "Default value")
}

func main() {
	flag.Parse()

	// Read config from system environment
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err, "could not get env conf parms")
	}

	// Read given file or from STDIN
	var logReader io.Reader

	if logFile == "dummy" {
		logReader = strings.NewReader(`1.2.3.4 - - [18/Feb/2018:16:54:10 +0600] "GET /1px.png?stat=preroll&name=PACMAN&act=click&pid=bQFrkd9E HTTP/1.1" 200 923 "https://localhost/embed/video/183223?mesto=1&ismobile=1" "Mozilla/5.0 (Linux; Android 6.0.1; SM-G532F Build/MMB29T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Mobile Safari/537.36" "-"
		1.2.3.4 - - [20/Feb/2018:04:22:24 +0600] "GET /1px.png?stat=preroll&name=PACMAN&act=show&pid=SeD3e3Ra HTTP/1.1" 200 923 "https://localhost/embed/video/210837?mesto=0&ismobile=1" "Mozilla/5.0 (Linux; U; Android 6.0.1; en-US; ASUS_Z010DD Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/40.0.2214.89 UCBrowser/11.5.0.1015 Mobile Safari/537.36" "-"`)
	} else if logFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		logReader = file
		defer file.Close()
	}

	// Use nginx config file to extract format by the name
	var nginxConfig io.Reader
	nginxConfig = strings.NewReader(`
            http {
                log_format   main  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"';
            }
		`)

	testRes, err := stats.ParseLogFile(logReader, nginxConfig, format)
	if err != nil {
		fmt.Println(err)
	}

	s, err := stats.NewServerDB(c)
	if err != nil {
		fmt.Println(err)
	}

	if mesto == "kg" {
		geo = 1
	} else {
		geo = 0
	}
	for k, v := range testRes {
		for i, j := range v {
			//fmt.Printf("%s => %s = %v, %v\n", k, i, j.Showcnt, j.Clickcnt)
			_ = s.WriteToDb(k, i, j.Showcnt, j.Clickcnt, geo)

		}
	}
	fmt.Println("All Done!")
}
