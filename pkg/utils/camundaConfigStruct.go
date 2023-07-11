package utils

type Camunda struct {
	URL      string
	User     string
	Password string
}

func (camunda Camunda) String() string {
	return "{ URL:" + camunda.URL + " User:" + camunda.User + " Password:" + camunda.Password + " }"
}
