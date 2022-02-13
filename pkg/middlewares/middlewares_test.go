package middlewares_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ozonLinkShortener/pkg/helpers"
	"ozonLinkShortener/pkg/middlewares"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestAddLogger(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(helpers.ContextLoggerKey)
		require.NotNil(t, val)
		logger, ok := r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry)
		require.NotEqual(t, false, ok)
		require.NotEqual(t, nil, logger)
	})
	handlerToTest := middlewares.AddLogger(nextHandler)
	req := httptest.NewRequest("Get", "http://localhost:15001", nil)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestManageLinks(t *testing.T) {
	{ // correct body
		input := helpers.Input{
			Link: "https:/asdaslkjas;dlklk;asdjflk;asdfas;dfjasdjf",
		}
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := r.Context().Value(helpers.ContextInputKey).(helpers.Input)
			require.True(t, ok)
			//require.Equal(t, inputIn.Link, input.Link)
		})
		handlerToTest := middlewares.ParseInputJson(nextHandler)

		toSend, err := json.Marshal(input)
		require.Nil(t, err)
		reader := bytes.NewReader(toSend)
		req := httptest.NewRequest("Get", "http://localhost:15001", reader)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
	{ // emtry body
		input := helpers.Input{
			Link: "https:/asdaslkjas;dlklk;asdjflk;asdfas;dfjasdjf",
		}
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := r.Context().Value(helpers.ContextInputKey).(helpers.Input)
			require.True(t, ok)
			//require.Equal(t, inputIn.Link, input.Link)
		})
		handlerToTest := middlewares.ParseInputJson(nextHandler)

		_, err := json.Marshal(input)
		require.Nil(t, err)
		reader := bytes.NewReader(nil)
		req := httptest.NewRequest("Get", "http://localhost:15001", reader)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}
