package models

type ShortURLCreateDTO struct {
	URL string `json:"url"`
}

type ShortURLReadDTO struct {
	Result string `json:"result"`
}

type SerializeData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchShortURLCreateDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string
}

type BatchShortURLReadDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
