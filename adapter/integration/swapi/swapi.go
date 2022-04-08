package swapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kaduartur/go-planet/core/port"
	"github.com/kaduartur/go-planet/pkg/log"
)

var (
	ErrPlanetNotFound = errors.New("the planet was not found")
	ErrFilmNotFound   = errors.New("the film was not found")
)

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string, timeout time.Duration) port.SwapiClient {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		url:    url,
		client: client,
	}
}

func (c *Client) FindPlanetByName(ctx context.Context, name string) (*port.SwapiPlanet, error) {
	planetURL := fmt.Sprintf("%s/api/planets", c.url)
	req, err := http.NewRequest(http.MethodGet, planetURL, nil)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("search", name)
	req.URL.RawQuery = params.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		log.Error(ctx, "error to find planet by name ", err, log.Event{"name": name})
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Error(ctx, "unknown error ", nil)
		return nil, fmt.Errorf("an unexpected error has occurred. STATUS=%s", res.Status)
	}

	defer res.Body.Close()
	var body port.SwapiResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	planet := body.Planets.First()
	if planet == nil {
		log.Error(ctx, "error not found planet ", nil, log.Event{"planet": name})
		return nil, ErrPlanetNotFound
	}

	return planet, err
}

func (c *Client) FindFilmByID(ctx context.Context, id string) (*port.SwapiFilm, error) {
	filmURL := fmt.Sprintf("%s/api/films/%s", c.url, id)
	req, err := http.NewRequest(http.MethodGet, filmURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		log.Error(ctx, "error to find film by id ", err, log.Event{"film_id": id})
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		log.Error(ctx, "error not found planet ", nil, log.Event{"film_id": id})
		return nil, ErrFilmNotFound
	}

	if res.StatusCode != http.StatusOK {
		log.Error(ctx, "unknown error ", nil)
		return nil, fmt.Errorf("an unexpected error has occurred. STATUS=%s", res.Status)
	}

	defer res.Body.Close()
	var body port.SwapiFilm
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, err
}
