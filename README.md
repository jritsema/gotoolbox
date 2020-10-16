# gotoolbox

A kitchen sink of Go tools that I've found useful. Uses only the standard library, no external dependencies.

### contents

- [super lightweight http server framework](web)
- [exponential backoff retry](retry.go)
- [working with JSON](json.go)
- [working with the file system](fs.go)
- [working with slices](slice.go)
- [working with CLIs](cli.go)

### example usage

```
go get github.com/jritsema/gotoolbox
```

### utils

```go
package main

import utils "github.com/jritsema/gotoolbox"

func main() {

	s := []string{"a", "b", "c"}
	if utils.SliceContains(&s, "b") {
		fmt.Println("b exists")
	}

	err := utils.Retry(3, 1, func() error {
		return callBrittleAPI()
	})
	if err != nil {
		fmt.Println("callBrittleAPI failed after 3 retries: %w", err)
	}

	config, err := utils.ReadJSONFile("config.json")
	if err != nil {
		fmt.Println("error reading json file: %w", err)
	}
}
```

#### web framework

```go
package main

import "github.com/jritsema/gotoolbox/web"

type Data struct {
	Hello string `json:"hello"`
}

func hello(r *http.Request) *web.Response {
	return web.DataJSON(http.StatusOK, Data{Hello: "world"}, nil)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", web.Action(hello))
	http.ListenAndServe(":8080", mux)
}
```

### development

```

Choose a make command to run

vet vet code
test run unit tests
build build a binary
autobuild auto build when source files change
start build and run local project

```
