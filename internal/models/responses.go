package models

type JSONSuccessResultData struct {
	TotalItemCount int         `json:"total_item_count"`
	Data           interface{} `json:"data"`
}

type JSONSuccessResultId struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}
