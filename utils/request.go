package utils

import (
	"context"
	"fmt"
	"net/http"
)

const CancelRequestKey = "cancel_request"

func CancelRequest(r *http.Request) error {
	value := r.Context().Value(CancelRequestKey)

	cancelRequest, ok := value.(context.CancelFunc)

	if !ok {
		return fmt.Errorf("invalid type of %T", cancelRequest)
	}

	cancelRequest()
	return nil
}
