package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `
		{
			"id": "random-uuid", 
			"file_path": "convite.mp4", 
			"status":"pending"
		}
	`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `wes`
	err = utils.IsJson(json)
	require.Error(t, err)
}

func TestGenerateUUIDString(t *testing.T) {
	id, err := utils.GenerateUUIDString()
	require.Nil(t, err)
	require.NotEqual(t, id, "")
}
