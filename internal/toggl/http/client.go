package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"time"

	"github.com/philipf/gt/internal/cache"
	"github.com/philipf/gt/internal/settings"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/viper"
)

type TogglClientGateway struct {
}

func (t *TogglClientGateway) Get(filter *toggl.GetClientOpts) (toggl.TogglClients, error) {
	cache := cache.JsonFileCache[toggl.TogglClients, toggl.GetClientOpts]{}

	cacheDir, err := settings.GetGtConfigPath()
	if err != nil {
		return nil, err
	}

	cacheFilePath := path.Join(cacheDir, "toggl-clients.json")
	cacheMaxAge := time.Duration(time.Minute * 5)

	if filter == nil {
		filter = &toggl.GetClientOpts{}
	}

	cr, err := cache.Get(*filter, cacheFilePath, cacheMaxAge)

	if err != nil {
		log.Println("Cache miss")
	} else {
		log.Println("Cache hit")
		return *cr, nil
	}

	r, err := t.getFromToggl(filter)
	if err != nil {
		return nil, err
	}

	err = cache.Save(*filter, cacheFilePath, &r, cacheMaxAge)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (t *TogglClientGateway) getFromToggl(filter *toggl.GetClientOpts) (toggl.TogglClients, error) {
	uri, err := getApiClientsListUri()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if filter != nil {
		q := u.Query()
		q.Add("name", filter.Name)
		u.RawQuery = q.Encode()
	}

	log.Println("URI:", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r toggl.TogglClients
	err = json.Unmarshal(body, &r)

	if err != nil {
		return nil, err
	}

	sort.Sort(toggl.ClientsByName(r))

	return r, nil
}
