package user_mgr

type UserModel struct {
	Id               int64             `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement"`
	Email            string            `json:"email" gorm:"type:varchar(128);uniqueIndex;"`
	Password         string            `json:"password" gorm:"type:varchar(128)"`
	Token            string            `json:"token" gorm:"type:varchar(64);uniqueIndex;"`
	Forbidden        bool              `json:"forbidden" gorm:"type:tinyint(1);"`
	Roles            string            `json:"roles" gorm:"type:longtext;"`
	Roles_map        map[string]string `json:"roles_map" gorm:"-"`
	Permissions      string            `json:"permissions" gorm:"type:longtext;"`
	Permissions_map  map[string]string `json:"permissions_map" gorm:"-"`
	Register_ipv4    string            `json:"register_ipv4" gorm:"type:varchar(64);index;"`
	Update_unixtime  int64             `json:"update_unixtime" gorm:"autoUpdateTime;type:bigint(20)"`
	Created_unixtime int64             `json:"created_unixtime" gorm:"autoCreateTime;type:bigint(20)"`
}

const TABLE_NAME_USER = "user"

func (model *UserModel) TableName() string {
	return TABLE_NAME_USER
}
