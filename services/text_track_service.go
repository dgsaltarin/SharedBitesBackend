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

func Detectitems(svc *textract.Textract, objectname string) *textract.GetExpenseAnalysisOutput {
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

	time.Sleep(12 * time.Second)

	output, err := svc.GetExpenseAnalysis(&textract.GetExpenseAnalysisInput{
		JobId: result.JobId,
	})

	if err != nil {
		fmt.Println(err)
	}

	//expensesR := extractExpensesFromResults(output.Blocks)

	return output
}

func extractExpensesFromResults(blocks []*textract.Block) []Expense {
	var expenses []Expense

	// Example: Extract expenses from table cells
	for _, block := range blocks {
		if *block.BlockType == "CELL" && *block.RowSpan == 1 && *block.ColumnSpan == 1 {
			// Extract item and price from table cells
			if len(block.Relationships) > 0 && *block.Relationships[0].Type == "CHILD" {
				item := ""
				price := 0.0

				for _, id := range block.Relationships[0].Ids {
					childBlock := findBlockByID(id, blocks)
					if childBlock != nil && *childBlock.BlockType == "WORD" {
						if item == "" {
							item = *childBlock.Text
						} else {
							price, _ = parsePrice(*childBlock.Text)
						}
					}
				}

				if item != "" && price > 0.0 {
					expenses = append(expenses, Expense{Item: item, Price: price})
				}
			}
		}
	}

	return expenses
}

func findBlockByID(id *string, blocks []*textract.Block) *textract.Block {
	for _, block := range blocks {
		if block.Id == id {
			return block
		}
	}
	return nil
}

// parse price from string
func parsePrice(text string) (float64, error) {
	// Remove non-numeric characters from the text
	text = strings.ReplaceAll(text, "$", "")
	text = strings.ReplaceAll(text, ",", "")

	// Parse the text as a float64
	price, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}
