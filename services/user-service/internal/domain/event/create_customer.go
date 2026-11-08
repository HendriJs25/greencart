package event

type CreateCustomerRequested struct {
	EventID          string `json:"event_id"`
	ReceiverID       int64  `json:"receiver_id"`
	ReceiverEmail    string `json:"receiver_email"`
	Subject          string `json:"subject"`
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
}
