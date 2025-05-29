package services

import (
	"github.com/eryalito/http-file-share/internal/listener"
)

type HttpFileServerService struct {
	HttpFileServer *listener.HttpFileServer
}

func (ds *HttpFileServerService) GetAddresses() []string {
	if ds.HttpFileServer == nil {
		return []string{}
	}
	return ds.HttpFileServer.Addresses()
}

func (ds *HttpFileServerService) SetFileToServe(filePath string) {
	if ds.HttpFileServer == nil {
		return
	}
	ds.HttpFileServer.SetFileToServe(filePath)
}

func (ds *HttpFileServerService) GetFileToServe() string {
	if ds.HttpFileServer == nil {
		return ""
	}
	return ds.HttpFileServer.FileToServe
}
