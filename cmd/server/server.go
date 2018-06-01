package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jberlinsky/faxman-server/service"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(c *cli.Context) (service.Config, error) {
	config := service.Config{
		SvcHost:            fmt.Sprintf(":%s", c.GlobalString("port")),
		TwilioAccountSID:   c.GlobalString("twilio-account-sid"),
		TwilioAccountToken: c.GlobalString("twilio-account-token"),
		S3Bucket:           c.GlobalString("s3-bucket"),
		S3Region:           c.GlobalString("s3-region"),
	}

	return config
}

func main() {
	app := cli.NewApp()
	app.Name = "faxman"
	app.Version = "0.0.1" // TODO read from elsewhere/git
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "port, p",
			Value:  "8080",
			Usage:  "Port to serve on",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "twilio-account-sid, s",
			Usage:  "Twilio Account SID",
			EnvVar: "TWILIO_ACCOUNT_SID",
		},
		cli.StringFlag{
			Name:   "twilio-account-token, t",
			Usage:  "Twilio Account Token",
			EnvVar: "TWILIO_ACCOUNT_TOKEN",
		},
		cli.StringFlag{
			Name:   "s3-bucket, b",
			Usage:  "S3 bucket to store files in",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "s3-region, r",
			Value:  "us-east-1",
			Usage:  "S3 Region",
			EnvVar: "S3_REGION",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the HTTP server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}
				svc := service.FaxmanService{}
				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
