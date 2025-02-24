package miniblog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBInit(t *testing.T) {
	db, err := ConnectDatabase(":memory:", "")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	defaultPost, err := IsPostExist(db, defaultPostSlug)
	assert.Nil(t, err)
	assert.Equal(t, defaultPost.Title, defaultPostTitle)
}
func TestPostCRUD(t *testing.T) {
	db, err := ConnectDatabase(":memory:", "")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	p, _ := IsPostExist(db, defaultPostSlug)

	t.Run("CreatePost", func(t *testing.T) {
		p.Content = "new content"
		err = UpdatePost(db, &p)
		assert.Nil(t, err)

		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, p.Content, "new content")

		_, err = GetPostBySlug(db, "not-exist")
		assert.NotNil(t, err)
	})

	t.Run("GetPosts", func(t *testing.T) {
		ps, err := GetPosts(db, 0, 10)
		assert.Nil(t, err)
		assert.Equal(t, len(ps), 1)

		count := GetPostsCount(db)
		assert.Equal(t, count, 1)

	})

	t.Run("GetPagination", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			p := Post{
				Title: "title",
				Slug:  fmt.Sprintf("slug-%d", i),
			}
			err = CreatePost(db, &p)
			assert.Nil(t, err)
		}
		{
			pa := GetPagination(db, 0, 10)
			assert.Equal(t, pa.Total, 101)
			assert.Equal(t, pa.Page, 1)
			assert.Equal(t, pa.Limit, 10)
			assert.Equal(t, pa.Prev, 1)
			assert.Equal(t, pa.Next, 2)
			assert.Equal(t, pa.Last, 11)
		}
		{
			pa := GetPagination(db, 2, 10)
			assert.Equal(t, pa.Total, 101)
			assert.Equal(t, pa.Page, 2)
			assert.Equal(t, pa.Limit, 10)
			assert.Equal(t, pa.Prev, 1)
			assert.Equal(t, pa.Next, 3)
			assert.Equal(t, pa.Last, 11)
		}
	})

	t.Run("DeletePostByID", func(t *testing.T) {
		err = DeletePostByID(db, p.ID)
		assert.Nil(t, err)

		_, err = GetPostByID(db, p.ID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "record not found")
	})
}

func TestCommentCRUD(t *testing.T) {
	db, err := ConnectDatabase(":memory:", "")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	p, _ := IsPostExist(db, defaultPostSlug)

	t.Run("CreateComment", func(t *testing.T) {
		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, len(p.Comments), 0)
		err = CreateComment(db, p.ID, "bob@ruzhila.cn", "bob", "Hi, nice to meet you", "127.0.0.1")

		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, len(p.Comments), 1)
		assert.Equal(t, p.Comments[0].Email, "bob@ruzhila.cn")

		err = DeleteCommentByID(db, p.Comments[0].ID)
		assert.Nil(t, err)
		p, _ = GetPostByID(db, p.ID)
		assert.Equal(t, len(p.Comments), 0)
	})
}
