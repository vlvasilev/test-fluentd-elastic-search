package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type WorkerStats struct {
	WorkerName                string             `json:worker_name`
	TimeToWaitAfterLoggingSec uint               `json:"time_waited_after_logging_sec"`
	TimeFirstLogIsFoundMs     uint64             `json"elapsed_time_after_firs_record"`
	LoggingState              types.LoggingState `json:logging_state`
	AnalyseState              types.AnalyseState `json:"analyse_state`
}

type ServerStats struct {
	listeningPort string
	address       string
}

type TestStatus struct {
	NumberOfWorkers         uint `json:"number_of_workers"`
	NumberOfEnroledWorkers  uint `json:"number_of_enroled_workers`
	NumberOfFinishedWorkers uint `json:"number_of_finished_workers`
	Done                    bool `json:"done"`
}

type ServerOneHandler struct {
	hadlerFunctions map[string]func(http.ResponseWriter, *http.Request)
	workerStats     map[string]WorkerStats
	kubeClient      KubeClient
	serverStats     ServerStats
	testStatus      TestStatus
	mutex           sync.Mutex
}

type TestInitResources struct {
	Pods                     uint   `json:"pods"`
	ElasticAPI               string `json:"elastic_api"`
	MsgCount                 uint   `json:"msgcount"`
	LogtimeMs                uint   `json:"logtime_ms"`
	TimeToWaitAterLoggingSec uint   `json:"time_to_wait_after_logging_sec"`
	Namespace                string `json:"namespace"`
}

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

	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "worker", "") {
		return
	}

	type EnrolData struct {
		WorkerName                string `json:"worker"`
		TimeToWaitAfterLoggingSec uint   `json:"time_waited_after_logging_sec"`
	}

	var enroll EnrolData
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &enroll) {
		return
	}

	entry := r.Form.Get("worker")

	soh.mutex.Lock()
	worker, ok := soh.workerStats[entry]
	defer soh.mutex.Unlock()

	if ok {
		setErrorToResponse("Worker already enroled", &w, http.StatusConflict)
		return
	}

	worker.WorkerName = enroll.WorkerName
	worker.TimeToWaitAfterLoggingSec = enroll.TimeToWaitAfterLoggingSec

	soh.workerStats[entry] = worker
	setOKResponse(w, []byte("Enroled!"))
}

func (soh *ServerOneHandler) SetLogginStatus(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/logging_status") {
		return
	}

	if !checkRequestForKeyValueCorrectnesAndSetResponseOnFailure(w, r, "worker", "") {
		return
	}

	var loggingStatus types.LoggingState
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &loggingStatus) {
		return
	}

	entry := r.Form.Get("worker")

	soh.mutex.Lock()
	worker, ok := soh.workerStats[entry]
	defer soh.mutex.Unlock()

	if !ok {
		setErrorToResponse("Worker not enroled", &w, http.StatusUnprocessableEntity)
		return
	}

	worker.LoggingState = loggingStatus

	soh.workerStats[entry] = worker
	setOKResponse(w, []byte("Loggins Status Saves!"))
}

func (soh *ServerOneHandler) GetAllStatus(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "GET", "/status/all") {
		return
	}

	var index int

	soh.mutex.Lock()
	statuses := make([]WorkerStats, len(soh.workerStats))

	for _, status := range soh.workerStats {
		statuses[index] = status
		index++
	}
	soh.mutex.Unlock()

	responseData, err := json.Marshal(&statuses)
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

	var analyseStatus types.AnalyseState
	var worker WorkerStats

	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &analyseStatus) {
		return
	}

	entry := r.Form.Get("worker")

	soh.mutex.Lock()
	worker, ok := soh.workerStats[entry]
	defer soh.mutex.Unlock()

	if !ok {
		setErrorToResponse("Worker not enroled", &w, http.StatusUnprocessableEntity)
		return
	}

	worker.AnalyseState = analyseStatus
	soh.workerStats[entry] = worker

	setOKResponse(w, []byte("Analyse Status Saves!"))
}

func (soh *ServerOneHandler) clearStatus() {
	soh.workerStats = make(map[string]WorkerStats)
	soh.testStatus.NumberOfEnroledWorkers = 0
	soh.testStatus.NumberOfFinishedWorkers = 0
	soh.testStatus.NumberOfWorkers = 0
	soh.testStatus.Done = false
}

func (soh *ServerOneHandler) stopCurrentTest() {
	//TODO:find a way to implement this method
}

