package posts

import (
	"errors"
	"testing"
)

type MockRepository struct {
	GetAllFn  func() ([]PostRead, error)
	GetByIDFn func(id int) (PostRead, error)
	CreateFn  func(data PostCreateUpdate) (PostRead, error)
	UpdateFn  func(id int, data PostCreateUpdate) (PostRead, error)
	DeleteFn  func(id int) error
}

func (m *MockRepository) GetAll() ([]PostRead, error) {
	return m.GetAllFn()
}

func (m *MockRepository) GetByID(id int) (PostRead, error) {
	return m.GetByIDFn(id)
}

func (m *MockRepository) Create(data PostCreateUpdate) (PostRead, error) {
	return m.CreateFn(data)
}

func (m *MockRepository) Update(id int, data PostCreateUpdate) (PostRead, error) {
	return m.UpdateFn(id, data)
}

func (m *MockRepository) Delete(id int) error {
	return m.DeleteFn(id)
}

var testPostsData = []PostRead{
	{ID: 1, Title: "Test Post 1", Content: "Content 1", Author: "Author 1"},
	{ID: 2, Title: "Test Post 2", Content: "Content 2", Author: "Author 2"},
}

func TestServiceGetAllPosts(t *testing.T) {
	tests := []struct {
		name          string
		mockGetAllFn  func() ([]PostRead, error)
		expectedPosts []PostRead
		expectedError bool
	}{
		{
			name: "Success",
			mockGetAllFn: func() ([]PostRead, error) {
				return testPostsData, nil
			},
			expectedPosts: testPostsData,
			expectedError: false,
		},
		{
			name: "Repository Error",
			mockGetAllFn: func() ([]PostRead, error) {
				return nil, errors.New("repository error")
			},
			expectedPosts: nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetAllFn: tc.mockGetAllFn,
			}

			service := NewPostService(mockRepo)

			posts, err := service.GetAllPosts()

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(posts) != len(tc.expectedPosts) {
				t.Errorf("Expected %d posts, got %d", len(tc.expectedPosts), len(posts))
			}

			for i, post := range posts {
				if post.ID != tc.expectedPosts[i].ID {
					t.Errorf("Expected post ID %d, got %d", tc.expectedPosts[i].ID, post.ID)
				}
				if post.Title != tc.expectedPosts[i].Title {
					t.Errorf("Expected post title %s, got %s", tc.expectedPosts[i].Title, post.Title)
				}
				if post.Content != tc.expectedPosts[i].Content {
					t.Errorf("Expected post content %s, got %s", tc.expectedPosts[i].Content, post.Content)
				}
				if post.Author != tc.expectedPosts[i].Author {
					t.Errorf("Expected post author %s, got %s", tc.expectedPosts[i].Author, post.Author)
				}
			}
		})
	}
}

func TestServiceGetPostByID(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		mockGetByIDFn func(id int) (PostRead, error)
		expectedPost  *PostRead
		expectedError bool
	}{
		{
			name: "Success",
			id:   1,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return testPostsData[0], nil
			},
			expectedPost:  &testPostsData[0],
			expectedError: false,
		},
		{
			name: "Invalid ID",
			id:   0,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name: "Post Not Found",
			id:   999,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, ErrPostNotFound
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name: "Repository Error",
			id:   1,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, errors.New("repository error")
			},
			expectedPost:  nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetByIDFn: tc.mockGetByIDFn,
			}

			service := NewPostService(mockRepo)

			post, err := service.GetPostByID(tc.id)

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tc.expectedPost == nil {
				return
			}

			if post.ID != tc.expectedPost.ID {
				t.Errorf("Expected post ID %d, got %d", tc.expectedPost.ID, post.ID)
			}
			if post.Title != tc.expectedPost.Title {
				t.Errorf("Expected post title %s, got %s", tc.expectedPost.Title, post.Title)
			}
			if post.Content != tc.expectedPost.Content {
				t.Errorf("Expected post content %s, got %s", tc.expectedPost.Content, post.Content)
			}
			if post.Author != tc.expectedPost.Author {
				t.Errorf("Expected post author %s, got %s", tc.expectedPost.Author, post.Author)
			}
		})
	}
}

