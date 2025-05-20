package texttrack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"

	"github.com/dgsaltarin/SharedBitesBackend/internal/ports" // Corrected import path
)

type AWSTextractAdapter struct {
	textractClient *textract.Client
}

func NewAWSTextractAdapter(cfg aws.Config) *AWSTextractAdapter {
	return &AWSTextractAdapter{
		textractClient: textract.NewFromConfig(cfg),
	}
}

func (a *AWSTextractAdapter) AnalyzeDocument(ctx context.Context, storagePath string) (*ports.ParsedTextractData, error) {
	// storagePath is expected to be in the format "s3://bucket-name/key"
	bucket, key, err := parseS3Path(storagePath)
	if err != nil {
		return nil, fmt.Errorf("invalid S3 path: %w", err)
	}

	input := &textract.AnalyzeExpenseInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
	}

	output, err := a.textractClient.AnalyzeExpense(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze document with Textract: %w", err)
	}

	return parseTextractOutput(output)
}

// parseS3Path extracts bucket and key from an S3 path string (e.g., "s3://bucket/key")
func parseS3Path(s3Path string) (string, string, error) {
	if !strings.HasPrefix(s3Path, "s3://") {
		return "", "", fmt.Errorf("path must start with s3://")
	}
	trimmedPath := strings.TrimPrefix(s3Path, "s3://")
	parts := strings.SplitN(trimmedPath, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid S3 path format, expected s3://bucket/key")
	}
	return parts[0], parts[1], nil
}

// parseTextractOutput will convert the AWS Textract AnalyzeExpenseOutput to our internal ParsedTextractData
func parseTextractOutput(output *textract.AnalyzeExpenseOutput) (*ports.ParsedTextractData, error) {
	parsedData := &ports.ParsedTextractData{
		LineItems: []ports.ParsedLineItem{},
	}
	var rawTextBuilder strings.Builder

	for _, expenseDoc := range output.ExpenseDocuments {
		// Collect text from summary fields for RawTextOutput
		for _, summaryField := range expenseDoc.SummaryFields {
			if summaryField.LabelDetection != nil && summaryField.LabelDetection.Text != nil {
				rawTextBuilder.WriteString(*summaryField.LabelDetection.Text)
				rawTextBuilder.WriteString(": ")
			}
			if summaryField.ValueDetection != nil && summaryField.ValueDetection.Text != nil {
				rawTextBuilder.WriteString(*summaryField.ValueDetection.Text)
				rawTextBuilder.WriteString("\n") // Newline after each summary field value
			}

			fieldType := summaryField.Type
			fieldLabel := summaryField.LabelDetection
			fieldValue := summaryField.ValueDetection

			if fieldValue == nil || fieldValue.Text == nil {
				continue
			}
			valueText := *fieldValue.Text

			switch {
			case fieldType != nil && fieldType.Text != nil && strings.ToUpper(*fieldType.Text) == "VENDOR_NAME":
				parsedData.VendorName = aws.String(valueText)
			case fieldType != nil && fieldType.Text != nil && (strings.ToUpper(*fieldType.Text) == "INVOICE_RECEIPT_DATE" || strings.ToUpper(*fieldType.Text) == "EXPENSE_DATE" || strings.ToUpper(*fieldType.Text) == "DATE" || (fieldLabel != nil && fieldLabel.Text != nil && strings.Contains(strings.ToUpper(*fieldLabel.Text), "DATE"))):
				parsedTime, err := parseDate(valueText)
				if err == nil {
					parsedData.TransactionDate = parsedTime
				}
			case fieldType != nil && fieldType.Text != nil && (strings.ToUpper(*fieldType.Text) == "TOTAL" || strings.ToUpper(*fieldType.Text) == "INVOICE_TOTAL" || strings.ToUpper(*fieldType.Text) == "AMOUNT_DUE"):
				amount, err := parseFloat(valueText)
				if err == nil {
					parsedData.TotalAmount = aws.Float64(amount)
				}
			}
		}

		// Collect text from line items for RawTextOutput and parse line items
		for _, itemGroup := range expenseDoc.LineItemGroups {
			for _, lineItem := range itemGroup.LineItems {
				parsedLineItem := ports.ParsedLineItem{}
				var lineItemTextParts []string
				for _, field := range lineItem.LineItemExpenseFields {
					fieldType := field.Type
					fieldValue := field.ValueDetection

					if fieldValue == nil || fieldValue.Text == nil {
						continue
					}
					valueText := *fieldValue.Text

					// Add to raw text builder
					if field.Type != nil && field.Type.Text != nil {
						lineItemTextParts = append(lineItemTextParts, fmt.Sprintf("%s: %s", *field.Type.Text, valueText))
					}

					switch {
					case fieldType != nil && fieldType.Text != nil && (strings.ToUpper(*fieldType.Text) == "ITEM" || strings.ToUpper(*fieldType.Text) == "PRODUCT_CODE" || strings.ToUpper(*fieldType.Text) == "DESCRIPTION"):
						parsedLineItem.Description = valueText
					case fieldType != nil && fieldType.Text != nil && strings.ToUpper(*fieldType.Text) == "QUANTITY":
						qty, err := parseFloat(valueText)
						if err == nil {
							parsedLineItem.Quantity = aws.Float64(qty)
						}
					case fieldType != nil && fieldType.Text != nil && (strings.ToUpper(*fieldType.Text) == "PRICE" || strings.ToUpper(*fieldType.Text) == "UNIT_PRICE"):
						price, err := parseFloat(valueText)
						if err == nil {
							parsedLineItem.UnitPrice = aws.Float64(price)
						}
					case fieldType != nil && fieldType.Text != nil && (strings.ToUpper(*fieldType.Text) == "TOTAL_PRICE" || strings.ToUpper(*fieldType.Text) == "LINE_ITEM_TOTAL"):
						totalPrice, err := parseFloat(valueText)
						if err == nil {
							parsedLineItem.TotalPrice = aws.Float64(totalPrice)
						}
					}
				}
				if parsedLineItem.Description != "" { // Only add if we have a description
					parsedData.LineItems = append(parsedData.LineItems, parsedLineItem)
					if len(lineItemTextParts) > 0 {
						rawTextBuilder.WriteString(strings.Join(lineItemTextParts, ", "))
						rawTextBuilder.WriteString("\n")
					}
				}
			}
		}
	}
	parsedData.RawTextOutput = strings.TrimSpace(rawTextBuilder.String())
	return parsedData, nil
}

