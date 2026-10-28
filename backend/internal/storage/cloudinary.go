package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary(cloudName, apiKey, apiSecret string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("create cloudinary client: %w", err)
	}
	return &Cloudinary{cld: cld}, nil
}

func (c *Cloudinary) Upload(ctx context.Context, file io.Reader) (string, string, error) {
	result, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "seraphim/listings",
	})
	if err != nil {
		return "", "", fmt.Errorf("upload photo: %w", err)
	}
	return result.SecureURL, result.PublicID, nil
}
