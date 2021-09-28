// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/c0llinn/ebook-store/config/aws"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/internal/api"
	http2 "github.com/c0llinn/ebook-store/internal/auth/delivery/http"
	"github.com/c0llinn/ebook-store/internal/auth/email"
	"github.com/c0llinn/ebook-store/internal/auth/helper"
	"github.com/c0llinn/ebook-store/internal/auth/middleware"
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/c0llinn/ebook-store/internal/auth/usecase"
	http3 "github.com/c0llinn/ebook-store/internal/catalog/delivery/http"
	helper2 "github.com/c0llinn/ebook-store/internal/catalog/helper"
	repository2 "github.com/c0llinn/ebook-store/internal/catalog/repository"
	"github.com/c0llinn/ebook-store/internal/catalog/storage"
	usecase2 "github.com/c0llinn/ebook-store/internal/catalog/usecase"
	"net/http"
)

// Injectors from wire.go:

func CreateWebServer() *http.Server {
	engine := api.NewRouter()
	gormDB := db.NewConnection()
	userRepository := repository.NewUserRepository(gormDB)
	hmacSecret := helper.NewHMACSecret()
	jwtWrapper := helper.NewJWTWrapper(hmacSecret)
	ses := aws.NewSNSService()
	client := email.NewEmailClient(ses)
	passwordGenerator := helper.NewPasswordGenerator()
	bcryptWrapper := helper.NewBcryptWrapper()
	authUseCase := usecase.NewAuthUseCase(userRepository, jwtWrapper, client, passwordGenerator, bcryptWrapper)
	uuidGenerator := helper.NewUUIDGenerator()
	authHandler := http2.NewAuthHandler(authUseCase, uuidGenerator)
	bookRepository := repository2.NewBookRepository(gormDB)
	s3 := aws.NewS3Service()
	bucket := aws.NewBucket()
	s3Client := storage.NewS3Client(s3, bucket)
	filenameGenerator := helper2.NewFilenameGenerator()
	catalogUseCase := usecase2.NewCatalogUseCase(bookRepository, s3Client, filenameGenerator)
	catalogHandler := http3.NewCatalogHandler(catalogUseCase)
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(jwtWrapper)
	server := api.NewHttpServer(engine, authHandler, catalogHandler, authenticationMiddleware)
	return server
}
