package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`error loading .env file`)
	}
}

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDbTest()
	defer db.Close()

	id := uuid.NewV4().String()

	video := domain.NewVideo()
	video.ID = id
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDB{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)
	if err != nil {
		log.Fatalf("video %v not found", video.ID)
	}

	return v, repo
}

func TestVideoServiceDownload(t *testing.T) {
	v, r := prepare()

	vs := services.NewVideoService()
	vs.Video = v
	vs.VideoRepository = r

	err := vs.Download("ferdev-video-encoder")
	require.Nil(t, err)

	err = vs.Fragment()
	require.Nil(t, err)

	err = vs.Encode()
	require.Nil(t, err)

	err = vs.Finish()
	require.Nil(t, err)

}

func TestVideoServiceFragment(t *testing.T) {
	v, r := prepare()

	vs := services.NewVideoService()
	vs.Video = v
	vs.VideoRepository = r

	err := vs.Download("ferdev-video-encoder")
	require.Nil(t, err)

	err = vs.Fragment()
	require.Nil(t, err)
}
