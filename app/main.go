package main

import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/enj/osa/osa"
)

func init() {
	if _, err := osa.RegisterService(); err != nil {
		panic(err.Error())
	}
	endpoints.HandleHTTP()
}
