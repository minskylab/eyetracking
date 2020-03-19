package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const MinskyPublicSpace = "minskypublic"
const projectFolder = "eyetraking"

type StorageInfo struct {
	Images     []string `json:"images"`
	LastUpdate time.Time `json:"lastUpdate"`
}

var minioClient *minio.Client

func minioClientInit() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Error(errors.Wrap(err, "error loading .env file"))
	}

	accessKey := os.Getenv("DO_ACCESS_KEY")
	secKey := os.Getenv("DO_SECRET_ACCESS_KEY")

	minioClient, err = minio.New("nyc3.digitaloceanspaces.com", accessKey, secKey, true)
	if err != nil {
		err = errors.Wrap(err, "problem at create minio client")
		log.Error(err)
		panic(err)
	}
}

func storageInfo() (*StorageInfo, error) {
	exists, err := minioClient.BucketExists(MinskyPublicSpace)
	if err != nil {
		return nil, errors.Wrap(err, "error at check if bucket exists")
	}

	if !exists {
		err = minioClient.MakeBucket(MinskyPublicSpace, "nyc3")
		if err != nil {
			return nil, errors.Wrap(err, "error at try to make a new bucket")
		}
	}
	done := make(chan struct{}, 1)
	totalImages := make([]string, 0)
	info := minioClient.ListObjectsV2WithMetadata(MinskyPublicSpace, projectFolder, true, done)

	for i := range info {
		log.WithField("content-type", i.ContentType).Infof("object: %s", i.Key)
		if i.Key == projectFolder + "/" {
			continue
		}
		totalImages = append(totalImages, i.Key)
	}

	close(done)

	return &StorageInfo{
		Images:     totalImages,
		LastUpdate: time.Now(),
	}, nil
}
