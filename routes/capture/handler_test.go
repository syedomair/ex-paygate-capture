package capture

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/mockserver"
)

const (
	valid_approve_key   = "06F3BCC1C3B836B1AA6D"
	invalid_approve_key = "1D754E20948F3EB8589A9"
)

func TestCaptureAction(t *testing.T) {
	c := Controller{
		Logger: logger.New("DEBUG", "TEST#", os.Stdout),
		Repo:   &mockDB{},
		Pay:    &mockPay{}}

	method := "POST"
	url := "/capture"

	type TestResponse struct {
		Data   string
		Result string
	}

	//Invalid approve_key
	res, req := mockserver.MockTestServer(method, url, []byte(`{"amount":"2", "approve_key":"`+invalid_approve_key+`"}`))
	c.CaptureAction(res, req)
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected := "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Valid approve_key
	res, req = mockserver.MockTestServer(method, url, []byte(`{"amount":"10", "approve_key":"`+valid_approve_key+`"}`))
	c.CaptureAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected = "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
}

type mockPay struct {
}

func (mdb *mockPay) CapturePayment(approveObj *models.Approve, captureAmount string) error {
	return nil
}

type mockDB struct {
}

func (mdb *mockDB) SetRequestID(requestID string) {
}

func (mdb *mockDB) CaptureApprove(inputApproveKey map[string]interface{}) (*models.Approve, error) {
	approveKey := ""
	if approveKeyValue, ok := inputApproveKey["approve_key"]; ok {
		approveKey = approveKeyValue.(string)
	}
	if approveKey != valid_approve_key {
		return nil, errors.New("invalid approve_key")
	}
	return &models.Approve{}, nil
}
