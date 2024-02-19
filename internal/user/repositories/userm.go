package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/glu/shopvui/internal/userm/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgxutil"
)

type UserRepository struct {
}

func (r *UserRepository) GetByEmail(ctx context.Context, db models.DBTX, email string) (*models.User, error) {
	q := models.New(db)
	result, err := q.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) Create(ctx context.Context, db models.DBTX, user *models.User) error {
	q := models.New(db)

	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	var params models.CreateUserParams
	if err := json.Unmarshal(b, &params); err != nil {
		return err
	}
	_, err = q.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUserV2(ctx context.Context, db models.DBTX, user *models.User) error {
	params := models.CreateUserParams{}
	m := mapCommonFields2(*user, params)
	result := make([]map[string]interface{}, 0)
	result = append(result, m)
	_, err := pgxutil.Insert(ctx, db, "users", result)
	if err != nil {
		return err
	}
	return nil
}

func mapCommonFields2(model1 interface{}, model2 interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val1 := reflect.ValueOf(model1)
	val2 := reflect.ValueOf(model2)

	typ1 := reflect.TypeOf(model1)
	typ2 := reflect.TypeOf(model2)

	for i := 0; i < val1.NumField(); i++ {
		field1 := typ1.Field(i)
		tag1 := field1.Tag.Get("json")

		for j := 0; j < val2.NumField(); j++ {
			field2 := typ2.Field(j)
			tag2 := field2.Tag.Get("json")

			if tag1 == tag2 {
				result[tag1] = val1.Field(i).Interface()
			}
		}
	}
	return result
}

func (r *UserRepository) GetUserByID(ctx context.Context, db models.DBTX, userID string) (*models.User, error) {
	query := fmt.Sprintf(`select * from %s where user_id = $1`, "users")

	user, err := pgxutil.SelectRow(ctx, db, query, []any{userID}, pgx.RowToAddrOfStructByPos[models.User])
	if err != nil {
		return nil, err
	}
	return user, nil
}
