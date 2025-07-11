package domain

import "errors"

// General Domain Errors (examples)
var (
	ErrInvalidInput = errors.New("invalid input provided")
	ErrNotFound     = errors.New("requested resource not found")
	// ... other general errors
)

// Group Specific Errors
var (
	ErrGroupNotFound              = errors.New("group not found")
	ErrGroupNameEmpty             = errors.New("group name cannot be empty")
	ErrGroupMembersEmpty          = errors.New("group must have at least one member")
	ErrGroupMemberNameEmpty       = errors.New("group member name cannot be empty")
	ErrUserNotFound               = errors.New("user not found")
	ErrAlreadyMember              = errors.New("user is already a member of the group")
	ErrNotMember                  = errors.New("user is not a member of the group")
	ErrCannotRemoveOwner          = errors.New("cannot remove the group owner")
	ErrPermissionDenied           = errors.New("permission denied for this operation")
	ErrUserNameEmpty              = errors.New("user name cannot be empty")
	ErrUserEmailEmpty             = errors.New("user email cannot be empty")
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrFirebaseUserCreationFailed = errors.New("failed to create user in Firebase")
	ErrDatabaseUserCreationFailed = errors.New("failed to create user in database")
	ErrUserPasswordTooShort       = errors.New("user password is too short")
	ErrFirebaseIDEmpty            = errors.New("firebase ID cannot be empty")
	ErrFirebaseUserAlreadyExists  = errors.New("firebase user already exists")
	ErrBillNotFound               = errors.New("bill not found")
	ErrUserIDEmpty                = errors.New("user ID cannot be empty for a bill")
	ErrBillIDEmpty                = errors.New("bill ID cannot be empty for a line item")
	ErrLineItemDescriptionEmpty   = errors.New("line item description cannot be empty")
	ErrBillProcessingFailed       = errors.New("bill processing failed")
	ErrFileUploadFailed           = errors.New("file upload failed")
	ErrFileDownloadFailed         = errors.New("file download failed")
	ErrTextractAnalysisFailed     = errors.New("Textract analysis failed")
	ErrTextractDataExtraction     = errors.New("failed to extract data from Textract result")
)

// Text Analysis Errors (as previously defined)
var (
	ErrTextAnalysisFailed = errors.New("text analysis failed")
	// ... other text analysis errors
)

// Queue/Filestore Errors (if added)
var (
	ErrQueueFailed   = errors.New("failed to interact with the queue")
	ErrStorageFailed = errors.New("failed to interact with the file store")
)

// Helper function (optional) for checking specific error types if needed elsewhere
func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, ErrGroupNotFound) || errors.Is(err, ErrUserNotFound)
}
