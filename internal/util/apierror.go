package util

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	InternalServerErrorCode                   = "internal_server_error"
	BadRequestErrorCode                       = "bad_request"
	ConflictErrorCode                         = "conflict_error"
	TransactionInsertErrorCode                = "transaction_insert_failed"
	EventInsertErrorCode                      = "event_insert_failed"
	AccountMintErrorCode                      = "account_mint_failed"
	AccountNotFoundErrorCode                  = "account_not_found"
	AccountStarkKeyAlreadyRegisteredErrorCode = "account_stark_key_already_registered"
	ClaimNotFoundErrorCode                    = "claim_not_found"
	TokenNotFoundErrorCode                    = "token_not_found"
	InvalidEtherKeyErrorCode                  = "invalid_ether_key"
	AssetInsertFailedErrorCode                = "asset_insert_failed"
	AssetNotFoundErrorCode                    = "asset_not_found"
	VaultNotFoundErrorCode                    = "vault_not_found"
	VaultInsertErrorCode                      = "vault_insert_failed"
	InsufficientBalanceErrorCode              = "insufficient_balance"
	MintValidationErrorCode                   = "mint_validation_failed"
	MintCountErrorCode                        = "mint_count_error"
	EndpointNotSupportedCode                  = "endpoint_not_supported_code"
	VersionNotSupportedCode                   = "version_not_supported_code"
	OrderHasBeenFilled                        = "order_has_been_filled"
	InsufficientBalanceForSellAmountErrorCode = "insufficient_balance_for_sell_amount"
	FieldUnparseableError                     = "field_not_parseable"
	InvalidEthSignatureErrorCode              = "invalid_eth_signature"
	InvalidStarkSignatureErrorCode            = "invalid_stark_signature"
	MintUnwithdrawableCode                    = "mint_unwithdrawable"
	ContractAlreadyExistsErrorCode            = "contract_address_already_exists"
	ContractInvalidCode                       = "contract_address_invalid"
	EthClientErrorCode                        = "ethereum_client_error"
	RoyaltiesAddressNotRegisteredErrorCode    = "royalties_address_not_registered"
	QuantizationErrorCode                     = "token_quantization_failed"
	ResourceNotFoundCode                      = "resource_not_found_code"
	InputValidationFailedCode                 = "input_validation_failed"
	InvalidExpirationTimestampCode            = "invalid_expiration_timestamp"
	UnableToConvertAmountToIntegerCode        = "unable_to_convert_amount_to_integer"
)

