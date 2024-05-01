package cmd

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutputResultToNewFile(t *testing.T) {
	assert := assert.New(t)

	// setup flags
	f.Decode = false
	f.Output = fmt.Sprint(time.Now().UnixNano()) + ".txt"
	input.Exists = false

	result := []byte("test")

	if err := outputResult(result); err != nil {
		t.Fatal(err)
	}

	fileContent, err := os.ReadFile(f.Output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(result, fileContent)

	// remove the file (cleanup)
	if err := os.Remove(f.Output); err != nil {
		t.Fatal(err)
	}
}

func TestOutputResultToExistingFile(t *testing.T) {
	assert := assert.New(t)

	// setup flags
	f.Decode = false
	f.Output = fmt.Sprint(time.Now().UnixNano()) + ".txt"
	input.Exists = false

	result := []byte("test")

	// create the file so it already exists
	if err := os.WriteFile(f.Output, []byte("some other text"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := outputResult(result); err != nil {
		t.Fatal(err)
	}

	fileContent, err := os.ReadFile(f.Output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(result, fileContent)

	// remove the file (cleanup)
	if err := os.Remove(f.Output); err != nil {
		t.Fatal(err)
	}
}
func TestOutputResultToInplaceFile(t *testing.T) {
	assert := assert.New(t)

	// setup flags
	f.Decode = false
	f.Output = ""
	f.Inplace = true
	input.Exists = true
	input.Path = fmt.Sprint(time.Now().UnixNano()) + ".txt"

	result := []byte("test")

	// create the file so it already exists
	if err := os.WriteFile(input.Path, []byte("some other text"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := outputResult(result); err != nil {
		t.Fatal(err)
	}

	fileContent, err := os.ReadFile(input.Path)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(result, fileContent)

	// remove the file (cleanup)
	if err := os.Remove(input.Path); err != nil {
		t.Fatal(err)
	}
}
