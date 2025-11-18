<div align="center">
<a href="#">
<img width="312" height="335" alt="igo-logo" src="assets/igo.svg" />
</a>
</div>

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/mukailasam/igo)](https://goreportcard.com/report/github.com/mukailasam/igo)

# Igo - Go Web Framework

Igo is a simple and lightweight micro web-framework for Go - inspired by [Bottle](https://bottlepy.org).  
It is distributed as a **single file (`igo.go`)** package and depends only on the Go standard library.

## Background

I built it because I wanted a Bottle-like microframework for my <a href="https://github.com/mukailasam/codelab"> CodeLab </a> - something personal, minimal, and easy to use for experimenting and creating **small Go projects** without setup overhead.

> The name **“Igo”** comes from the **Yoruba** word **“Ìgò”**, which means **“bottle.”** \
> ma pa igo lo ri e😂

## Install

```go
go get -u github.com/mukailasam/igo
```

## Features

- **Routing:** Map HTTP requests to handler functions with dynamic URL support.
- **Utilities:** Easy access to URL parameters, form data, headers, and other HTTP features.
- **Server:** Built-in development server powered by Go’s `net/http`.

Homepage and documentation: <a href="">https://mukailasam.github.io/mukailasam/igo/docs/docs.html</a>

## Example Usage:

```go
package main

import (
	"fmt"

	"github.com/mukailasam/igo"
)

func main() {
	app := igo.NewIgo()

	app.GET("/<name>", func(ctx *igo.Context) {
		name := ctx.GetParam("name")
		ctx.Text(200, fmt.Sprintf("Hello, %s!", name))
	})

	app.Run(":1337")
}
```

## Contributing

Contributions are welcome!
To contribute:

- Fork the repository

- Create a new branch

- Implement your feature or fix

- Open a pull request

Make sure to follow Go’s code formatting and keep your commits clean and descriptive.
