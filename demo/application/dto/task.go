package dto

type Task struct {
	ID                 uint32 `json:"id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	EstimationInSecond int    `json:"estimation_in_second"`
}
