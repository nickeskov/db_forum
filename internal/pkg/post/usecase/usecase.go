package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/post/repository"
	"github.com/nickeskov/db_forum/internal/pkg/thread"
	"github.com/pkg/errors"
	"strconv"
)

type UseCase struct {
	repository repository.Repository
	threadRepo thread.Repository
}

func NewUseCase(repository repository.Repository, threadRepo thread.Repository) UseCase {
	return UseCase{
		repository: repository,
		threadRepo: threadRepo,
	}
}

func (useCase UseCase) CreatePostsByThreadSlugOrID(threadSlugOrID string,
	posts models.Posts) (models.Posts, error) {

	if len(posts) == 0 {
		return make(models.Posts, 0), nil
	}

	var err error
	var postsThread models.Thread

	if id, convertErr := strconv.Atoi(threadSlugOrID); convertErr != nil {
		postsThread, err = useCase.threadRepo.GetBySlug(threadSlugOrID)
	} else {
		postsThread, err = useCase.threadRepo.GetByID(int32(id))
	}

	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		return nil, models.ErrDoesNotExist
	case err != nil:
		return nil, errors.WithStack(err)
	}

	return useCase.repository.CreatePostsInThread(postsThread, posts)
}
