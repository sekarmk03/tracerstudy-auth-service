package repository

import (
	"context"
	"errors"
	"log"
	"tracerstudy-auth-service/modules/user/entity"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserRepositoryUseCase interface {
	FindAll(ctx context.Context, req any) ([]*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, id uint64) (*entity.User, error)
	Create(ctx context.Context, req *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User, updatedFields map[string]interface{}) (*entity.User, error)
	Delete(ctx context.Context, id uint64) error
}

func (u *UserRepository) FindAll(ctx context.Context, req any) ([]*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindAll")
	defer span.End()

	var users []*entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Find(&users).Error; err != nil {
		log.Println("ERROR: [UserRepository - FindAll] Internal server error:", err)
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindByUsername")
	defer span.End()

	var user entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [UserRepository - FindByUsername] Record not found for username", username)
			return nil, status.Errorf(codes.NotFound, "record not found for username %s", username)
		}
		log.Println("ERROR: [UserRepository - FindByUsername] Internal server error:", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindByEmail")
	defer span.End()

	var user entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [UserRepository - FindByEmail] Record not found for email", email)
			return nil, status.Errorf(codes.NotFound, "record not found for email %s", email)
		}
		log.Println("ERROR: [UserRepository - FindByEmail] Internal server error:", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindById(ctx context.Context, id uint64) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindById")
	defer span.End()

	var user entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [UserRepository - FindById] Record not found for id", id)
			return nil, status.Errorf(codes.NotFound, "record not found for id %d", id)
		}
		log.Println("ERROR: [UserRepository - FindById] Internal server error:", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Create(ctx context.Context, req *entity.User) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - Create")
	defer span.End()

	if err := u.db.Debug().WithContext(ctxSpan).Create(req).Error; err != nil {
		log.Println("ERROR: [UserRepository - Create] Internal server error:", err)
		return nil, err
	}

	return req, nil
}

func (u *UserRepository) Update(ctx context.Context, user *entity.User, updatedFields map[string]interface{}) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - Update")
	defer span.End()

	if err := u.db.Debug().WithContext(ctxSpan).Model(user).Updates(updatedFields).Error; err != nil {
		log.Println("ERROR: [UserRepository - Update] Internal server error:", err)
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Delete(ctx context.Context, id uint64) error {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - Delete")
	defer span.End()

	if err := u.db.Debug().WithContext(ctxSpan).Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		log.Println("ERROR: [UserRepository - Delete] Internal server error:", err)
		return err
	}

	return nil
}
