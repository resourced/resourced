package executors

type AgentLoglinePayload struct {
	Created int64
	Content string
}

type AgentLogPayload struct {
	Host struct {
		Name string
		Tags map[string]string
	}
	Data struct {
		Loglines []AgentLoglinePayload
		Filename string
	}
}
