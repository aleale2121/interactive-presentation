package intiator

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/aleale2121/interactive-presentation/internal/glue/routing"
	"github.com/aleale2121/interactive-presentation/internal/handler/rest"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"

	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/aleale2121/interactive-presentation/util"
)

func Init() {
	logger := logrus.New()
	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Fatalf("%s cannot load config", err.Error())
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Fatal("conn error: ", err)
	}
	defer conn.Close()

	m, err := migrate.New(
		config.MigrationURL,
		config.DBSource,
	)
	if err != nil {
		logger.Fatal("migration error: ", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Println("No migrations to apply.")
		} else {
			logger.Fatal(err)
		}
	}

	logger.Println("Migrations applied successfully.")
	store := db.NewStore(conn)

	presentationHandler := rest.NewPresentationHandler(logger, store)
	presentationRouting := routing.PresentationRouting(presentationHandler)

	pollHandler := rest.NewPollsHandler(logger, store)
	pollRouting := routing.PollRouting(pollHandler)

	voteHandler := rest.NewVoteHandler(logger, store)
	voteRouting := routing.VoteRouting(voteHandler)

	var routersList []routers.Router
	routersList = append(routersList, presentationRouting...)
	routersList = append(routersList, pollRouting...)
	routersList = append(routersList, voteRouting...)

	server := routers.NewRouting("localhost", "8080", routersList)
	server.Serve()
}
