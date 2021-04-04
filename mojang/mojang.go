package mojang

import (
	"bytes"
	"encoding/json"
	"github.com/rotisserie/eris"
	"io"
	"net/http"
)

// https://wiki.vg/Mojang_AP
type (
	ApiClient struct {
		Client http.Client
	}

	PlayerUuidRequest  []string
	PlayerUuidResponse struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
)

func (c *ApiClient) GetPlayerUuid(request PlayerUuidRequest) ([]PlayerUuidResponse, error) {
	req, err := json.Marshal(request)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal request")
	}
	resp, err := c.Client.Post("https://api.mojang.com/profiles/minecraft", "application/json", bytes.NewReader(req))
	if err != nil {
		return nil, eris.Wrap(err, "failed to hit Mojang profiles endpoint")
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read response body")
	}

	var respObj []PlayerUuidResponse
	err = json.Unmarshal(respBytes, &respObj)
	return respObj, err
}
