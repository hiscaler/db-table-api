package response

type Error struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Success bool `json:"success"`
}

// {'success': false, 'error': {'message': "error message"}}
type FailResponse struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

type SuccessListData struct {
	Items []interface{}    `json:"items"`
	Meta  map[string]int64 `json:"_meta"`
}

// {'success': true, data: {"items": [{...}, {...}], "_meta": {"totalPage": 1}}}
type SuccessListResponse struct {
	Success bool            `json:"success"`
	Data    SuccessListData `json:"data"`
}

type SuccessOneResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}
