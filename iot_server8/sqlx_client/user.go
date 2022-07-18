package sqlx_client

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

//判断用户是否注册
func IsUserExisted(name string) (bool, error) {
	var res int
	rows, err := db.Queryx("select count(*) from user where username=?", name)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			return false, err
		}
	}
	if res == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

//获取用户密码
func GetUserInfo(name string) (*User, error) {
	var user User
	err := db.Get(&user, "select username,password from user where username=?", name)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}
