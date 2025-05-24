package posts

import (
	"sync"
	"testing"
)

func TestMapRepositoryGetAll(t *testing.T) {
	repo := setupTestRepository()

	posts, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	expectedIDs := []int{1, 2}
	for _, post := range posts {
		found := false
		for _, id := range expectedIDs {
			if post.ID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected post ID: %d", post.ID)
		}
	}
}

func TestMapRepositoryGetByID(t *testing.T) {
	repo := setupTestRepository()

	tests := []struct {
		name          string
		id            int
		expectedError bool
		expectedPost  *PostRead
	}{
		{
			name:          "Existing Post",
			id:            1,
			expectedError: false,
			expectedPost: &PostRead{
				ID:      1,
				Title:   "Test Post 1",
				Content: "Test Content 1",
				Author:  "Test Author 1",
			},
		},
		{
			name:          "Non-existent Post",
			id:            999,
			expectedError: true,
			expectedPost:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			post, err := repo.GetByID(tc.id)

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

func TestMapRepositoryCreate(t *testing.T) {
	repo := setupTestRepository()

	newPost := PostCreateUpdate{
		Title:   "New Post",
		Content: "New Content",
		Author:  "New Author",
	}

	createdPost, err := repo.Create(newPost)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if createdPost.ID <= 0 {
		t.Errorf("Expected a positive ID, got %d", createdPost.ID)
	}

	if createdPost.Title != newPost.Title {
		t.Errorf("Expected title %s, got %s", newPost.Title, createdPost.Title)
	}

	if createdPost.Content != newPost.Content {
		t.Errorf("Expected content %s, got %s", newPost.Content, createdPost.Content)
	}

	if createdPost.Author != newPost.Author {
		t.Errorf("Expected author %s, got %s", newPost.Author, createdPost.Author)
	}

	retrievedPost, err := repo.GetByID(createdPost.ID)
	if err != nil {
		t.Errorf("Expected no error when retrieving created post, got %v", err)
	}

	if retrievedPost.ID != createdPost.ID {
		t.Errorf("Expected retrieved post ID %d, got %d", createdPost.ID, retrievedPost.ID)
	}
}

func TestMapRepositoryUpdate(t *testing.T) {
	repo := setupTestRepository()

	updatedData := PostCreateUpdate{
		Title:   "Updated Post",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	tests := []struct {
		name          string
		id            int
		data          PostCreateUpdate
		expectedError bool
		expectedPost  *PostRead
	}{
		{
			name:          "Existing Post",
			id:            1,
			data:          updatedData,
			expectedError: false,
			expectedPost: &PostRead{
				ID:      1,
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
		},
		{
			name:          "Non-existent Post",
			id:            999,
			data:          updatedData,
			expectedError: true,
			expectedPost:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			post, err := repo.Update(tc.id, tc.data)

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

			retrievedPost, err := repo.GetByID(tc.id)
			if err != nil {
				t.Errorf("Expected no error when retrieving updated post, got %v", err)
			}

			if retrievedPost.Title != tc.expectedPost.Title {
				t.Errorf("Expected retrieved post title %s, got %s", tc.expectedPost.Title, retrievedPost.Title)
			}
		})
	}
}

func TestMapRepositoryDelete(t *testing.T) {
	repo := setupTestRepository()

	tests := []struct {
		name          string
		id            int
		expectedError bool
	}{
		{
			name:          "Existing Post",
			id:            1,
			expectedError: false,
		},
		{
			name:          "Non-existent Post",
			id:            999,
			expectedError: false, // Delete is idempotent, so no error is expected
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(tc.id)

			if tc.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tc.expectedError {
				_, err := repo.GetByID(tc.id)
				if err == nil {
					t.Errorf("Expected post with ID %d to be deleted", tc.id)
				}
				if err != ErrPostNotFound {
					t.Errorf("Expected ErrPostNotFound, got %v", err)
				}
			}
		})
	}
}

func setupTestRepository() *MapRepository {
	repo := &MapRepository{
		posts:  make(map[int]PostRead),
		mutex:  sync.RWMutex{},
		nextID: 3,
	}

	repo.posts[1] = PostRead{
		ID:      1,
		Title:   "Test Post 1",
		Content: "Test Content 1",
		Author:  "Test Author 1",
	}

	repo.posts[2] = PostRead{
		ID:      2,
		Title:   "Test Post 2",
		Content: "Test Content 2",
		Author:  "Test Author 2",
	}

	return repo
}
