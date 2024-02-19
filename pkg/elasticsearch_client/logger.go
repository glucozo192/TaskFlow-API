package elastic_client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ColorLogger struct {
	Output             io.Writer
	EnableRequestBody  bool
	EnableResponseBody bool
}

func (l *ColorLogger) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	query, _ := url.QueryUnescape(req.URL.RawQuery)
	if query != "" {
		query = "?" + query
	}

	var (
		status string
		color  string
	)

	status = res.Status
	switch {
	case res.StatusCode > 0 && res.StatusCode < 300:
		color = "\x1b[32m"
	case res.StatusCode > 299 && res.StatusCode < 500:
		color = "\x1b[33m"
	case res.StatusCode > 499:
		color = "\x1b[31m"
	default:
		status = "ERROR"
		color = "\x1b[31;4m"
	}

	fmt.Fprintf(l.Output, "%6s \x1b[1;4m%s://%s%s\x1b[0m%s %s%s\x1b[0m \x1b[2m%s\x1b[0m\n",
		req.Method,
		req.URL.Scheme,
		req.URL.Host,
		req.URL.Path,
		query,
		color,
		status,
		dur.Truncate(time.Millisecond),
	)

	if l.RequestBodyEnabled() && req != nil && req.Body != nil && req.Body != http.NoBody {
		var buf bytes.Buffer
		if req.GetBody != nil {
			b, _ := req.GetBody()
			buf.ReadFrom(b)
		} else {
			buf.ReadFrom(req.Body)
		}
		fmt.Fprint(l.Output, "\x1b[2m")
		logBodyAsText(l.Output, &buf, "       »")
		fmt.Fprint(l.Output, "\x1b[0m")
	}

	if l.ResponseBodyEnabled() && res != nil && res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		fmt.Fprint(l.Output, "\x1b[2m")
		logBodyAsText(l.Output, &buf, "       «")
		fmt.Fprint(l.Output, "\x1b[0m")
	}

	if err != nil {
		fmt.Fprintf(l.Output, "\x1b[31;1m» ERROR \x1b[31m%v\x1b[0m\n", err)
	}

	if l.RequestBodyEnabled() || l.ResponseBodyEnabled() {
		fmt.Fprintf(l.Output, "\x1b[2m%s\x1b[0m\n", strings.Repeat("─", 80))
	}
	return nil
}

// RequestBodyEnabled returns true when the request body should be logged.
func (l *ColorLogger) RequestBodyEnabled() bool { return l.EnableRequestBody }

// ResponseBodyEnabled returns true when the response body should be logged.
func (l *ColorLogger) ResponseBodyEnabled() bool { return l.EnableResponseBody }

func logBodyAsText(dst io.Writer, body io.Reader, prefix string) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		s := scanner.Text()
		if s != "" {
			fmt.Fprintf(dst, "%s %s\n", prefix, s)
		}
	}
}
