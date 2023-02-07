package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
)

// WebhookRequest matches the format of OME admission webhook request
// See: https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks#request
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

// WebhookResponse matches the format of OME admission webhook response
// See: https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks#format-1
type WebhookResponse struct {
	Allowed  bool   `json:"allowed"`
	NewURL   string `json:"new_url,omitempty"`
	Lifetime int64  `json:"lifetime,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// AdmissionWebhook is called via /api/admission by OME to verify if a client is
// allowed to stream or watch the stream. We accept all viewers, but for
// streaming the stream key is checked.
func (a *API) AdmissionWebhook(w http.ResponseWriter, r *http.Request) {
	var req WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("parsing JSON body failed: %s", err)
		return
	}

	log.Printf("admission webhook: %s, %s", req.Request.Direction, req.Request.Status)

	res := json.NewEncoder(w)

	if req.Request.Direction == "outgoing" {
		// allow everyone to play the stream
		res.Encode(WebhookResponse{Allowed: true})
		return
	}

	if req.Request.Status == "closing" {
		res.Encode(map[string]bool{}) // "ta? to zajebiście"

		// Notify the state updater that a stream is about to end
		a.admissionWebhookSignal <- false
		return
	}

	// Extract the streamKey from the stream URL
	u, err := url.Parse(req.Request.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("url parsing failed: %s", err)
		return
	}

	_, streamKey := path.Split(u.Path)

	// Lock the config to check the streamKey
	a.Config.RLock()
	defer a.Config.RUnlock()

	if len(a.StreamKey) == 0 {
		res.Encode(WebhookResponse{Allowed: false})
		return
	}

	if streamKey != a.StreamKey {
		res.Encode(WebhookResponse{Allowed: false})
	}

	// Notify the state updater that a stream is about to start
	a.admissionWebhookSignal <- true

	// Redirect the stream to a known path that doesn't contain the streamKey -
	// otherwise the viewer's would need to know the key to watch the stream
	u.Path = "/live/stream"
	res.Encode(WebhookResponse{
		Allowed: true,
		NewURL:  u.String(),
	})
}
