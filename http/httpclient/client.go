package httpclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

	_ "github.com/JoveYu/zgo/http/httplog/patch"
)

var DefaultClient = http.DefaultClient

var Do = DefaultClient.Do
var Get = DefaultClient.Get
var Post = DefaultClient.Post
var PostForm = DefaultClient.PostForm
var Head = DefaultClient.Head

func PostJson(url string, v interface{}) (*http.Response, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Post(url, "application/json", bytes.NewReader(data))
}

func PostXml(url string, v interface{}) (*http.Response, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Post(url, "application/xml", bytes.NewReader(data))
}
