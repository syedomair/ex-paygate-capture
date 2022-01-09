package capture

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Repository interface
type Repository interface {
	SetRequestID(requestID string)
	CaptureApprove(inputApproveKey map[string]interface{}) (*models.Approve, error)
}
