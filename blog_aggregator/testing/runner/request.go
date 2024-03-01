package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/testing/utils"
)

type RequestParameters struct {
	URL           string
	Method        string
	Endpoint      string
	PathParameter string
	ContentType   string
	Headers       map[string]string
	Body          interface{}
}

func (p *RequestParameters) ConstructRequest(caller string) *http.Request {
	var buf bytes.Buffer
	if p.Body != nil {
		err := json.NewEncoder(&buf).Encode(p.Body)
		if err != nil {
			utils.EncodingError(err, caller)
		}
	}

	var url string
	if url = p.URL + p.Endpoint; p.PathParameter != "" {
		url = url + "/" + p.PathParameter
	}
	// log.Println(url)

	req, err := http.NewRequest(p.Method, url, &buf)
	if err != nil {
		utils.RequestFormationError(err, caller)
	}

	for key, val := range p.Headers {
		req.Header.Set(key, val)
	}

	if p.ContentType != "" {
		req.Header.Set("Content-Type", p.ContentType)
	}

	return req
}

func SendRequest(
	rp RequestParameters,
	caller string,
	c *http.Client,
) ([]byte, int) {
	req := rp.ConstructRequest(caller)
	resp, err := c.Do(req)
	if err != nil {
		utils.BadResponseError(err, caller)
	}

	body, err := io.ReadAll(resp.Body)
	// log.Println(string(body))
	if err != nil {
		utils.BodyResponseError(err, caller)
	}
	return body, resp.StatusCode
}
