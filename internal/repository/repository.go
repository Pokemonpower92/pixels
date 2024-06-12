package repository

type Repository[O any] interface {
    Get(id int) (*O, error)
    Create(obj *O) error
    Update(id int, obj *O) (*O, error)
    Delete(id int) (*O, error)
}

