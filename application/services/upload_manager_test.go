package services_test

import (
	"encoder/application/services"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}
}

func TestUploadManager(t *testing.T) {
	bktName := "codeeducationtest"
	v, r := prepare()

	vs := services.NewVideoService()
	vs.Video = v
	vs.VideoRepository = r

	err := vs.Download(bktName)
	require.Nil(t, err)

	err = vs.Fragment()
	require.Nil(t, err)

	vu := services.NewVideoUpload()
	vu.OutputBucket = bktName
	vu.VideoPath = os.Getenv("localStoragePath") + "/" + v.ID

	doneUpload := make(chan string)
	go vu.ProcessUpload(50, doneUpload)

	result := <-doneUpload

	require.Equal(t, result, "upload completed")
}
