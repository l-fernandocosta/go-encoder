package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIsNotAValidUUID(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "INVALID_UUID"
	video.ResourceID = "uuid"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	uid := uuid.NewV4()
	video.ID = uid.String()
	video.ResourceID = "uuid"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Nil(t, err)
}
