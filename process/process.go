package process

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/NoahHakansson/yml64/model"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// metadata keys to keep
var metadataKeyWhitelist = []string{"name", "namespace"}

func cleanMetadata(metadata map[string]interface{}, whitelist []string) {
	// Remove unnecessary metadata
	log.Debug("Removing unnecessary metadata")
	for k := range metadata {
		if !slices.Contains(whitelist, k) {
			log.Debugf("Removing metadata key: [%s]", k)
			delete(metadata, k)
		}
	}
	log.Debugf("Cleaned Metadata: %#v", metadata)
}

// EncodeDataProps base64 encodes the data properties of a Kubernetes Secret
func EncodeDataProps(input []byte, f model.Flags) ([]byte, error) {
	ymlData := model.KubeSecret{}
	err := yaml.Unmarshal(input, &ymlData)
	if err != nil {
		return nil, err
	}

	if err := validateInput(ymlData); err != nil {
		return nil, err
	}

	// Remove unnecessary metadata if the metadata flag is not set
	if !f.Metadata {
		cleanMetadata(ymlData.Metadata, metadataKeyWhitelist)
	}

	// encode each data child property in base64
	for k, v := range ymlData.Data {
		ymlData.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}

	return yaml.Marshal(ymlData)
}

// DecodeDataProps base64 decodes the data properties of a Kubernetes Secret
func DecodeDataProps(input []byte, f model.Flags) ([]byte, error) {
	ymlData := model.KubeSecret{}
	err := yaml.Unmarshal(input, &ymlData)
	if err != nil {
		return nil, err
	}

	if err := validateInput(ymlData); err != nil {
		return nil, err
	}

	// Remove unnecessary metadata if the metadata flag is not set
	if !f.Metadata {
		cleanMetadata(ymlData.Metadata, metadataKeyWhitelist)
	}

	// decode each data child property from base64
	for k, v := range ymlData.Data {
		res, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		ymlData.Data[k] = string(res)
	}

	return yaml.Marshal(ymlData)
}

// validateInput validates the input against the KubeSecret struct
// and returns an error if the input is not a valid KubeSecret
func validateInput(s model.KubeSecret) error {
	err := validate.Struct(s)
	if err != nil {
		verr := err.(validator.ValidationErrors)
		log.Debugf("Validation errors: %v", verr)
		missingKeys := []string{}
		for _, e := range verr {
			key := lowercaseFirst(e.Field())
			missingKeys = append(missingKeys, key)
		}
		return fmt.Errorf("input is not a valid kube secret: Missing keys: [%v]", strings.Join(missingKeys, ", "))
	}
	return nil
}

func lowercaseFirst(s string) string {
	if len(s) < 1 {
		return s
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
