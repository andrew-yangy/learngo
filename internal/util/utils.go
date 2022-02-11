package util

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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
