package service

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	SvcHost            string
	TwilioAccountSID   string
	TwilioAccountToken string
	S3Bucket           string
	S3Region           string
}

type FaxmanService struct{}

func (s *FaxmanService) Run(cfg Config) error {
	faxResource := &FaxResource{
		TwilioAccountSID:   cfg.TwilioAccountSID,
		TwilioAccountToken: cfg.TwilioAccountToken,
		S3Bucket:           cfg.S3Bucket,
		S3Region:           cfg.S3Region,
	}

	r := gin.Default()

	r.GET("/fax", faxResource.GetAllFaxes)
	r.GET("/fax/:id", faxResource.GetFax)
	r.POST("/fax", faxResource.CreateFax)

	r.Run(cfg.SvcHost)

	return nil
}
