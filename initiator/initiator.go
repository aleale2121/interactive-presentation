package intiator

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/aleale2121/interactive-presentation/internal/glue/routing"
	presentationHttpV1 "github.com/aleale2121/interactive-presentation/internal/handler/presentation/http/v1"
	voteHttpV1 "github.com/aleale2121/interactive-presentation/internal/handler/vote/http/v1"

	"github.com/aleale2121/interactive-presentation/internal/module/presentation"
	"github.com/aleale2121/interactive-presentation/internal/module/vote"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"

	"github.com/aleale2121/interactive-presentation/pkg/config"
	"github.com/aleale2121/interactive-presentation/platform/routers"
)

func Init() {
	logger := logrus.New()
	config, err := config.LoadConfig(".")
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

	presentationUseCase := presentation.Initialize(store)
	presentationHandler := presentationHttpV1.NewPresentationHandler(logger, presentationUseCase)
	presentationRouting := routing.PresentationRouting(presentationHandler)



	voteUseCase := vote.Initialize(store)
	voteHandler := voteHttpV1.NewVoteHandler(logger, voteUseCase)
	voteRouting := routing.VoteRouting(voteHandler)

	var routersList []routers.Router
	routersList = append(routersList, presentationRouting...)
	routersList = append(routersList, voteRouting...)

	server := routers.NewRouting(config.ServerAddress, routersList)
	server.Serve()
}
