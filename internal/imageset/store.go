package imageset

import (
	"context"
	"image"
	_ "image/jpeg"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pokemonpower92/imagesetservice/config"
)

type Store interface {
	GetImageSet() (*image.RGBA, error)
}

type S3Store struct {
	Bucket string
	l      *log.Logger
	client s3.Client
}

func NewS3Store(bucket string) *S3Store {
	conf := config.NewS3Config()

	cfg, err := awsconf.LoadDefaultConfig(
		context.TODO(),
		awsconf.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.AccessKeyID,
			conf.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return &S3Store{
		Bucket: bucket,
		client: *client,
	}
}

func (s *S3Store) GetImageSet() ([]*image.YCbCr, error) {
	output, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Images in the bucket:")
	var images []*image.YCbCr
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
		// Load the image and append it to the images slice
		img, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    object.Key,
		})
		if err != nil {
			return nil, err
		}
		// Decode the image
		decodedImage, _, err := image.Decode(img.Body)
		if err != nil {
			return nil, err
		}
		// Convert the image to a YCbCr image
		rgbaImage := decodedImage.(*image.YCbCr)
		images = append(images, rgbaImage)
	}

	return images, nil
}
