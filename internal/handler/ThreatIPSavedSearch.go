package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type ThreatIPSavedSearch struct {
	Section util.SectionAPI
}

type ThreatIPPayload struct {
	Results []map[string]interface{} `json:"results"`
}

func (t *ThreatIPSavedSearch) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var p ThreatIPPayload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var ips util.SectionIpRestrictionSchema

	if len(p.Results) > 0 {
		for _, r := range p.Results {
			ips.IpBlacklist = append(ips.IpBlacklist, r["message.request.remote_addr"].(string))
		}
	} else {
		ips = util.SectionIpRestrictionSchema{IpBlacklist: []string{}}
	}

	// go t.Section.AddIPBlocklist(ips)
	w.WriteHeader(http.StatusOK)
}
