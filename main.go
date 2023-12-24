package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/aleale2121/interactive-presentation/api"
	db "github.com/aleale2121/interactive-presentation/db/sqlc"
	"github.com/aleale2121/interactive-presentation/util"
)


func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("%s cannot load config", err.Error())
	}
	conn, err := sql.Open(config.DBDriver,  config.DBSource)
	if err != nil {
		log.Fatal("conn error: ",err)
	}
	defer conn.Close()


	m, err := migrate.New(
		config.DBSource,
		config.DBSource, 
	)
	if err != nil {
		log.Fatal("migration error: ", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply.")
		} else {
			log.Fatal(err)
		}
	}

	log.Println("Migrations applied successfully.")
	store := db.NewStore(conn)

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = server.Start(":8080")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}
