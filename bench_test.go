package benchmarks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/cors"

	"github.com/jub0bs/fcors"
)

const hostMaxLen = 253

var dummyHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

var reqHeaders = []string{
	"Accept",
	"Content-Type",
	"X-Requested-With",
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
			name: "without CORS_________________vs_actual",
			mw:   identity[http.Handler],
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors_______________single_vs_actual",
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
			name: "jub0bs_fcors__________single_vs_actual",
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
			name: "rs_cors_____________multiple_vs_actual",
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
			name: "jub0bs_fcors________multiple_vs_actual",
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
			name: "rs_cors_________pathological_vs_actual",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{
					"https://a" + strings.Repeat(".a", hostMaxLen/2),
					"https://b" + strings.Repeat(".a", hostMaxLen/2),
				},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
				},
			),
		}, {
			name: "jub0bs_fcors____pathological_vs_actual",
			mw: mustAllowAccess(
				fcors.FromOrigins(
					"https://a"+strings.Repeat(".a", hostMaxLen/2),
					"https://b"+strings.Repeat(".a", hostMaxLen/2),
				),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
				},
			),
		}, {
			name: "rs_cors_________________many_vs_actual",
			mw: cors.New(cors.Options{
				AllowedOrigins: append(manyOrigins, dummyOrigin),
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "jub0bs_fcors____________many_vs_actual",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin, manyOrigins...),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodGet,
				http.Header{
					headerOrigin: {dummyOrigin},
				},
			),
		}, {
			name: "rs_cors__________________any_vs_actual",
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
			name: "jub0bs_fcors_____________any_vs_actual",
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
			name: "rs_cors____________single_vs_preflight",
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
			name: "jub0bs_fcors_______single_vs_preflight",
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
			name: "rs_cors__________multiple_vs_preflight",
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
			name: "jub0bs_fcors_____multiple_vs_preflight",
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
			name: "rs_cors______pathological_vs_preflight",
			mw: cors.New(cors.Options{
				AllowedOrigins: []string{
					"https://a" + strings.Repeat(".a", hostMaxLen/2),
					"https://b" + strings.Repeat(".a", hostMaxLen/2),
				},
				AllowedHeaders: reqHeaders,
			}).Handler,
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "jub0bs_fcors_pathological_vs_preflight",
			mw: mustAllowAccess(
				fcors.FromOrigins(
					"https://a"+strings.Repeat(".a", hostMaxLen/2),
					"https://b"+strings.Repeat(".a", hostMaxLen/2),
				),
				withRequestHeaders(reqHeaders...),
			),
			req: newRequest(
				http.MethodOptions,
				http.Header{
					headerOrigin: {"https://c" + strings.Repeat(".a", hostMaxLen/2)},
					headerACRM:   {http.MethodGet},
				},
			),
		}, {
			name: "rs_cors______________many_vs_preflight",
			mw: cors.New(cors.Options{
				AllowedOrigins: append(manyOrigins, dummyOrigin),
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
			name: "jub0bs_fcors_________many_vs_preflight",
			mw: mustAllowAccess(
				fcors.FromOrigins(dummyOrigin, manyOrigins...),
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
			name: "rs_cors_______________any_vs_preflight",
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
			name: "jub0bs_fcors__________any_vs_preflight",
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
			name: "rs_cors_______any_1header_vs_preflight",
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
			name: "jub0bs_fcors__any_1header_vs_preflight",
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
