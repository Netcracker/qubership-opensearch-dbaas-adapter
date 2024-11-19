package api

import (
	"context"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func newCreateUserFunc(t opensearchapi.Transport) CreateUser {
	return func(username string, o ...func(request *CreateUserRequest)) (*opensearchapi.Response, error) {
		var r = CreateUserRequest{Username: username}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// CreateUser creates a user
type CreateUser func(username string, o ...func(request *CreateUserRequest)) (*opensearchapi.Response, error)

// CreateUserRequest configures the User API request.
type CreateUserRequest struct {
	Username string

	Body io.Reader

	WaitForCompletion *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do function executes the request and returns response or error.
func (r CreateUserRequest) Do(ctx context.Context, transport opensearchapi.Transport) (*opensearchapi.Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = http.MethodPut
	path.Grow(1 + len("_plugins/_security/api/internalusers") + 1 + len(r.Username))
	path.WriteString("/_plugins/_security/api/internalusers")
	path.WriteString("/")
	path.WriteString(r.Username)

	params = make(map[string]string)

	if r.WaitForCompletion != nil {
		params["wait_for_completion"] = strconv.FormatBool(*r.WaitForCompletion)
	}

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := http.NewRequest(method, path.String(), r.Body)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := opensearchapi.Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithUsername sets the request username.
func (f CreateUser) WithUsername(v string) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.Username = v
	}
}

// WithBody sets the request body.
func (f CreateUser) WithBody(v io.Reader) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.Body = v
	}
}

// WithContext sets the request context.
func (f CreateUser) WithContext(v context.Context) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
func (f CreateUser) WithPretty() func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
func (f CreateUser) WithHuman() func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
func (f CreateUser) WithErrorTrace() func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
func (f CreateUser) WithFilterPath(v ...string) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
func (f CreateUser) WithHeader(h map[string]string) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
func (f CreateUser) WithOpaqueID(s string) func(*CreateUserRequest) {
	return func(r *CreateUserRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}