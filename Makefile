build:
	go build -o stats_from_nginx cmd/main.go

clean:
	rm stats_from_nginx

mysqlenv:
	export DATABASE_URL='user:passw@tcp(localhost:3306)/dbname?multiStatements=true&parseTime=True'