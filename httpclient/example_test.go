package httpclient_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/pallat/wtf/httpclient"
)

func ExampleGet() {
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"how to use httpclient.Get"}`))
	}))
	defer serv.Close()

	type Result struct {
		Data string `json:"data"`
	}

	// example of using httpclient.Get
	resp, err := httpclient.Get[Result](context.Background(), httpclient.NewHTTPClient(), serv.URL)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	fmt.Fprintf(os.Stdout, "%d\n", resp.Code)
	fmt.Fprintln(os.Stdout, resp.Response.Data)

	// Output:
	// 200
	// how to use httpclient.Get
}

func ExamplePost() {
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"how to use httpclient.Post"}`))
	}))
	defer serv.Close()

	type Result struct {
		Data string `json:"data"`
	}

	type Payload struct {
		Payload string `json:"payload"`
	}

	pl := Payload{Payload: "request paly load"}

	// example of using httpclient.Post
	resp, err := httpclient.Post[Payload, Result](context.Background(), httpclient.NewHTTPClient(), serv.URL, pl)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	fmt.Fprintf(os.Stdout, "%d\n", resp.Code)
	fmt.Fprintln(os.Stdout, resp.Response.Data)

	// Output:
	// 200
	// how to use httpclient.Post
}
