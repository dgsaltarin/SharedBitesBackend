package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/domain/entity"
)

const (
	AWS_TOPIC_ARN string = "arn:aws:sns:us-east-1:123456789012:TextractTopic"
	AWS_ROLE_ARN  string = "arn:aws:iam::955650001607:role/texttrack-sns"
)

type TextTrackService struct{}

func NewTextTrackService() *TextTrackService {
	return &TextTrackService{}
}

func (tt *TextTrackService) TextTrackSesson(session *session.Session) *textract.Textract {
	svc := textract.New(session)
	return svc
}

// Detectitems detects start the expenses analysis job and get teh result once the job is complete
func (tt *TextTrackService) Detectitems(svc *textract.Textract, objectname string) []entity.Item {
	input := &textract.StartExpenseAnalysisInput{
		NotificationChannel: &textract.NotificationChannel{
			SNSTopicArn: aws.String(AWS_TOPIC_ARN),
			RoleArn:     aws.String(AWS_ROLE_ARN),
		},
		DocumentLocation: &textract.DocumentLocation{
			S3Object: &textract.S3Object{
				Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
				Name:   aws.String(objectname),
			},
		},
		OutputConfig: &textract.OutputConfig{
			S3Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
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

	fmt.Printf("result %s", output.ExpenseDocuments[0].LineItemGroups)

	expensesR := tt.extractExpensesFromResults(output.ExpenseDocuments[0].LineItemGroups)

	fmt.Printf("expenses %s", expensesR)

	return expensesR
}

// extractExpensesFromResults extracts the expenses item and price from the results of the expense analysis
func (tt *TextTrackService) extractExpensesFromResults(itemsGroup []*textract.LineItemGroup) []entity.Item {
	var expenses []entity.Item

	// Example: Extract expenses from table cells
	for _, itemRow := range itemsGroup[0].LineItems {
		var name string
		var price string

		if *itemRow.LineItemExpenseFields[0].Type.Text == "ITEM" {
			name = *itemRow.LineItemExpenseFields[0].ValueDetection.Text
		}

		if *itemRow.LineItemExpenseFields[1].Type.Text == "PRICE" {
			price = *itemRow.LineItemExpenseFields[1].ValueDetection.Text
		}

		parsedPrice, err := tt.parsePrice(price)
		if err != nil {
			fmt.Println(err)
		}

		expenses = append(expenses, entity.Item{Name: name, Price: parsedPrice})
	}

	return expenses
}

// parse price from string
func (tt *TextTrackService) parsePrice(text string) (float64, error) {
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
