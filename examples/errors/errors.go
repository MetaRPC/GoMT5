package errors

/*
══════════════════════════════════════════════════════════════════════════════
FILE: errors.go - MT5 Error Handling
══════════════════════════════════════════════════════════════════════════════

This file provides Go-native error handling for MT5 API operations.

ERROR TYPES:
  1. ErrNotConnected - Sentinel error when calling methods before Connect()
  2. ApiError - Wraps protobuf Error with convenient Go methods
  3. Trade return codes - Constants for checking trading operation results

══════════════════════════════════════════════════════════════════════════════
*/

import (
	"errors"
	"fmt"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
)

// ══════════════════════════════════════════════════════════════════════════════
// SENTINEL ERRORS
// ══════════════════════════════════════════════════════════════════════════════

// ErrNotConnected is returned when attempting to call API methods before Connect().
// Use errors.Is(err, mt5.ErrNotConnected) to check for this error.
var ErrNotConnected = errors.New("MT5Account not connected - call Connect() or ConnectEx() first")

// ══════════════════════════════════════════════════════════════════════════════
// API ERROR - Wrapper for protobuf Error
// ══════════════════════════════════════════════════════════════════════════════

// ApiError wraps protobuf Error with convenient Go methods.
// This error is returned when MT5 server responds with an error.
//
// WHEN THIS ERROR OCCURS:
//   - MT5 server responds with an error in the reply message
//   - Trade operation fails (invalid parameters, insufficient margin, etc.)
//   - MQL script execution error on the server side
//   - Invalid request parameters or forbidden operations

type ApiError struct {
	err *pb.Error
}

// NewApiError creates a new ApiError from protobuf Error.
func NewApiError(err *pb.Error) *ApiError {
	if err == nil {
		return nil
	}
	return &ApiError{err: err}
}

// Error implements the error interface.
// Returns a clean, user-friendly error message.
func (e *ApiError) Error() string {
	if e.err == nil {
		return "unknown API error"
	}

	// Prioritize most specific error: Trade Error > MQL Error > API Error
	// Trade errors are the most specific and user-friendly
	if e.MqlErrorTradeIntCode() != 0 {
		return fmt.Sprintf("%s (%s)",
			e.MqlErrorTradeDescription(),
			e.MqlErrorTradeCode())
	}

	// MQL errors are more specific than generic API errors
	if e.MqlErrorIntCode() != 0 {
		return fmt.Sprintf("%s (%s)",
			e.MqlErrorDescription(),
			e.MqlErrorCode())
	}

	// Generic API error
	if e.err.ErrorCode != "" {
		return fmt.Sprintf("API error: %s", e.err.ErrorCode)
	}

	return "unknown API error"
}

// ErrorCode returns the API-level error code (e.g., "INVALID_SYMBOL", "TRADE_DISABLED").
func (e *ApiError) ErrorCode() string {
	if e.err == nil {
		return ""
	}
	return e.err.ErrorCode
}

// ErrorMessage returns the human-readable error description.
func (e *ApiError) ErrorMessage() string {
	if e.err == nil {
		return ""
	}
	return e.err.ErrorMessage
}

// ErrorType returns the error type from the protobuf Error.
func (e *ApiError) ErrorType() pb.ErrorType {
	if e.err == nil {
		return 0 // Default error type
	}
	return e.err.Type
}

// MqlErrorCode returns the MQL5 error code enum.
func (e *ApiError) MqlErrorCode() pb.MqlErrorCode {
	if e.err == nil {
		return pb.MqlErrorCode_ERR_SUCCESS
	}
	return e.err.MqlErrorCode
}

// MqlErrorIntCode returns the MQL5 error code as integer.
func (e *ApiError) MqlErrorIntCode() int32 {
	if e.err == nil {
		return 0
	}
	return e.err.MqlErrorIntCode
}

