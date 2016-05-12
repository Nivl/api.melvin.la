package blog

import "github.com/Nivl/api.melvin.la/blog/article"

// EnsureIndexes sets the indexes for all the documents in the blog
func EnsureIndexes() {
	article.EnsureIndexes()
}
