package elastic

import (
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/requests"
)

type ElasticSearchClient struct {
	api string
}

func (e *ElasticSearchClient) Init(api string) {
	e.api = api
}

func (e *ElasticSearchClient) IsActive() bool {
	status, _, err := requests.MakeHTTPGetRequest(e.api)
	if status != requests.OK || err != nil {
		return false
	}
	return true
}

func (e *ElasticSearchClient) IsSuchIndex(index string) bool {
	status, _, err := requests.MakeHTTPGetRequest(e.api + "/" + index)
	if status != requests.OK || err != nil {
		return false
	}
	return true
}

func (e *ElasticSearchClient) GetIndexRecords(index, logtype, key, value string) ([]byte, error) {
	searchQuery := []byte("{\"query\":{\"term\":{\"" + key + "\":\"" + value + "\"}},\"_source\":[\"worker\",\"message\",\"sentence\",\"text\"],\"sort\":[{\"message\":{\"order\":\"asc\"}},{\"sentence\":{\"order\":\"asc\"}}],\"from\":0,\"size\":10000}")
	rawurl := e.api + "/" + index + "/" + logtype + "/_search"
	url, _ := url.Parse(rawurl)

	status, body, err := requests.MakeJsonHTTPPOSTRequest(url.String(), searchQuery)
	if status != requests.OK || err != nil {
		return []byte{}, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return []byte{}, err
	}

	numhit, ok := data["hits"].(map[string]interface{})["total"].(float64)
	if !ok {
		return []byte{}, errors.New("can't extract hits from response")
	}

	numberOfhits := int(numhit)

	if numberOfhits >= 10000 {
		numberOfhits = 10000
	}

	if numberOfhits < 1 {
		return []byte{}, errors.New("no records avaliable")
	}

	hits, ok := data["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return []byte{}, errors.New("can't extract hits from response")
	}

	maps := make([]map[string]interface{}, numberOfhits)

	for index, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		delete(source, key)
		maps[index] = source
	}

	result, err := json.Marshal(maps)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

func (e *ElasticSearchClient) GetRegordsFromServer(index, logtype, key, value string) ([]byte, error) {
	if !e.IsActive() {
		return []byte{}, errors.New("elastich-search server is unreachable")
	}
	if !e.IsSuchIndex(index) {
		return []byte{}, errors.New("there in no such index")
	}

	return e.GetIndexRecords(index, logtype, key, value)
}

func (e *ElasticSearchClient) GetIndexRecordNumbers(index, logtype, key, value string) (int, error) {
	rawurl := e.api + "/" + index + "/" + logtype + "/_search"
	url, _ := url.Parse(rawurl)
	searchQuery := []byte("{\"query\":{\"term\":{\"" + key + "\":\"" + value + "\"}},\"size\":0}")
	status, body, err := requests.MakeJsonHTTPPOSTRequest(url.String(), searchQuery)
	if status != requests.OK || err != nil {
		return -1, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return -1, err
	}

	hits, ok := data["hits"].(map[string]interface{})["total"].(float64)
	if !ok {
		return -1, errors.New("can't extract hits from response")
	}

	return int(hits), nil
}

func (e *ElasticSearchClient) GetLogstashIndex() string {
	current_time := time.Now().Local()
	return string(current_time.AppendFormat([]byte(`logstash-`), "2006.01.02"))
}

func (e *ElasticSearchClient) GetCurrenLogstashRecords(logtype, key, value string) ([]byte, error) {
	index := e.GetLogstashIndex()
	return e.GetRegordsFromServer(index, logtype, key, value)
}
