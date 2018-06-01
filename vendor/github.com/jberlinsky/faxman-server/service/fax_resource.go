package service

import (
	"github.com/sfreiberg/gotwilio"
	"github.com/jberlinsky/faxman-server/api"
	"github.com/gin-gonic/gin"
	"log"
)

type FaxResource struct {
	*gotwilio.FaxResource
	TwilioAccountSID string
	TwilioAccountToken string
	S3Bucket string
	S3Region string
}

func (r *FaxResource) CreateFax(c *gin.Context) {
	faxMediaLocations := make([]string)
	multipart, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(400, api.NewError("Failed to create MultipartReader", err)
	}

	for {
		mimePart, err := multipart.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(400, api.NewError("Error reading multipart section: %v", err))
			break
		}
		disposition, params, err := mime.ParseMediaType(mimePart.Header.Get("Content-Disposition"))
		if err != nil {
			c.JSON(400, api.NewError("Invalid Content Disposition: %v", err))
			break
		}

		uploader := s3manager.NewUploader(
			session.New(
				&aws.Config{
					Region: aws.String(r.S3Region),
				},
			),
		)
		result, err := uploader.Upload(
			&s3manager.UploadInput{
				Body: mimePart,
				Bucket: aws.String(r.S3Bucket),
				Key: aws.String(params["filename"]),
				ContentType: aws.String(mimePart.Header.Get("Content-Type")),
				ACL: aws.String("public-read"),
			},
		)
		if err != nil {
			c.JSON(500, api.NewError("Failed to upload to S3"))
			return
		}
		faxMediaLocations = append(faxMediaLocations, result.Location)
	}

	if len(faxMediaLocations) != 1 {
		c.JSON(400, api.NewError("You must upload exactly one file"))
		return
	}

	var fax gotwilio.FaxResource

	if c.Bind(&fax) != nil {
		c.JSON(400, api.NewError("Error decoding body"))
		return
	}

	fax.MediaUrl = faxMediaLocations[0]

	fax, exception, err := twilio.SendFax(
		fax.To,
		fax.From,
		fax.MediaUrl,
		fax.Quality,
		"", // ... // TODO status callback
		false,
	)

	if err != nil {
		c.JSON(500, api.NewError(fmt.Sprintf("Something went wrong sending fax: %v", err.Error())))
		return
	} else if exception != nil {
		c.JSON(500, api.NewError("Something went wrong sending fax"))
		return
	} else {
		c.JSON(201, fax)
	}
}

func (r *FaxResource) GetAllFaxes(c *gin.Context) {
	var faxes []*gotwilio.FaxResource

	faxes, exception, err := twilio.GetFaxes("", "", "", "")
	if err != nil {
		c.JSON(500, api.NewError(fmt.Sprintf("Something went wrong retrieving faxes: %v", err.Error())))
		return
	} else if exception != nil {
		c.JSON(500, api.NewError("Something went wrong sending fax"))
		return
	} else {
		c.JSON(200, faxes)
	}
}

func (r *FaxResource) GetFax(c *gin.Context) {
	id, err := r.getId(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	faxResource, exception, err := twilio.GetFax(id)
	if err != nil {
		c.JSON(500, api.NewError(fmt.Sprintf("Something went wrong retrieving fax: %v", err.Error())))
		return
	} else if exception != nil {
		c.JSON(500, api.NewError("Something went wrong sending fax"))
		return
	} else {
		c.JSON(200, faxResource)
	}
}

func (r *FaxResource) getId(c *gin.Context) string {
	id = c.Params.ByName("id")
	return id
}

func (r *FaxResource) twilioClient() *gotwilio.Twilio {
	return gotwilio.NewTwilioClient(
		r.TwilioAccountSID,
		r.TwilioAccountToken,
	)
}
