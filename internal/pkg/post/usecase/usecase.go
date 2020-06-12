package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/post"
	"github.com/nickeskov/db_forum/internal/pkg/thread"
	"github.com/nickeskov/db_forum/internal/pkg/user"
	"github.com/pkg/errors"
	"strconv"
)

type UseCase struct {
	repository post.Repository
	userRepo   user.Repository
	forumRepo  forum.Repository
	threadRepo thread.Repository
}

func NewUseCase(repository post.Repository,
	userRepo user.Repository, forumRepo forum.Repository, threadRepo thread.Repository) UseCase {
	return UseCase{
		repository: repository,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
		threadRepo: threadRepo,
	}
}

func (useCase UseCase) CreatePostsByThreadSlugOrID(threadSlugOrID string,
	posts models.Posts) (models.Posts, error) {

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

func (useCase UseCase) GetPostInfoByID(id int64,
	related []string) (postFullInfo models.PostFullInfo, err error) {

	postModel, err := useCase.repository.GetPostByID(id)
	switch {
	case errors.Is(err, models.ErrDoesNotExist):
		return models.PostFullInfo{}, models.ErrDoesNotExist
	case err != nil:
		return models.PostFullInfo{}, errors.WithStack(err)
	}

	postFullInfo.Post = &postModel

	for _, entityName := range related {
		switch entityName {
		case "user":
			relatedAuthor, relatedErr := useCase.userRepo.GetByNickname(postModel.Author)
			if relatedErr != nil {
				err = relatedErr
			} else {
				postFullInfo.Author = &relatedAuthor
			}
		case "forum":
			relatedForum, relatedErr := useCase.forumRepo.GetBySlug(postModel.Forum)
			if relatedErr != nil {
				err = relatedErr
			} else {
				postFullInfo.Forum = &relatedForum
			}
		case "thread":
			relatedThread, relatedErr := useCase.threadRepo.GetByID(postModel.Thread)
			if relatedErr != nil {
				err = relatedErr
			} else {
				postFullInfo.Thread = &relatedThread
			}
		default:
			return models.PostFullInfo{}, models.ErrInvalid
		}

		switch {
		case errors.Is(err, models.ErrDoesNotExist):
			return models.PostFullInfo{}, models.ErrDoesNotExist
		case err != nil:
			return models.PostFullInfo{}, errors.WithStack(err)
		}
	}

	return postFullInfo, nil
}

func (useCase UseCase) UpdatePostByID(post models.Post) (models.Post, error) {
	return useCase.repository.UpdatePostByID(post)
}
