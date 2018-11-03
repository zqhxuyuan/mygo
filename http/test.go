package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	m := make(map[string]string)
	a1("http://localhost:8080", "", m)
}

// POST 带参数 + 带文件 ... 请求
func a1(u string, ff string, params map[string]string) ([]byte, error) {
	f, err := os.Open(ff)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	file, err := w.CreateFormFile("File", filepath.Base(ff))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, f)
	content_type := w.FormDataContentType()
	for k, v := range params {
		w.WriteField(k, v)
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", content_type)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GET 带参数请求 ...
func a(u string, params map[string]string) ([]byte, error) {
	p, _ := url.Parse(u)
	q := p.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	p.RawQuery = q.Encode()
	resp, err := http.Get(p.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// POST 带参数请求 ...
func a2(u string, params map[string]string) ([]byte, error) {
	values := make(url.Values)
	for k, v := range params {
		values.Add(k, v)
	}
	resp, err := http.PostForm(u, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