func (soh *ServerOneHandler) deployPods(testData TestInitResources) ([]byte, error) {
	pods := int32(testData.Pods)
	elep := testData.ElasticAPI
	master := soh.getListeningEndPoint()
	msg := strconv.Itoa(int(testData.MsgCount))
	logtime := strconv.Itoa(int(testData.LogtimeMs))
	ttlaw := strconv.Itoa(int(testData.TimeToWaitAterLoggingSec))
	namespace := testData.Namespace
	return soh.kubeClient.DeployJob(pods, namespace, logtime, msg, ttlaw, elep, master)
	//return util.Exe_cmd("./deploy_job.sh " + pods + " " + elep + " " + master + " " + msg + " " + logtime + " " + ttlaw + " " + namespace)
	//TODO: set soh.testStatus.NumberOfPod = testDAta.Pods
	//TODO: use command to verify the succseed of the deployment
}

func (soh *ServerOneHandler) startTest(w http.ResponseWriter, r *http.Request) {
	if !checkRequestCorrectnessAndSetResponseOnFailure(w, r, "POST", "/test/start") {
		return
	}

	var resources TestInitResources
	if !loadOjectFromRequestBodyAndSetResponseOnFailure(w, r, &resources) {
		return
	}

	soh.stopCurrentTest()
	soh.clearStatus()

	data, err := soh.deployPods(resources)
	if err != nil {
		setErrorToResponse("Can't deploy workers: "+string(err.Error())+" "+string(data), &w, http.StatusUnprocessableEntity)
		return
	}

	setOKResponse(w, []byte("Test is runing!"+string(data)))
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
	soh.workerStats = make(map[string]WorkerStats)
	soh.kubeClient.Init(kubeconfig)
	soh.hadlerFunctions["/enrol"] = soh.Enrol
	soh.hadlerFunctions["/logging_status"] = soh.SetLogginStatus
	soh.hadlerFunctions["/analyse_status"] = soh.SetAnalyseStatus
	soh.hadlerFunctions["/status/all"] = soh.GetAllStatus
	soh.hadlerFunctions["/test/start"] = soh.startTest
}

func (soh *ServerOneHandler) getListeningEndPoint() string {
	return soh.serverStats.address + ":" + soh.serverStats.listeningPort
}

func main() {
	var myIP = flag.String("address", "0.0.0.0", "ip on wich the server will listen")
	var myPort = flag.String("port", "8000", "port on wich the server will listen")
	var kubeconfig = flag.String("kubeconfig", "./conf/kubeconfig.yaml", "The config file for the kubernetes kluster")
	flag.Parse()

	serverOneHandler := ServerOneHandler{}
	serverOneHandler.Init(*myIP, *myPort, *kubeconfig)

	server := http.Server{
		Addr:    serverOneHandler.getListeningEndPoint(),
		Handler: &serverOneHandler,
	}
	log.Println("Server Start To Listen On: " + serverOneHandler.getListeningEndPoint())
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Server Exit Normaly")
	}

}

type KubeClient struct {
	clientset *kubernetes.Clientset
}

func (k *KubeClient) Init(kubeconfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		k.clientset = nil
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *KubeClient) DeployJob(numberOfPods int32, namespace, logtime, msgcount, TimeToWaitAfterLoggingSec, alasticAPI, master string) ([]byte, error) {
	if k.clientset == nil {
		return []byte{}, errors.New("missing or unvalid kubeconfig file")
	}
	jobClient := k.clientset.BatchV1().Jobs(namespace)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flood-and-analyse",
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "flood-and-analyse",
				"role":    "test",
				"section": namespace,
			},
		},
		Spec: batchv1.JobSpec{
			// Selector: &metav1.LabelSelector{
			// 	MatchLabels: map[string]string{
			// 		"app":     "flood-and-analyse",
			// 		"role":    "test",
			// 		"section": namespace,
			// 	},
			// },
			Parallelism: int32Ptr(numberOfPods),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     "flood-and-analyse",
						"role":    "test",
						"section": namespace,
					},
				},
				Spec: apiv1.PodSpec{
					RestartPolicy: apiv1.RestartPolicyNever,
					Containers: []apiv1.Container{
						{
							Name:  "flood-and-anlyse",
							Image: "hisshadow85/flood-and-analyse:1.1",
							Env: []apiv1.EnvVar{
								{
									Name: "POD_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							Command: []string{
								"./flood_and_analyse",
								"--workername=$(POD_NAME)",
								"--logtime=" + logtime,
								"--msgcount=" + msgcount,
								"--time_to_wait_after_logging_sec=" + TimeToWaitAfterLoggingSec,
								"--elastic_end_point=" + alasticAPI,
								"--master=" + master,
							},
						},
					},
				},
			},
		},
	}
	result, err := jobClient.Create(job)
	if err != nil {
		return []byte{}, err
	}
	return []byte("Created job " + result.GetObjectMeta().GetName()), nil
}

func int32Ptr(i int32) *int32 { return &i }

// func main() {
// 	server := http.Server{
// 		Addr:    ":8000",
// 		Handler: &myHandler{},
// 	}

// 	mux = make(map[string]func(http.ResponseWriter, *http.Request))
// 	mux["/"] = hello

// 	server.ListenAndServe()
// }

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

var mux map[string]func(http.ResponseWriter, *http.Request)
