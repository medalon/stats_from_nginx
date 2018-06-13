build:
	go build -o logNginxParse cmd/main.go

clean:
	rm logNginxParse

mysqlenv:
	export DATABASE_URL='user:passw@tcp(localhost:3306)/dbname?multiStatements=true&parseTime=True'