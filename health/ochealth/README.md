<h1>ochealth</h1>

# Purpose

`ochealth` exports health.State data via OpenCensus.

`ochealth` allows you to take the state from a health server and export the data in it via OpenCensus metrics. 
For example, to track the readiness of your server over time you can use OpenCensus to export via Prometheus metrics where you can graph and alert on your service not being ready. 
It also allows you to export version information so you can easily correlate version changes with changes in other metrics making it easier to identify regressions.

For more information see [register.go](./register.go).

## Example Usage

```go
package ochealth_test

import (
	"fmt"
	"net/http/httptest"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/anz-bank/pkg/health"
	"github.com/anz-bank/pkg/health/ochealth"
)

func Example() {
	// normally set with go build linker option
	health.CommitHash = "0123456789abcdef0123456789abcdef01234567"
	health.Semver = "v1.2.3"

	// Real code handles errors.
	server, _ := health.NewHTTPServer()
	_ = ochealth.Register(server.State, ochealth.WithPrefix("myapp"))
	prom, _ := ocprom.NewExporter(ocprom.Options{})

	server.SetReady(true)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	prom.ServeHTTP(w, r)

	fmt.Print(w.Body.String())
	//nolint:lll
	// output:
	// # HELP myapp_ready Readiness state of server
	// # TYPE myapp_ready gauge
	// myapp_ready 1
	// # HELP myapp_version Version information
	// # TYPE myapp_version gauge
	// myapp_version{build_log_url="undefined",commit_hash="0123456789abcdef0123456789abcdef01234567",container_tag="undefined",repo_url="undefined",semver="v1.2.3"} 1
}
```
