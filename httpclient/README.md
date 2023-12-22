# httpclient

### The `httpclient` package implements the core functionality of a [coder](https://github.com/easy-techno-lab/proton/blob/main/coder/README.md)-based http client.

## Getting Started

```go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/easy-techno-lab/proton/coder"
	"github.com/easy-techno-lab/proton/httpclient"
	"github.com/easy-techno-lab/proton/logger"
)

func main() {
	cdrJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	clientJSON := httpclient.New(cdrJSON, http.DefaultClient)

	URL := "http://localhost:8080/example/"

	params := make(url.Values)
	params.Add("id", "1")

	// To add additional data to the request, use the optional function f(*http.Request)
	f := func(r *http.Request) {
		r.Header.Set("Accept", "application/json")
		r.URL.RawQuery = params.Encode()
	}

	resp, err := clientJSON.Request(context.TODO(), http.MethodGet, URL, nil, f)
	if err != nil {
		panic(err)
	}

	defer logger.Closer(resp.Body)

	res := &struct {
		// some fields
	}{}

	if err = clientJSON.Decode(resp.Body, res); err != nil {
		panic(err)
	}
}

```

```go
package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/easy-techno-lab/proton/coder"
	"github.com/easy-techno-lab/proton/httpclient"
	"github.com/easy-techno-lab/proton/logger"
)

func main() {
	cdrJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	clientJSON := httpclient.New(cdrJSON, http.DefaultClient)

	URL := "http://localhost:8080/v1/example/"

	req := &struct {
		ID int `json:"id"`
	}{ID: 1}

	resp, err := clientJSON.Request(context.TODO(), http.MethodPost, URL, req, nil)
	if err != nil {
		panic(err)
	}

	defer logger.Closer(resp.Body)

	res := &struct {
		// some fields
	}{}

	if err = clientJSON.Decode(resp.Body, res); err != nil {
		panic(err)
	}
}

```

### The `httpclient` package contains functions that are used as middleware on the http client side.

## Getting Started

```go
package main

import (
	"net/http"

	"github.com/easy-techno-lab/proton/httpclient"
	"github.com/easy-techno-lab/proton/logger"
)

func main() {
	transport := httpclient.RoundTripperSequencer(
		http.DefaultTransport,
		httpclient.DumpHttp(logger.LevelTrace),
		httpclient.Timer(logger.LevelInfo),
		httpclient.PanicCatcher,
	)

	hct := new(http.Client)
	hct.Transport = transport
}

```