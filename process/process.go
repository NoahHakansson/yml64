package process

import (
	"encoding/base64"
	"slices"

	"github.com/NoahHakansson/yml64/model"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// metadata keys to keep
var metadataKeys = []string{"name", "namespace"}

func cleanMetadata(metadata map[string]interface{}) {
	// Remove unnecessary metadata
	log.Debug("Removing unnecessary metadata")
	for k := range metadata {
		if !slices.Contains(metadataKeys, k) {
			log.Debugf("Removing metadata key: [%s]", k)
			delete(metadata, k)
		}
	}
	log.Debugf("Cleaned Metadata: %#v", metadata)
}

// EncodeDataProps base64 encodes the data property of a Kubernetes Secret
func EncodeDataProps(input []byte, f model.Flags) ([]byte, error) {

	ymlData := model.KubeSecret{}
	err := yaml.Unmarshal(input, &ymlData)
	if err != nil {
		return nil, err
	}

	// Remove unnecessary metadata if the metadata flag is not set
	if !f.Metadata {
		cleanMetadata(ymlData.Metadata)
	}

	// encode each data child property in base64
	for k, v := range ymlData.Data {
		ymlData.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}

	// Marshal the data back to YAML
	result, err := yaml.Marshal(ymlData)

	return result, err
}
