package capture

import (
	"net/http"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"
	"github.com/syedomair/ex-paygate-lib/lib/tools/response"
)

const (
	errorCodePrefix       = "03"
	Amount                = "amount"
	ApproveKey            = "approve_key"
	ApprovedAmountBalance = "approved_amount_balance"
	Currency              = "currency"
)

// Controller Public
type Controller struct {
	Logger logger.Logger
	Repo   Repository
	Pay    Payment
}

//var httpClient = &http.Client{}

// Ping Public
func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	methodName := "Ping"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()
	responseToken := map[string]string{"response": "captureController pong"}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseToken, http.StatusOK)
}

// CaptureAction Public
func (c *Controller) CaptureAction(w http.ResponseWriter, r *http.Request) {
	methodName := "CaptureAction"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()

	paramConf := make(map[string]models.ParamConf)
	paramConf[Amount] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[ApproveKey] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}

	paramMap, errCode, err := request.ValidateInputParameters(r, request.GetRequestID(r), c.Logger, paramConf, nil)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+errCode, err.Error(), http.StatusBadRequest)
		return
	}

	approveObj, err := c.Repo.CaptureApprove(paramMap)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"2", err.Error(), http.StatusBadRequest)
		return
	}

	responseAction := map[string]string{ApprovedAmountBalance: approveObj.AmountBalance, Currency: approveObj.Currency}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseAction, http.StatusOK)
}
