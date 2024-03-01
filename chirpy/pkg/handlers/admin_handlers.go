package handlers

import (
	"fmt"
	"net/http"
)

func (m *Repository) AdminHitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := `
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
  `
	bodyBytes := []byte(fmt.Sprintf(body, m.App.GetHits()))
	w.Write(bodyBytes)
}
