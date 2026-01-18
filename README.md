# framework

```
 __  __ ___ ___ ___  ___   ___ ___  ___ ___ 
|  \/  |_ _/ __| _ \/ _ \ / __/ _ \| _ \ __|
| |\/| || | (__|   / (_) | (_| (_) |   / _| 
|_|  |_|___\___|_|_\\___/ \___\___/|_|_\___|
```

Microcore is a high-performance, modular framework for building microservices in Go.

[![Go Reference](https://pkg.go.dev/badge/go.microcore.dev/framework.svg)](https://pkg.go.dev/go.microcore.dev/framework)

## Install
```bash
go get go.microcore.dev/framework
```

## Usage
```go
package main

import (
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/transport/http/server"
)

func main() {
	// Up http server
	go server.New().Up()

	// Graceful shutdown
	shutdown.Wait()
}
```

## Docs
[https://microcore.dev/framework](https://microcore.dev/framework)


## License

This project is licensed under the terms of the [MIT License](LICENSE).

Copyright © 2026 Microcore