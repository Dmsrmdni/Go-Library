package models

type Response struct {
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Paginate interface{} `json:"paginate"`
}

type ResponseDetail struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
