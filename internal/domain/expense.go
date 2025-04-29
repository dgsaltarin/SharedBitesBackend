package domain

//
//import (
//	"fmt"
//	"time"
//	// Use a decimal library for accurate currency calculations
//	// "github.com/shopspring/decimal"
//)
//
//// Amount represents a monetary value, typically in the smallest currency unit (e.g., cents).
//// Using int64 helps avoid floating-point inaccuracies. Consider using a dedicated decimal library
//// like github.com/shopspring/decimal for more robust currency handling.
//type Amount int64
//
//// Share represents how much a specific user's portion of an expense is.
//type Share struct {
//	UserName UserName
//	Amount   Amount // The amount this user owes/covered for this expense share
//}
//
//// Expense represents a single financial transaction within a group.
//type Expense struct {
//	ID          string    // Unique identifier for the expense (e.g., UUID)
//	GroupID     string    // ID of the group this expense belongs to
//	Description string    // Description of the expense
//	TotalAmount Amount    // The total amount of the expense
//	Currency    string    // Currency code (e.g., "USD", "EUR")
//	PayerName   UserName  // Name of the user who paid the expense
//	Date        time.Time // Date the expense occurred
//	// Participants maps UserName to the amount they are responsible for in this expense.
//	// The sum of amounts in Participants should equal TotalAmount.
//	Participants map[UserName]Amount
//	CreatedAt    time.Time // Timestamp of expense record creation
//	UpdatedAt    time.Time // Timestamp of last update
//}
//
//// NewExpense is a factory function to create a new Expense.
//// It performs validation, including ensuring the payer and participants are valid
//// and that the sum of participant shares equals the total amount.
//func NewExpense(
//	id string,
//	groupID string,
//	description string,
//	totalAmount Amount,
//	currency string,
//	payerName UserName,
//	date time.Time,
//	participants map[UserName]Amount, // Represents who owes what for this expense
//) (*Expense, error) {
//
//	// Basic validation
//	if id == "" {
//		return nil, fmt.Errorf("expense ID cannot be empty: %w", ErrInvalidInput)
//	}
//	if groupID == "" {
//		return nil, fmt.Errorf("group ID cannot be empty: %w", ErrInvalidInput)
//	}
//	if description == "" {
//		return nil, fmt.Errorf("expense description cannot be empty: %w", ErrInvalidInput)
//	}
//	if totalAmount <= 0 { // Usually expenses are positive
//		return nil, fmt.Errorf("expense amount must be positive: %w", ErrInvalidInput)
//	}
//	if currency == "" { // Basic currency check
//		return nil, fmt.Errorf("currency code cannot be empty: %w", ErrInvalidInput)
//	}
//	if payerName == "" {
//		return nil, fmt.Errorf("payer name cannot be empty: %w", ErrInvalidInput)
//	}
//	if len(participants) == 0 {
//		return nil, fmt.Errorf("expense must have at least one participant share: %w", ErrInvalidInput)
//	}
//
//	// Validate participant shares sum up to the total amount
//	var calculatedTotal Amount
//	hasPayerShare := false
//	for name, shareAmount := range participants {
//		if name == "" {
//			return nil, fmt.Errorf("participant name cannot be empty: %w", ErrInvalidInput)
//		}
//		if shareAmount <= 0 {
//			// Allow zero share? Maybe for tracking participation without cost? For now, require positive.
//			return nil, fmt.Errorf("participant share amount for %s must be positive: %w", name, ErrInvalidInput)
//		}
//		calculatedTotal += shareAmount
//		if name == payerName {
//			hasPayerShare = true
//		}
//	}
//
//	if calculatedTotal != totalAmount {
//		return nil, fmt.Errorf("sum of participant shares (%d) does not match total expense amount (%d): %w",
//			calculatedTotal, totalAmount, ErrMismatchedShares) // Add ErrMismatchedShares to errors.go
//	}
//
//	// Application Layer Responsibility: Ensure PayerName and all participant UserNames
//	// are actually members of the group identified by GroupID before calling NewExpense.
//	// The domain object itself cannot easily check this without access to the Group repository.
//
//	// Ensure the payer is included in the participants map if they are also sharing the cost.
//	// This validation might depend on your exact splitting logic (e.g., does the payer always pay for their share upfront?)
//	// For simplicity here, we assume the participants map defines the *final* responsibility split.
//	// If the payer isn't in the map, it means they paid, but their share is covered by others.
//	// If the payer *is* in the map, it reflects their portion of the cost.
//
//	now := time.Now().UTC()
//	return &Expense{
//		ID:           id,
//		GroupID:      groupID,
//		Description:  description,
//		TotalAmount:  totalAmount,
//		Currency:     currency,
//		PayerName:    payerName,
//		Date:         date.UTC(), // Store date in UTC
//		Participants: participants,
//		CreatedAt:    now,
//		UpdatedAt:    now,
//	}, nil
//}
//
//// UpdateExpenseDetails allows modifying certain fields of an existing expense.
//// Be cautious about which fields are updatable. Often Amount/Participants require recalculation logic.
//func (e *Expense) UpdateExpenseDetails(description string, date time.Time) error {
//	if description == "" {
//		return fmt.Errorf("expense description cannot be empty: %w", ErrInvalidInput)
//	}
//	// Add more validation as needed
//	e.Description = description
//	e.Date = date.UTC()
//	e.UpdatedAt = time.Now().UTC()
//	return nil
//}
//
//// Note: Modifying TotalAmount or Participants after creation is complex.
//// It often implies deleting and recreating the expense or having specific
//// "adjustment" logic, potentially involving recalculating balances if tracked separately.
