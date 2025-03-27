package httpclient

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type APIResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func TestHTTPClient(t *testing.T) {
	t.Run("should be able to GET request succesfuly", func(t *testing.T) {
		client := NewHTTPClient()
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Get[APIResponse](ctx, client, serv.URL)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to POST request succesfuly", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := RESQ{Name: "somebody name"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Post[RESQ, APIResponse](ctx, client, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("error when convert to json request", func(t *testing.T) {
		ctx := context.Background()
		client := NewHTTPClient()
		req, err := NewRequest(ctx, client, http.MethodPost, "http://0.0.0.0/api", make(chan int))

		if req != nil {
			t.Errorf("unsupported type expect expect nil request but actual %v\n", req)
		}
		if err == nil {
			t.Error("unsupported type expect not nil error")
		}
	})

	t.Run("error when create http request", func(t *testing.T) {
		ctx := context.Background()
		client := NewHTTPClient()
		req, err := NewRequest(ctx, client, "\\", "http://0.0.0.0/api", bytes.NewBufferString(`{}`))

		if req != nil {
			t.Errorf("invalid method expect expect nil request but actual %v\n", req)
		}
		if err == nil {
			t.Error("invalid method expect not nil error")
		}
	})
}

func TestDoRequest(t *testing.T) {
	t.Run("client do error", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api", bytes.NewBufferString("{}"))

		_, err := func() (Response[RESQ], error) {
			return DoRequest[RESQ](client, req)
		}()

		if err == nil {
			t.Error("invalid url expect not nil error")
		}
	})
	t.Run("error when decode json response", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()

		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "text/json")
			w.Write(nil)
		}))
		defer serv.Close()

		ctx := context.Background()
		payload := RESQ{Name: "somebody name"}
		req, _ := NewRequest(ctx, client, http.MethodPost, serv.URL, payload)

		_, err := func() (Response[RESQ], error) {
			return DoRequest[RESQ](client, req)
		}()

		if err == nil {
			t.Error("invalid response expect not nil error")
		}
	})
}

func TestDo(t *testing.T) {
	t.Run("POST method with request body and response", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		type RESQ struct {
			Name string `json:"name"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "POST"
		payload := RESQ{Name: "test"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		resp, err := do[RESQ, RESP](ctx, client, method, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("GET method with no request body", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "GET"
		var payload bytes.Buffer
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		reps, err := do[bytes.Buffer, RESP](ctx, client, method, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, reps)
		assert.Equal(t, 200, reps.Code)
		assert.Equal(t, "OK", reps.Response.Message)
		assert.Equal(t, "Success", reps.Response.Description)
	})

	t.Run("error newRequest", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "\\"
		var payload bytes.Buffer
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		_, err := do[bytes.Buffer, RESP](ctx, client, method, serv.URL, payload)

		if err == nil {
			t.Error("can not new http request expect not nil error")
		}
	})
}