type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"status_code"`
}

func (e *APIError) Error() string {
	return e.Message
}

func InternalApiError(err error) *APIError {
	log.Errorf("Internal server error: %s", err.Error())

	return &APIError{
		Code:       InternalServerErrorCode,
		Message:    "The server encountered an internal error and was unable to process your request",
		Details:    "",
		StatusCode: http.StatusInternalServerError,
	}
}

var InputValidationFailed = &APIError{
	Code:       InputValidationFailedCode,
	Message:    "There was a problem validating API input, please check the documentation and try again",
	Details:    retrieveDetailsByErrorCode(InputValidationFailedCode),
	StatusCode: http.StatusBadRequest,
}

var InvalidExpirationTimestamp = &APIError{
	Code:       InvalidExpirationTimestampCode,
	Message:    "Invalid expiration timestamp for order, must be at least 1 week in the future",
	Details:    retrieveDetailsByErrorCode(InvalidExpirationTimestampCode),
	StatusCode: http.StatusBadRequest,
}

func BadRequestApiError(msg string) *APIError {
	return &APIError{
		Code:       BadRequestErrorCode,
		Message:    msg,
		Details:    "",
		StatusCode: http.StatusBadRequest,
	}
}

//ApiConflictError returns a HTTP 409 CONFLICT error with the provided message msg
func ApiConflictError(msg string) *APIError {
	return &APIError{
		Code:       ConflictErrorCode,
		Message:    msg,
		Details:    "",
		StatusCode: http.StatusConflict,
	}
}

// TODO: For now, details is empty.
func retrieveDetailsByErrorCode(errCode string) string {
	return ""
}

func BadRequestWithCodeApiError(errorCode string, err error) *APIError {
	return &APIError{
		Code:       errorCode,
		Message:    err.Error(),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func AccountNotFoundApiError(address string) *APIError {
	errorCode := AccountNotFoundErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Account not found: %s", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusNotFound,
	}
}

func AccountStarkKeyAlreadyRegisteredError() *APIError {
	errorCode := AccountStarkKeyAlreadyRegisteredErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    "Stark key already registered",
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func ClaimNotFoundApiError(address string) *APIError {
	errorCode := ClaimNotFoundErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Claim not found for address: %s", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusNotFound,
	}
}

func TokenNotFoundApiError(address string) *APIError {
	errorCode := TokenNotFoundErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Token with address %s could not be found.", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusNotFound,
	}
}

func InvalidEtherKeyApiError(address string) *APIError {
	errorCode := InvalidEtherKeyErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Ethereum address provided is empty and/or invalid: %s", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func AssetIdNotFoundApiError(asset string) *APIError {
	errorCode := AssetNotFoundErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Asset not found: %s", asset),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func VaultNotFoundApiError(vault string) *APIError {
	errorCode := VaultNotFoundErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Vault not found: %s", vault),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func InsufficientBalanceApiError(amount uint64, balance uint64) *APIError {
	errorCode := InsufficientBalanceErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Insufficient balance. Transaction amount: %d, balance: %d", amount, balance),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func InsufficientBalanceForSellAmountError(details string) *APIError {
	errorCode := InsufficientBalanceForSellAmountErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    details,
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func MintLimitExceededApiError(numAssets, maxMintsInRequest int) *APIError {
	errorCode := MintCountErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("minting of %d assets exceeds single-request limit of %d", numAssets, maxMintsInRequest),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func VersionNotSupportedApiError() *APIError {
	errorCode := VersionNotSupportedCode
	return &APIError{
		Code:       errorCode,
		Message:    "version not supported",
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func EndpointNotSupportedApiError() *APIError {
	errorCode := EndpointNotSupportedCode
	return &APIError{
		Code:       errorCode,
		Message:    "endpoint not supported",
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func OrderHasBeenFilledApiError(orderID uint64) *APIError {
	errorCode := OrderHasBeenFilled
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Order %d has been filled already", orderID),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func FieldNotParseableAPIError(fieldName string, value string) *APIError {
	errorCode := FieldUnparseableError
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("could not parse field %s with value %s", fieldName, value),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func InvalidEthSignature(signature string) *APIError {
	errorCode := InvalidEthSignatureErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("unable to recover eth address from signature %s", signature),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func InvalidStarkSignature(signature string) *APIError {
	errorCode := InvalidStarkSignatureErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("unable to verify stark signature %s", signature),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func ContractAlreadyExists(address string) *APIError {
	errorCode := ContractAlreadyExistsErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("contract with address %s already exists", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func ContractInvalid(address string) *APIError {
	errorCode := ContractInvalidCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("invalid contract bytecode, address %s is not a contract", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func EthClientError() *APIError {
	errorCode := EthClientErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    "unable to process eth client request",
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusNotFound,
	}
}

func RoyaltiesAddressNotRegistered(address string) *APIError {
	errorCode := RoyaltiesAddressNotRegisteredErrorCode
	return &APIError{
		Code:       errorCode,
		Message:    fmt.Sprintf("Royalties address not registered: %s", address),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func QuantizationError(quantizationTarget string, amount string, quantum uint64, err error) *APIError {
	errorCode := QuantizationErrorCode
	return &APIError{
		Code: errorCode,
		//quantizationTarget is the thing we're trying to quantize. E.g. fee exclusive buy amount, sell amount, fee etc.
		Message:    fmt.Sprintf("error quantizing %s (%s) with quantum %d: %+v", quantizationTarget, amount, quantum, err),
		Details:    retrieveDetailsByErrorCode(errorCode),
		StatusCode: http.StatusBadRequest,
	}
}

func ResourceNotFound(resourceType string, resourceId string) *APIError {
	return &APIError{
		Code:       ResourceNotFoundCode,
		Message:    fmt.Sprintf("%s id '%s' not found", resourceType, resourceId),
		Details:    retrieveDetailsByErrorCode(ResourceNotFoundCode),
		StatusCode: http.StatusNotFound,
	}
}

func UnableToConvertAmountToInteger(amount string) *APIError {
	return &APIError{
		Code:       UnableToConvertAmountToIntegerCode,
		Message:    fmt.Sprintf("error when attempting to set amount (%s) to integer", amount),
		Details:    retrieveDetailsByErrorCode(UnableToConvertAmountToIntegerCode),
		StatusCode: http.StatusInternalServerError,
	}
}
