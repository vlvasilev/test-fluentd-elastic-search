package server_one

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
)

func setErrorToResponse(errMsg string, w *http.ResponseWriter, status int) {
	log.Printf(errMsg)
	http.Error(*w, errMsg, http.StatusMethodNotAllowed)
}

func checkRequestCorrectnessAndSetResponseOnFailure(w http.ResponseWriter, r *http.Request, method string, path string) bool {
	if r.Method != method {
		setErrorToResponse("Expected ‘"+method+"’ request, got "+r.Method, &w, http.StatusMethodNotAllowed)
		return false
	}

	if r.URL.EscapedPath() != path {
		setErrorToResponse("Expected request to ‘"+path+"’, got "+r.URL.EscapedPath(), &w, http.StatusNotFound)
		return false
	}
	return true
}

func checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w http.ResponseWriter, r *http.Request, key string, value string) bool {
	r.ParseForm()
	keyV := r.Form.Get(key)
	//check if the key exists in case we do not know the value
	if keyV == "" {
		if value == "" {
			value = "<some_value>"
		}
		setErrorToResponse("Expected request to have ‘"+key+"="+value+"’, got empty key/value instead", &w, http.StatusUnprocessableEntity)
		return false
	}
	//in case that we want exactly the pair key/value
	if value != "" && keyV != value {
		setErrorToResponse("Expected request to have ‘"+key+"="+value+"’, got value: "+keyV, &w, http.StatusUnprocessableEntity)
		return false
	}

	return true
}

func loadOjectFromRequestBodyAndSetResponseOnFailure(w http.ResponseWriter, r *http.Request, object interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		setErrorToResponse("Error reading body:"+err.Error(), &w, http.StatusBadRequest)
		return false
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		setErrorToResponse("Error reading body:"+err.Error(), &w, http.StatusBadRequest)
		return false
	}
	return true
}

func setOKResponse(w http.ResponseWriter, msg []byte) {
	w.WriteHeader(http.StatusOK)
	if len(msg) > 0 {
		w.Write(msg)
	}
}

func (soh *ServerOneHandler) Enrol(w http.ResponseWriter, r *http.Request) {

	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/enrol") {
		return
	}

	type EnrolData struct {
		WorkerName                string `json:"worker"`
		TestName                  string `json:"test_name"`
		TimeToWaitAfterLoggingSec uint   `json:"time_waited_after_logging_sec"`
	}

	var enroll EnrolData
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &enroll) {
		return
	}

	worker.WorkerName = enroll.WorkerName
	worker.TimeToWaitAfterLoggingSec = enroll.TimeToWaitAfterLoggingSec

	soh.mutex.Lock()
	test, ok := soh.tests[enroll.TestName]
	defer soh.mutex.Unlock()

	if !ok {
		setErrorToResponse("There is no such test!", &w, http.StatusNotFound)
		return
	}

	_, ok = test.WorkerStats[enroll.WorkerName]
	if ok {
		setErrorToResponse("Worker already enroled", &w, http.StatusConflict)
		return
	}

	test.NumberOfEnroledWorkers++
	test.WorkerStats[enroll.WorkerName] = worker
	soh.tests[enroll.TestName] = test
	setOKResponse(w, []byte("Enroled!"))
}

func (soh *ServerOneHandler) SetLogginStatus(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/logging_status") {
		return
	}

	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "worker", "") {
		return
	}
	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "test", "") {
		return
	}

	var loggingStatus types.LoggingState
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &loggingStatus) {
		return
	}

	workerName := r.Form.Get("worker")
	testName := r.Form.Get("test")

	soh.mutex.Lock()
	test, ok := soh.tests[testName]
	defer soh.mutex.Unlock()

	if !ok {
		setErrorToResponse("There is no shuch test: "+testName, &w, http.StatusUnprocessableEntity)
		return
	}

	worker, ok = test.WorkerStats[workerName]

	if !ok {
		setErrorToResponse("Worker not enroled", &w, http.StatusUnprocessableEntity)
		return
	}
	test.NumberOfFinishedWorkers++
	worker.LoggingState = loggingStatus
	test.WorkerStats[workerName] = worker
	soh.tests[testName] = test

	setOKResponse(w, []byte("Loggins Status Saves!"))
}

func (soh *ServerOneHandler) GetAllStatus(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "GET", "/status/all") {
		return
	}

	var index int

	soh.mutex.Lock()
	tests := make([]Test, len(soh.tests))

	for _, test := range soh.tests {
		tests[index] = test
		index++
	}
	soh.mutex.Unlock()

	responseData, err := json.Marshal(tests)
	if err != nil {
		setErrorToResponse("Server can't Marshal Status Data", &w, http.StatusInternalServerError)
		return
	}
	setOKResponse(w, responseData)
}

func (soh *ServerOneHandler) SetAnalyseStatus(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/analyse_status") {
		return
	}

	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "worker", "") {
		return
	}
	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "test", "") {
		return
	}

	var analyseStatus types.AnalyseState
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &analyseStatus) {
		return
	}

	workerName := r.Form.Get("worker")
	testName := r.Form.Get("test")

	soh.mutex.Lock()
	test, ok := soh.tests[testName]
	defer soh.mutex.Unlock()

	if !ok {
		setErrorToResponse("There is no shuch test: "+testName, &w, http.StatusUnprocessableEntity)
		return
	}

	worker, ok = test.WorkerStats[workerName]

	if !ok {
		setErrorToResponse("Worker not enroled", &w, http.StatusUnprocessableEntity)
		return
	}

	worker.AnalyseState = analyseStatus

	soh.tests[testName].WorkerStats[workerName] = worker

	setOKResponse(w, []byte("Analyse Status Saves!"))
}

