package mocks

import (
	"context"

	"github.com/glu/shopvui/internal/userm/entities"
	"github.com/glu/shopvui/internal/userm/golibs/database"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func NewMockUserRepo() *mockUserRepo {
	return &mockUserRepo{}
}

func (m *mockUserRepo) GetUser(ctx context.Context, db database.Ext, email pgtype.Text) (*entities.User, error) {
	args := m.Called(ctx, db, email)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, db database.Ext, e *entities.User) (*entities.User, error) {
	args := m.Called(ctx, db, e)
	return args.Get(0).(*entities.User), args.Error(1)
}
