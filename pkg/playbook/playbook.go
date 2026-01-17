package playbook

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Playbook struct {
	Plays []*Play
}

type Play struct {
	Name       string                 `yaml:"name"`
	Hosts      string                 `yaml:"hosts"`
	Become     bool                   `yaml:"become,omitempty"`
	BecomeUser string                 `yaml:"become_user,omitempty"`
	Vars       map[string]interface{} `yaml:"vars,omitempty"`
	Tasks      []*Task                `yaml:"tasks,omitempty"`
	Handlers   []*Task                `yaml:"handlers,omitempty"`
	Roles      []string               `yaml:"roles,omitempty"`
	PreTasks   []*Task                `yaml:"pre_tasks,omitempty"`
	PostTasks  []*Task                `yaml:"post_tasks,omitempty"`
	Tags       []string               `yaml:"tags,omitempty"`
}

func LoadPlaybook(filename string) (*Playbook, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var plays []*Play
	if err := yaml.Unmarshal(data, &plays); err != nil {
		return nil, err
	}

	return &Playbook{Plays: plays}, nil
}
