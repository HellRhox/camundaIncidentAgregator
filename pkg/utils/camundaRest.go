package utils

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type CamundaRest struct {
	apiClient       http.Client
	baseUrl         string
	apiUrlExtension string
	username        string
	password        string
	basicAuthString string
}

func (camundaRest CamundaRest) CreatClient(url string, user string, password string) CamundaRest {
	return CamundaRest{
		apiClient: http.Client{
			Jar:     &cookiejar.Jar{},
			Timeout: 10 * time.Second,
		},
		baseUrl:         url,
		apiUrlExtension: "/engineRest",
		username:        user,
		password:        password,
		basicAuthString: camundaRest.basicAuth(user, password),
	}

}

func (camundaRest CamundaRest) basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (camundaRest CamundaRest) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+camundaRest.basicAuthString)
	return nil
}

func (camundaRest CamundaRest) GetListOfIncidents(startDate string, enddate string) (error, ListResponse) {
	endpoint := "/incident"
	request, err := http.NewRequest("GET", camundaRest.baseUrl+camundaRest.apiUrlExtension+endpoint, nil)
	if err != nil {
		return err, ListResponse{}
	}
	request.Header.Add("Authorization", "Basic "+camundaRest.basicAuthString)
	q := request.URL.Query()
	q.Add("incidentTimestampBefore", enddate)
	q.Add("incidentTimestampAfter", startDate)
	response, err := camundaRest.apiClient.Do(request)
	if err != nil {
		return err, ListResponse{}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err, ListResponse{}
	}
	var listResponse ListResponse
	json.Unmarshal(body, &listResponse)
	return nil, listResponse
}

func (camundaRest CamundaRest) GetListOfIncidentsCount(startDate string, enddate string) (error, ListCountResponse) {
	endpoint := "/incident/count"
	request, err := http.NewRequest("GET", camundaRest.baseUrl+camundaRest.apiUrlExtension+endpoint, nil)
	if err != nil {
		return err, ListCountResponse{}
	}
	request.Header.Add("Authorization", "Basic "+camundaRest.basicAuthString)
	q := request.URL.Query()
	q.Add("incidentTimestampBefore", enddate)
	q.Add("incidentTimestampAfter", startDate)
	response, err := camundaRest.apiClient.Do(request)
	if err != nil {
		return err, ListCountResponse{}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err, ListCountResponse{}
	}
	var listCountResponse ListCountResponse
	json.Unmarshal(body, &listCountResponse)
	return nil, listCountResponse
}
