package http

import "io/ioutil"
import "strings"
import "testing"
import "net/http"
import "net/url"
import "encoding/xml"
import "github.com/JoveYu/zgo/log"

func TestClient(t *testing.T) {
	log.Install("stdout")
	client := NewClient()
	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	if err != nil {
		log.Error(err)
	}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	resp, _ = client.Get("http://httpbin.org/get")
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	resp, _ = client.Head("http://httpbin.org/get")
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	resp, _ = client.Post("http://httpbin.org/post", "application/json", strings.NewReader("{\"key\":\"value\"}"))
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	resp, _ = client.PostForm("http://httpbin.org/post", url.Values{
		"key": []string{"1", "2"},
	})
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	resp, _ = client.PostJson("http://httpbin.org/post", map[string]int{"key": 1})
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

	type User struct {
		XMLName xml.Name `xml:"xml"`
		Name    string   `xml:"name"`
		Id      int      `xml:"id,attr"`
	}
	user := User{
		Name: "test",
		Id:   1,
	}
	resp, _ = client.PostXml("http://httpbin.org/post", user)
	body, _ = ioutil.ReadAll(resp.Body)
	log.Debug("%s", body)

}
