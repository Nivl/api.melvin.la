package article

import (
	"encoding/json"

	"github.com/Nivl/api.melvin.la/src/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestSuite() {
	ginkgo.Describe("Article", func() {
		ginkgo.It("Should return a list of article", func() {
			articles := []Article{}
			response := test.Get("/blog/articles")

			json.NewDecoder(response.Body).Decode(&articles)
			defer response.Body.Close()

			gomega.Expect(len(articles)).To(gomega.Equal(0))
		})
	})
}