// parseFloat attempts to parse a string into a float64, handling common currency symbols and commas.
func parseFloat(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", "")
	// Consider a regex for more robust currency symbol removal
	s = strings.NewReplacer("$", "", "€", "", "£", "", "USD", "").Replace(s)
	s = strings.TrimSpace(s)
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f) // fmt.Sscanf is basic; consider strconv.ParseFloat for more control
	return f, err
}

// parseDate attempts to parse a date string.
// AWS Textract can return dates in various formats. This function tries a few common ones.
// For a production system, a more robust date parsing library might be preferable.
func parseDate(dateStr string) (*time.Time, error) {
	formats := []string{
		"2006-01-02",      // YYYY-MM-DD
		"01/02/2006",      // MM/DD/YYYY
		"01-02-2006",      // MM-DD-YYYY
		"02/01/2006",      // DD/MM/YYYY
		"02-Jan-2006",     // DD-Mon-YYYY
		"Jan 02, 2006",    // Mon DD, YYYY
		"02 January 2006", // DD Month YYYY
		"2006/01/02",      // YYYY/MM/DD
		"01-Jan-06",       // DD-Mon-YY
		"02.01.2006",      // DD.MM.YYYY
		"2006.01.02",      // YYYY.MM.DD
		"01/02/06",        // MM/DD/YY
		"02/01/06",        // DD/MM/YY
		time.RFC3339,      // Example: 2024-07-15T10:00:00Z
		"02-01-2006",      // DD-MM-YYYY
	}
	dateStr = strings.TrimSpace(dateStr)
	// Normalize separators, e.g., replace dots or slashes with hyphens if one format is preferred
	// dateStr = strings.ReplaceAll(dateStr, "/", "-")
	// dateStr = strings.ReplaceAll(dateStr, ".", "-")

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("unable to parse date: '%s'", dateStr)
}
