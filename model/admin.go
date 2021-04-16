package model

type Admin struct {
	Model
	Name string `json:"name"`
}

func GetFirstAdmin() (admin Admin, err error) {

	var a Admin

	res := Db.Debug().Where("id = ?", 1).First(&a)
	err = res.Error
	admin = a
	return
}
