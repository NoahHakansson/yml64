package model

type KubeSecret struct {
	ApiVersion string                 `yaml:"apiVersion" validate:"required"`
	Data       map[string]string      `yaml:"data" validate:"required"`
	Kind       string                 `yaml:"kind" validate:"required"`
	Metadata   map[string]interface{} `yaml:"metadata" validate:"required"`
	Type       string                 `yaml:"type" validate:"required"`
}

type Flags struct {
	Decode   bool
	Metadata bool
	Inplace  bool
	Output   string
}
