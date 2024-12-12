// Package web is a super minimalistic web/http server library
package web

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

// Headers is a map of string to string where the key is the key for the header
// And the value is the value for the header
type Headers map[string]string

// Response is a generic response object for our handlers
type Response struct {
	// StatusCode
	Status int
	// Content Type to writer
	ContentType string
	// Content to be written to the response writer
	Content io.Reader
	// Headers to be written to the response writer
	Headers Headers
}

// Write writes a response to an http.ResponseWriter
func (response *Response) Write(rw http.ResponseWriter) {
	if response != nil {
		if response.ContentType != "" {
			rw.Header().Set("Content-Type", response.ContentType)
		}
		for k, v := range response.Headers {
			rw.Header().Set(k, v)
		}
		rw.WriteHeader(response.Status)
		_, err := io.Copy(rw, response.Content)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

// Action represents a simplified http action
// that implements http.Handler
type Action func(r *http.Request) *Response

// Hyperlink represents a hyperlink
type Hyperlink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// Action's http.Handler implementation
func (a Action) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	response := a(r)
	response.Write(rw)
}

// Error returns an error response
func Error(status int, err error, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBufferString(err.Error()),
		Headers: headers,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

// ErrorJSON returns an error in json format
func ErrorJSON(status int, err error, headers Headers) *Response {
	errResp := errorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(errResp)
	if err != nil {
		return Error(http.StatusInternalServerError, err, headers)
	}

	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

// Data returns a data response
func Data(status int, content []byte, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBuffer(content),
		Headers: headers,
	}
}

// DataJSON returns a data response in json format
func DataJSON(status int, v interface{}, headers Headers) *Response {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ErrorJSON(http.StatusInternalServerError, err, headers)
	}

	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

// Empty returns an empty http response
func Empty(status int) *Response {
	return Data(status, []byte(""), nil)
}

// Sends an HTTP 302 redirect
func Redirect(location string) *Response {
	return &Response{
		Status:  http.StatusFound,
		Content: bytes.NewBuffer(nil),
		Headers: map[string]string{
			"Location": location,
		},
	}
}

// HTML renders an html template to a web response
func HTML(status int, t *template.Template, template string, data interface{}, headers Headers) *Response {

	//render template to buffer
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, template, data); err != nil {
		log.Println(err)
		return Empty(http.StatusInternalServerError)
	}
	return &Response{
		Status:      status,
		ContentType: "text/html",
		Content:     &buf,
		Headers:     headers,
	}
}

// TemplateParseFSRecursive recursively parses all templates in the FS with the given extension.
// File paths are used as template names to support duplicate file names.
// Use nonRootTemplateNames to exclude root directory from template names
// (e.g. index.html instead of templates/index.html)
func TemplateParseFSRecursive(
	templates fs.FS,
	ext string,
	nonRootTemplateNames bool,
	funcMap template.FuncMap) (*template.Template, error) {

	root := template.New("")
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}
			b, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}
			name := ""
			if nonRootTemplateNames {
				//name the template based on the file path (excluding the root)
				parts := strings.Split(path, string(os.PathSeparator))
				name = strings.Join(parts[1:], string(os.PathSeparator))
			}
			t := root.New(name).Funcs(funcMap)
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return root, err
}

// PathLast returns the last segment in the path
// and the number of path segments
func PathLast(r *http.Request) (string, int) {

	//remove trailing /
	path := r.URL.Path
	if path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}

	parts := strings.Split(path, "/")
	l := len(parts) - 1
	return parts[l], l
}
