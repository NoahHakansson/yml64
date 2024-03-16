package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanMetadata(t *testing.T) {
	assert := assert.New(t)

	whitelist := []string{"name", "namespace"}
	metadata := map[string]interface{}{"name": "some-name", "namespace": "my-namespace", "creationTimestamp": "2021-08-05T14:22:34Z"}
	cleanMetadata(metadata, whitelist)

	assert.Equal(map[string]interface{}{"name": "some-name", "namespace": "my-namespace"}, metadata)
}

func TestEncodeDataProps(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`
