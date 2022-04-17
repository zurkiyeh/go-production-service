package comment

import (
	"context"
	"errors"
	"fmt"
)

// To prevent leaving us exposed to disclusore attacks, it's better to define our own errors. This way the call stack won't be returned back to the calling func and eventually printed on user screen. So we define our own errors:
var (
	ErrFetchingComment = errors.New("error fetching comment")
	ErrNotImplemented  = errors.New("error not implemented")
)

type Comment struct {
	ID     string
	Slug   string
	Body   string
	Author string
}

// The reason we chose this design, is because we want to keep service decoupled from the implementation of store. For a store to use the service, all it has to do is to implement the functions defined in the interface store. Store will then be able to use the various methods defined in the interface. We do not care about the implementation of store, and how it reaches out to the backend/db/API to get the comment, rather we will be decoupled from it and call the functions Store defines.
// For example the GetComment -> the store has the freedom to return whatever it wants and whatever error it wants as well!

// What is also neat about this approach, is that we can also mock Store and this allows us to easily test Service independently from Store without actually implementing it.

// Service is just forwards the call
type Service struct {
	// This is necessary so that when we create a service and pass it to the repository layer service, it will match the interface "Store"
	// An alternative but more messy approach to design would be to have Store *db.Repository for example but then it's very difficult to test this because you will have to instatiate a db when testing
	Store Store
}

// This is essentially a constructor for our service struct
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("Getting comment")
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		// We print out the real error (logging) but we return a generic error that we predefined
		fmt.Println(err)
		return Comment{}, ErrFetchingComment
	}
	return cmt, nil
}

func (s *Service) UpdateComment(ctx context.Context, id string, updatedComment Comment) (Comment, error) {
	cmt, err := s.Store.UpdateComment(ctx, id, updatedComment)
	if err != nil {
		fmt.Println("error updating comment")
		return Comment{}, err
	}
	return cmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.Store.DeleteComment(ctx, id)
}

func (s *Service) PostComment(ctx context.Context, cmt Comment) (Comment, error) {
	insertedcmt, err := s.Store.PostComment(ctx, cmt)
	if err != nil {
		return Comment{}, err
	}

	return insertedcmt, nil
}

// Since we want to keep things decoupled, we dont want this service to be communicating with db directly for example. So we will create an interface that any service that wants to retrieve comments must implement in order for things to stay simple
// So if you want a repository layer service to be able to retrieve use the GetComment func, you must create an interface for that:

type Store interface {
	GetComment(context.Context, string) (Comment, error)
	PostComment(context.Context, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
	UpdateComment(context.Context, string, Comment) (Comment, error)
}
