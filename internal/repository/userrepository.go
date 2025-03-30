package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type UserRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewUserRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*UserRepository, error) {
	logger := log.New(
		log.Writer(),
		"UserRepository: ",
		log.LstdFlags,
	)
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &UserRepository{
		client: client,
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (ur *UserRepository) Close() {
	ur.client.Close()
}

func (ur *UserRepository) Get(id uuid.UUID) (*sqlc.User, error) {
	User, err := ur.q.GetUser(ur.ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (ur *UserRepository) GetByResourceId(userName string) ([]*sqlc.User, error) {
	User, err := ur.q.GetByUserName(ur.ctx, userName)
	if err != nil {
		return nil, err
	}
	return []*sqlc.User{User}, nil
}

func (ur *UserRepository) GetAll() ([]*sqlc.User, error) {
	User, err := ur.q.ListUsers(ur.ctx)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (ur *UserRepository) Create(
	req sqlc.CreateUserParams,
) (*sqlc.User, error) {
	imageset, err := ur.q.CreateUser(ur.ctx, req)
	if err != nil {
		return nil, err
	}
	return imageset, nil
}

func (ur *UserRepository) Update(
	id uuid.UUID,
	req uuid.UUID,
) (*sqlc.User, error) {
	ur.logger.Printf("Update not implemented for User")
	return nil, errors.New("Not implemented")
}

func (ur *UserRepository) Delete(id uuid.UUID) error {
	ur.logger.Printf("Delete not implemented for User")
	return errors.New("Not implemented")
}
