package models

type ApiParameter struct {
	ID        string `json:"id"`
	Context   string `json:"context"`
	Arguments JsonB  `json:"arguments"`
	Buffer    JsonB  `json:"buffer"`
}

type ApiResponse struct {
	Action Action           `json:"action,notnull"`
	Buffer *JsonB           `json:"buffer"`
	Data   *ApiResponseData `json:"data"`
}

type ApiResponseData struct {
	Step        *string `json:"step"`
	Interval    *int64  `json:"interval"`
	Comment     *string `json:"comment"`
	NoDecrement *bool   `json:"no_decrement"`
}