// MqlErrorDescription returns the MQL5 error description.
func (e *ApiError) MqlErrorDescription() string {
	if e.err == nil {
		return ""
	}
	return e.err.MqlErrorDescription
}

// MqlErrorTradeCode returns the MQL5 trade-specific error code enum.
func (e *ApiError) MqlErrorTradeCode() pb.MqlErrorTradeCode {
	if e.err == nil {
		return 0 // ERR_TRADE_SUCCESS
	}
	return e.err.MqlErrorTradeCode
}

// MqlErrorTradeIntCode returns the MQL5 trade error code as integer.
func (e *ApiError) MqlErrorTradeIntCode() int32 {
	if e.err == nil {
		return 0
	}
	return e.err.MqlErrorTradeIntCode
}

// MqlErrorTradeDescription returns the MQL5 trade error description.
func (e *ApiError) MqlErrorTradeDescription() string {
	if e.err == nil {
		return ""
	}
	return e.err.MqlErrorTradeDescription
}

// RemoteStackTrace returns the server-side stack trace (if available).
func (e *ApiError) RemoteStackTrace() string {
	if e.err == nil {
		return ""
	}
	return e.err.StackTrace
}

// CommandTypeName returns the gRPC method name that caused the error.
func (e *ApiError) CommandTypeName() string {
	if e.err == nil {
		return ""
	}
	return e.err.CommandTypeName
}

// CommandID returns the command ID from the protobuf Error.
func (e *ApiError) CommandID() int64 {
	if e.err == nil {
		return 0
	}
	return e.err.CommandId
}

// Unwrap returns the underlying protobuf Error for inspection.
func (e *ApiError) Unwrap() error {
	return nil // ApiError is a leaf error, doesn't wrap another error
}

// String returns a detailed string representation of the error.
func (e *ApiError) String() string {
	if e.err == nil {
		return "ApiError{nil}"
	}

	return fmt.Sprintf("API Exception: %s - %s\n"+
		"MQL: %s (%d) - %s\n"+
		"Trade: %s (%d) - %s\n"+
		"Command: %s (ID: %d)\n"+
		"Stack: %s",
		e.ErrorCode(), e.ErrorMessage(),
		e.MqlErrorCode(), e.MqlErrorIntCode(), e.MqlErrorDescription(),
		e.MqlErrorTradeCode(), e.MqlErrorTradeIntCode(), e.MqlErrorTradeDescription(),
		e.CommandTypeName(), e.CommandID(),
		e.RemoteStackTrace())
}

// ══════════════════════════════════════════════════════════════════════════════
// TRADE RETURN CODES - Constants for checking trading operation results
// ══════════════════════════════════════════════════════════════════════════════

// Trade operation return codes.
// After EVERY trading operation (OrderSend, OrderModify, etc.), you MUST check
// the ReturnedCode field. Only 10009 (TradeRetCodeDone) means success!
//
// CRITICAL: Trading success/failure is determined by ReturnCode, NOT by exceptions!

