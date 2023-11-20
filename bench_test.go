package benchmarks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jub0bs/fcors"
	"github.com/rs/cors"
)

const hostMaxLen = 253

var dummyHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

var reqHeaders = []string{
	"Accept",
	"Content-Type",
	"X-Requested-With",
}

const (
	headerOrigin = "Origin"
	headerACRM   = "Access-Control-Request-Method"
	headerACRH   = "Access-Control-Request-Headers"
)

const (
	rsCors      = "rs_cors"
	jub0bsFcors = "jub0bs_fcors"
)

const dummyOrigin = "https://jub0bs.com"

type Middleware struct {
	wrap func(http.Handler) http.Handler
	name string
}

type Case2 struct {
	middleware []Middleware
	desc       string
	req        *http.Request
}

func BenchmarkAll(b *testing.B) {
	cases := []Case2{
		{
			middleware: []Middleware{
				{
					wrap: identity[http.Handler],
					name: "identity",
				},
			},
			desc: "vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{dummyOrigin},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "single_vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		},
		{
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: append(otherOrigins, dummyOrigin),
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin, otherOrigins...),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "multiple_vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{
							"https://a" + strings.Repeat(".a", hostMaxLen/2),
							"https://b" + strings.Repeat(".a", hostMaxLen/2),
						},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(
							"https://a"+strings.Repeat(".a", hostMaxLen/2),
							"https://b"+strings.Repeat(".a", hostMaxLen/2),
						),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "pathological_vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: append(manyOrigins, dummyOrigin),
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin, manyOrigins...),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "many_vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{"*"},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromAnyOrigin(),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "any_vs_actual",
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{dummyOrigin},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "single_vs_preflight",
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: append(otherOrigins, dummyOrigin),
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin, otherOrigins...),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "multiple_vs_preflight",
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{
							"https://a" + strings.Repeat(".a", hostMaxLen/2),
							"https://b" + strings.Repeat(".a", hostMaxLen/2),
						},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(
							"https://a"+strings.Repeat(".a", hostMaxLen/2),
							"https://b"+strings.Repeat(".a", hostMaxLen/2),
						),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "pathological_vs_preflight",
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: append(manyOrigins, dummyOrigin),
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromOrigins(dummyOrigin, manyOrigins...),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "many_vs_preflight",
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{"*"},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromAnyOrigin(),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "any_vs_preflight",
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {dummyOrigin},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			middleware: []Middleware{
				{
					wrap: cors.New(cors.Options{
						AllowedOrigins: []string{"*"},
						AllowedHeaders: reqHeaders,
					}).Handler,
					name: rsCors,
				}, {
					wrap: mustAllowAccess(
						fcors.FromAnyOrigin(),
						withRequestHeaders(reqHeaders...),
					),
					name: jub0bsFcors,
				},
			},
			desc: "any_1header_vs_preflight",
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
		for _, middleware := range c.middleware {
			handler := middleware.wrap(dummyHandler)
			f := func(b *testing.B) {
				recs := make([]*httptest.ResponseRecorder, b.N)
				for i := 0; i < b.N; i++ {
					recs[i] = httptest.NewRecorder()
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					handler.ServeHTTP(recs[i], c.req)
				}
			}
			const padding = 38
			name := middleware.name + strings.Repeat("_", padding-len(middleware.name)-len(c.desc)) + c.desc
			b.Run(name, f)
		}
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

var otherOrigins = []string{
	"https://example.com",
	"https://*.example.com",
	"https://google.com",
	"https://*.google.com",
}

var manyOrigins []string

func init() {
	const n = 100
	for i := 0; i < n; i++ {
		manyOrigins = append(
			manyOrigins,
			// https
			fmt.Sprintf("https://example%d.com", i),
			fmt.Sprintf("https://example%d.com:7070", i),
			fmt.Sprintf("https://example%d.com:8080", i),
			fmt.Sprintf("https://example%d.com:9090", i),
			// one subdomain deep
			fmt.Sprintf("https://foo.example%d.com", i),
			fmt.Sprintf("https://foo.example%d.com:6060", i),
			fmt.Sprintf("https://foo.example%d.com:7070", i),
			fmt.Sprintf("https://foo.example%d.com:9090", i),
			// two subdomains deep
			fmt.Sprintf("https://foo.bar.example%d.com", i),
			fmt.Sprintf("https://foo.bar.example%d.com:6060", i),
			fmt.Sprintf("https://foo.bar.example%d.com:7070", i),
			fmt.Sprintf("https://foo.bar.example%d.com:9090", i),
			// arbitrary subdomains
			fmt.Sprintf("https://*.foo.bar.example%d.com", i),
			fmt.Sprintf("https://*.foo.bar.example%d.com:6060", i),
			fmt.Sprintf("https://*.foo.bar.example%d.com:7070", i),
			fmt.Sprintf("https://*.foo.bar.example%d.com:9090", i),
		)
	}
}
