package models

type ShortURLCreateDTO struct {
	URL string `json:"url"`
}

type ShortURLReadDTO struct {
	Result string `json:"result"`
}
