package api

import (
	"github.com/kararnab/authdemo/pkg/iam"
	internalprov "github.com/kararnab/authdemo/pkg/iam/provider/inhouse"
)

type Handlers struct {
	IAM       iam.Service
	UserStore internalprov.UserStore
}

func NewHandlers(
	iamSvc iam.Service,
	userStore internalprov.UserStore,
) *Handlers {
	return &Handlers{
		IAM:       iamSvc,
		UserStore: userStore,
	}
}
