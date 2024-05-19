package internal

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrateNewUserPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			registrateNewUserPage(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
func TestAuthentificateUserPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/api/user/login"
	m[1] = "http://localhost:8080/api/user/login"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			authentificateUserPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestUploadNewOrderPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/api/user/orders"
	m[1] = "http://localhost:8080/api/user/orders"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			uploadNewOrderPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestGetUserBalancePage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/api/user/balance"
	m[1] = "http://localhost:8080/api/user/balance"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			getUserBalancePage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestDropBalancePage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/api/user/balance/withdraw"
	m[1] = "http://localhost:8080/api/user/balance/withdraw"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			dropBalancePage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestGetAllOrdersBalanceDropPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 500,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/api/user/balance/withdraw"
	m[1] = "http://localhost:8080/api/user/balance/withdraw"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			getAllOrdersBalanceDropPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}
