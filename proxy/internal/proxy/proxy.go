package proxy

import (
	"io"
	"net/http"
)

func ProxyRequest(w http.ResponseWriter, r *http.Request, plBearer string) {
	upstreamURL := "https://test.com" + r.URL.Path //TODO: replace url loaded from env

	req, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plBearer = "Bearer " + plBearer

	req.Header = r.Header.Clone()
	req.Header.Del("X-Api-Key")
	req.Header.Set("Authorization", plBearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, value := range v {
			w.Header().Add(k, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
