package app

import (
	"context"
	"log"
	"ordermngmt/internal/config"
	"ordermngmt/internal/handler/http/api"
	"ordermngmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg    config.Config
	c      *Container
	logger *log.Logger
}

func NewApp(cfgPath string) (*App, error) {
	cfg, err := config.NewConfig(cfgPath)

	if err != nil {
		log.Printf("config reading error: %v\n", err)
		return nil, err
	}

	pgsqlConn, err := newPostgresqlConnection(cfg.Pgsql)

	if err != nil {
		log.Printf("pgsql connection error: %v\n", err)
		return nil, err
	}
	log.Println("pgsql connection established")

	if err := initDB(pgsqlConn); err != nil {
		log.Printf("initDB error: %v\n", err)
		return nil, err
	}
	log.Println("initDB executed")

	app := &App{
		cfg: cfg,
		c:   NewContainer(cfg, pgsqlConn),
	}

	return app, nil
}

func (app *App) Run(errChan chan error) {
	useCase := app.c.GetUseCase()

	go func() {
		errChan <- app.startHTTPServer(useCase)
		log.Println("http server returned error")
	}()

	go func() {
		errChan <- app.runStan(useCase)
		log.Println("stan server returned error")
	}()
}

func (app *App) runStan(uc *usecase.UseCase) error {
	for data := range app.c.msgStream {
		err := uc.AddOrder(context.Background(), data)
		if err != nil {
			log.Println(err)
		}
	}
	app.c.subscription.Close()
	return nil
}

func (app *App) startHTTPServer(uc *usecase.UseCase) (err error) {
	gengine := gin.Default()

	handler := api.NewHandler(uc)

	handler.AddRoutes(gengine)

	err = gengine.Run(":8080")
	app.logger.Println(err)

	return
}
