package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstName, lastName, email, phone *string) error
	Count(filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *domain.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("user created with id: ", user.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	var u []domain.User

	tx := repo.db.Model(&u)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&u)

	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (repo *repo) Get(id string) (*domain.User, error) {
	user := domain.User{ID: id}

	result := repo.db.First(&user)
	// result := repo.db.First(&user, "id = ?", id) ----> FUNCIONA DE LA MISMA MANERA pero &user es = User{}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Delete físico y completo
func (repo *repo) Delete(id string) error {
	user := domain.User{
		ID: id,
	}
	result := repo.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repo) Update(id string, firstName, lastName, email, phone *string) error {
	values := make(map[string]interface{})

	if firstName != nil {
		values["firstName"] = firstName
	}
	if lastName != nil {
		values["lastName"] = lastName
	}
	if email != nil {
		values["email"] = email
	}
	if phone != nil {
		values["phone"] = phone
	}

	res := repo.db.Model(&domain.User{}).Where("id = ?", id).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.User{})
	tx = applyFilters(tx, filters)

	result := tx.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}
	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}
