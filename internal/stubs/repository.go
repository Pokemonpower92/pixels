package stubs

type RepositoryStub[O any] struct {
	GetFunc    func(id int) (*O, bool)
	CreateFunc func(obj *O) error
	UpdateFunc func(id int, obj *O) (*O, error)
	DeleteFunc func(id int) error
}

func (r *RepositoryStub[O]) Get(id int) (*O, bool) {
	return r.GetFunc(id)
}

func (r *RepositoryStub[O]) Create(obj *O) error {
	return r.CreateFunc(obj)
}

func (r *RepositoryStub[O]) Update(id int, obj *O) (*O, error) {
	return r.UpdateFunc(id, obj)
}

func (r *RepositoryStub[O]) Delete(id int) error {
	return r.DeleteFunc(id)
}