const (
	// Success codes
	TradeRetCodeDone        uint32 = 10009 // Request completed successfully
	TradeRetCodeDonePartial uint32 = 10010 // Only part of the request was completed
	TradeRetCodePlaced      uint32 = 10008 // Order placed (pending order activated)

	// Requote codes
	TradeRetCodeRequote      uint32 = 10004 // Requote (price changed, need to retry)
	TradeRetCodePriceChanged uint32 = 10020 // Prices changed (requote)

	// Request rejection codes
	TradeRetCodeReject            uint32 = 10006 // Request rejected
	TradeRetCodeCancel            uint32 = 10007 // Request canceled by trader
	TradeRetCodeInvalidRequest    uint32 = 10013 // Invalid request
	TradeRetCodeInvalidVolume     uint32 = 10014 // Invalid volume in the request
	TradeRetCodeInvalidPrice      uint32 = 10015 // Invalid price in the request
	TradeRetCodeInvalidStops      uint32 = 10016 // Invalid stops in the request (SL/TP too close)
	TradeRetCodeInvalidExpiration uint32 = 10022 // Invalid order expiration date in the request
	TradeRetCodeInvalidFill       uint32 = 10030 // Invalid order filling type
	TradeRetCodeInvalidOrder      uint32 = 10035 // Incorrect or prohibited order type
	TradeRetCodeInvalidCloseVolume uint32 = 10038 // Invalid close volume (exceeds position volume)

	// Trading restriction codes
	TradeRetCodeTradeDisabled    uint32 = 10017 // Trade is disabled
	TradeRetCodeMarketClosed     uint32 = 10018 // Market is closed
	TradeRetCodeServerDisablesAt uint32 = 10026 // Autotrading disabled by server
	TradeRetCodeClientDisablesAt uint32 = 10027 // Autotrading disabled by client terminal
	TradeRetCodeOnlyReal         uint32 = 10032 // Operation is allowed only for live accounts
	TradeRetCodeLongOnly         uint32 = 10042 // Only long positions allowed
	TradeRetCodeShortOnly        uint32 = 10043 // Only short positions allowed
	TradeRetCodeCloseOnly        uint32 = 10044 // Only position close operations allowed
	TradeRetCodeFifoClose        uint32 = 10045 // Position close only by FIFO rule
	TradeRetCodeHedgeProhibited  uint32 = 10046 // Opposite positions on same symbol prohibited (hedging disabled)

	// Resource limit codes
	TradeRetCodeNoMoney        uint32 = 10019 // Not enough money to complete the request (insufficient margin)
	TradeRetCodeLimitOrders    uint32 = 10033 // The number of pending orders has reached the limit
	TradeRetCodeLimitVolume    uint32 = 10034 // The volume of orders and positions for the symbol has reached the limit
	TradeRetCodeLimitPositions uint32 = 10040 // The number of open positions has reached the limit

	// Technical issue codes
	TradeRetCodeError           uint32 = 10011 // Request processing error
	TradeRetCodeTimeout         uint32 = 10012 // Request canceled by timeout
	TradeRetCodeNoQuotes        uint32 = 10021 // No quotes to process the request
	TradeRetCodeTooManyRequests uint32 = 10024 // Too frequent requests
	TradeRetCodeLocked          uint32 = 10028 // Request locked for processing
	TradeRetCodeFrozen          uint32 = 10029 // Order or position frozen
	TradeRetCodeNoConnection    uint32 = 10031 // No connection with the trade server

	// State management codes
	TradeRetCodeOrderChanged    uint32 = 10023 // Order state changed
	TradeRetCodeNoChanges       uint32 = 10025 // No changes in request
	TradeRetCodePositionClosed  uint32 = 10036 // Position with the specified identifier already closed
	TradeRetCodeCloseOrderExist uint32 = 10039 // A close order already exists for a specified position
	TradeRetCodeRejectCancel    uint32 = 10041 // Pending order activation rejected and canceled
)

// IsRetCodeSuccess checks if the return code indicates success.
// Only 10009 (TradeRetCodeDone) is considered success.
func IsRetCodeSuccess(retCode uint32) bool {
	return retCode == TradeRetCodeDone
}

// IsRetCodeRequote checks if the return code indicates a requote (price changed).
// When this happens, you should retry the operation with updated price.
func IsRetCodeRequote(retCode uint32) bool {
	return retCode == TradeRetCodeRequote || retCode == TradeRetCodePriceChanged
}

// IsRetCodeRetryable checks if the error is retryable (temporary).
// These errors might succeed if you retry with exponential backoff.
func IsRetCodeRetryable(retCode uint32) bool {
	switch retCode {
	case TradeRetCodeTimeout,
		TradeRetCodeNoConnection,
		TradeRetCodeFrozen,
		TradeRetCodeLocked,
		TradeRetCodeTooManyRequests,
		TradeRetCodeNoQuotes:
		return true
	}
	return false
}

