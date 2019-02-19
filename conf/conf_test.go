package conf

import "io/ioutil"
import "testing"
import "github.com/JoveYu/zgo/log"

type Config struct {
	Num      int               `json:"num" yaml:"num" toml:"num"`
	Text     string            `json:"text" yaml:"text" toml:"text"`
	NumList  []int             `json:"num_list" yaml:"num_list" toml:"num_list"`
	TextDict map[string]string `json:"text_dict" yaml:"text_dict" toml:"text_dict"`
}

func TestJson(t *testing.T) {
	log.Install("stdout")
	data := `{"num":123, "text":"hello", "num_list": [1, 2], "text_dict":{"key1": "value1", "key2": "value2"}}`
	ioutil.WriteFile("/tmp/zgo_conf.json", []byte(data), 0755)

	c := Config{}

	err := Install("/tmp/zgo_conf.json", &c)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("json %+v", c)
}

func TestYaml(t *testing.T) {
	log.Install("stdout")
	data := `
---
num: 123
num_list:
- 1
- 2
text: hello
text_dict:
  key1: value1
  key2: value2
`
	ioutil.WriteFile("/tmp/zgo_conf.yaml", []byte(data), 0755)

	c := Config{}

	err := Install("/tmp/zgo_conf.yaml", &c)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("yaml %+v", c)
}

func TestToml(t *testing.T) {
	log.Install("stdout")
	data := `
num = 123
text = "hello"
num_list = [ 1, 2,]
[text_dict]
key1 = "value1"
key2 = "value2"
`
	ioutil.WriteFile("/tmp/zgo_conf.toml", []byte(data), 0755)

	c := Config{}

	err := Install("/tmp/zgo_conf.toml", &c)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("toml %+v", c)
}
