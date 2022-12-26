package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
)

type WebhookRequest struct {
	Client struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"client"`

	Request struct {
		Direction string `json:"direction"`
		Protocol  string `json:"protocol"`
		Status    string `json:"status"`
		URL       string `json:"url"`
		Time      string `json:"time"`
	} `json:"request"`
}

type WebhookResponse struct {
	Allowed  bool   `json:"allowed"`
	NewURL   string `json:"new_url,omitempty"`
	Lifetime int64  `json:"lifetime,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func (a *API) AdmissionWebhook(w http.ResponseWriter, r *http.Request) {
	var req WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("parsing JSON body failed: %s", err)
		return
	}

	res := json.NewEncoder(w)

	if req.Request.Status == "closing" {
		res.Encode(map[string]bool{})
		return
	}

	if req.Request.Direction == "outgoing" {
		// allow everyone to play the stream
		res.Encode(WebhookResponse{Allowed: true})
		return
	}

	a.Config.RLock()
	defer a.Config.RUnlock()

	if len(a.StreamKey) == 0 {
		res.Encode(WebhookResponse{Allowed: false})
		return
	}

	u, err := url.Parse(req.Request.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("url parsing failed: %s", err)
		return
	}

	_, streamKey := path.Split(u.Path)

	if streamKey != a.StreamKey {
		res.Encode(WebhookResponse{Allowed: false})
	}

	u.Path = "/live/stream"
	res.Encode(WebhookResponse{
		Allowed: true,
		NewURL:  u.String(),
	})
}
