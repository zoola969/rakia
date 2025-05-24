package posts

import (
	"encoding/json"
	"errors"
	"maps"
	"os"
	"slices"
	"sync"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Repository interface {
	GetAll() ([]PostRead, error)
	GetByID(id int) (PostRead, error)
	Create(data PostCreateUpdate) (PostRead, error)
	Update(id int, data PostCreateUpdate) (PostRead, error)
	Delete(id int) error
}

type MapRepository struct {
	posts  map[int]PostRead
	nextID int
	mutex  sync.RWMutex
}

func NewMapRepository() *MapRepository {
	data, err := os.ReadFile("blog_data.json")
	if err != nil {
		panic(err)
	}

	var jsonData struct {
		Posts []PostRead `json:"posts"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		panic(err)
	}
	posts := jsonData.Posts

	repo := &MapRepository{
		posts:  make(map[int]PostRead),
		mutex:  sync.RWMutex{},
		nextID: 1,
	}

	maxID := 0
	for _, post := range posts {
		repo.posts[post.ID] = post
		if post.ID > maxID {
			maxID = post.ID
		}
	}
	repo.nextID = maxID + 1
	return repo
}

func (r *MapRepository) GetAll() ([]PostRead, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return slices.Collect(maps.Values(r.posts)), nil
}

func (r *MapRepository) GetByID(id int) (PostRead, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	val, ok := r.posts[id]
	if ok {
		return val, nil
	}
	return PostRead{}, ErrPostNotFound
}

func (r *MapRepository) Create(data PostCreateUpdate) (PostRead, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	createdPost := PostRead{
		ID:      r.nextID,
		Title:   data.Title,
		Content: data.Content,
		Author:  data.Author,
	}
	r.posts[r.nextID] = createdPost
	r.nextID += 1
	return createdPost, nil
}

func (r *MapRepository) Update(id int, data PostCreateUpdate) (PostRead, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.posts[id]
	if !ok {
		return PostRead{}, ErrPostNotFound
	}
	updatedPost := PostRead{
		ID:      id,
		Title:   data.Title,
		Content: data.Content,
		Author:  data.Author,
	}
	r.posts[id] = updatedPost
	return updatedPost, nil
}

func (r *MapRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.posts, id)
	return nil
}
