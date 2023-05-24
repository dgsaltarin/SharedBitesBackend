package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

const (
	AWS_TOPIC_ARN string = "arn:aws:sns:us-east-1:123456789012:TextractTopic"
	AWS_ROLE_ARN  string = "arn:aws:iam::955650001607:role/texttrack-sns"
)

type Expense struct {
	Item  string
	Price float64
}

func TextTrackSesson(session *session.Session) *textract.Textract {
	svc := textract.New(session)

	fmt.Println("Textract session created")

	return svc
}

// Detectitems detects start the expenses analysis job and get teh result once the job is complete
func Detectitems(svc *textract.Textract, objectname string) []Expense {
	input := &textract.StartExpenseAnalysisInput{
		NotificationChannel: &textract.NotificationChannel{
			SNSTopicArn: aws.String(AWS_TOPIC_ARN),
			RoleArn:     aws.String(AWS_ROLE_ARN),
		},
		DocumentLocation: &textract.DocumentLocation{
			S3Object: &textract.S3Object{
				Bucket: aws.String(AWS_S3_BUCKET),
				Name:   aws.String(objectname),
			},
		},
		OutputConfig: &textract.OutputConfig{
			S3Bucket: aws.String(AWS_S3_BUCKET),
			S3Prefix: aws.String("output"),
		},
	}

	result, err := svc.StartExpenseAnalysis(input)
	if err != nil {
		fmt.Println(err)
	}

	output, err := svc.GetExpenseAnalysis(&textract.GetExpenseAnalysisInput{
		JobId: result.JobId,
	})

	if err != nil {
		fmt.Println(err)
	}

	// repeat the request until the analysis job is complete
	for *output.JobStatus == "IN_PROGRESS" {
		time.Sleep(2 * time.Second)
		output, err = svc.GetExpenseAnalysis(&textract.GetExpenseAnalysisInput{
			JobId: result.JobId,
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	expensesR := extractExpensesFromResults(output.ExpenseDocuments[0].LineItemGroups)

	return expensesR
}

// extractExpensesFromResults extracts the expenses item and price from the results of the expense analysis
func extractExpensesFromResults(itemsGroup []*textract.LineItemGroup) []Expense {
	var expenses []Expense

	// Example: Extract expenses from table cells
	for _, itemRow := range itemsGroup[0].LineItems {
		var item string
		var price string

		if *itemRow.LineItemExpenseFields[0].Type.Text == "ITEM" {
			item = *itemRow.LineItemExpenseFields[0].ValueDetection.Text
		}

		if *itemRow.LineItemExpenseFields[1].Type.Text == "PRICE" {
			price = *itemRow.LineItemExpenseFields[1].ValueDetection.Text
		}

		parsedPrice, err := parsePrice(price)
		if err != nil {
			fmt.Println(err)
		}

		expenses = append(expenses, Expense{Item: item, Price: parsedPrice})
	}

	return expenses
}

// parse price from string
func parsePrice(text string) (float64, error) {
	// Remove non-numeric characters from the text
	text = strings.ReplaceAll(text, "$", "")
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, ".", "")

	// Parse the text as a float64
	price, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}
