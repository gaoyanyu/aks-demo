package request

type CreateBody struct {
	Master  string `json:"master"`
	Version string `json:"version"`
}

type UpdateBody struct {
	CreateBody
}
