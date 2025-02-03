package resources

const (
	FINISHED = "completed"
)

type AsyncJob struct {
	Url string `header:"Location"`
}

type AsyncJobStatus struct {
	Status             string  `json:"status,omitempty"`
	PercentageComplete float64 `json:"percentageComplete,omitempty"`
	Operation          string  `json:"operation,omitempty"`
	ResourceId         string  `json:"resourceId,omitempty"`
}
