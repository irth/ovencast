package ome

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	Addr  string
	Token string
}

func NewAPI(addr string, token string) (*API, error) {
	ome := API{
		Addr:  addr,
		Token: token,
	}

	_, err := ome.Stats()
	if err != nil {
		return nil, fmt.Errorf("OMEAPI.Stats: %w", err)
	}

	return &ome, nil
}

var ErrAPIFailed = fmt.Errorf("status code")

// We want errors.Is(NotFoundErr, APIErr) to be true. Other errors are generated
// ad-hoc, and still can be checked using errors.Is.

// Basically, errors other than 404 usually signify a problem, but 404 can just
// mean that, for example, a stream hasn't been started yet, so this allows
// handling that in an easier way.

var ErrNotFound = fmt.Errorf("%w 404", ErrAPIFailed)

type Response[T any] struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Response   T      `json:"response"`
}

type Metrics struct {
	CreatedTime            time.Time `json:"createdTime"`
	LastUpdatedTime        time.Time `json:"lastUpdatedTime"`
	TotalBytesIn           int64     `json:"totalBytesIn"`
	TotalBytesOut          int64     `json:"totalBytesOut"`
	TotalConnectiosn       int64     `json:"totalConnections"`
	MaxTotalConnections    int64     `json:"maxTotalConnections"`
	MaxTotalConnectionTime time.Time `json:"maxTotalConnectionTime"`
	LastRecvTime           time.Time `json:"lastRecvTime"`
	LastSentTime           time.Time `json:"lastSentTime"`
}

func (o *API) request(endpointParts []string, data interface{}) error {
	endpoint, err := url.JoinPath(o.Addr, endpointParts...)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	token := base64.StdEncoding.EncodeToString([]byte(o.Token))
	header := fmt.Sprintf("Basic %s", token)

	req.Header.Add("Authorization", header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}

	var response Response[json.RawMessage]

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("JSON decode: %w", err)
	}

	if response.StatusCode == 404 {
		return fmt.Errorf("%w: %s", ErrNotFound, response.Message)
	}

	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		return fmt.Errorf("%w %d: %s", ErrAPIFailed, response.StatusCode, response.Message)
	}

	if data != nil {
		err = json.Unmarshal(response.Response, data)
		if err != nil {
			return fmt.Errorf("JSON decode (data): %w", err)
		}
		return nil
	}

	return nil
}

func (o *API) stats(endpoint []string) (*Metrics, error) {
	var response Metrics
	err := o.request(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("OMEAPI.request: %w", err)
	}

	return &response, nil
}

func (o *API) Stats() (*Metrics, error) {
	return o.stats([]string{"v1/stats/current"})
}

func (o *API) VHostStats(vhost string) (*Metrics, error) {
	return o.stats([]string{"v1/stats/current", "vhosts", vhost})
}

func (o *API) AppStats(vhost string, app string) (*Metrics, error) {
	return o.stats([]string{"v1/stats/current", "vhosts", vhost, "apps", app})
}

func (o *API) StreamStats(vhost string, app string, stream string) (*Metrics, error) {
	return o.stats([]string{"v1/stats/current", "vhosts", vhost, "apps", app, "streams", stream})
}

func (o *API) StreamExists(vhost string, app string, stream string) (bool, error) {
	_, err := o.StreamStats(vhost, app, stream)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("OMEAPI.StreamStats: %w", err)
	}

	return true, nil
}
