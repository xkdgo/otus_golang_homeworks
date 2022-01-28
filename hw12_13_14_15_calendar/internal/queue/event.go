package queue

type NotifyEvent struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	DateTimeStart string `json:"datetimestart"`
	UserID        string `json:"userid"`
}
