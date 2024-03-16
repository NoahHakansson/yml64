package model

type KubeSecret struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Data       map[string]string      `yaml:"data"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Type       string                 `yaml:"type"`
}

type Flags struct {
	Decode   bool
	Metadata bool
	Inplace  bool
	Output   string
}
