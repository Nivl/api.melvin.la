package main

import (
	"github.com/Nivl/api.melvin.la/src/blog"
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("Api", func() {
	blog.TestSuite()
})
