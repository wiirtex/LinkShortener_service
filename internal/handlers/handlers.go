package handlers

import (
	"encoding/json"
	"net/http"
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/pkg/helpers"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var Memory memory.Memory

func sendSuccess(w http.ResponseWriter, logger *logrus.Entry, output []byte) {
	_, errSending := w.Write(output)
	if errSending != nil {
		logger.Error("Error sending non-error answer: ", errSending)
	}
}

func sendError(w http.ResponseWriter, logger *logrus.Entry, output []byte, code int) {
	w.WriteHeader(code)
	_, errSending := w.Write(output)
	if errSending != nil {
		logger.Error("Error sending error: ", errSending)
	}
}

func ManageLongLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := Memory
	logger := r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry)
	input := r.Context().Value(helpers.ContextInputKey).(helpers.Input)
	authorIP := r.Context().Value(helpers.ContextAuthorIpKey).(string)

	mResponse, err := db.Storage.GetEntryByLongLink(input.Link)
	if err != nil {
		if err.Error() == "No such entry" {
			mResponse = memory.MemoryEntry{
				Author:    authorIP,
				CreatedAt: time.Now().UTC(),
				LongLink:  input.Link,
			}
			mResponse.GenerateUniqueShortLink(db)

			err = db.Storage.AddEntry(mResponse)
			var output helpers.Output
			if err != nil {
				output = helpers.Output{
					InputLink:   input.Link,
					OutputLink:  helpers.GetConfig().ShortLinkBase + mResponse.ShortId,
					Error:       true,
					ErrorString: err.Error(),
				}
				bytes, _ := json.Marshal(output)
				sendError(w, logger, bytes, http.StatusInternalServerError)

			} else {
				output = helpers.Output{
					InputLink:  input.Link,
					OutputLink: helpers.GetConfig().ShortLinkBase + mResponse.ShortId,
					Error:      false,
				}
				bytes, err := json.Marshal(output)
				if err != nil {
					sendError(w, logger, bytes, http.StatusInternalServerError)
				} else {
					sendSuccess(w, logger, bytes)
				}
			}

		} else {
			output := helpers.Output{
				InputLink:   input.Link,
				OutputLink:  "",
				Error:       true,
				ErrorString: err.Error(),
			}
			bytes, _ := json.Marshal(output)
			sendError(w, logger, bytes, http.StatusInternalServerError)
		}
	} else {
		output := helpers.Output{
			InputLink:   input.Link,
			OutputLink:  helpers.GetConfig().ShortLinkBase + mResponse.ShortId,
			Error:       false,
			ErrorString: "",
		}
		bytes, err := json.Marshal(output)
		if err != nil {
			sendError(w, logger, bytes, http.StatusInternalServerError)
		} else {
			sendSuccess(w, logger, bytes)
		}
	}
}

func ManageShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger := r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry)
	db := Memory.Storage
	input := r.Context().Value(helpers.ContextInputKey).(helpers.Input)
	var output helpers.Output
	shortLink := input.Link

	// Check if the provided link is link, that is created on the server
	if strings.HasPrefix(shortLink, helpers.GetConfig().ShortLinkBase) {
		shortLink = shortLink[len(helpers.GetConfig().ShortLinkBase):]
	} else {
		output = helpers.Output{
			InputLink:   shortLink,
			OutputLink:  "",
			Error:       true,
			ErrorString: "Not acceptible input link",
		}
		bytes, err := json.Marshal(output)
		if err != nil {
			sendError(w, logger, bytes, http.StatusInternalServerError)
		} else {
			sendError(w, logger, bytes, http.StatusBadRequest)
		}
		return
	}

	if mResponse, err := db.GetEntryByShortId(shortLink); err != nil {
		if err.Error() == "No such entry" {
			output = helpers.Output{
				InputLink:   helpers.GetConfig().ShortLinkBase + shortLink,
				OutputLink:  "",
				Error:       true,
				ErrorString: "Not Found",
			}
			bytes, err := json.Marshal(output)
			if err != nil {
				sendError(w, logger, bytes, http.StatusInternalServerError)
			} else {
				sendError(w, logger, bytes, http.StatusNotFound)
			}
		} else {
			output = helpers.Output{
				InputLink:   helpers.GetConfig().ShortLinkBase + shortLink,
				OutputLink:  "",
				Error:       true,
				ErrorString: "Internal error",
			}
			logger.Error("Error in db: ", err)
			bytes, _ := json.Marshal(output)
			sendError(w, logger, bytes, http.StatusInternalServerError)
		}
	} else {
		output = helpers.Output{
			InputLink:   helpers.GetConfig().ShortLinkBase + mResponse.ShortId,
			OutputLink:  mResponse.LongLink,
			Error:       false,
			ErrorString: "",
		}
		bytes, err := json.Marshal(output)
		if err != nil {
			sendError(w, logger, bytes, http.StatusInternalServerError)
		} else {
			sendSuccess(w, logger, bytes)
		}
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {

	logger := r.Context().Value(helpers.ContextLoggerKey).(*logrus.Entry)
	db := Memory.Storage
	shortId := chi.URLParam(r, "shortId")

	if memoryResponse, err := db.GetEntryByShortId(shortId); err == nil {
		http.Redirect(w, r, memoryResponse.LongLink, http.StatusSeeOther)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		output := helpers.Output{
			InputLink:   "",
			OutputLink:  "",
			Error:       true,
			ErrorString: "Not Found",
		}
		bytes, err := json.Marshal(output)
		if err != nil {
			sendError(w, logger, bytes, http.StatusInternalServerError)
		} else {
			sendError(w, logger, bytes, http.StatusNotFound)
		}
	}
}
