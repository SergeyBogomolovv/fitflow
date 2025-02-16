package content

import "github.com/SergeyBogomolovv/fitflow/pkg/httpx"

type GenerateContentResponse struct {
	Status  httpx.Status `json:"status"`
	Content string       `json:"content"`
}
