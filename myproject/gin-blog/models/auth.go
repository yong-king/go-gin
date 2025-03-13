package models

type Auth struct {
	ID       uint   `gorm:"primary_key"`
	UserName string `valid:"Required;MaxSize(50)" json:"user_name"`
	Password string `valid:"Required;MaxSize(50)" json:"password"`
}

func CheckAuth(userName, password string) (bool, error) {
	var auth Auth
	// 获取用户名和密码
	err := db.Select("id").Where("username = ? and password = ?", userName, password).First(&auth).Error
	if err != nil {
		return false, err
	}
	return auth.ID > 0, nil
}
