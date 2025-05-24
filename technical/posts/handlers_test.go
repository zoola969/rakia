package posts

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	GetAllPostsFn func() ([]PostRead, error)
	GetPostByIDFn func(id int) (PostRead, error)
	CreatePostFn  func(req PostCreateUpdate) (PostRead, error)
	UpdatePostFn  func(id int, req PostCreateUpdate) (PostRead, error)
	DeletePostFn  func(id int) error
}

func (m *MockService) GetAllPosts() ([]PostRead, error) {
	return m.GetAllPostsFn()
}

func (m *MockService) GetPostByID(id int) (PostRead, error) {
	return m.GetPostByIDFn(id)
}

func (m *MockService) CreatePost(req PostCreateUpdate) (PostRead, error) {
	return m.CreatePostFn(req)
}

func (m *MockService) UpdatePost(id int, req PostCreateUpdate) (PostRead, error) {
	return m.UpdatePostFn(id, req)
}

func (m *MockService) DeletePost(id int) error {
	return m.DeletePostFn(id)
}

var testPosts = []PostRead{
	{ID: 1, Title: "Test Post 1", Content: "Content 1", Author: "Author 1"},
	{ID: 2, Title: "Test Post 2", Content: "Content 2", Author: "Author 2"},
}

func setupTestRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func TestGetAllPosts(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAllFn   func() ([]PostRead, error)
		expectedStatus int
		expectedBody   []PostRead
	}{
		{
			name: "Success",
			mockGetAllFn: func() ([]PostRead, error) {
				return testPosts, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testPosts,
		},
		{
			name: "Service Error",
			mockGetAllFn: func() ([]PostRead, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockService{
				GetAllPostsFn: tc.mockGetAllFn,
			}

			handler := NewHandler(mockService)

			req, err := setupTestRequest(http.MethodGet, "/posts", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.GetAllPosts(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var response []PostRead
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(response) != len(tc.expectedBody) {
					t.Errorf("Expected %d posts, got %d", len(tc.expectedBody), len(response))
				}

				for i, post := range response {
					if post.ID != tc.expectedBody[i].ID {
						t.Errorf("Expected post ID %d, got %d", tc.expectedBody[i].ID, post.ID)
					}
					if post.Title != tc.expectedBody[i].Title {
						t.Errorf("Expected post title %s, got %s", tc.expectedBody[i].Title, post.Title)
					}
				}
			}
		})
	}
}

func TestGetPostByID(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockGetByIDFn  func(id int) (PostRead, error)
		expectedStatus int
		expectedBody   *PostRead
	}{
		{
			name:   "Success",
			postID: "1",
			mockGetByIDFn: func(id int) (PostRead, error) {
				return testPosts[0], nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   &testPosts[0],
		},
		{
			name:   "Invalid ID",
			postID: "invalid",
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:   "Negative ID",
			postID: "-1",
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, InvalidPostIDError
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:   "Post Not Found",
			postID: "999",
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, ErrPostNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:   "Service Error",
			postID: "1",
			mockGetByIDFn: func(id int) (PostRead, error) {
				return PostRead{}, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockService{
				GetPostByIDFn: tc.mockGetByIDFn,
			}

			handler := NewHandler(mockService)

			req, err := setupTestRequest(http.MethodGet, "/posts/"+tc.postID, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.GetPostByID(rr, req, tc.postID)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			if tc.expectedStatus == http.StatusOK && tc.expectedBody != nil {
				var response PostRead
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.ID != tc.expectedBody.ID {
					t.Errorf("Expected post ID %d, got %d", tc.expectedBody.ID, response.ID)
				}
				if response.Title != tc.expectedBody.Title {
					t.Errorf("Expected post title %s, got %s", tc.expectedBody.Title, response.Title)
				}
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockCreateFn   func(req PostCreateUpdate) (PostRead, error)
		expectedStatus int
		expectedBody   *PostRead
	}{
		{
			name: "Success",
			requestBody: PostCreateUpdate{
				Title:   "New Post",
				Content: "New Content",
				Author:  "New Author",
			},
			mockCreateFn: func(req PostCreateUpdate) (PostRead, error) {
				return PostRead{
					ID:      3,
					Title:   req.Title,
					Content: req.Content,
					Author:  req.Author,
				}, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &PostRead{
				ID:      3,
				Title:   "New Post",
				Content: "New Content",
				Author:  "New Author",
			},
		},
		{
			name:        "Invalid Request Body",
			requestBody: "invalid json",
			mockCreateFn: func(req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "Validation Error",
			requestBody: PostCreateUpdate{
				Title:  "New Post",
				Author: "New Author",
			},
			mockCreateFn: func(req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, errors.New("validation error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockService{
				CreatePostFn: tc.mockCreateFn,
			}

			handler := NewHandler(mockService)

			req, err := setupTestRequest(http.MethodPost, "/posts", tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.CreatePost(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			if tc.expectedStatus == http.StatusCreated && tc.expectedBody != nil {
				var response PostRead
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.ID != tc.expectedBody.ID {
					t.Errorf("Expected post ID %d, got %d", tc.expectedBody.ID, response.ID)
				}
				if response.Title != tc.expectedBody.Title {
					t.Errorf("Expected post title %s, got %s", tc.expectedBody.Title, response.Title)
				}
				if response.Content != tc.expectedBody.Content {
					t.Errorf("Expected post content %s, got %s", tc.expectedBody.Content, response.Content)
				}
				if response.Author != tc.expectedBody.Author {
					t.Errorf("Expected post author %s, got %s", tc.expectedBody.Author, response.Author)
				}
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		requestBody    interface{}
		mockUpdateFn   func(id int, req PostCreateUpdate) (PostRead, error)
		expectedStatus int
		expectedBody   *PostRead
	}{
		{
			name:   "Success",
			postID: "1",
			requestBody: PostCreateUpdate{
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			mockUpdateFn: func(id int, req PostCreateUpdate) (PostRead, error) {
				return PostRead{
					ID:      id,
					Title:   req.Title,
					Content: req.Content,
					Author:  req.Author,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: &PostRead{
				ID:      1,
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
		},
		{
			name:   "Invalid ID",
			postID: "invalid",
			requestBody: PostCreateUpdate{
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			mockUpdateFn: func(id int, req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "Invalid Request Body",
			postID:      "1",
			requestBody: "invalid json",
			mockUpdateFn: func(id int, req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:   "Post Not Found",
			postID: "999",
			requestBody: PostCreateUpdate{
				Title:   "Updated Post",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			mockUpdateFn: func(id int, req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, ErrPostNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:   "Validation Error",
			postID: "1",
			requestBody: PostCreateUpdate{
				Title:  "Updated Post",
				Author: "Updated Author",
			},
			mockUpdateFn: func(id int, req PostCreateUpdate) (PostRead, error) {
				return PostRead{}, errors.New("validation error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockService{
				UpdatePostFn: tc.mockUpdateFn,
			}

			handler := NewHandler(mockService)

			req, err := setupTestRequest(http.MethodPut, "/posts/"+tc.postID, tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.UpdatePost(rr, req, tc.postID)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			if tc.expectedStatus == http.StatusOK && tc.expectedBody != nil {
				var response PostRead
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.ID != tc.expectedBody.ID {
					t.Errorf("Expected post ID %d, got %d", tc.expectedBody.ID, response.ID)
				}
				if response.Title != tc.expectedBody.Title {
					t.Errorf("Expected post title %s, got %s", tc.expectedBody.Title, response.Title)
				}
				if response.Content != tc.expectedBody.Content {
					t.Errorf("Expected post content %s, got %s", tc.expectedBody.Content, response.Content)
				}
				if response.Author != tc.expectedBody.Author {
					t.Errorf("Expected post author %s, got %s", tc.expectedBody.Author, response.Author)
				}
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockDeleteFn   func(id int) error
		expectedStatus int
	}{
		{
			name:   "Success",
			postID: "1",
			mockDeleteFn: func(id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Invalid ID",
			postID: "invalid",
			mockDeleteFn: func(id int) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Service Error",
			postID: "1",
			mockDeleteFn: func(id int) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockService{
				DeletePostFn: tc.mockDeleteFn,
			}

			handler := NewHandler(mockService)

			req, err := setupTestRequest(http.MethodDelete, "/posts/"+tc.postID, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.DeletePost(rr, req, tc.postID)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}
