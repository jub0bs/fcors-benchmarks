package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/cors"

	"github.com/jub0bs/fcors"
)

var dummyHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

var reqHeaders = []string{"Accept", "Content-Type", "X-Requested-With"}

var otherOrigins = []string{
	"https://example.com",
	"https://*.example.com",
	"https://google.com",
	"https://*.google.com",
}

type Middleware = func(http.Handler) http.Handler

const (
	headerOrigin = "Origin"
	headerACRM   = "Access-Control-Request-Method"
	headerACRH   = "Access-Control-Request-Headers"
)

const dummyOrigin = "https://jub0bs.com"

func BenchmarkAll(b *testing.B) {
	type Case struct {
		name string
		mw   Middleware
		req  *http.Request
	}
	cases := []Case{
		{
			name: "without CORS",
			mw:   identity[http.Handler],
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors single origin vs actual request",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{dummyOrigin},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors multiple origins vs actual request",
			mw: cors.New(cors.Options{
				AllowedOrigins: append(otherOrigins, dummyOrigin),
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors any origin vs actual request",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors single origin vs preflight request",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{dummyOrigin},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "rs_cors multiple origins vs preflight request",
			mw: cors.New(cors.Options{
				AllowedOrigins: append(otherOrigins, dummyOrigin),
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "rs_cors any origin vs preflight request",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "rs_cors any origin vs preflight request with one header",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
					headerACRH:   {"Accept"},
				},
			),
		}, {
			name: "jub0bs_fcors single origin vs actual request",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "jub0bs_fcors multiple origins vs actual request",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin, otherOrigins...),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "jub0bs_fcors any origin vs actual request",
			mw: mustAllowAccess(
				fcors.FromAnyOrigin(),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "jub0bs_fcors single origin vs preflight request",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "jub0bs_fcors multiple origins vs preflight request",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin, otherOrigins...),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "jub0bs_fcors any origin vs preflight request",
			mw: mustAllowAccess(
				fcors.FromAnyOrigin(),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "jub0bs_fcors any origin vs preflight request with one header",
			mw: mustAllowAccess(
				fcors.FromAnyOrigin(),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
					headerACRH:   {"Accept"},
				},
			),
		},
	}
	for _, c := range cases {
		handler := c.mw(dummyHandler)
		f := func(b *testing.B) {
			rec := httptest.NewRecorder()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				clear(rec.Header())
				handler.ServeHTTP(rec, c.req)
			}
		}
		b.Run(c.name, f)
	}
}

func withRequestHeaders(names ...string) fcors.OptionAnon {
	if len(names) == 0 {
		panic("at least one header is required")
	}
	return fcors.WithRequestHeaders(names[0], names[1:]...)
}

func mustAllowAccess(one fcors.OptionAnon, others ...fcors.OptionAnon) fcors.Middleware {
	cors, err := fcors.AllowAccess(one, others...)
	if err != nil {
		panic("invalid policy")
	}
	return cors
}

func newRequest(method string, headers http.Header) *http.Request {
	const dummyEndpoint = "https://example.com/whatever"
	req := httptest.NewRequest(method, dummyEndpoint, nil)
	req.Header = headers
	return req
}

func identity[T any](t T) T { return t }

func clear(h http.Header) {
	for k := range h {
		delete(h, k)
	}
}
