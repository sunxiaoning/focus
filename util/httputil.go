package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ApplicationJSONUtf8Value = "application/json;charset=utf-8"
)

type SimpleHttpClient interface {
	PostJsonWithHeader(url string, headers map[string]string, reqParam interface{}, timeout time.Duration) (code int, respBody []byte, err error)
	PostJson(url string, reqParam interface{}, timeout time.Duration) (code int, respBody []byte, err error)
	GetWithHeader(url string, header map[string]string, timeout time.Duration) (code int, res []byte, err error)
	Get(url string, timeout time.Duration) (code int, res []byte, err error)
	DownloadFileWithHeader(url string, destPath string, header map[string][]string, timeout time.Duration) (err error)
	DownLoadFile(url string, destPath string, header map[string]string) (err error, timeout time.Duration)
}

var DefaultHttpClient = &defaultHttpClient{}

type defaultHttpClient struct {
}

func (c *defaultHttpClient) PostJsonWithHeader(url string, headers map[string]string, reqParam interface{}, timeout time.Duration) (code int, respBody []byte, err error) {
	var jsonBytes []byte
	if reqParam != nil {
		jsonBytes, err = json.Marshal(reqParam)
		if err != nil {
			return fasthttp.StatusBadRequest, nil, err
		}
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.SetRequestURI(url)
	req.SetBody(jsonBytes)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(ApplicationJSONUtf8Value)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if timeout > 0 {
		err = fasthttp.DoTimeout(req, resp, timeout)
	} else {
		err = fasthttp.Do(req, resp)
	}
	return resp.StatusCode(), resp.Body(), err
}

func (c *defaultHttpClient) PostJson(url string, reqParam interface{}, timeout time.Duration) (code int, res []byte, err error) {
	return DefaultHttpClient.PostJsonWithHeader(url, nil, reqParam, timeout)
}

func (c *defaultHttpClient) GetWithHeader(url string, headers map[string]string, timeout time.Duration) (code int, respBody []byte, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	if timeout > 0 {
		err = fasthttp.DoTimeout(req, res, timeout)
	} else {
		err = fasthttp.Do(req, res)
	}
	return res.StatusCode(), res.Body(), err
}

func (c *defaultHttpClient) Get(url string, timeout time.Duration) (code int, res []byte, err error) {
	return DefaultHttpClient.GetWithHeader(url, nil, timeout)
}

func convertHeaders(headers map[string][]string) map[string]string {
	if len(headers) == 0 {
		return nil
	}
	hds := make(map[string]string)
	for k, v := range headers {
		hds[k] = strings.Join(v, ",")
	}
	return hds
}

type httpReqError struct {
	error
	msg string
}

func (c *defaultHttpClient) DownloadFileWithHeader(url string, destPath string, header map[string][]string, timeout time.Duration) (err error) {
	code, content, err := DefaultHttpClient.GetWithHeader(url, convertHeaders(header), timeout)
	if code != http.StatusOK || len(content) <= 0 || err != nil {
		return httpReqError{err, fmt.Sprintf("req source error! code:%d, content-len:%s", code, content)}
	}
	fileHelper := DefaultFileHelper
	downloadFile, err := fileHelper.OpenFile(destPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer downloadFile.Close()
	buf := bufio.NewWriterSize(downloadFile, 4*1024*1024)
	_, err = io.Copy(buf, bytes.NewBuffer(content))
	buf.Flush()
	return err
}
func (c *defaultHttpClient) DownLoadFile(url string, destPath string, header map[string]string, timeout time.Duration) (err error) {
	return DefaultHttpClient.DownloadFileWithHeader(url, destPath, nil, timeout)
}
