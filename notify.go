package laundry

// NotificationType represents the different kind of notifications
// that can be sent
type NotificationType struct {
	Id          int    `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
}
