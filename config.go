package config

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config map[string]interface{}

func New(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v := map[interface{}]interface{}{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	c := Config{}
	err = c.visit(v, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c Config) Run(symbol string, ss ...*string) (interface{}, error) {
	s := c[symbol]
	if s == nil {
		return nil, fmt.Errorf("%v not found", symbol)
	}
	return s.(func(map[string]interface{}, ...*string) (interface{}, error))(c, ss...)
}

func (c Config) visit(m map[interface{}]interface{}, path []string) error {
	for k, v := range m {
		if mm, ok := v.(map[interface{}]interface{}); ok {
			err := c.visit(mm, append(path, k.(string)))
			if err != nil {
				return err
			}
		} else if k == "PluginFile" {
			symbol := path[len(path)-1:][0]
			p, err := plugin.Open(v.(string))
			if err != nil {
				return err
			}
			s, err := p.Lookup(symbol)
			if err != nil {
				return err
			}
			c[strings.Join(path, "/")] = s
		} else {
			c[strings.Join(append(path, k.(string)), "/")] = v
		}
	}
	return nil
}
