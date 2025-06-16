package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/BeehiveInnovations/zen-mcp-server-go/internal/server"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetOutput(os.Stderr)

	srv := server.NewServer()

	if err := srv.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
		logrus.Fatalf("server error: %v", err)
	}
}
