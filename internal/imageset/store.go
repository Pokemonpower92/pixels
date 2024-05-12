package imageset

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

type Store interface {
	GetImageSet() ([]*image.RGBA, error)
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
		awsconf.WithRegion(conf.Region),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return &S3Store{
		Bucket: bucket,
		client: *client,
		l:      log.New(log.Writer(), "s3store ", log.LstdFlags),
	}
}

func (s *S3Store) GetImageSet() ([]*image.RGBA, error) {
	output, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		s.l.Printf("Failed to list objects in bucket: %s: %s", s.Bucket, err)
	}

	s.l.Printf("Found %d images in bucket %s", len(output.Contents), s.Bucket)

	var images []*image.RGBA
	for _, object := range output.Contents {
		img, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    object.Key,
		})
		if err != nil {
			s.l.Printf("Failed to get image: %s", err)
			return nil, err
		}

		decodedImage, _, err := image.Decode(img.Body)
		if err != nil {
			s.l.Printf("Failed to decode image: %s", err)
			return nil, err
		}

		YCbCrImage := decodedImage.(*image.YCbCr)

		// Convert YCbCr image to RGBA
		b := YCbCrImage.Bounds()
		RGBAImage := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(RGBAImage, RGBAImage.Bounds(), YCbCrImage, b.Min, draw.Src)

		images = append(images, RGBAImage)
	}

	s.l.Printf("Loaded %d images", len(images))

	return images, nil
}
