package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type HealthcheckHandler struct {
	Section util.Section
}

type HealthCheckResponse struct {
	Reason string `json:"reason"`
	Status string `json:"status"`
}

func (t *HealthcheckHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var b HealthCheckResponse

	b.Status = "ok"
	b.Reason = "All checks passed"
	body, _ := json.Marshal(b)
	w.Write(body)
}
