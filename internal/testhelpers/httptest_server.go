package testhelpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

func NewHttpTestServer() *httptest.Server {
	r := chi.NewRouter()

	channels := make(map[int64]*adzerk.Channel)

	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Post("/v1/channel", func(w http.ResponseWriter, r *http.Request) {
		var rb adzerk.CreateChannelJSONRequestBody
		if err := decodeJsonRequestBody(r, &rb); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		channelId := int64(1234)
		engine := "0"
		channels[channelId] = &adzerk.Channel{Id: int32(channelId), Engine: &engine}

		writeJsonMarshalable(w, applyChannelCreate(channels[channelId], rb))
	})

	r.Get("/v1/channel/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		channelId, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJsonMarshalable(w, channels[channelId])
	})

	r.Put("/v1/channel/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		channelId, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var rb adzerk.UpdateChannelJSONRequestBody
		if err := decodeJsonRequestBody(r, &rb); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJsonMarshalable(w, applyChannelUpdate(channels[channelId], rb))
	})

	r.Get("/v1/channel/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		channelId, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		delete(channels, channelId)
		w.Write([]byte(""))
	})

	return httptest.NewServer(r)
}

func decodeJsonRequestBody(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJsonMarshalable(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func applyChannelCreate(channel *adzerk.Channel, update adzerk.CreateChannelJSONRequestBody) *adzerk.Channel {
	channel.Title = update.Title
	if update.AdTypes != nil {
		channel.AdTypes = *update.AdTypes
	}
	return channel
}

func applyChannelUpdate(channel *adzerk.Channel, update adzerk.UpdateChannelJSONRequestBody) *adzerk.Channel {
	channel.Id = update.Id
	channel.Title = update.Title
	if update.AdTypes != nil {
		channel.AdTypes = *update.AdTypes
	}
	return channel
}
