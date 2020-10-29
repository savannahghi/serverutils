package base

// MailgunEMailMessage used to create a message used in requests
type MailgunEMailMessage struct {
	Subject string   `json:"subject,omitempty"`
	Text    string   `json:"text,omitempty"`
	To      []string `json:"to,omitempty"`
}
