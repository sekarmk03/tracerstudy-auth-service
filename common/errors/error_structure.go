package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RestError struct {
	Code    codes.Code
	Message string
}

func (re *RestError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", re.Code, re.Message)
}

func ConvertToRestError(err error) error {
	if err == nil {
		return nil
	}

	statusErr, ok := status.FromError(err)
	if !ok {
		// If it's not a gRPC status error, create a generic RestError
		return &RestError{
			Code:    codes.Unknown,
			Message: err.Error(),
		}
	}

	// If it's a gRPC status error, convert it to a RestError
	return &RestError{
		Code:    statusErr.Code(),
		Message: statusErr.Message(),
	}
}