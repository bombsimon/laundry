package laundry

// NotificationTypes represents the different kind of notifications
// that can be sent
type NotificationTypes struct {
	Id          int    `db:"id"          json:"-"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
}

type Notifications struct {
}
