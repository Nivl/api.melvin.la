package blog

import "github.com/Nivl/api.melvin.la/api/components/blog/articles"

// EnsureIndexes sets the indexes for all the documents in the blog
func EnsureIndexes() {
	articles.EnsureIndexes()
}
