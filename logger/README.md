# logger

### The `logger` package implements the core logging functionality.

By default, the `logger` is initialized:

- with the time format `2006/01/02 15:04:05.000`;
- with logging format `FormatText`;
- at the highest logging level `LevelTrace`;
- printing the name of the calling function `true`.

You can set the required time format, level, log message format and additional output, as in the following example:

```go
package main

import (
	"os"

	"github.com/easy-techno-lab/proton/logger"
)

func main() {
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer logger.Closer(file)

	logger.Info("Hello World!")
	// 2006/01/02 15:04:05.000 INFO [1] /../../main.go:17 main() Hello World!

	logger.SetTimeFormat("2006-01-02 15:04:05")
	logger.SetFormat(logger.FormatJSON)
	logger.SetAdditionalOut(file)
	logger.SetLevel(logger.LevelInfo)
	logger.SetFuncNamePrinting(false)

	logger.Info("Hello World!")
	// {"time":"2006-01-02 15:04:05","level":"INFO","routine":1,"message":"Hello World!"}
}

```
