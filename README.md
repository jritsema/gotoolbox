# gotoolbox

A kitchen sink of Go tools that I've found useful. Uses only the standard library, no external dependencies.

### contents

- [super lightweight http server library](web)
- [exponential backoff retry](retry.go)
- [working with JSON](json.go)
- [working with the file system](fs.go)
- [working with slices](slice.go)
- [working with CLIs](cli.go)

### example usage

```
go get github.com/jritsema/gotoolbox
```

### utilities

```go
package main

import "github.com/jritsema/gotoolbox"

func main() {

	s := []string{"a", "b", "c"}
	if gotoolbox.SliceContainsLike(&s, "b") {
		fmt.Println("b exists")
	}

	err := gotoolbox.Retry(3, 1, func() error {
		return callBrittleAPI()
	})
	if err != nil {
		fmt.Println("callBrittleAPI failed after 3 retries: %w", err)
	}

	f := "config.json"
	if !gotoolbox.IsDirectory(f) && gotoolbox.FileExists(f) {
		config, err := gotoolbox.ReadJSONFile(f)
		if err != nil {
			fmt.Println("error reading json file: %w", err)
		}
	}

	value := gotoolbox.GetEnvWithDefault("MY_ENVVAR", "true")

	command := exec.Command("docker", "build", "-t", "foo", ".")
	err = gotoolbox.ExecCmd(command, true)
	if err != nil {
		fmt.Println("error executing command: %w", err)
	}

	var data interface{}
	err = gotoolbox.HttpGetJSON("https://api.example.com/data.json", &data)
	err = gotoolbox.HttpPostJSON("https://api.example.com/data.json", data, http.StatusOK)
}
```

#### web package

```go
package main

import (
	"embed"
	"html/template"
	"net/http"
	"github.com/jritsema/gotoolbox/web"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS
	html *template.Template
)

type Data struct {
	Hello string `json:"hello"`
}

func index(r *http.Request) *web.Response {
	return HTML(http.StatusOK, html, "index.html", Data{Hello: "world"}, nil)
}

func api(r *http.Request) *web.Response {
	return web.DataJSON(http.StatusOK, Data{Hello: "world"}, nil)
}

func main() {
	html, _ = web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	mux := http.NewServeMux()
	mux.Handle("/api", web.Action(api))
	mux.Handle("/", web.Action(index))
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
