package model

type CertModel struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type CertRequest struct {
	FileName string `json:"fileName"`
}
