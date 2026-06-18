package proxy

import (
	"io"
	"net/http"
	"strings"
)

func ProxyRequest(w http.ResponseWriter, r *http.Request, plBearer string, baseURL string) {
	proxyPath := strings.TrimPrefix(r.URL.Path, "/api/pl")
	if proxyPath == "" {
		proxyPath = "/"
	}
	if !strings.HasPrefix(proxyPath, "/") {
		proxyPath = "/" + proxyPath
	}

	upstreamURL := strings.TrimRight(baseURL, "/") + proxyPath
	if r.URL.RawQuery != "" {
		upstreamURL += "?" + r.URL.RawQuery
	}

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
