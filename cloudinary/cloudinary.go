package cloudinary

import (
	"fmt"
	"io"
	"time"

	"github.com/HomesNZ/go-common/env"
	"github.com/nicday/go-cloudinary"
)

var (
	// Service is a singleton Cloudinary Service instance that be modified during testing.
	service *CDNService
)

// Service returns a Cloudinary Service singleton.
func Service() (*CDNService, error) {
	if service != nil {
		return service, nil
	}

	key := env.GetString("CLOUDINARY_API_KEY", "")
	if key == "" {
		return nil, fmt.Errorf("CLOUDINARY_API_KEY is not set")
	}
	secret := env.GetString("CLOUDINARY_API_SECRET", "")
	if secret == "" {
		return nil, fmt.Errorf("CLOUDINARY_API_SECRET is not set")
	}
	name := env.GetString("CLOUDINARY_CLOUD_NAME", "")
	if name == "" {
		return nil, fmt.Errorf("CLOUDINARY_CLOUD_NAME is not set")
	}

	c, err := cloudinary.Dial(fmt.Sprintf(
		"cloudinary://%s:%s@%s",
		key,
		secret,
		name,
	))
	if err != nil {
		return nil, err
	}

	service = &CDNService{
		service: c,
	}

	return service, nil
}

// CDNService is a Cloudinary concrete implementation of cdn.Interface
type CDNService struct {
	service *cloudinary.Service
}

// UploadImage uploads a new image to the Cloudinary CDN, with the supplied name and reader data. The public URL will be
// returned when successful, otherwise an error will be returned.
func (c CDNService) UploadImage(name string, data io.Reader) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s_%d", name, now.UnixNano())

	_, err := c.service.UploadStaticImage(fileName, data, "")
	if err != nil {
		return "", err
	}

	return c.service.Url(fileName, cloudinary.ImageType), nil
}

// UploadURI sets the URI used when uploading assets to the CDN.
func (c *CDNService) UploadURI(uri string) error {
	return c.service.UploadURI(uri)
}

// RemoveAsset removes an asset from Cloudinary.
func (c CDNService) RemoveAsset(uri string) error {
	id, err := c.service.PublicID(uri)
	if err != nil {
		return err
	}

	return c.service.Delete(id, "", cloudinary.ImageType)
}
