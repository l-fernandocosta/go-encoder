package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	uid := uuid.NewV4()

	video := domain.NewVideo()
	video.ID = uid.String()
	video.CreatedAt = time.Now()
	video.FilePath = "path"
	video.ResourceID = "a"

	job, err := domain.NewJob("OUTPUT", "CONVERTED", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
