package stubs

import "github.com/google/uuid"

type RepositoryStub[O, R any] struct {
	GetFunc    func(id uuid.UUID) (*O, error)
	GetAllFunc func() ([]*O, error)
	CreateFunc func(req R) (*O, error)
	UpdateFunc func(id uuid.UUID, req R) (*O, error)
	DeleteFunc func(id uuid.UUID) error
}

func (r *RepositoryStub[O, R]) Get(id uuid.UUID) (*O, error) {
	return r.GetFunc(id)
}

func (r *RepositoryStub[O, R]) GetAll() ([]*O, error) {
	return r.GetAllFunc()
}

func (r *RepositoryStub[O, R]) Create(req R) (*O, error) {
	return r.CreateFunc(req)
}

func (r *RepositoryStub[O, R]) Update(id uuid.UUID, req R) (*O, error) {
	return r.UpdateFunc(id, req)
}

func (r *RepositoryStub[O, R]) Delete(id uuid.UUID) error {
	return r.DeleteFunc(id)
}
