package main

import (
	"context"
	"fmt"

	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/texttrack"
	"github.com/google/uuid"
)

// Example of how to use enhanced text detection for Spanish bills
func main() {
	// This is an example of how to use the enhanced text detection features
	// In a real application, you would inject these dependencies through your DI container

	ctx := context.Background()

	// Example 1: Using default configuration (good for most cases)
	fmt.Println("=== Example 1: Default Configuration ===")
	analyzeWithDefaultConfig(ctx)

	// Example 2: Using Spanish-optimized configuration
	fmt.Println("\n=== Example 2: Spanish-Optimized Configuration ===")
	analyzeWithSpanishConfig(ctx)

	// Example 3: Using high accuracy configuration for important documents
	fmt.Println("\n=== Example 3: High Accuracy Configuration ===")
	analyzeWithHighAccuracyConfig(ctx)

	// Example 4: Using low confidence configuration for poor quality images
	fmt.Println("\n=== Example 4: Low Confidence Configuration ===")
	analyzeWithLowConfidenceConfig(ctx)

	// Example 5: Custom configuration for specific needs
	fmt.Println("\n=== Example 5: Custom Configuration ===")
	analyzeWithCustomConfig(ctx)
}

func analyzeWithDefaultConfig(ctx context.Context) {
	// This uses the default configuration which is optimized for Spanish bills
	// but still supports English as a fallback

	// In a real application, you would get the billService from your DI container
	// billService := container.GetBillService()

	// billID := uuid.MustParse("your-bill-id")
	// err := billService.AnalyzeBill(ctx, billID) // Uses default config automatically

	fmt.Println("Using default configuration:")
	fmt.Println("- Languages: Spanish first, then English")
	fmt.Println("- Confidence threshold: 70%")
	fmt.Println("- Supported currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU")
}

func analyzeWithSpanishConfig(ctx context.Context) {
	// This configuration is specifically optimized for Spanish bills
	// It uses a lower confidence threshold since Spanish OCR might be less accurate

	// In a real application:
	// config := texttrack.SpanishOptimizedConfig()
	// billID := uuid.MustParse("your-bill-id")
	// err := billService.AnalyzeBillWithConfig(ctx, billID, config)

	fmt.Println("Using Spanish-optimized configuration:")
	fmt.Println("- Languages: Spanish only")
	fmt.Println("- Confidence threshold: 60% (lower for Spanish)")
	fmt.Println("- Supported currencies: EUR, MXN, ARS, CLP, COP, PEN, UYU, USD")
}

func analyzeWithHighAccuracyConfig(ctx context.Context) {
	// This configuration requires higher confidence for more accurate results
	// Useful for important documents where accuracy is critical

	// In a real application:
	// config := texttrack.HighAccuracyConfig()
	// billID := uuid.MustParse("your-bill-id")
	// err := billService.AnalyzeBillWithConfig(ctx, billID, config)

	fmt.Println("Using high accuracy configuration:")
	fmt.Println("- Languages: Spanish first, then English")
	fmt.Println("- Confidence threshold: 85% (higher accuracy)")
	fmt.Println("- Supported currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU")
}

func analyzeWithLowConfidenceConfig(ctx context.Context) {
	// This configuration accepts lower confidence results
	// Useful for poor quality images or unclear text

	// In a real application:
	// config := texttrack.LowConfidenceConfig()
	// billID := uuid.MustParse("your-bill-id")
	// err := billService.AnalyzeBillWithConfig(ctx, billID, config)

	fmt.Println("Using low confidence configuration:")
	fmt.Println("- Languages: Spanish first, then English")
	fmt.Println("- Confidence threshold: 40% (accepts lower quality)")
	fmt.Println("- Supported currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU")
}

func analyzeWithCustomConfig(ctx context.Context) {
	// Create a custom configuration for specific needs
	// In a real application:
	// config := texttrack.TextDetectionConfig{
	//     Languages:     []string{"es", "en", "fr"}, // Spanish, English, French
	//     MinConfidence: 0.75,                       // 75% confidence threshold
	//     CurrencyCodes: []string{"EUR", "USD", "MXN", "ARS"}, // Specific currencies
	// }
	// billID := uuid.MustParse("your-bill-id")
	// err := billService.AnalyzeBillWithConfig(ctx, billID, config)

	fmt.Println("Using custom configuration:")
	fmt.Println("- Languages: Spanish, English, French")
	fmt.Println("- Confidence threshold: 75%")
	fmt.Println("- Supported currencies: EUR, USD, MXN, ARS")
}

// Example of how to handle different types of Spanish bills
func handleDifferentSpanishBills(ctx context.Context) {
	// Example: Restaurant bill from Spain
	restaurantConfig := texttrack.TextDetectionConfig{
		Languages:     []string{"es"},
		MinConfidence: 0.65,
		CurrencyCodes: []string{"EUR"},
	}

	// Example: Mexican restaurant bill
	mexicanConfig := texttrack.TextDetectionConfig{
		Languages:     []string{"es"},
		MinConfidence: 0.6,
		CurrencyCodes: []string{"MXN", "USD"},
	}

	// Example: Argentine restaurant bill
	argentineConfig := texttrack.TextDetectionConfig{
		Languages:     []string{"es"},
		MinConfidence: 0.6,
		CurrencyCodes: []string{"ARS", "USD"},
	}

	fmt.Println("Different configurations for different Spanish-speaking countries:")
	fmt.Printf("Spain: %+v\n", restaurantConfig)
	fmt.Printf("Mexico: %+v\n", mexicanConfig)
	fmt.Printf("Argentina: %+v\n", argentineConfig)
}

// Example of how to retry with different configurations if analysis fails
func retryWithDifferentConfigs(ctx context.Context, billID uuid.UUID) error {
	// Try with default configuration first
	// err := billService.AnalyzeBill(ctx, billID)
	// if err == nil {
	//     return nil
	// }

	// If default fails, try with Spanish-optimized config
	// spanishConfig := texttrack.SpanishOptimizedConfig()
	// err = billService.AnalyzeBillWithConfig(ctx, billID, spanishConfig)
	// if err == nil {
	//     return nil
	// }

	// If that fails, try with low confidence config
	// lowConfidenceConfig := texttrack.LowConfidenceConfig()
	// err = billService.AnalyzeBillWithConfig(ctx, billID, lowConfidenceConfig)
	// if err == nil {
	//     return nil
	// }

	// If all configurations fail, return the last error
	return fmt.Errorf("all text detection configurations failed for bill %s", billID)
}