func TestServiceCreatePost(t *testing.T) {
	validPostData := PostCreateUpdate{
		Title:   "New Post",
		Content: "New Content",
		Author:  "New Author",
	}

	tests := []struct {
		name          string
		postData      PostCreateUpdate
		mockCreateFn  func(data PostCreateUpdate) (PostRead, error)
		expectedPost  *PostRead
		expectedError bool
	}{
		{
			name:     "Success",
			postData: validPostData,
			mockCreateFn: func(data PostCreateUpdate) (PostRead, error) {
				return PostRead{
					ID:      3,
					Title:   data.Title,
					Content: data.Content,
					Author:  data.Author,
				}, nil
			},
			expectedPost: &PostRead{
				ID:      3,
				Title:   "New Post",
				Content: "New Content",
				Author:  "New Author",
			},
			expectedError: false,
		},
		{
			name: "Validation Error",
			postData: PostCreateUpdate{
				Title:  "New Post",
				Author: "New Author",
			},
			mockCreateFn: func(data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name:     "Repository Error",
			postData: validPostData,
			mockCreateFn: func(data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, errors.New("repository error")
			},
			expectedPost:  nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				CreateFn: tc.mockCreateFn,
			}

			service := NewPostService(mockRepo)

			post, err := service.CreatePost(tc.postData)

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tc.expectedPost == nil {
				return
			}

			if post.ID != tc.expectedPost.ID {
				t.Errorf("Expected post ID %d, got %d", tc.expectedPost.ID, post.ID)
			}
			if post.Title != tc.expectedPost.Title {
				t.Errorf("Expected post title %s, got %s", tc.expectedPost.Title, post.Title)
			}
			if post.Content != tc.expectedPost.Content {
				t.Errorf("Expected post content %s, got %s", tc.expectedPost.Content, post.Content)
			}
			if post.Author != tc.expectedPost.Author {
				t.Errorf("Expected post author %s, got %s", tc.expectedPost.Author, post.Author)
			}
		})
	}
}

func TestServiceUpdatePost(t *testing.T) {
	validPostData := PostCreateUpdate{
		Title:   "Updated Post",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	tests := []struct {
		name          string
		id            int
		postData      PostCreateUpdate
		mockGetByIDFn func(id int) (PostRead, error)
		mockUpdateFn  func(id int, data PostCreateUpdate) (PostRead, error)
		expectedPost  *PostRead
		expectedError bool
	}{
		{
			name:     "Success",
			id:       1,
			postData: validPostData,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return testPostsData[0], nil
			},
			mockUpdateFn: func(id int, data PostCreateUpdate) (PostRead, error) {
				return PostRead{
					ID:      id,
					Title:   data.Title,
					Content: data.Content,
					Author:  data.Author,
				}, nil
			},
			expectedPost: &PostRead{
				ID:      1,
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			expectedError: false,
		},
		{
			name:     "Invalid ID",
			id:       0,
			postData: validPostData,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, nil
			},
			mockUpdateFn: func(id int, data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name: "Validation Error",
			id:   1,
			postData: PostCreateUpdate{
				Title:  "Updated Post",
				Author: "Updated Author",
			},
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, nil
			},
			mockUpdateFn: func(id int, data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name:     "Post Not Found",
			id:       999,
			postData: validPostData,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, ErrPostNotFound
			},
			mockUpdateFn: func(id int, data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedPost:  nil,
			expectedError: true,
		},
		{
			name:     "Repository Error",
			id:       1,
			postData: validPostData,
			mockGetByIDFn: func(id int) (PostRead, error) {
				return testPostsData[0], nil
			},
			mockUpdateFn: func(id int, data PostCreateUpdate) (PostRead, error) {
				return PostRead{}, errors.New("repository error")
			},
			expectedPost:  nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetByIDFn: tc.mockGetByIDFn,
				UpdateFn:  tc.mockUpdateFn,
			}

			service := NewPostService(mockRepo)

			post, err := service.UpdatePost(tc.id, tc.postData)

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tc.expectedPost == nil {
				return
			}

			if post.ID != tc.expectedPost.ID {
				t.Errorf("Expected post ID %d, got %d", tc.expectedPost.ID, post.ID)
			}
			if post.Title != tc.expectedPost.Title {
				t.Errorf("Expected post title %s, got %s", tc.expectedPost.Title, post.Title)
			}
			if post.Content != tc.expectedPost.Content {
				t.Errorf("Expected post content %s, got %s", tc.expectedPost.Content, post.Content)
			}
			if post.Author != tc.expectedPost.Author {
				t.Errorf("Expected post author %s, got %s", tc.expectedPost.Author, post.Author)
			}
		})
	}
}

func TestServiceDeletePost(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		mockDeleteFn  func(id int) error
		expectedError bool
	}{
		{
			name: "Success",
			id:   1,
			mockDeleteFn: func(id int) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "Invalid ID",
			id:   0,
			mockDeleteFn: func(id int) error {
				return nil
			},
			expectedError: true,
		},
		{
			name: "Repository Error",
			id:   1,
			mockDeleteFn: func(id int) error {
				return errors.New("repository error")
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				DeleteFn: tc.mockDeleteFn,
			}

			service := NewPostService(mockRepo)

			err := service.DeletePost(tc.id)

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
