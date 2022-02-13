package middlewares

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"ozonLinkShortener/pkg/helpers"

	"github.com/sirupsen/logrus"
)

var requestId = 0

func newRequest(logEntry *logrus.Entry) *logrus.Entry {
	requestId += 1
	return logEntry.WithField("rqId", requestId)
}

func AddLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.NewEntry(logrus.New())
		logger = newRequest(logger)
		ctx := context.WithValue(r.Context(), helpers.ContextLoggerKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseInputJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)

		if err != nil || len(bytes) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			output := helpers.Output{
				Error:       true,
				ErrorString: "Body should not be empty",
			}
			bytes, err := json.Marshal(output)
			if err != nil {
				r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry).Error("Error marshaling error: ", err)
			}
			_, err = w.Write(bytes)
			if err != nil {
				r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry).Error("Error sending error: ", err)
			}
			return
		}
		input := helpers.Input{}
		err = json.Unmarshal(bytes, &input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			output := helpers.Output{
				Error:       true,
				ErrorString: "Body should be json",
			}
			bytes, err := json.Marshal(output)
			if err != nil {
				r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry).Error("Error marshaling error: ", err)
			}
			_, err = w.Write(bytes)
			if err != nil {
				r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry).Error("Error sending error: ", err)
			}
			return
		}
		ctx := context.WithValue(r.Context(), helpers.ContextInputKey, input)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AddClientInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry).Error("Error parsing remoteAddr: ", r.RemoteAddr)
			ip = r.RemoteAddr
		}
		ctx := context.WithValue(r.Context(), helpers.ContextAuthorIpKey, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
