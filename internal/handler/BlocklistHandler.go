package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	sectionio "github.com/dpc-sdp/go-section-io"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type BlocklistWebhookPayload struct {
	Results []BlocklistWebhookItem `json:"results"`
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

	err = json.Unmarshal(body, &p)
	if err != nil {
		t.Section.Logger.Error().Err(err).Msg("error parsing request body")
		t.Section.Logger.Debug().Msg(string(body))
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	var ips sectionio.IpRestrictions
	blocklist := make([]string, 0)
	for _, r := range p.Results {
		blocklist = append(blocklist, r.Cidr)
	}
	ips = sectionio.IpRestrictions{
		IpBlacklist: blocklist,
	}

	go t.Section.AddIpRestrictionsToAllApplications(ips)
	w.WriteHeader(http.StatusOK)
}
