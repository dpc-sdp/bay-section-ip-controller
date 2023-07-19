package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/sectionio"
	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type ThreatIPSavedSearch struct {
	Section util.Section
}

type ThreatIPPayload struct {
	Results string `json:"results"`
}

type ThreatIPResult struct {
	Count      int    `json:"Count"`
	RemoteAddr string `json:"message.request.remote_addr"`
}

func (t *ThreatIPSavedSearch) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var p ThreatIPPayload
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	var ips sectionio.IpRestrictions
	var results []ThreatIPResult
	err = json.Unmarshal([]byte(p.Results), &results)

	if len(results) > 0 {
		for _, r := range results {
			ips.IpBlacklist = append(ips.IpBlacklist, r.RemoteAddr)
		}
	} else {
		ips = sectionio.IpRestrictions{IpBlacklist: []string{}}
	}

	go t.Section.AddIpRestrictionsToAllApplications(ips)
	w.WriteHeader(http.StatusOK)
}
