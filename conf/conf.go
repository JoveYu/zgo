package conf

import (
    "errors"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "github.com/go-yaml/yaml"
    "github.com/BurntSushi/toml"
)

type formater struct {
    Marshal func(v interface{}) ([]byte, error)
    Unmarshal func(data []byte, v interface{}) error
}

var (
    formaters = map[string]formater{
        "json": formater{
            Marshal: json.Marshal,
            Unmarshal: json.Unmarshal,
        },
        "yaml": formater{
            Marshal: yaml.Marshal,
            Unmarshal: yaml.Unmarshal,
        },
        "yml": formater{
            Marshal: yaml.Marshal,
            Unmarshal: yaml.Unmarshal,
        },
        "toml": formater{
            Marshal: func(v interface{}) ([]byte, error) {
                b := bytes.Buffer{}
                err := toml.NewEncoder(&b).Encode(v)
                return b.Bytes(), err
            },
            Unmarshal: toml.Unmarshal,
        },
    }
)

func Install(path string, v interface{}) error {

    ext := ""
    // get extension
    for i:=len(path) - 1; i>=0; i-- {
        if path[i] == '.' {
            ext = path[i+1:]
            break
        }
    }
    if ext == "" {
        return errors.New("invalid file extension")
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }

    if formater, ok := formaters[ext]; ok {
        err := formater.Unmarshal(data, v)
        if err != nil {
            return err
        }

    } else {
        return errors.New("no support extension")
    }
    return nil
}


