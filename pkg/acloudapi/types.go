package acloudapi

type Pageable struct {
	Sort       Sort `json:"sort"`
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Offset     int  `json:"offset"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
}

type Sort struct {
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
	Empty    bool `json:"empty"`
}

type Error struct {
	Message string `json:"message"`
}

type PagedResult struct {
	Content          []interface{} `json:"content"`
	Pageable         interface{}   `json:"pageable"`
	Last             bool          `json:"last"`
	TotalPages       int           `json:"totalPages"`
	TotalElements    int           `json:"totalElements"`
	NumberOfElements int           `json:"numberOfElements"`
	Number           int           `json:"number"`
	Sort             Sort          `json:"sort"`
	First            bool          `json:"first"`
	Size             int           `json:"size"`
	Empty            bool          `json:"empty"`
}
