package http

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	//"context"
)

func Get(urlstr string) ([]byte, error) {
	return Request(urlstr, nil, nil, nil, 3*time.Second)
}

func Post(url string, post []byte) ([]byte, error) {
	return Request(url, nil, nil, post, 3*time.Second)
}

func Request(urlstr string, params, headers map[string]string, post []byte, timeout time.Duration) ([]byte, error) {
	if len(params) > 0 {
		up := url.Values{}
		for k, v := range params {
			up.Add(k, v)
		}
		urlstr += "?" + up.Encode()
	}
	//fmt.Printf("url encode: %s\n", urlstr)

	//ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	client := &http.Client{
		Timeout: timeout,
	}
	var req *http.Request
	var err error
	if len(post) > 0 {
		req, err = http.NewRequest("POST", urlstr, bytes.NewReader(post))
	} else {
		req, err = http.NewRequest("GET", urlstr, nil)
	}
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s return %d", urlstr, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
