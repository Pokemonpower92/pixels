package datastore

import (
	"context"
	"image"
	"image/draw"
	_ "image/jpeg"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pokemonpower92/imagesetservice/config"
)

// Store is an interface that defines the methods for retrieving image sets.
type Store interface {
	GetImageSet() ([]*image.RGBA, error)
}

// S3Api is an interface that defines the methods for interacting with Amazon S3.
type S3Api interface {
	ListObjectsV2(
		ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options),
	) (*s3.ListObjectsV2Output, error)
	GetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.GetObjectOutput, error)
}

// S3Store is a struct that implements the Store interface and provides methods for retrieving image sets from Amazon S3.
type S3Store struct {
	Bucket string
	logger *log.Logger
	api    S3Api
}

func NewS3Store(bucket string) *S3Store {
	s3Config := config.NewS3Config()

	awsConfig, err := awsconf.LoadDefaultConfig(
		context.TODO(),
		awsconf.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s3Config.AccessKeyID,
			s3Config.SecretAccessKey,
			"",
		)),
		awsconf.WithRegion(s3Config.Region),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsConfig)
	return &S3Store{
		Bucket: bucket,
		api:    client,
		logger: log.New(log.Writer(), "s3store ", log.LstdFlags),
	}
}

func (store *S3Store) GetImageSet() ([]*image.RGBA, error) {
	output, err := store.api.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(store.Bucket),
	})
	if err != nil {
		store.logger.Printf("Failed to list objects in bucket: %s: %s", store.Bucket, err)
		return nil, err
	}
	store.logger.Printf("Found %d images in bucket %s", len(output.Contents), store.Bucket)

	var images []*image.RGBA
	for _, object := range output.Contents {
		imageObject, err := store.api.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(store.Bucket),
			Key:    object.Key,
		})
		if err != nil {
			store.logger.Printf("Failed to get image: %s", err)
			return nil, err
		}

		decodedImage, _, err := image.Decode(imageObject.Body)
		if err != nil {
			store.logger.Printf("Failed to decode image: %s", err)
			return nil, err
		}

		YCbCrImage := decodedImage.(*image.YCbCr)

		b := YCbCrImage.Bounds()
		RGBAImage := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(RGBAImage, RGBAImage.Bounds(), YCbCrImage, b.Min, draw.Src)

		images = append(images, RGBAImage)
	}

	store.logger.Printf("Loaded %d images", len(images))

	return images, nil
}
