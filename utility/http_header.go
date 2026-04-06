package utility

import "net/http"

func WithBrowserUA(header http.Header) {
	header.Set("User-Agent",
		`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36`,
	)
}
