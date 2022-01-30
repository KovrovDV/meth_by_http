package methbyhttp

import (
	"net/http"
)

type TBaseRequest struct {
	pRequest http.Request
	arParams []TReqParam
}
