package req

import (
	"demo/purpleSchool/pkg/res"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.Json(*w, err.Error(), 402)
		return nil, err
	}
	err = IsValide(body)
	if err != nil {
		res.Json(*w, err.Error(), 402)
		return nil, err
	}
	return &body, nil
}