// GetRetCodeMessage returns a human-readable description for a return code.
func GetRetCodeMessage(retCode uint32) string {
	switch retCode {
	// Success codes
	case TradeRetCodeDone:
		return "Request completed successfully"
	case TradeRetCodeDonePartial:
		return "Only part of the request was completed"
	case TradeRetCodePlaced:
		return "Order placed (pending order activated)"

	// Requote codes
	case TradeRetCodeRequote:
		return "Requote (price changed, need to retry)"
	case TradeRetCodePriceChanged:
		return "Prices changed (requote)"

	// Request rejection codes
	case TradeRetCodeReject:
		return "Request rejected"
	case TradeRetCodeCancel:
		return "Request canceled by trader"
	case TradeRetCodeInvalidRequest:
		return "Invalid request"
	case TradeRetCodeInvalidVolume:
		return "Invalid volume in the request"
	case TradeRetCodeInvalidPrice:
		return "Invalid price in the request"
	case TradeRetCodeInvalidStops:
		return "Invalid stops in the request (SL/TP too close)"
	case TradeRetCodeInvalidExpiration:
		return "Invalid order expiration date in the request"
	case TradeRetCodeInvalidFill:
		return "Invalid order filling type"
	case TradeRetCodeInvalidOrder:
		return "Incorrect or prohibited order type"
	case TradeRetCodeInvalidCloseVolume:
		return "Invalid close volume (exceeds position volume)"

	// Trading restriction codes
	case TradeRetCodeTradeDisabled:
		return "Trade is disabled"
	case TradeRetCodeMarketClosed:
		return "Market is closed"
	case TradeRetCodeServerDisablesAt:
		return "Autotrading disabled by server"
	case TradeRetCodeClientDisablesAt:
		return "Autotrading disabled by client terminal"
	case TradeRetCodeOnlyReal:
		return "Operation is allowed only for live accounts"
	case TradeRetCodeLongOnly:
		return "Only long positions allowed"
	case TradeRetCodeShortOnly:
		return "Only short positions allowed"
	case TradeRetCodeCloseOnly:
		return "Only position close operations allowed"
	case TradeRetCodeFifoClose:
		return "Position close only by FIFO rule"
	case TradeRetCodeHedgeProhibited:
		return "Opposite positions on same symbol prohibited (hedging disabled)"

	// Resource limit codes
	case TradeRetCodeNoMoney:
		return "Not enough money to complete the request (insufficient margin)"
	case TradeRetCodeLimitOrders:
		return "The number of pending orders has reached the limit"
	case TradeRetCodeLimitVolume:
		return "The volume of orders and positions for the symbol has reached the limit"
	case TradeRetCodeLimitPositions:
		return "The number of open positions has reached the limit"

	// Technical issue codes
	case TradeRetCodeError:
		return "Request processing error"
	case TradeRetCodeTimeout:
		return "Request canceled by timeout"
	case TradeRetCodeNoQuotes:
		return "No quotes to process the request"
	case TradeRetCodeTooManyRequests:
		return "Too frequent requests"
	case TradeRetCodeLocked:
		return "Request locked for processing"
	case TradeRetCodeFrozen:
		return "Order or position frozen"
	case TradeRetCodeNoConnection:
		return "No connection with the trade server"

	// State management codes
	case TradeRetCodeOrderChanged:
		return "Order state changed"
	case TradeRetCodeNoChanges:
		return "No changes in request"
	case TradeRetCodePositionClosed:
		return "Position with the specified identifier already closed"
	case TradeRetCodeCloseOrderExist:
		return "A close order already exists for a specified position"
	case TradeRetCodeRejectCancel:
		return "Pending order activation rejected and canceled"

	default:
		return fmt.Sprintf("Unknown return code: %d", retCode)
	}
}
