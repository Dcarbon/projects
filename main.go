package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/projects/internal/service"
	"google.golang.org/grpc"
)

var config = gutils.Config{
	Port:   utils.IntEnv("PORT", 4012),
	DbUrl:  utils.StringEnv("DB_URL", "postgres://admin:hellosecret@dev01.dcarbon.org/projects"), //,localhost
	Name:   "ProjectService",
	JwtKey: utils.StringEnv("JWT", ""),
	Options: map[string]string{
		"AMQP_URL":        utils.StringEnv("AMQP_URL", "amqp://rbuser:hellosecret@localhost"),
		"REDIS_URL":       utils.StringEnv("REDIS_URL", "redis://localhost:6379"),
		gutils.ISVStorage: utils.StringEnv(gutils.ISVStorage, "http://localhost:4100"),
	},
	AuthConfig: map[string]*gutils.ARConfig{
		"/pb.ProjectService/Create": {
			Require:    true,
			Permission: "project-info-create",
			PermDesc:   "Create project",
		},
		"/pb.ProjectService/UpdateDesc": {
			Require:    true,
			Permission: "project-info-update-desc",
			PermDesc:   "Update project description",
		},
		"/pb.ProjectService/UpdateSpecs": {
			Require:    true,
			Permission: "project-info-update-specs",
			PermDesc:   "Update project specification",
		},
		"/pb.ProjectService/AddImage": {
			Require:    true,
			Permission: "project-info-add-image",
			PermDesc:   "Add image to project",
		},
		"/pb.ProjectService/GetById": {
			Require:    false,
			Permission: "project-info-get-by-id",
			PermDesc:   "",
		},
		"/pb.ProjectService/GetList": {
			Require:    false,
			Permission: "project-info-get-list",
			PermDesc:   "",
		},
	},
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	utils.PanicError(config.Name+" open port", err)

	logger := gutils.NewLogInterceptor()
	auth, err := gutils.NewAuthInterceptor(
		// config.GetIAM(),
		config.JwtKey,
		config.AuthConfig,
	)
	utils.PanicError("Authen init", err)

	handler, err := service.NewProjectService(&config)
	utils.PanicError(config.Name+" init", err)

	var sv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			gutils.UnaryPreventPanic,
			logger.Intercept,
			auth.Intercept,
		),
	)
	pb.RegisterProjectServiceServer(sv, handler)
	log.Println(config.Name+" listen and serve at ", config.Port)
	err = sv.Serve(listen)
	if nil != err {
		log.Fatal(config.Name+" listen and serve error: ", err)
	}
}
