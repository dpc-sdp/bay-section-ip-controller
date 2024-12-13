package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	sectionio "github.com/dpc-sdp/go-section-io"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

// Sumo posts the results payload as a json string inside a json object. This
// requires unmarshaling it twice.
type BlocklistWebhookPayload struct {
	Results string `json:"results"`
}
type BlocklistWebhookItem struct {
	Cidr string `json:"cidr"`
}

type BlocklistHandler struct {
	Section util.Section
}

func (t *BlocklistHandler) Serve(w http.ResponseWriter, r *http.Request) {
	t.Section.Logger.Debug().Msg("serving http request")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var p BlocklistWebhookPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Section.Logger.Error().Err(err).Msg("error reading request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	// Unmarshal the outer json wrapper.
	err = json.Unmarshal(body, &p)
	if err != nil {
		t.Section.Logger.Error().Err(err).Msg("error parsing request outer wrapper")
		t.Section.Logger.Debug().Msg(string(body))
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	// Unmarshal the inner json.
	results := make([]BlocklistWebhookItem, 0)
	err = json.Unmarshal([]byte(p.Results), &results)
	if err != nil {
		t.Section.Logger.Error().Err(err).Msg("error parsing request inner json")
		t.Section.Logger.Debug().Msg(p.Results)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	var ips sectionio.IpRestrictions
	blocklist := make([]string, 0)
	for _, r := range results {
		blocklist = append(blocklist, r.Cidr)
	}
	ips = sectionio.IpRestrictions{
		IpBlacklist: blocklist,
	}

	go t.Section.AddIpRestrictionsToAllApplications(ips)
	w.WriteHeader(http.StatusOK)
}
