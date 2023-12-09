package mail

import (
	"errors"
	"io"
	"net/http"
)

func githubGetResBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+githubToken)
	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not 200")
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not 200")
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
