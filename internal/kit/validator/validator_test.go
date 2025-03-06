package validator_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AddMile/backend/internal/kit/validator"
)

var (
	//go:embed testdata/dynamic_content_ok.json
	dynamicContentOK []byte

	//go:embed testdata/dynamic_content_error.json
	dynamicContentError []byte
)

func TestValidateDynamicContent_OK(t *testing.T) {
	type data struct {
		Content []map[string]any `json:"content" validate:"required,custom_content"`
	}

	var d data
	err := json.Unmarshal(dynamicContentOK, &d)
	assert.NoError(t, err)

	v := validator.New()
	err = v.Validate(d)
	assert.NoError(t, err)
}

func TestValidateDynamicContent_Error(t *testing.T) {
	type data struct {
		Content []map[string]any `json:"content" validate:"required,custom_content"`
	}

	var d data
	err := json.Unmarshal(dynamicContentError, &d)
	assert.NoError(t, err)

	v := validator.New()
	err = v.Validate(d)
	assert.Error(t, err)
}
