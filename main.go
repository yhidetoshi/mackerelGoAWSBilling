package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/mackerelio/mackerel-client-go"
)

var (
	mkrKey = os.Getenv("MKRKEY")
	client = mackerel.NewClient(mkrKey)
	config = aws.Config{Region: aws.String(region)}
	cwt    = cloudwatch.New(session.New(&config))
)

const (
	region      = "us-east-1"
	serviceName = "AWS"
	timezone    = "Asia/Tokyo"
	offset      = 9 * 60 * 60
)

func main() {
	lambda.Start(Handler)
}

// Handler Lambda
func Handler() {
	var cost float64

	jst := time.FixedZone(timezone, offset)
	nowTime := time.Now().In(jst)

	input := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("Currency"),
				Value: aws.String("USD"),
			},
		},
		StartTime:  aws.Time(time.Now().Add(time.Hour * -24)),
		EndTime:    aws.Time(time.Now()),
		Period:     aws.Int64(86400),
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
	}
	response, err := cwt.GetMetricStatistics(input)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range response.Datapoints {
		cost = *v.Maximum
	}
	fmt.Println(cost)

	errMkr := PostValuesToMackerel(cost, nowTime)
	if errMkr != nil {
		fmt.Println(errMkr)
	}
}

// PostValuesToMackerel Post Metrics to Mackerel
func PostValuesToMackerel(cost float64, nowTime time.Time) error {
	err := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Cost.cost",
			Time:  nowTime.Unix(),
			Value: cost,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
