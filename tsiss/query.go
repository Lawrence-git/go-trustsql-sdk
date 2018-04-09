package tsiss

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

var (
	client   *http.Client
	validate *validator.Validate
)

const (
	// QueryIssURI 共享信息查询
	QueryIssURI = "https://baas.trustsql.qq.com/cgi-bin/v1.0/trustsql_iss_query_v1.cgi"
)

func init() {
	client = &http.Client{
		Timeout: 1 * time.Second,
	}
	validate = validator.New()
}

// QueryIss 共享信息查询
func QueryIss(iss *IssQuery) (*IssResponse, error) {
	// 校验common是否符合标准
	err := validate.Struct(iss)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(iss)

	log.Printf("trustsql request data is %s", string(data))

	// send http request
	req, err := http.NewRequest("POST", QueryIssURI, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	// log.Printf("trustsql response body is %s", string(body))
	_ = resp.Body.Close()

	issResponse := IssResponse{}
	err = json.Unmarshal(body, &issResponse)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(issResponse)
	if err != nil {
		return nil, err
	}
	return &issResponse, nil
}
