package main

import (
	"os"

	"github.com/intelsdi-x/pulse-plugin-processor-movingaverage/movingaverage"
	"github.com/intelsdi-x/pulse/control/plugin"
)

func main() {
	meta := movingaverage.Meta()
	plugin.Start(meta, movingaverage.NewMovingaverageProcessor(), os.Args[1])
}
