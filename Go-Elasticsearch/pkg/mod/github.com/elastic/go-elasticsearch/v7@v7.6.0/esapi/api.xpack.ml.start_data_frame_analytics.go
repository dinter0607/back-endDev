// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 7.5.0: DO NOT EDIT

package esapi

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func newMLStartDataFrameAnalyticsFunc(t Transport) MLStartDataFrameAnalytics {
	return func(id string, o ...func(*MLStartDataFrameAnalyticsRequest)) (*Response, error) {
		var r = MLStartDataFrameAnalyticsRequest{ID: id}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// MLStartDataFrameAnalytics -
//
// See full documentation at http://www.elastic.co/guide/en/elasticsearch/reference/current/start-dfanalytics.html.
//
type MLStartDataFrameAnalytics func(id string, o ...func(*MLStartDataFrameAnalyticsRequest)) (*Response, error)

// MLStartDataFrameAnalyticsRequest configures the ML Start Data Frame Analytics API request.
//
type MLStartDataFrameAnalyticsRequest struct {
	ID string

	Body io.Reader

	Timeout time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r MLStartDataFrameAnalyticsRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "POST"

	path.Grow(1 + len("_ml") + 1 + len("data_frame") + 1 + len("analytics") + 1 + len(r.ID) + 1 + len("_start"))
	path.WriteString("/")
	path.WriteString("_ml")
	path.WriteString("/")
	path.WriteString("data_frame")
	path.WriteString("/")
	path.WriteString("analytics")
	path.WriteString("/")
	path.WriteString(r.ID)
	path.WriteString("/")
	path.WriteString("_start")

	params = make(map[string]string)

	if r.Timeout != 0 {
		params["timeout"] = formatDuration(r.Timeout)
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

	req, err := newRequest(method, path.String(), r.Body)
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

	if r.Body != nil {
		req.Header[headerContentType] = headerContentTypeJSON
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

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f MLStartDataFrameAnalytics) WithContext(v context.Context) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.ctx = v
	}
}

// WithBody - The start data frame analytics parameters.
//
func (f MLStartDataFrameAnalytics) WithBody(v io.Reader) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.Body = v
	}
}

// WithTimeout - controls the time to wait until the task has started. defaults to 20 seconds.
//
func (f MLStartDataFrameAnalytics) WithTimeout(v time.Duration) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.Timeout = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f MLStartDataFrameAnalytics) WithPretty() func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f MLStartDataFrameAnalytics) WithHuman() func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f MLStartDataFrameAnalytics) WithErrorTrace() func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f MLStartDataFrameAnalytics) WithFilterPath(v ...string) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f MLStartDataFrameAnalytics) WithHeader(h map[string]string) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f MLStartDataFrameAnalytics) WithOpaqueID(s string) func(*MLStartDataFrameAnalyticsRequest) {
	return func(r *MLStartDataFrameAnalyticsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
