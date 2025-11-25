package link

type LinkCreateRequest struct {
	Url string `json:"URL" validate:"required,url"`
}
