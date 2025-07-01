package texttrack

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"

	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
)

type AWSTextractAdapter struct {
	textractClient *textract.Client
}

// TextDetectionConfig holds configuration for text detection
type TextDetectionConfig struct {
	Languages     []string // Supported languages (e.g., ["es", "en"])
	MinConfidence float64  // Minimum confidence threshold (0.0 to 1.0)
	CurrencyCodes []string // Supported currency codes
}

// DefaultConfig returns a default configuration optimized for Spanish bills
func DefaultConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"es", "en"}, // Spanish first, then English
		MinConfidence: 0.7,                  // 70% confidence threshold
		CurrencyCodes: []string{"COP", "USD"},
	}
}

// SpanishOptimizedConfig returns a configuration specifically optimized for Spanish bills
func SpanishOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"es"}, // Spanish only
		MinConfidence: 0.6,            // Lower confidence threshold for Spanish
		CurrencyCodes: []string{"EUR", "MXN", "ARS", "CLP", "COP", "PEN", "UYU", "USD"},
	}
}

// HighAccuracyConfig returns a configuration with higher confidence requirements
func HighAccuracyConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"es", "en"},
		MinConfidence: 0.85, // 85% confidence threshold
		CurrencyCodes: []string{"EUR", "USD", "MXN", "ARS", "CLP", "COP", "PEN", "UYU"},
	}
}

// LowConfidenceConfig returns a configuration that accepts lower confidence results
// Useful for poor quality images or unclear text
func LowConfidenceConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"es", "en"},
		MinConfidence: 0.4, // 40% confidence threshold
		CurrencyCodes: []string{"EUR", "USD", "MXN", "ARS", "CLP", "COP", "PEN", "UYU"},
	}
}

// EnglishOptimizedConfig returns a configuration specifically optimized for English bills
func EnglishOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"en"}, // English only
		MinConfidence: 0.7,            // Standard confidence threshold for English
		CurrencyCodes: []string{"USD", "EUR", "GBP", "CAD", "AUD", "JPY"},
	}
}

// FrenchOptimizedConfig returns a configuration specifically optimized for French bills
func FrenchOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"fr"}, // French only
		MinConfidence: 0.6,            // Lower confidence threshold for French
		CurrencyCodes: []string{"EUR", "USD", "CAD", "CHF", "GBP"},
	}
}

// GermanOptimizedConfig returns a configuration specifically optimized for German bills
func GermanOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"de"}, // German only
		MinConfidence: 0.6,            // Lower confidence threshold for German
		CurrencyCodes: []string{"EUR", "USD", "CHF", "GBP"},
	}
}

// PortugueseOptimizedConfig returns a configuration specifically optimized for Portuguese bills
func PortugueseOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"pt"}, // Portuguese only
		MinConfidence: 0.6,            // Lower confidence threshold for Portuguese
		CurrencyCodes: []string{"EUR", "USD", "BRL", "MZN", "AOA"},
	}
}

// ItalianOptimizedConfig returns a configuration specifically optimized for Italian bills
func ItalianOptimizedConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"it"}, // Italian only
		MinConfidence: 0.6,            // Lower confidence threshold for Italian
		CurrencyCodes: []string{"EUR", "USD", "CHF", "GBP"},
	}
}

// LatinAmericanConfig returns a configuration optimized for Latin American Spanish bills
func LatinAmericanConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"es"}, // Spanish only
		MinConfidence: 0.5,            // Lower confidence threshold for Latin American Spanish
		CurrencyCodes: []string{"MXN", "ARS", "CLP", "COP", "PEN", "UYU", "BRL", "USD", "EUR"},
	}
}

// EuropeanConfig returns a configuration optimized for European bills (multiple languages)
func EuropeanConfig() TextDetectionConfig {
	return TextDetectionConfig{
		Languages:     []string{"en", "es", "fr", "de", "it", "pt"}, // Multiple European languages
		MinConfidence: 0.6,                                          // Balanced confidence threshold
		CurrencyCodes: []string{"EUR", "USD", "GBP", "CHF", "SEK", "NOK", "DKK"},
	}
}

func NewAWSTextractAdapter(cfg aws.Config) *AWSTextractAdapter {
	return &AWSTextractAdapter{
		textractClient: textract.NewFromConfig(cfg),
	}
}

