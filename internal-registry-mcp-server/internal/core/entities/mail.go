package entities

type MailAttachment struct {
	Filename    string
	Content     []byte
	ContentType string
}

type MailEmbedded struct {
	CID         string // Content-ID for embedding (e.g., "logo")
	Filename    string
	Content     []byte
	ContentType string
}

type Mail struct {
	To          string
	Subject     string
	Body        string
	IsHTML      bool
	Bcc         []string
	Cc          []string
	Attachments []MailAttachment
	Embedded    []MailEmbedded
}
