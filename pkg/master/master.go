package master

import (
	"encoding/json"
	"errors"

	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/requests"
	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
)

type MasterClient struct {
	api        string
	workerName string
	testName   string
}

func (m *MasterClient) Init(api, workerName, testName string) {
	m.api = api
	m.workerName = workerName
	m.testName = testName
}

func sendStructAsHTTPPOSTRequest(url string, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	status, body, err := requests.MakeJsonHTTPPOSTRequest(url, data)
	if err != nil {
		return err
	} else if status != requests.OK {
		return errors.New(string(body))
	}
	return nil
}

func (m *MasterClient) Enrol(timeToWaitAfterLoggingSec uint) error {
	type EnrolData struct {
		WorkerName                string `json:"worker"`
		TestName                  string `json:"test_name"`
		TimeToWaitAfterLoggingSec uint   `json:"time_waited_after_logging_sec"`
	}
	enrol := EnrolData{m.workerName, m.testName, timeToWaitAfterLoggingSec}
	url := m.api + "/enrol"
	return sendStructAsHTTPPOSTRequest(url, &enrol)
}

func (m *MasterClient) SendLoggingStatus(loggingStatus *types.LoggingState) error {
	url := m.api + "/logging_status?worker=" + m.workerName + "&&test=" + m.testName
	return sendStructAsHTTPPOSTRequest(url, loggingStatus)
}

func (m *MasterClient) SendAnalyseStatus(analyseStatus *types.AnalyseState) error {
	url := m.api + "/analyse_status?worker=" + m.workerName + "&&test=" + m.testName
	return sendStructAsHTTPPOSTRequest(url, analyseStatus)
}
