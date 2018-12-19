package pkg

// NotificationType represents the different kind of notifications
// that can be sent
type NotificationType struct {
	ID          int    `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
}
