package types

type ExecSuccessResponse struct {
	Result                ContainerResponseData
	ExecutionTime         int64
	InternalExecutionTime int64
	Meta                  FunctionExecutionMetaData
}

type ContainerResponse struct {
	StartTime int64                   `json:"startTime"`
	EndTime   int64                   `json:"endTime"`
	Data      ContainerResponseData   `json:"result"`
}

type ContainerResponseData struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type FunctionExecutionMetaData struct {
	ImageBuilt                 bool
	UsingPooledContainer       bool
	UsingExistingRestContainer bool
	ContainerName              string
	ImageName                  string
}
