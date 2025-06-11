package util

import "net/http"

var GlobalHeader = http.Header{
	"User-Agent": []string{
		"Mozilla/5.0 (X11; Linux 86_64) AppleWebKit/537.36 " +
			"(KHTML, like Gecko) Version/4.0 Chrome/130.0.0.0 Safari/537.36"},
}
