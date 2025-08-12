package response

import (
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	formProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/form"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
)

func FormList(arr []models.Form, pagResp pagination.Response) *formProto.FormToCreateCityList {
	formList := make([]*formProto.FormToCreateCity, len(arr))
	for i, form := range arr {
		formList[i] = Form(form)
	}

	return &formProto.FormToCreateCityList{
		Forms: formList,
		Pagination: &pagProto.Response{
			Page:  pagResp.Page,
			Size:  pagResp.Size,
			Total: pagResp.Total,
		},
	}
}
