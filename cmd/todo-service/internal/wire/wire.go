package wire

import (
	"fmt"

	goConfig "github.com/liampulles/go-config"
)

// Run is the main entrypoint for todo-service
func Run(source goConfig.Source) int {
	fmt.Println("Hello world!")
	return 0
}