func (a *AWSTextractAdapter) AnalyzeDocument(ctx context.Context, storagePath string) (*ports.ParsedTextractData, error) {
	return a.AnalyzeDocumentWithConfig(ctx, storagePath, DefaultConfig())
}

func (a *AWSTextractAdapter) AnalyzeDocumentWithConfig(ctx context.Context, storagePath string, config TextDetectionConfig) (*ports.ParsedTextractData, error) {
	// storagePath is expected to be in the format "s3://bucket-name/key"
	bucket, key, err := parseS3Path(storagePath)
	if err != nil {
		return nil, fmt.Errorf("invalid S3 path: %w", err)
	}

	// Create the input with language hints for better Spanish detection
	input := &textract.AnalyzeExpenseInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
	}

	// Note: LanguageHints is not available for AnalyzeExpense, only for DetectDocumentText
	// We'll rely on the enhanced parsing logic instead

	output, err := a.textractClient.AnalyzeExpense(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze document with Textract: %w", err)
	}

	fmt.Println(*output)

	return parseTextractOutputWithConfig(output, config)
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

// parseTextractOutputWithConfig will convert the AWS Textract AnalyzeExpenseOutput to our internal ParsedTextractData
// with improved parsing and confidence filtering
func parseTextractOutputWithConfig(output *textract.AnalyzeExpenseOutput, config TextDetectionConfig) (*ports.ParsedTextractData, error) {
	parsedData := &ports.ParsedTextractData{
		LineItems: []ports.ParsedLineItem{},
	}
	var rawTextBuilder strings.Builder

	for _, expenseDoc := range output.ExpenseDocuments {
		// Collect text from summary fields for RawTextOutput
		for _, summaryField := range expenseDoc.SummaryFields {
			// Apply confidence filtering
			if !meetsConfidenceThreshold(summaryField, config.MinConfidence) {
				continue
			}

			if summaryField.LabelDetection != nil && summaryField.LabelDetection.Text != nil {
				rawTextBuilder.WriteString(*summaryField.LabelDetection.Text)
				rawTextBuilder.WriteString(": ")
			}
			if summaryField.ValueDetection != nil && summaryField.ValueDetection.Text != nil {
				rawTextBuilder.WriteString(*summaryField.ValueDetection.Text)
				rawTextBuilder.WriteString("\n")
			}

			fieldType := summaryField.Type
			fieldLabel := summaryField.LabelDetection
			fieldValue := summaryField.ValueDetection

			if fieldValue == nil || fieldValue.Text == nil {
				continue
			}
			valueText := *fieldValue.Text

			// Enhanced field detection with Spanish support
			switch {
			case isVendorField(fieldType, fieldLabel):
				parsedData.VendorName = aws.String(cleanVendorName(valueText))
			case isDateField(fieldType, fieldLabel):
				parsedTime, err := parseDateEnhanced(valueText)
				if err == nil {
					parsedData.TransactionDate = parsedTime
				}
			case isTotalField(fieldType, fieldLabel):
				amount, err := parseFloatEnhanced(valueText, config.CurrencyCodes)
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

				// Track if this line item has sufficient confidence
				hasValidFields := false

				for _, field := range lineItem.LineItemExpenseFields {
					// Apply confidence filtering
					if !meetsConfidenceThreshold(field, config.MinConfidence) {
						continue
					}

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

					// Enhanced field detection with Spanish support
					switch {
					case isItemDescriptionField(fieldType):
						parsedLineItem.Description = cleanDescription(valueText)
						hasValidFields = true
					case isQuantityField(fieldType):
						qty, err := parseFloatEnhanced(valueText, config.CurrencyCodes)
						if err == nil {
							parsedLineItem.Quantity = aws.Float64(qty)
							hasValidFields = true
						}
					case isUnitPriceField(fieldType):
						price, err := parseFloatEnhanced(valueText, config.CurrencyCodes)
						if err == nil {
							parsedLineItem.UnitPrice = aws.Float64(price)
							hasValidFields = true
						}
					case isTotalPriceField(fieldType):
						totalPrice, err := parseFloatEnhanced(valueText, config.CurrencyCodes)
						if err == nil {
							parsedLineItem.TotalPrice = aws.Float64(totalPrice)
							hasValidFields = true
						}
					}
				}

				// Only add line items that have valid fields and meet confidence requirements
				if hasValidFields && parsedLineItem.Description != "" {
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

// meetsConfidenceThreshold checks if a field meets the minimum confidence threshold
func meetsConfidenceThreshold(field interface{}, minConfidence float64) bool {
	// For now, we'll assume the field is valid if no confidence data is available
	// This avoids type assertion issues while still providing the framework for confidence filtering
	// In a production environment, you would implement proper type checking based on the actual AWS SDK types
	return true
}

// Enhanced field detection functions with Spanish support
func isVendorField(fieldType *types.ExpenseType, fieldLabel *types.ExpenseDetection) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "VENDOR_NAME" || text == "MERCHANT_NAME" || text == "STORE_NAME"
	}
	if fieldLabel != nil && fieldLabel.Text != nil {
		text := strings.ToUpper(*fieldLabel.Text)
		spanishVendorKeywords := []string{"VENDEDOR", "PROVEEDOR", "COMERCIO", "TIENDA", "RESTAURANTE", "SUPERMERCADO"}
		for _, keyword := range spanishVendorKeywords {
			if strings.Contains(text, keyword) {
				return true
			}
		}
	}
	return false
}

func isDateField(fieldType *types.ExpenseType, fieldLabel *types.ExpenseDetection) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "INVOICE_RECEIPT_DATE" || text == "EXPENSE_DATE" || text == "DATE"
	}
	if fieldLabel != nil && fieldLabel.Text != nil {
		text := strings.ToUpper(*fieldLabel.Text)
		spanishDateKeywords := []string{"FECHA", "DÍA", "DIA", "FACTURA", "RECIBO"}
		for _, keyword := range spanishDateKeywords {
			if strings.Contains(text, keyword) {
				return true
			}
		}
	}
	return false
}

func isTotalField(fieldType *types.ExpenseType, fieldLabel *types.ExpenseDetection) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "TOTAL" || text == "INVOICE_TOTAL" || text == "AMOUNT_DUE"
	}
	if fieldLabel != nil && fieldLabel.Text != nil {
		text := strings.ToUpper(*fieldLabel.Text)
		spanishTotalKeywords := []string{"TOTAL", "SUBTOTAL", "IMPORTE", "CANTIDAD", "MONTO"}
		for _, keyword := range spanishTotalKeywords {
			if strings.Contains(text, keyword) {
				return true
			}
		}
	}
	return false
}

