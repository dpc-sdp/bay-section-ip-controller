package handler

import (
	"net/http"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

type ThreatIPList struct {
	Section util.Section
}

func (t *ThreatIPList) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//jsonData, err := json.Marshal(t.Section.IPTracker.GetAll())
	jsonData, err := func() ([]byte, error) { return []byte("{}"), nil }()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Section.Logger.Error().Err(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
