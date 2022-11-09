package dto

type Task struct {
	ID                 uint32 `json:"id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Status             string `json:"status"`
	EstimationInSecond int    `json:"estimation_in_second"`
}
