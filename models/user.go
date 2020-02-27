package models

import (
	"github.com/jinzhu/gorm"
	"github.com/wenxian2012/go-rbac-admin/dto"
)

type AuthSwag struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Model
	Username *string `gorm:"column:username"`
	Password *string `gorm:"column:password"`
	Nickname *string `gorm:"column:nickname"`
	Disabled *int    `gorm:"column:disabled"`
	Phone    *string `gorm:"column:phone"`
	Email    *string `gorm:"column:email"`
	Roles    []*Role `gorm:"many2many:user_role;"`
}

func (a User) TableName() string {
	return a.Model.TableName("user")
}

func (a *User) ToDtoUser() *dto.User {
	user := dto.User{
		ID:        a.ID,
		Username:  *a.Username,
		Password:  *a.Password,
		Nickname:  *a.Nickname,
		Roles:     ToDtoRoles(a.Roles),
		UpdatedAt: a.UpdatedAt,
	}
	if a.Phone != nil {
		user.Phone = *a.Phone
	}
	if a.Email != nil {
		user.Email = *a.Email
	}
	if a.Disabled != nil {
		user.Disabled = *a.Disabled
	}
	return &user
}

type UserRole struct {
	Model
	UserID int `gorm:"column:user_id"`
	RoleID int `gorm:"column:role_id"`
}

func (a UserRole) TableName() string {
	return a.Model.TableName("user_role")
}

func CheckUser(username, password string) (*User, error) {
	var user User
	err := db.Model(&User{}).Where(User{Username: &username, Password: &password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user.ID > 0 {
		return &user, nil
	}
	return nil, nil
}

func CheckUsername(username string) (bool, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetUsersAll() ([]*User, error) {
	var user []*User
	err := db.Preload("Roles").Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return user, nil
}

func (a *User) Query(params *dto.UserQueryParams) ([]*User, int, error) {
	var list = []*User{}
	var count int
	db := db.Model(&User{})
	if v := params.Username; v != "" {
		db = db.Where("username LIKE ?", "%"+v+"%")
	}
	if v := params.Nickname; v != "" {
		db = db.Where("nickname LIKE ?", "%"+v+"%")
	}
	if v := params.Phone; v != "" {
		db = db.Where("phone LIKE ?", "%"+v+"%")
	}
	if v := params.Email; v != "" {
		db = db.Where("email LIKE ?", "%"+v+"%")
	}
	if v := params.Disabled; v != nil {
		db = db.Where("disabled = ?", v)
	}
	if v := params.RoleId; v != 0 {
		subQuery := db.Model(&UserRole{}).Select("user_id").Where("role_id = ?", v).SubQuery()
		db = db.Where("id IN (?)", subQuery)
	}
	db.Count(&count)
	err := db.Preload("Roles").Offset(params.PageNum).Limit(params.PageSize).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, count, err
	}
	return list, count, nil
}

func (a *User) Find(id int) (*User, error) {
	var user User
	err := db.Preload("Roles").Where("id = ? ", id).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func (a *User) Update(roleIds []int) error {
	var role []Role
	db.Where("id in (?)", roleIds).Find(&role)
	if err := db.Model(&User{}).Updates(a).Association("Roles").Replace(role).Error; err != nil {
		return err
	}
	return nil
}

func (a *User) Create(roleIds []int) error {
	var role []Role
	db.Where("id in (?)", roleIds).Find(&role)
	if err := db.Create(a).Error; err != nil {
		return err
	}
	if err := db.Model(a).Association("Roles").Append(role).Error; err != nil {
		return err
	}
	return nil
}

func (a *User) Delete(id int) error {
	if err := db.Delete(User{Model: Model{ID: id}}).Association("Roles").Clear().Error; err != nil {
		return err
	}

	return nil
}
