package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/ckkb_api/internal/adapter/postgresql"
	"github.com/ckkb_api/internal/controller/handlers"
	"github.com/ckkb_api/internal/controller/handlers/issue"
	"github.com/ckkb_api/internal/controller/router"
	"github.com/ckkb_api/internal/controller/server"
	issueAction "github.com/ckkb_api/internal/domain/issue"
	"github.com/sirupsen/logrus"
)

func main() {

	host := ":8005"
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	postgresqlConn, err := postgresql.NewDB()
	if err != nil {
		logrus.Fatalf("failed to connect to postgreSQL: %v", err)
	}
	postgreSQLStorage := postgresql.NewPostgreSQLStorage(postgresqlConn)

	var registeredHandlers []handlers.Handler

	issueAction := issueAction.NewIssueAction(postgreSQLStorage)
	issueHandler := issue.NewIssueHandler(issueAction)
	registeredHandlers = append(registeredHandlers, issueHandler)

	appRouter := router.NewEchoRouter(registeredHandlers)
	HTTPserver := server.NewHTTPServer(host, 30, appRouter.Echo)

	log.Printf("starting download http server: %s", host)
	go HTTPserver.Start(ctx)

	<-ctx.Done()

	HTTPserver.Stop(ctx)
}
