package stubs

import "github.com/google/uuid"

type RepositoryStub[O any] struct {
	GetFunc    func(id uuid.UUID) (*O, bool)
	GetAllFunc func() ([]*O, bool)
	CreateFunc func(obj *O) error
	UpdateFunc func(id uuid.UUID, obj *O) (*O, error)
	DeleteFunc func(id uuid.UUID) error
}

func (r *RepositoryStub[O]) Get(id uuid.UUID) (*O, bool) {
	return r.GetFunc(id)
}

func (r *RepositoryStub[O]) GetAll() ([]*O, bool) {
	return r.GetAllFunc()
}

func (r *RepositoryStub[O]) Create(obj *O) error {
	return r.CreateFunc(obj)
}

func (r *RepositoryStub[O]) Update(id uuid.UUID, obj *O) (*O, error) {
	return r.UpdateFunc(id, obj)
}

func (r *RepositoryStub[O]) Delete(id uuid.UUID) error {
	return r.DeleteFunc(id)
}
