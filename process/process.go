package process

import (
	"encoding/base64"

	"github.com/NoahHakansson/yml64/model"
	"gopkg.in/yaml.v3"
)

// EncodeDataProps base64 encodes the data property of a Kubernetes Secret
func EncodeDataProps(input []byte) ([]byte, error) {

	ymlData := model.KubeSecret{}
	err := yaml.Unmarshal(input, &ymlData)
	if err != nil {
		return nil, err
	}
	// encode each data child property in base64
	for k, v := range ymlData.Data {
		ymlData.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}
	result, err := yaml.Marshal(ymlData)

	return result, err
}
