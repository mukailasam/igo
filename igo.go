// Igo is a simple and lightweight micro web-framework for Go.

package igo

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// ----------------------------------------------------------------------------
// Context
// ----------------------------------------------------------------------------

// Context represents the context of the current HTTP request.
// It holds information about the request, response, and any route parameters.
// This struct is passed to every handler so it can access and modify request/response data.
type Context struct {
	responseWriter http.ResponseWriter // Used to send responses back to the client
	request        *http.Request       // The incoming HTTP request
	params         map[string]string   // Dynamic URL parameters (e.g. /user/:id → {"id": "42"})
}

// Header is a shortcut for ctx.Writer.Header().Set(key, value).
// It writes a header in the response.
func (c *Context) SetHeader(key, value string) {
	c.responseWriter.Header().Set(key, value)
}

// Redirect returns an HTTP redirect to the specific url
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.responseWriter, c.request, url, code)
}

// GetHeader returns value from request header.
// It gets the first value associated with the given key.
// If there are no values associated with the key, it returns ""
func (c *Context) GetHeader(key string) string {
	return c.request.Header.Get(key)
}

// WriteHeader is a shortcut for ctx.ResponseWriter.WriteHeader(code).
// It sends an HTTP response header with the provided status code.
func (c *Context) WriteHeader(statusCode int) {
	c.responseWriter.WriteHeader(statusCode)
}

// GetParam retrieves a URL path parameter by its key.
// For example, if route is "/users/<id>" and URL is "/users/123",
// GetParam("id") will return "123"
func (c *Context) GetParam(key string) string {
	return c.params[key]
}

// GetContentType detects and returns the MIME type of the provided data bytes using the http.DetectContentType function.
// The content type is determined by examining the initial bytes of the data.
// If the content type cannot be determined, it defaults to "application/octet-stream".
//
// Parameters:
//   - data: []byte - The byte slice containing the data to analyze
//
// Returns:
//   - string - The detected MIME type of the data
func (c *Context) GetContentType(data []byte) string {
	return http.DetectContentType(data)
}

// GetQuery retrieves a query parameter from the URL by its key.
// For example, if URL is "/?name=sam", GetQuery("name") will return "sam"
func (c *Context) GetQuery(key string) string {
	return c.request.URL.Query().Get(key)
}

// Abort sends an HTTP error response with the specified status code and message.
// It uses http.Error, which writes a plain-text error message, sets the proper
// Content-Type and headers, and sends the given status code.
func (c *Context) Abort(code int, msg string) {
	http.Error(c.responseWriter, msg, code)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.responseWriter, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// GetCookie retrieves a cookie from the request by its name.
// Returns the cookie value or an error if the cookie does not exist.
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetFormValue retrieves the value associated with the given key from the
// form data of the request. If the key does not exist, it returns an empty string.
func (c *Context) GetFormValue(key string) string {
	return c.request.Form.Get(key)
}

// Sends plain text to the client with the given status code
// It also sets the Content-Type as "text/plain".
func (c *Context) Text(statusCode int, data string) {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.WriteHeader(statusCode)
	c.responseWriter.Write([]byte(data))
}

// ----------------------------------------------------------------------------
// Igo
// ----------------------------------------------------------------------------

// HandlerFunc defines the function type for route handlers.
// Each route in Igo must have a handler with this signature.
type HandlerFunc func(ctx *Context)

// Route represents a single registered route in the Igo framework.
// It stores information about the HTTP method, pattern, handler, and
// compiled regex for matching incoming requests.
type Route struct {
	Method     string         // HTTP method (GET, POST, etc.)
	pattern    string         // Original route pattern (e.g. /users/<id>)
	Handler    HandlerFunc    // Function to handle matching requests
	regex      *regexp.Regexp // Compiled regex used for matching paths
	ParamNames []string       // Names of dynamic route parameters (e.g. ["id"])
}

// Igo is the main framework struct.
// It holds all registered routes and static file mappings.
type Igo struct {
	routes []*Route
}

// NewIgo creates and returns a new instance of the Igo framework.
func NewIgo() *Igo {
	return &Igo{
		routes: []*Route{},
	}
}

// addRoute registers a new route in the router with its method, pattern, and handler.
// It converts the pattern to a regex and extracts parameter names for dynamic routes.
func (i *Igo) addRoute(method, pattern string, handler HandlerFunc) {
	regex, paramNames := parsePattern(pattern)
	i.routes = append(i.routes, &Route{
		Method:     method,
		pattern:    pattern,
		Handler:    handler,
		regex:      regex,
		ParamNames: paramNames,
	})
}

// GET registers a new HTTP GET route.
func (i *Igo) GET(pattern string, handler HandlerFunc) {
	i.addRoute("GET", pattern, handler)
}

// POST registers a new HTTP POST route.
func (i *Igo) POST(pattern string, handler HandlerFunc) {
	i.addRoute("POST", pattern, handler)
}

// PUT registers a new HTTP PUT route.
func (i *Igo) PUT(pattern string, handler HandlerFunc) {
	i.addRoute("PUT", pattern, handler)
}

// DELETE registers a new HTTP DELETE route.
func (i *Igo) DELETE(pattern string, handler HandlerFunc) {
	i.addRoute("DELETE", pattern, handler)
}

// PATCH registers a new HTTP PATCH route.
func (i *Igo) PATCH(pattern string, handler HandlerFunc) {
	i.addRoute("PATCH", pattern, handler)
}

// OPTIONS registers a new HTTP OPTIONS route.
func (i *Igo) OPTIONS(pattern string, handler HandlerFunc) {
	i.addRoute("OPTIONS", pattern, handler)
}

// ServeHTTP is the main request handler for the Igo framework.
// It is automatically called by the HTTP server for each incoming request.
// This method handles both static files and registered routes.
func (i *Igo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Routes
	for _, route := range i.routes {
		if r.Method != route.Method {
			continue
		}
		matches := route.regex.FindStringSubmatch(path)
		if matches != nil {
			params := map[string]string{}
			for i, name := range route.ParamNames {
				params[name] = matches[i+1]
			}
			ctx := &Context{
				responseWriter: w,
				request:        r,
				params:         params,
			}

			route.Handler(ctx)
			return
		}
	}

	http.NotFound(w, r)
}

// parsePattern converts a route pattern into a regular expression and extracts parameter names.
//
// Example:
//
//	Input:  "/users/<id>/posts/*slug"
//	Output: Regex "^/users/([^/]+)/posts/(.*)$"
//	         ParamNames ["id", "slug"]
//
// Supported formats:
//
//	<name>  - named parameter that matches a single path segment
//	*name   - wildcard parameter that matches the rest of the path
func parsePattern(pattern string) (*regexp.Regexp, []string) {
	parts := strings.Split(pattern, "/")
	paramNames := []string{}
	for i, part := range parts {
		if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			name := part[1 : len(part)-1]
			paramNames = append(paramNames, name)
			parts[i] = "([^/]+)"
		}
		if strings.HasPrefix(part, "*") {
			name := part[1:]
			paramNames = append(paramNames, name)
			parts[i] = "(.*)"
		}
	}

	regex := regexp.MustCompile("^" + strings.Join(parts, "/") + "$")

	return regex, paramNames
}

// Run starts the HTTP server and begins listening on the specified address.
func (i *Igo) Run(addr string) error {
	log.Println("server started")
	return http.ListenAndServe(addr, i)
}
