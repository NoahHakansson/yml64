package model

// apiVersion: v1
// data:
//   something: aHR0
// kind: Secret
// metadata:
//   name: some-api-config
//   namespace: project-dev
// type: Opaque

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type KubeSecret struct {
	ApiVersion string            `yaml:"apiVersion"`
	Data       map[string]string `yaml:"data"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Type       string            `yaml:"type"`
}
