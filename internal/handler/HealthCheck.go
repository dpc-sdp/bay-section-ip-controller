package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type HealthCheck struct {
	Section util.SectionAPI
}

type HealthCheckResponse struct {
	Reason string `json:"reason"`
	Status string `json:"status"`
}

func (t *HealthCheck) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	a := t.Section.GetAccount()

	var b HealthCheckResponse

	if a.Name == "" {
		b.Reason = "Unable to fetch account details"
		b.Status = "fail"
		body, _ := json.Marshal(b)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	b.Status = "ok"
	b.Reason = "All checks passed"
	body, _ := json.Marshal(b)
	w.Write(body)
}
