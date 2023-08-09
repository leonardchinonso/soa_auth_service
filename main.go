package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leonardchinonso/auth_service_cmp7174/datasource"
	"github.com/leonardchinonso/auth_service_cmp7174/injection"
)

func main() {
    log.Println("Starting Server...")

    // initialize data sources
    dataSource, err := datasource.InitDataSource()
    if err != nil {
        log.Fatalf("Failed to initialize data sources: %v", err)
    }

    // release resources when main function returns
    defer dataSource.Close()

    // initialize dependency injection
    router, err := injection.Inject(dataSource)
    if err != nil {
        log.Fatalf("Failed to inject data sources: %v", err)
    }

    srv := &http.Server{
        Addr: ":8080",
        Handler: router,
    }

    // Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
    // listening to the server in a goroutine so it does not block the graceful
    // shutdown after
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to initialize server. Error: %v\n", err)
        }
    }()

    log.Printf("Listening on port %v\n", srv.Addr)

    // wait for kill signal in channel
    quit := make(chan os.Signal)

    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // this blocks until a signal is passed into the quit channel
    <-quit

    // use context to communicate to server it has 5 seconds
    // to finish the current requests its handling
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // shutting down the server
    log.Println("Shutting down Server...")
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Failed to shutdown server. Error: %v\n", err)
    }

    log.Println("Server exiting")
}
