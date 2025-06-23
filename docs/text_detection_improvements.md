# Enhanced Text Detection for Spanish Bills

This document describes the improvements made to the text detection system to better handle Spanish bills and receipts.

## Overview

The text detection system has been enhanced with several improvements specifically designed to improve accuracy for Spanish-language bills and receipts:

1. **Language-specific optimizations**
2. **Enhanced currency parsing**
3. **Improved date parsing for Spanish formats**
4. **Better field detection with Spanish keywords**
5. **Configurable confidence thresholds**
6. **Vendor name cleaning**

## Key Improvements

### 1. Language Support

The system now prioritizes Spanish language detection and includes Spanish-specific field keywords:

- **Vendor fields**: `VENDEDOR`, `PROVEEDOR`, `COMERCIO`, `TIENDA`, `RESTAURANTE`, `SUPERMERCADO`
- **Date fields**: `FECHA`, `DÍA`, `DIA`, `FACTURA`, `RECIBO`
- **Total fields**: `TOTAL`, `SUBTOTAL`, `IMPORTE`, `CANTIDAD`, `MONTO`

### 2. Enhanced Currency Parsing

Improved support for Spanish and Latin American currencies:

- **Currency symbols**: €, $, ₱, ₦, ₹, ₪, ₩, ₨
- **Currency codes**: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU
- **Number formatting**: Handles both comma-as-decimal (1,50) and comma-as-thousands (1,500) separators

### 3. Spanish Date Formats

Support for common Spanish date formats:

- `DD/MM/YYYY` (Spanish standard)
- `DD-MM-YYYY` (Spanish standard)
- `DD.MM.YYYY` (Spanish standard)
- Spanish month names: `Enero`, `Febrero`, `Marzo`, etc.
- Abbreviated months: `Ene`, `Feb`, `Mar`, etc.

### 4. Vendor Name Cleaning

Automatic removal of common business suffixes:

- `S.A.`, `S.A. DE C.V.`, `S.L.`, `INC.`, `LLC.`, `LTD.`

## Configuration Options

### Default Configuration

```go
config := texttrack.DefaultConfig()
// Languages: ["es", "en"] (Spanish first, then English)
// Confidence: 70%
// Currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU
```

### Spanish-Optimized Configuration

```go
config := texttrack.SpanishOptimizedConfig()
// Languages: ["es"] (Spanish only)
// Confidence: 60% (lower threshold for Spanish)
// Currencies: EUR, MXN, ARS, CLP, COP, PEN, UYU, USD
```

### High Accuracy Configuration

```go
config := texttrack.HighAccuracyConfig()
// Languages: ["es", "en"]
// Confidence: 85% (higher accuracy)
// Currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU
```

### Low Confidence Configuration

```go
config := texttrack.LowConfidenceConfig()
// Languages: ["es", "en"]
// Confidence: 40% (accepts lower quality)
// Currencies: EUR, USD, MXN, ARS, CLP, COP, PEN, UYU
```

## Usage Examples

### Basic Usage (Default Configuration)

```go
// Uses the default configuration automatically
err := billService.AnalyzeBill(ctx, billID)
```

### Custom Configuration

```go
// Use Spanish-optimized configuration
config := texttrack.SpanishOptimizedConfig()
err := billService.AnalyzeBillWithConfig(ctx, billID, config)
```

### Retry Strategy

```go
// Try with different configurations if analysis fails
err := billService.AnalyzeBill(ctx, billID)
if err != nil {
    // Try with Spanish-optimized config
    spanishConfig := texttrack.SpanishOptimizedConfig()
    err = billService.AnalyzeBillWithConfig(ctx, billID, spanishConfig)
    if err != nil {
        // Try with low confidence config
        lowConfidenceConfig := texttrack.LowConfidenceConfig()
        err = billService.AnalyzeBillWithConfig(ctx, billID, lowConfidenceConfig)
    }
}
```

## Country-Specific Configurations

### Spain

```go
config := texttrack.TextDetectionConfig{
    Languages:     []string{"es"},
    MinConfidence: 0.65,
    CurrencyCodes: []string{"EUR"},
}
```

### Mexico

```go
config := texttrack.TextDetectionConfig{
    Languages:     []string{"es"},
    MinConfidence: 0.6,
    CurrencyCodes: []string{"MXN", "USD"},
}
```

### Argentina

```go
config := texttrack.TextDetectionConfig{
    Languages:     []string{"es"},
    MinConfidence: 0.6,
    CurrencyCodes: []string{"ARS", "USD"},
}
```

## Best Practices

### 1. Image Quality

- Ensure good lighting and contrast
- Avoid blurry or low-resolution images
- Make sure the entire bill is visible

### 2. Configuration Selection

- Use `SpanishOptimizedConfig()` for Spanish-only bills
- Use `HighAccuracyConfig()` for important documents
- Use `LowConfidenceConfig()` for poor quality images

### 3. Error Handling

- Implement retry logic with different configurations
- Log analysis failures for debugging
- Provide user feedback on image quality issues

### 4. Testing

- Test with various Spanish bill formats
- Test with different currencies and date formats
- Test with different image qualities

## Troubleshooting

### Common Issues

1. **Low accuracy on Spanish bills**

   - Try `SpanishOptimizedConfig()` with lower confidence threshold
   - Check image quality and lighting

2. **Currency parsing errors**

   - Verify the currency is supported in the configuration
   - Check for mixed number formatting (comma vs dot)

3. **Date parsing failures**

   - Ensure the date format is supported
   - Check for mixed date formats in the same document

4. **Missing vendor information**
   - The system automatically cleans vendor names
   - Check if the vendor name contains unusual characters

### Debugging

Enable detailed logging to see:

- Confidence scores for each field
- Parsed values before and after cleaning
- Field detection results

## Performance Considerations

- Higher confidence thresholds may result in fewer detected fields
- Lower confidence thresholds may include more false positives
- Language-specific configurations may be slower but more accurate
- Consider caching results for frequently analyzed bills

## Future Improvements

1. **Machine Learning**: Train custom models on Spanish bill datasets
2. **Template Matching**: Add support for common Spanish bill templates
3. **Regional Variations**: Add country-specific optimizations
4. **Real-time Feedback**: Provide immediate feedback on image quality
5. **Batch Processing**: Optimize for processing multiple bills simultaneously
