package middleware

import "net/http"

func TestHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//Add a "server" header to the response
		writer.Header().Set("Server", "Apache")
		writer.Header().Set("Strict-Transport-Security", "max-age=15768000; includeSubDomains; preload;")
		writer.Header().Set("Expires", "Thu, 19 Nov 1981 08:52:00 GMT")
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Content-Security-Policy", "default-src 'none'")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Vary", "Brief,Prefer,Accept-Encoding")
		writer.Header().Set("DAV", "1, 3, extended-mkcol, access-control, calendarserver-principal-property-search, nc-calendar-search, nc-enable-birthday-calendar")
		writer.Header().Set("Referrer-Policy", "no-referrer")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Download-Options", "noopen")
		writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		writer.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		writer.Header().Set("X-Robots-Tag", "none")
		writer.Header().Set("X-XSS-Protection", "1; mode=block")

		handler.ServeHTTP(writer, request)
	})
}
