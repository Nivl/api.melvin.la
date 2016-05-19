package blog

import (
	"github.com/Nivl/api.melvin.la/src/blog/article"
	. "github.com/onsi/ginkgo"
)

func TestSuite() {
	Describe("Blog", func() {
		article.TestSuite()
	})
}
