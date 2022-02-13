package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"ozonLinkShortener/internal/handlers"
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/internal/memory/inMemory"
	"ozonLinkShortener/pkg/helpers"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
}

func (m MockStorage) AddEntry(entry memory.MemoryEntry) error {
	return errors.New("Don't work")
}

func (m MockStorage) GetEntryByShortId(shortId string) (entry memory.MemoryEntry, err error) {
	return entry, errors.New("Don't work")
}
func (m MockStorage) GetEntryByLongLink(longLink string) (entry memory.MemoryEntry, err error) {
	return entry, errors.New("Don't work")
}
func (m MockStorage) Clear() error {
	return errors.New("Don't work")
}

func TestManageLonkLink(t *testing.T) {
	r := &http.Request{}
	ctx := context.WithValue(r.Context(), helpers.ContextLoggerKey, logrus.NewEntry(logrus.New()))
	ctx = context.WithValue(ctx, helpers.ContextAuthorIpKey, "123.123.123.11")
	handlers.Memory = memory.Memory{
		Storage: inMemory.GetMemoryInstance(),
	}
	output := helpers.Output{}

	{ // Correct add and get
		var shortLink string
		input := helpers.Input{
			Link: "https://google.com",
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w := httptest.NewRecorder()
		handlers.ManageLongLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, resp.StatusCode, 200)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, "https://google.com", output.InputLink)
		require.NotEqual(t, "", output.OutputLink)
		shortLink = output.OutputLink
		require.False(t, output.Error)
		require.Equal(t, "", output.ErrorString)

		input = helpers.Input{
			Link: "https://google.com",
		}
		ctx2 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w = httptest.NewRecorder()
		handlers.ManageLongLink(w, r.WithContext(ctx2))
		resp = w.Result()
		require.Equal(t, resp.StatusCode, 200)

		bytes, err = ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, "https://google.com", output.InputLink)
		require.Equal(t, shortLink, output.OutputLink)
		require.False(t, output.Error)
		require.Equal(t, "", output.ErrorString)

		err = handlers.Memory.Clear()
		require.Nil(t, err)
	}
	{ // incorrect add (database error)
		handlers.Memory = memory.Memory{
			Storage: MockStorage{},
		}
		input := helpers.Input{
			Link: "https://google.com",
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w := httptest.NewRecorder()
		handlers.ManageLongLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		err := handlers.Memory.Storage.Clear()
		require.NotNil(t, err)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, "https://google.com", output.InputLink)
		require.Equal(t, "", output.OutputLink)
		require.True(t, output.Error)
		handlers.Memory = memory.Memory{
			Storage: inMemory.GetMemoryInstance(),
		}
	}
}

func TestManageShortLink(t *testing.T) {
	r := &http.Request{}
	ctx := context.WithValue(r.Context(), helpers.ContextLoggerKey, logrus.NewEntry(logrus.New()))
	handlers.Memory = memory.Memory{
		Storage: inMemory.GetMemoryInstance(),
	}
	output := helpers.Output{}

	{ // not acceptible input link
		input := helpers.Input{
			Link: "ljk;safdlkjsfdljk;sadfljk;;",
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w := httptest.NewRecorder()
		handlers.ManageShortLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, "ljk;safdlkjsfdljk;sadfljk;;", output.InputLink)
		require.Equal(t, "", output.OutputLink)
		require.True(t, output.Error)
		require.Equal(t, "Not acceptible input link", output.ErrorString)
		err = handlers.Memory.Clear()
		require.Nil(t, err)
	}

	{ // Not Found
		shortLink := helpers.GetConfig().ShortLinkBase + "helloworld"
		input := helpers.Input{
			Link: shortLink,
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w := httptest.NewRecorder()
		handlers.ManageShortLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, resp.StatusCode, http.StatusNotFound)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, shortLink, output.InputLink)
		require.Equal(t, "", output.OutputLink)
		require.True(t, output.Error)
		require.Equal(t, "Not Found", output.ErrorString)
		err = handlers.Memory.Clear()
		require.Nil(t, err)
	}
	{ // incorrect get (database error)
		handlers.Memory = memory.Memory{
			Storage: MockStorage{},
		}
		shortLink := helpers.GetConfig().ShortLinkBase + "helloworld"
		input := helpers.Input{
			Link: shortLink,
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w := httptest.NewRecorder()
		handlers.ManageShortLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		err := handlers.Memory.Storage.Clear()
		require.NotNil(t, err)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, shortLink, output.InputLink)
		require.Equal(t, "", output.OutputLink)
		require.True(t, output.Error)
		handlers.Memory = memory.Memory{
			Storage: inMemory.GetMemoryInstance(),
		}
	}

	{ // correct add and get
		var shortLink string
		input := helpers.Input{
			Link: "https://google.com",
		}
		ctx1 := context.WithValue(ctx, helpers.ContextInputKey, input)
		ctx1 = context.WithValue(ctx1, helpers.ContextAuthorIpKey, "102.122.112.212")
		w := httptest.NewRecorder()
		handlers.ManageLongLink(w, r.WithContext(ctx1))
		resp := w.Result()
		require.Equal(t, resp.StatusCode, 200)
		bytes, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, "https://google.com", output.InputLink)
		require.NotEqual(t, "", output.OutputLink)
		shortLink = output.OutputLink
		require.False(t, output.Error)
		require.Equal(t, "", output.ErrorString)

		input = helpers.Input{
			Link: shortLink,
		}
		ctx2 := context.WithValue(ctx, helpers.ContextInputKey, input)
		w = httptest.NewRecorder()
		handlers.ManageShortLink(w, r.WithContext(ctx2))
		resp = w.Result()
		require.Equal(t, resp.StatusCode, 200)

		bytes, err = ioutil.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bytes, &output)
		require.Nil(t, err)
		require.Equal(t, shortLink, output.InputLink)
		require.Equal(t, "https://google.com", output.OutputLink)
		require.False(t, output.Error)
		require.Equal(t, "", output.ErrorString)

		err = handlers.Memory.Clear()
		require.Nil(t, err)
	}
}
