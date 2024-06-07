package dto

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Paging struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalRows  int `json:"totalRows"`
	TotalPages int `json:"totalPages"`
}

type ManyResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data"`
	Paging Paging        `json:"paging"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}
