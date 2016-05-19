package main

import (
	"time"

	"github.com/Nivl/api.melvin.la/src/app"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"testing"
)

func TestSrc(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "api-melvin-la Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	argv := &app.Args{
		Port:     "5001",
		MongoURI: "mongodb://localhost/api-melvin-test",
	}

	app.InitContextWithParams(argv)

	ensureIndexes()
	go start()
	// TODO find a better way to enure the server has started
	time.Sleep(3 * time.Second)
})

var _ = ginkgo.AfterSuite(func() {
	app.GetContext().Destroy()
})