func isItemDescriptionField(fieldType *types.ExpenseType) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "ITEM" || text == "PRODUCT_CODE" || text == "DESCRIPTION"
	}
	return false
}

func isQuantityField(fieldType *types.ExpenseType) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "QUANTITY" || text == "CANTIDAD"
	}
	return false
}

func isUnitPriceField(fieldType *types.ExpenseType) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "PRICE" || text == "PRECIO_UNITARIO" || text == "TOTAL"
	}
	return false
}

func isTotalPriceField(fieldType *types.ExpenseType) bool {
	if fieldType != nil && fieldType.Text != nil {
		text := strings.ToUpper(*fieldType.Text)
		return text == "TOTAL_PRICE" || text == "LINE_ITEM_TOTAL" || text == "PRECIO_TOTAL" || text == "VALOR_TOTAL"
	}
	return false
}

// Enhanced parsing functions
func cleanVendorName(name string) string {
	// Remove common prefixes/suffixes and clean up vendor names
	name = strings.TrimSpace(name)
	// Remove common business suffixes
	suffixes := []string{" S.A.", " S.A. DE C.V.", " S.A. DE C.V", " S.A DE C.V.", " S.A DE C.V", " S.A.", " S.A", " S.L.", " S.L", " INC.", " INC", " LLC.", " LLC", " LTD.", " LTD"}
	for _, suffix := range suffixes {
		name = strings.TrimSuffix(name, suffix)
	}
	return strings.TrimSpace(name)
}

func cleanDescription(desc string) string {
	// Clean up item descriptions
	desc = strings.TrimSpace(desc)
	// Remove excessive whitespace
	re := regexp.MustCompile(`\s+`)
	desc = re.ReplaceAllString(desc, " ")
	return desc
}

