package utils

type ListResponseEntre struct {
	StatusCode          int
	Id                  string `json:"id"`
	ProcessDefinitionId string `json:"processDefinitionId"`
	ProcessName         string
	ProcessInstanceId   string `json:"processInstanceId"`
	ExecutionId         string `json:"executionId"`
	IncidentTimestamp   string `json:"incidentTimestamp"`
	IncidentType        string `json:"incidentType"`
	ActivityId          string `json:"activityId"`
	FailedActivityId    string `json:"failedActivityId"`
	CauseIncidentId     string `json:"causeIncidentId"`
	RootCauseIncidentId string `json:"rootCauseIncidentId"`
	Configuration       string `json:"configuration"`
	TenantId            string `json:"tenantId"`
	IncidentMessage     string `json:"incidentMessage"`
	JobDefinitionId     string `json:"jobDefinitionId"`
	Annotation          string `json:"annotation"`
}

type ListResponse []ListResponseEntre

func (listResponse ListResponse) isEmpty() bool {
	if listResponse == nil {
		return true
	}
	return false
}

type ListCountResponse struct {
	StatusCode int
	Count      int `json:"count"`
}

type Definition struct {
	Id                  string      `json:"id"`
	Key                 string      `json:"key"`
	Category            string      `json:"category"`
	Description         interface{} `json:"description"`
	Name                string      `json:"name"`
	Version             int         `json:"version"`
	Resource            string      `json:"resource"`
	DeploymentId        string      `json:"deploymentId"`
	Diagram             interface{} `json:"diagram"`
	Suspended           bool        `json:"suspended"`
	TenantId            interface{} `json:"tenantId"`
	VersionTag          interface{} `json:"versionTag"`
	HistoryTimeToLive   int         `json:"historyTimeToLive"`
	StartableInTasklist bool        `json:"startableInTasklist"`
}
