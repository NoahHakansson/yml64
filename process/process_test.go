package process

import (
	"testing"

	"github.com/NoahHakansson/yml64/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

func TestCleanMetadata(t *testing.T) {
	assert := assert.New(t)

	whitelist := []string{"name", "namespace"}
	metadata := map[string]interface{}{"name": "some-name", "namespace": "my-namespace", "creationTimestamp": "2021-08-05T14:22:34Z"}
	cleanMetadata(metadata, whitelist)

	assert.Equal(map[string]interface{}{"name": "some-name", "namespace": "my-namespace"}, metadata)
}

type EncodeDecodeDataTestSuite struct {
	suite.Suite
	DecodedInputWithMetadata []byte
	EncodedInputWithMetadata []byte
	DecodedInputNoMetadata   []byte
	EncodedInputNoMetadata   []byte
}

func (suite *EncodeDecodeDataTestSuite) SetupTest() {
	// With Metadata
	// input
	suite.DecodedInputWithMetadata = []byte(`apiVersion: v1
data:
  DB: mydb
  JSON_DATA: |-
    {
    	"key1": "value1",
    	"key2": "value2"
    }
  PASS: 9x34jh8a3#cf9$asdf
  USER: admin
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"password":"OTMzNGpoOGEzI2NmOSRhZHNkZg==","username":"YWRtaW4","database":"bXlkYg=="},"kind":"Secret","metadata":{"annotations":{},"name":"my-secret","namespace":"my-namespace"}}
  creationTimestamp: "2021-08-05T14:22:34Z"
  name: my-secret
  namespace: my-namespace
  resourceVersion: "123456"
  uid: edc15e05-22c0-4e7a-9878-73f6c61f4b1f
type: Opaque
`)
	// base64 encoded input
	suite.EncodedInputWithMetadata = []byte(`apiVersion: v1
data:
    DB: bXlkYg==
    JSON_DATA: ewoJImtleTEiOiAidmFsdWUxIiwKCSJrZXkyIjogInZhbHVlMiIKfQ==
    PASS: OXgzNGpoOGEzI2NmOSRhc2Rm
    USER: YWRtaW4=
kind: Secret
metadata:
    annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
            {"apiVersion":"v1","data":{"password":"OTMzNGpoOGEzI2NmOSRhZHNkZg==","username":"YWRtaW4","database":"bXlkYg=="},"kind":"Secret","metadata":{"annotations":{},"name":"my-secret","namespace":"my-namespace"}}
    creationTimestamp: "2021-08-05T14:22:34Z"
    name: my-secret
    namespace: my-namespace
    resourceVersion: "123456"
    uid: edc15e05-22c0-4e7a-9878-73f6c61f4b1f
type: Opaque
`)

	// No Metadata
	// input (no metadata)
	suite.DecodedInputNoMetadata = []byte(`apiVersion: v1
data:
  DB: mydb
  JSON_DATA: |-
    {
    	"key1": "value1",
    	"key2": "value2"
    }
  PASS: 9x34jh8a3#cf9$asdf
  USER: admin
kind: Secret
metadata:
  name: my-secret
  namespace: my-namespace
type: Opaque
`)
	// base64 encoded input (no metadata)
	suite.EncodedInputNoMetadata = []byte(`apiVersion: v1
data:
    DB: bXlkYg==
    JSON_DATA: ewoJImtleTEiOiAidmFsdWUxIiwKCSJrZXkyIjogInZhbHVlMiIKfQ==
    PASS: OXgzNGpoOGEzI2NmOSRhc2Rm
    USER: YWRtaW4=
kind: Secret
metadata:
    name: my-secret
    namespace: my-namespace
type: Opaque
`)
}

// TestEncodeDataPropsNoMetadata() tests the EncodeDataProps function
// should return a byte slice with the data properties base64 encoded
// no metadata flag set, so metadata should be removed
func (suite *EncodeDecodeDataTestSuite) TestEncodeDataPropsNoMetadata() {
	assert := assert.New(suite.T())

	flags := model.Flags{}
	output, err := EncodeDataProps(suite.DecodedInputWithMetadata, flags)
	assert.Nil(err)

	// unmarshal the expected and actual output to interface{} to make comparison easier
	var expected interface{}
	err = yaml.Unmarshal(suite.EncodedInputNoMetadata, &expected)
	assert.Nil(err)

	var actual interface{}
	err = yaml.Unmarshal(output, &actual)
	assert.Nil(err)

	assert.Equal(expected, actual)
}

// TestDecodeDataPropsWithMetadata() tests the DecodeDataProps function
// should return a byte slice with the data properties base64 decoded
// no metadata flag set, so metadata should be removed
func (suite *EncodeDecodeDataTestSuite) TestDecodeDataPropsNoMetadata() {
	assert := assert.New(suite.T())

	flags := model.Flags{}
	output, err := DecodeDataProps(suite.EncodedInputWithMetadata, flags)
	assert.Nil(err)

	// unmarshal the expected and actual output to interface{} to make comparison easier
	var expected interface{}
	err = yaml.Unmarshal(suite.DecodedInputNoMetadata, &expected)
	assert.Nil(err)

	var actual interface{}
	err = yaml.Unmarshal(output, &actual)
	assert.Nil(err)

	assert.Equal(expected, actual)
}

// TestEncodeDataPropsWithMetadata() tests the EncodeDataProps function
// should return a byte slice with the data properties base64 encoded
// metadata flag set, so metadata should be kept
func (suite *EncodeDecodeDataTestSuite) TestEncodeDataPropsWithMetadata() {
	assert := assert.New(suite.T())

	flags := model.Flags{Metadata: true} // metadata flag set
	output, err := EncodeDataProps(suite.DecodedInputWithMetadata, flags)
	assert.Nil(err)

	// unmarshal the expected and actual output to interface{} to make comparison easier
	var expected interface{}
	err = yaml.Unmarshal(suite.EncodedInputWithMetadata, &expected)
	assert.Nil(err)

	var actual interface{}
	err = yaml.Unmarshal(output, &actual)
	assert.Nil(err)

	assert.Equal(expected, actual)
}

// TestDecodeDataPropsWithMetadata() tests the DecodeDataProps function
// should return a byte slice with the data properties base64 decoded
// metadata flag set, so metadata should be kept
func (suite *EncodeDecodeDataTestSuite) TestDecodeDataPropsWithMetadata() {
	assert := assert.New(suite.T())

	flags := model.Flags{Metadata: true} // metadata flag set
	output, err := DecodeDataProps(suite.EncodedInputWithMetadata, flags)
	assert.Nil(err)

	// unmarshal the expected and actual output to interface{} to make comparison easier
	var expected interface{}
	err = yaml.Unmarshal(suite.DecodedInputWithMetadata, &expected)
	assert.Nil(err)

	var actual interface{}
	err = yaml.Unmarshal(output, &actual)
	assert.Nil(err)

	assert.Equal(expected, actual)
}

func (suite *EncodeDecodeDataTestSuite) TestEncodeDataPropsInvalidInput() {
	assert := assert.New(suite.T())

	flags := model.Flags{}
	_, err := EncodeDataProps([]byte("invalid input"), flags)
	assert.NotNil(err)
}

func (suite *EncodeDecodeDataTestSuite) TestDecodeDataPropsInvalidInput() {
	assert := assert.New(suite.T())

	flags := model.Flags{}
	_, err := DecodeDataProps([]byte("invalid input"), flags)
	assert.NotNil(err)
}

func TestEncodeDataPropsSuite(t *testing.T) {
	suite.Run(t, new(EncodeDecodeDataTestSuite))
}
