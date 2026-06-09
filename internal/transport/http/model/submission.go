package model

type ResponseModel struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type SubmitRequest struct {
	Language       string `json:"language"`
	CodeSubmission string `json:"code_submission"`
}
