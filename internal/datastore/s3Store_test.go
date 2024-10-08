package datastore

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mockS3Api struct{}

func (m *mockS3Api) ListObjectsV2(
	ctx context.Context,
	params *s3.ListObjectsV2Input,
	optFns ...func(*s3.Options),
) (*s3.ListObjectsV2Output, error) {
	key := "test-image-1.jpg"
	return &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: &key,
			},
			{
				Key: &key,
			},
		},
	}, nil
}

func (m *mockS3Api) GetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	mockImage := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var buf bytes.Buffer
	jpeg.Encode(&buf, mockImage, nil)
	mockJPEGData := buf.Bytes()
	return &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader(mockJPEGData)),
	}, nil
}

func TestS3Store(t *testing.T) {
	tests := []struct {
		name           string
		expectedImages []*image.RGBA
	}{
		{
			name: "Test case 1",
			expectedImages: []*image.RGBA{
				image.NewRGBA(image.Rect(0, 0, 2, 2)),
				image.NewRGBA(image.Rect(0, 0, 2, 2)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s3Store := S3Store{
				Bucket: "test-bucket",
				logger: log.New(log.Writer(), "test: ", log.LstdFlags),
				api:    &mockS3Api{},
			}
			result, err := s3Store.GetImages("test")
			if err != nil {
				t.Errorf("Error occurred: %v", err)
			}
			if len(result) != len(tt.expectedImages) {
				t.Errorf("Image length mismatch. Expected: %d, Got: %d",
					len(tt.expectedImages),
					len(result),
				)
			}
		})
	}
}
