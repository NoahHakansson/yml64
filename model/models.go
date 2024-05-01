package model

type KubeSecret struct {
	APIVersion string                 `yaml:"apiVersion" validate:"required"`
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

type Input struct {
	Exists bool
	Path   string
}

func FieldToKey(s string) string {
	switch s {
	case "APIVersion":
		return "apiVersion"
	case "Data":
		return "data"
	case "Kind":
		return "kind"
	case "Metadata":
		return "metadata"
	case "Type":
		return "type"
	default:
		return "unknown field"
	}
}
