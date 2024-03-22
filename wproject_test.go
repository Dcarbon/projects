package main

import (
	"context"
	"testing"
	"time"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
	"github.com/Dcarbon/go-shared/libs/utils"
)

var reqCtx = context.TODO()
var host = utils.StringEnv(gutils.ISVProjects, "localhost:4012")

var projectClient pb.ProjectServiceClient

func init() {

	cc, err := gutils.GetCCTimeout(host, 5*time.Second)
	utils.PanicError("Connect to service timeout", err)

	projectClient = pb.NewProjectServiceClient(cc)
}

func TestGetById(t *testing.T) {
	data, err := projectClient.GetById(reqCtx, &pb.RPGetById{
		ProjectId: 1,
		Lang:      "vi",
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestProjectGetList(t *testing.T) {
	data, err := projectClient.GetList(reqCtx, &pb.RPGetList{
		Skip:  0,
		Limit: 3,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}
