package server_one

import (
	"net/http"
	"sync"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
	"k8s.io/client-go/kubernetes"
)

//WorkerStats represents the logging and analysing status one worker(pod)
type WorkerStats struct {
	WorkerName                string             `json:"worker_name"`
	TimeToWaitAfterLoggingSec uint               `json:"time_waited_after_logging_sec"`
	TimeFirstLogIsFoundMs     uint64             `json:"elapsed_time_after_firs_record"`
	LoggingState              types.LoggingState `json:"logging_state"`
	AnalyseState              types.AnalyseState `json:"analyse_state"`
}

//ServerStats represents the endpoin of the server
type ServerStats struct {
	listeningPort string
	address       string
}

//TestStatus represents the common status of a test
type TestStatus struct {
	NumberOfWorkers         uint `json:"number_of_workers"`
	NumberOfEnroledWorkers  uint `json:"number_of_enroled_workers"`
	NumberOfFinishedWorkers uint `json:"number_of_finished_workers"`
	Done                    bool `json:"done"`
	Passed                  bool `json:"passed"`
}

//Test represents a test which takes in account the stsatus of every worker involved in the test
type Test struct {
	Name      string    `json:"test_name"`
	Namespace string    `json:"namespace"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	TestStatus
	WorkerStats map[string]WorkerStats `json:"workers_status"`
}

//ServerOneHandler merges hadlers functions, test statuses, k8s clinet, server endpoint and a mutex.
//The mutex is used to syncronize the tests map. One must use mutex when dealing with that map
type ServerOneHandler struct {
	hadlerFunctions map[string]func(http.ResponseWriter, *http.Request)
	//workerStats     map[string]WorkerStats
	tests       map[string]Test
	kubeClient  KubeClient
	serverStats ServerStats
	testStatus  TestStatus
	mutex       sync.Mutex
}

//TestInitResources represent the information needed to be start a new test
type TestInitResources struct {
	TestName                 string `json:"test_name"`
	Pods                     uint   `json:"pods"`
	ElasticAPI               string `json:"elastic_api"`
	MsgCount                 uint   `json:"msgcount"`
	LogtimeMs                uint   `json:"logtime_ms"`
	TimeToWaitAterLoggingSec uint   `json:"time_to_wait_after_logging_sec"`
	Namespace                string `json:"namespace"`
}

//KubeClient holds a kubernetes.Clientset and is used to interact with the k8s API
type KubeClient struct {
	clientset *kubernetes.Clientset
}
