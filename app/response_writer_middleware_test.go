package app

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pallat/wtf/serror"
)

func TestWriter(t *testing.T) {
	t.Run("response message json decode fail", func(t *testing.T) {
		handler := func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, "\\")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Add("X-REF-ID", "12345")

		c.Request = req

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)
		buf := bytes.NewBuffer([]byte{})
		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(l)

		e.Use(AutoLoggingMiddleware)
		e.POST("/api", handler)

		e.HandleContext(c)

		t.Run("log should contains DEBUG log", func(t *testing.T) {
			if !bytes.Contains(buf.Bytes(), []byte(`level=DEBUG`)) && !bytes.Contains(buf.Bytes(), []byte("invalid character")) {
				t.Errorf(`"level=DEBUG" and "invalid character" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
		t.Run("not standard should contains Warn log, standard is important", func(t *testing.T) {
			if !bytes.Contains(buf.Bytes(), []byte(`level=WARN`)) && !bytes.Contains(buf.Bytes(), []byte("response not standard")) {
				t.Errorf(`"level=Warn" and "standard is important" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
	})

	t.Run("log include meta data, refID", func(t *testing.T) {
		handler := func(c *gin.Context) {
			err := serror.New("test message")
			c.JSON(http.StatusInternalServerError, Response{
				Message: err.Error(),
			})
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Add("X-REF-ID", "12345")

		c.Request = req

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)
		buf := bytes.NewBuffer([]byte{})
		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		e.Use(AutoLoggingMiddleware)
		e.POST("/api", handler)

		e.HandleContext(c)

		expect := fmt.Sprintf(`%s=%s`, string(refIDKey), "12345")
		if !bytes.Contains(buf.Bytes(), []byte(expect)) {
			t.Errorf("%q should contain, actual\n%q\n\n", expect, buf.Bytes())
		}
	})
	t.Run("log should contain code when response/code not 0", func(t *testing.T) {
		handler := func(c *gin.Context) {
			err := serror.New("test message")
			c.JSON(http.StatusInternalServerError, Response{
				Code:    5000,
				Message: err.Error(),
			})
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Add("X-REF-ID", "12345")

		c.Request = req

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)
		buf := bytes.NewBuffer([]byte{})
		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		e.Use(AutoLoggingMiddleware)
		e.POST("/api", handler)

		e.HandleContext(c)

		expect := fmt.Sprintf(`code=%s`, "5000")
		if !bytes.Contains(buf.Bytes(), []byte(expect)) {
			t.Errorf("%q should contain, actual\n%q\n\n", expect, buf.Bytes())
		}
	})
	t.Run("log level selected", func(t *testing.T) {
		// httpw := httptest.NewRecorder()

		// ctx, _ := gin.CreateTestContext(httpw)
		// w := newWriter(ctx.Writer, map[string]string{
		// 	string(refIDKey): "12345",
		// })

		// defaultLogger := slog.Default()
		// defer slog.SetDefault(defaultLogger)

		// buf := bytes.NewBuffer([]byte{})

		// l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		// slog.SetDefault(l)

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)
		buf := bytes.NewBuffer([]byte{})
		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		e.Use(AutoLoggingMiddleware)

		t.Run("ERROR level for Internal Server Error", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/internal", nil)
			req.Header.Add("X-REF-ID", "12345")
			c.Request = req

			e.POST("/internal", func(c *gin.Context) {
				err := serror.New("test message")
				c.JSON(http.StatusInternalServerError, Response{
					Code:    5000,
					Message: err.Error(),
				})
			})

			e.HandleContext(c)

			if !bytes.Contains(buf.Bytes(), []byte(`level=ERROR`)) {
				t.Errorf(`"level=ERROR" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
		t.Run("ERROR level for Bad Request", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/badrequest", nil)
			req.Header.Add("X-REF-ID", "12345")
			c.Request = req

			e.POST("/badrequest", func(c *gin.Context) {
				err := serror.New("test message")
				c.JSON(http.StatusBadRequest, Response{
					Code:    5000,
					Message: err.Error(),
				})
			})

			e.HandleContext(c)
			if !bytes.Contains(buf.Bytes(), []byte(`level=ERROR`)) {
				t.Errorf(`"level=ERROR" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
	})
	t.Run("response message from serror correct logging", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newWriter(ctx.Writer, map[string]string{})

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("test message")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="test message"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=app.TestWriter`)) {
			t.Errorf("%q should contain", `func=app.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((test message+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "test message"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
	t.Run("response message with escape string", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newWriter(ctx.Writer, map[string]string{})

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("invalid character 'a' looking for beginning of object key string")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="invalid character 'a' looking for beginning of object key string"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=app.TestWriter`)) {
			t.Errorf("%q should contain", `func=app.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((invalid character 'a' looking for beginning of object key string+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "invalid character 'a' looking for beginning of object key string"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
	t.Run("response message with escape string '\u003c' and '\u003e;", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newWriter(ctx.Writer, map[string]string{})

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("invalid character 'a' looking for beginning of object key string")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "\u003ctest\u003e%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="<test>invalid character 'a' looking for beginning of object key string"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=app.TestWriter`)) {
			t.Errorf("%q should contain", `func=app.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((\u003ctest\u003einvalid character 'a' looking for beginning of object key string+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "<test>invalid character 'a' looking for beginning of object key string"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
}

func BenchmarkWriterWriteSerror(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	httpw := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(httpw)
	ctx.Writer = newWriter(ctx.Writer, map[string]string{})

	defaultLogger := slog.Default()
	defer slog.SetDefault(defaultLogger)

	l := slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), &slog.HandlerOptions{}))
	slog.SetDefault(l)

	for range b.N {
		ctx.JSON(500, Response{
			Code:    50001,
			Message: serror.New("testing error message").Error(),
		})
	}
}
