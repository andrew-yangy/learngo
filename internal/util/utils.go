package util

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"strings"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

func DecodeRequest(req string, dest interface{}) error {
	if err := json.Unmarshal([]byte(req), &dest); err != nil {
		return fmt.Errorf("Cannot decode json: %w", err)
	}
	return nil
}

func DecodeResponse(res string, dest interface{}) error {
	var er V1EngineResponse
	if err := json.Unmarshal([]byte(res), &er); err != nil {
		return err
	}
	bytes, err := json.Marshal(er.Result)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return err
	}
	return nil
}

func toErrResponse(err error) string {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make([]string, len(fieldErrors)),
		}
		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf(`%s is a required field`, err.Field())
			case "bitlength22":
				resp.Errors[i] = fmt.Sprintf(`%s cannot exceed 22 bits`, err.Field())
			case "bitlength31":
				resp.Errors[i] = fmt.Sprintf(`%s cannot exceed 31 bits`, err.Field())
			case "valid_eco_fees":
				resp.Errors[i] = fmt.Sprintf(`%s information is invalid`, err.Field())
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}
		res := fmt.Sprintf("[%s]", strings.Join(resp.Errors, ", "))
		return res
	}
	return ""
}

func HandlePanic(logEntry *logrus.Entry, cb func(interface{})) {
	if r := recover(); r != nil {
		logEntry.Println("Recovering from panic:", r)
		logEntry.Println("Stack Trace:")
		logEntry.Println(string(debug.Stack()))
		if cb != nil {
			cb(r)
		}
	}
}

func InternalError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       InternalApiError(err).Error(),
	}, nil
}

type V1EngineResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
}

func EncodeLambdaError(err error) (events.APIGatewayProxyResponse, error) {
	// Handle error cases
	apiError, ok := err.(*APIError)
	if ok {
		return EncodeError(apiError)
	}
	return InternalError(err)
}

func EncodeLambdaResponse(resp interface{}, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return EncodeLambdaError(err)
	}

	// No error; handle response
	var v1 V1EngineResponse
	v1.Result = resp

	body, err := json.Marshal(v1)
	if err != nil {
		return InternalError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func EncodeError(apiError *APIError) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(apiError)
	if err != nil {
		return InternalError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: apiError.StatusCode,
		Body:       string(body),
	}, nil
}
