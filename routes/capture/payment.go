package capture

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Payment Interface
type Payment interface {
	CapturePayment(*models.Approve) (string, error)
}
