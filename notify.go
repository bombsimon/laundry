package laundry

type NotificationTypes struct {
	Id          int    `db:"id"          json:"-"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
}

type Notifications struct {
}
