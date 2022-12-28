package entity

type EventLog struct {
	Level        string `json:"level" `
	ProviderName string `json:"provider_name" `
	Msg          string `json:"msg"`
	Created      int64  `json:"created" gorm:"UNIQUE"`
}
