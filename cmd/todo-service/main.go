package main

import (
	"github.com/liampulles/go-config"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/wire"
)

func main() {
	wire.Run(config.NewEnvSource())
}
