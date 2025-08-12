package response

import (
	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/city"
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
)

func FormList(arr []models.Form, pagResp pagination.Response) *cityProto.FormToCreateCityList {
	formList := make([]*cityProto.FormToCreateCity, len(arr))
	for i, form := range arr {
		formList[i] = Form(form)
	}

	return &cityProto.FormToCreateCityList{
		Forms: formList,
		Pagination: &pagProto.Response{
			Page:  pagResp.Page,
			Size:  pagResp.Size,
			Total: pagResp.Total,
		},
	}
}
