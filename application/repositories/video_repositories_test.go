package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepository_INSERT(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	id, _ := uuid.NewV4()
	video := domain.NewVideo()
	video.ID = id.String()
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDB{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.Nil(t, err)
	require.NotEmpty(t, v.ID)
	require.Equal(t, v.ID, video.ID)
}
