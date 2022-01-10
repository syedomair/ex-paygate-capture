package capture

import (
	"errors"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type PaymentService struct {
	logger    logger.Logger
	requestID string
}

const (
	CaptureFailureCCNumber = "4000000000000259"
)

// NewPaymentService Public.
func NewPaymentService(logger logger.Logger) Payment {
	return &PaymentService{logger: logger}
}

// CapturePayment Public.
func (payWrap *PaymentService) CapturePayment(approveObj *models.Approve, captureAmount string) error {
	methodName := "CapturePayment"
	payWrap.logger.Debug(payWrap.requestID, "M:%v start", methodName)
	start := time.Now()

	if approveObj.CCNumber == CaptureFailureCCNumber {
		return errors.New("capture failure")
	}

	payWrap.logger.Debug(payWrap.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return nil
}
