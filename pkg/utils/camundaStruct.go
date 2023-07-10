package utils

type Camunda struct {
	URL  string
	User string
	pw   string
}

func (camunda Camunda) String() string {
	return "{ URL:" + camunda.URL + " User:" + camunda.User + " pw" + camunda.pw + " }"
}
