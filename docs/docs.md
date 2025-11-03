# Igo: Go Micro Web-Framework

Igo is a simple and lightweight micro web-framework for Go, inspired by Bottle.py.
It is distributed as a single Go source file and depends only on the Go standard library.

- Routing - requests map to handler functions, with support for dynamic path parameters.

- Context & Utilities - form data, cookies, headers, file uploads, query parameters, etc.

- Responses - built-in HTML templating, JSON, plain text.

- Static files - easy static file serving.

- Middleware - logging, recovery, etc.

- Error handling - graceful failures and error logging.

# Table of Contents

[Getting Started](#getting-started)

[Context & Request Handling](#context--request-handling)

[Routing](#routing)

[Response Helpers](#response-helpers)

[Templates](#templates)

[Static Files](#static-files)

[Middleware](#middleware)

[Error Handling & Recovery](#error-handling--recovery)

[Utilities & Helpers](#utilities--helpers)

[Examples](#examples)

[Roadmap](#roadmap)

# Getting Started

### Installation

Igo does not depend on any external libraries. You can just <a href="https://github.com/mukailasam/snip">Snip</a> **Igo.go** into your project directory and start coding:

```sh
$ snip repo github.com/mukailasam/igo --file igo.go
```

OR

Using go get command

```go
$ go get github.com/mukailasam/igo@latest
```

OR

With <a href="https://go.dev/wiki/Modules#how-to-use-modules">Go's module support</a>, simply import igo in your code and Go will automatically fetch it during build:

```go
import "github.com/mukailasam/igo"
```

### Basic Example

```go
package main

import "github.com/mukailasam/igo"

type Post struct {
	Title string
	Body  string
}

func main() {
	app := igo.Default()

	app.GET("/hello", func(ctx *igo.Context) {
		ctx.JSON(200, Post{Title: "Title 1", Body: "Body 1"})
	})

	app.Run(":1337")
}

```

Running the application:

1. Save the code above as main.go
2. Run the application:
   ```go
   go run main.go
   ```
3. Open your browser and visit http://localhost:1337/hello

You should see:

```json
{
  "Title": "Title 1",
  "Body": "Body 1"
}
```

## Context & Request Handling

The Context object is passed to every handler and provides:

- Access to ResponseWriter and Request
- URL parameters via GetParam(key)
- Query parameters via GetQuery(key)
- Form and file uploads via FormFile and GetFormValue
- Cookies via SetCookie / GetCookie
- Header manipulation: SetHeader, GetHeader
- Error capturing: ctx.Error(err)
- Aborting a request early: ctx.Abort(code, message)

Example:

```go
func handler(ctx *igo.Context) {
    user := ctx.GetParam("user")
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.Abort(401, "Missing token")
        return
    }

    ctx.Text(200, "Hello "+user)
}
```

## Routing

Routes are registered with HTTP methods and path patterns.

```go
app.GET("/users/<id>", handler)
app.POST("/upload", uploadHandler)
app.GET("/assets/*path", staticHandler)
```

## Dynamic Segments

"\<name>" - matches a single path segment.

"/\*" - path is a wildcard matching the rest of the path (including "/").

Under the hood, routes are converted to regular expressions and parameter names are extracted.

## Response Helpers

- Text(status, string)
  Returns plain text with Content-Type: text/plain.

- JSON(status, data)
  Writes JSON response (you pass a struct, map, slice).
  Sets Content-Type: application/json.

- HTML(status, templateName, data)
  Renders HTML templates (using Go's html/template) by name or file path.
  Supports templates in a views/ folder or custom path.

## Templates

Templates may use {{.Field}} notation, loops, conditionals, etc.

You can name a template file (e.g. index) and igo will look for views/index.html or files matching pattern views/\*.html.

You can also provide a full file path, e.g. "templates/home.html".

Errors in parsing or execution are captured (via ctx.Error) and result in a 500 Internal Server Error.

## Static Files

Use ServeStatic(urlPrefix, directory):

```go
app.ServeStatic("/static/", "./public")
```

When a request URL starts with /static/, Igo will attempt to serve the corresponding file under ./public.
If file not found or directory, it returns 404 or 403 as appropriate.

Static serving also handles correct MIME types and caching headers (ETag / Last-Modified) when implemented.

## Middleware

Middleware functions allow interception before or after handlers:

```go
func Auth(c *igo.Context, next igo.HandlerFunc) {
	if c.GetHeader("X-Token") == "" {
		c.Abort(403, "Forbidden")
		return
	}
	next(c)
}

app.Use(Auth)
```

## Built-in middleware:

- Logger - logs requests, duration, and captured errors

- Recovery recovers from panics and responds 500

You can create your own and add with app.Use.

## Error Handling & Recovery

When a handler panics, the Recovery middleware intercepts it, logs, and returns 500 Internal Server Error.
You can also call ctx.Error(err) inside any handler; errors get collected and can be logged or handled in middleware (e.g. Logger prints them).

## Examples

JSON API

```go
app.GET("/api/users/<id>", func(ctx *igo.Context) {
    user := getUserFromDB(ctx.GetParam("id"))
    ctx.JSON(200, user)
})
```

Template + Query

```go
app.GET("/search", func(ctx *igo.Context) {
    q := ctx.GetQuery("q")
    results := doSearch(q)
    ctx.HTML(200, "search", results)
})
```

File Upload

```go
app.POST("/upload", func(ctx *igo.Context) {
    f, fh, err := ctx.FormFile("file")
    if err != nil {
        ctx.Abort(400, "Upload failed")
        return
    }
    defer f.Close()
    saveFile(f, fh.Filename)
    ctx.Text(200, "Upload successful")
})
```

## Roadmap

- Support for a radix tree router (quicker route matching)
- Hot-reload
- Support for template caching

_Community contributions welcome_
