package handler

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/google/uuid"
	"golang.org/x/net/http2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Collect anonymous data from the request:
	//
	//  - URL path (/blog/article-name, /apps, ...)
	//  - Referrer (google, twitter, ...)
	//  - utm_source query parameter (hackernewsletter, wsjnews, ...)
	//  - utm_medium query parameter (email, post, ...)
	//
	// Only the URL path is guaranteed to be available

	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form data: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pathValue := r.PostForm.Get("path")
	referrerValue := r.PostForm.Get("referrer")
	sourceValue := r.PostForm.Get("source")
	mediumValue := r.PostForm.Get("medium")
	visitIDValue := uuid.NewString()

	// Collect environment variables

	awsRegion := os.Getenv("MY_AWS_REGION")
	awsAccessKey := os.Getenv("MY_AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("MY_AWS_SECRET_ACCESS_KEY")
	tsDatabaseName := os.Getenv("MY_AWS_TIMESTREAM_DATABASE_NAME")
	tsTableName := os.Getenv("MY_AWS_TIMESTREAM_TABLE_NAME")

	// Initialize the AWS SDK
	// Based on https://github.com/awslabs/amazon-timestream-tools/blob/mainline/sample_apps/go/crud-ingestion-sample.go

	tr := &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		// Using DefaultTransport values for other parameters: https://golang.org/pkg/net/http/#RoundTripper
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			DualStack: true,
			Timeout:   30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	http2.ConfigureTransport(tr)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey,
			awsSecretKey,
			"",
		),
	})
	if err != nil {
		log.Printf("error creating AWS session: %s", err.Error())
		w.WriteHeader(500)
		return
	}

	// Write the visit information to AWS Timestream

	writeSvc := timestreamwrite.New(sess)

	dimensions := []*timestreamwrite.Dimension{
		{Name: aws.String("path"), Value: &pathValue},
		{Name: aws.String("visit_id"), Value: &visitIDValue},
	}
	if referrerValue != "" {
		dimensions = append(dimensions, &timestreamwrite.Dimension{Name: aws.String("referrer"), Value: &referrerValue})
	}
	if sourceValue != "" {
		dimensions = append(dimensions, &timestreamwrite.Dimension{Name: aws.String("source"), Value: &sourceValue})
	}
	if mediumValue != "" {
		dimensions = append(dimensions, &timestreamwrite.Dimension{Name: aws.String("medium"), Value: &mediumValue})
	}

	writeRecordsInput := timestreamwrite.WriteRecordsInput{
		DatabaseName: &tsDatabaseName,
		TableName:    &tsTableName,
		Records: []*timestreamwrite.Record{
			{
				Dimensions:       dimensions,
				MeasureName:      aws.String("visit"),
				MeasureValue:     aws.String("1"),
				MeasureValueType: aws.String("BIGINT"),
				Time:             aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
				TimeUnit:         aws.String("SECONDS"),
			},
		},
	}

	result, err := writeSvc.WriteRecords(&writeRecordsInput)
	if err != nil {
		log.Printf("error writing AWS timestream record: %s", err.Error())
		w.WriteHeader(500)
		return
	}

	if result.RecordsIngested == nil || *result.RecordsIngested.Total != 1 {
		log.Printf("no AWS timestream records written")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