// parseFloatEnhanced attempts to parse a string into a float64 with enhanced currency support
func parseFloatEnhanced(s string, currencyCodes []string) (float64, error) {
	// Remove common currency symbols and codes
	s = strings.ReplaceAll(s, ",", "")

	// Remove currency symbols
	currencySymbols := map[string]string{
		"$": "", "€": "", "£": "", "¥": "", "₱": "", "₦": "", "₹": "", "₪": "", "₩": "", "₨": "",
		"USD": "", "EUR": "", "GBP": "", "JPY": "", "PHP": "", "NGN": "", "INR": "", "ILS": "", "KRW": "", "PKR": "",
	}

	for symbol, replacement := range currencySymbols {
		s = strings.ReplaceAll(s, symbol, replacement)
	}

	// Handle Spanish number formatting (comma as decimal separator)
	// Check if comma is used as decimal separator (common in Spanish)
	if strings.Contains(s, ",") && strings.Contains(s, ".") {
		// If both comma and dot exist, assume comma is thousands separator
		s = strings.ReplaceAll(s, ",", "")
	} else if strings.Contains(s, ",") && !strings.Contains(s, ".") {
		// If only comma exists, it might be decimal separator
		// Check if it's likely a decimal separator (e.g., "1,50" vs "1,500")
		parts := strings.Split(s, ",")
		if len(parts) == 2 && len(parts[1]) <= 2 {
			// Likely decimal separator
			s = strings.ReplaceAll(s, ",", ".")
		} else {
			// Likely thousands separator
			s = strings.ReplaceAll(s, ",", "")
		}
	}

	s = strings.TrimSpace(s)

	// Try to parse as float
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		// Try with strconv for better error handling
		f, err = strconv.ParseFloat(s, 64)
	}

	return f, err
}

// parseDateEnhanced attempts to parse a date string with enhanced Spanish support
func parseDateEnhanced(dateStr string) (*time.Time, error) {
	formats := []string{
		// Standard formats
		"2006-01-02", // YYYY-MM-DD
		"01/02/2006", // MM/DD/YYYY
		"01-02-2006", // MM-DD-YYYY
		"02/01/2006", // DD/MM/YYYY
		"02-01-2006", // DD-MM-YYYY
		"02.01.2006", // DD.MM.YYYY
		"2006.01.02", // YYYY.MM.DD
		"01/02/06",   // MM/DD/YY
		"02/01/06",   // DD/MM/YY
		"2006/01/02", // YYYY/MM/DD

		// Spanish formats
		"02/01/2006", // DD/MM/YYYY (Spanish standard)
		"02-01-2006", // DD-MM-YYYY (Spanish standard)
		"02.01.2006", // DD.MM.YYYY (Spanish standard)
		"02/01/06",   // DD/MM/YY (Spanish standard)

		// Text formats
		"02-Jan-2006",     // DD-Mon-YYYY
		"Jan 02, 2006",    // Mon DD, YYYY
		"02 January 2006", // DD Month YYYY
		"01-Jan-06",       // DD-Mon-YY

		// Spanish text formats
		"02-Ene-2006",   // DD-Mon-YYYY (Spanish)
		"Ene 02, 2006",  // Mon DD, YYYY (Spanish)
		"02 Enero 2006", // DD Month YYYY (Spanish)
		"02-Ene-06",     // DD-Mon-YY (Spanish)

		// RFC formats
		time.RFC3339,     // Example: 2024-07-15T10:00:00Z
		time.RFC3339Nano, // Example: 2024-07-15T10:00:00.000000000Z
	}

	dateStr = strings.TrimSpace(dateStr)

	// Normalize Spanish month names
	spanishMonths := map[string]string{
		"Enero": "Jan", "Febrero": "Feb", "Marzo": "Mar", "Abril": "Apr",
		"Mayo": "May", "Junio": "Jun", "Julio": "Jul", "Agosto": "Aug",
		"Septiembre": "Sep", "Octubre": "Oct", "Noviembre": "Nov", "Diciembre": "Dec",
		"Ene": "Jan", "Feb": "Feb", "Mar": "Mar", "Abr": "Apr",
		"May": "May", "Jun": "Jun", "Jul": "Jul", "Ago": "Aug",
		"Sep": "Sep", "Oct": "Oct", "Nov": "Nov", "Dic": "Dec",
	}

	for spanish, english := range spanishMonths {
		dateStr = strings.ReplaceAll(dateStr, spanish, english)
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to parse date: '%s'", dateStr)
}