func (soh *ServerOneHandler) stopCurrentTest(testName, namespace string) {
	//TODO:find a way to implement this method
	soh.kubeClient.DeleteJob(namespace, testName)
}

func (soh *ServerOneHandler) deployPods(testData TestInitResources) ([]byte, error) {
	pods := int32(testData.Pods)
	elep := testData.ElasticAPI
	master := soh.GetListeningEndPoint()
	msg := strconv.Itoa(int(testData.MsgCount))
	logtime := strconv.Itoa(int(testData.LogtimeMs))
	ttlaw := strconv.Itoa(int(testData.TimeToWaitAterLoggingSec))
	namespace := testData.Namespace
	testName := testData.TestName
	return soh.kubeClient.DeployJob(pods, testName, namespace, logtime, msg, ttlaw, elep, master)
	//return util.Exe_cmd("./deploy_job.sh " + pods + " " + elep + " " + master + " " + msg + " " + logtime + " " + ttlaw + " " + namespace)
	//TODO: set soh.testStatus.NumberOfPod = testDAta.Pods
	//TODO: use command to verify the succseed of the deployment
}

func (soh *ServerOneHandler) isSuchTest(name string) bool {
	if _, ok := soh.tests[name]; ok {
		return true
	}
	return false
}

func (soh *ServerOneHandler) allocateTest(resources TestInitResources) error {
	newTest := Test{}
	newTest.Name = resources.TestName
	newTest.WorkerStats = make(map[string]WorkerStats)
	newTest.NumberOfWorkers = resources.Pods
	newTest.Namespace = resources.Namespace
	newTest.StartTime = time.Now()
	maxTestPeriod := uint64(resources.LogtimeMs/1000) + uint64(resources.TimeToWaitAterLoggingSec) + uint64(time.Second*300)
	newTest.EndTime = newTest.StartTime.Add(time.Duration(maxTestPeriod * uint64(time.Second)))

	soh.mutex.Lock()
	defer soh.mutex.Unlock()
	if _, ok := soh.tests[newTest.Name]; ok {
		return errors.New("such test already exists")
	}

	soh.tests[newTest.Name] = newTest
	return nil
}

func (soh *ServerOneHandler) destroyTest(testName, namespace string) error {
	soh.stopCurrentTest(testName, namespace)
	soh.mutex.Lock()
	defer soh.mutex.Unlock()
	delete(soh.tests, testName)
	return nil
}

func (soh *ServerOneHandler) startTest(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/test/start") {
		return
	}

	var resources TestInitResources
	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &resources) {
		return
	}

	if err := soh.allocateTest(resources); err != nil {
		setErrorToResponse("Such test already exists. Please stop it first and than try again!", &w, http.StatusUnprocessableEntity)
		return
	}

	data, err := soh.deployPods(resources)
	if err != nil {
		setErrorToResponse("Can't deploy workers: "+string(err.Error())+" "+string(data), &w, http.StatusUnprocessableEntity)
		soh.destroyTest(resources.TestName, resources.Namespace)
		return
	}

	setOKResponse(w, []byte("Test is runing!\n"+string(data)))
}

func (soh *ServerOneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := soh.hadlerFunctions[r.URL.EscapedPath()]; ok {
		handler(w, r)
		return
	}

	setErrorToResponse("Unknow resource path: "+r.URL.EscapedPath(), &w, http.StatusNotFound)
}

func (soh *ServerOneHandler) Init(address, port, kubeconfig string) {
	soh.serverStats.address = address
	soh.serverStats.listeningPort = port
	soh.hadlerFunctions = make(map[string]func(http.ResponseWriter, *http.Request))
	soh.tests = make(map[string]Test)
	soh.kubeClient.Init(kubeconfig)
	soh.hadlerFunctions["/enrol"] = soh.Enrol
	soh.hadlerFunctions["/logging_status"] = soh.SetLogginStatus
	soh.hadlerFunctions["/analyse_status"] = soh.SetAnalyseStatus
	soh.hadlerFunctions["/status/all"] = soh.GetAllStatus
	soh.hadlerFunctions["/test/start"] = soh.startTest
}

func (soh *ServerOneHandler) GetListeningEndPoint() string {
	return soh.serverStats.address + ":" + soh.serverStats.listeningPort
}

//TODO remove this function
// func (soh *ServerOneHandler) clearStatus(testName string) {
// 	test, ok := soh.tests[testName]
// 	if !ok {
// 		return
// 	}
// 	test.workerStats = make(map[string]WorkerStats)
// 	test.NumberOfEnroledWorkers = 0
// 	test.NumberOfFinishedWorkers = 0
// 	test.NumberOfWorkers = 0
// 	test.Done = false
// 	test.passed = true
// 	test.startTime = time.Now()
// 	test.endTime = test.startTime
// }
