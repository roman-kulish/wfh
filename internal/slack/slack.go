package slack

type CommandRequest struct {
	Token          string
	TeamID         string
	TeamDomain     string
	EnterpriseID   string
	EnterpriseName string
	ChannelID      string
	ChannelName    string
	UserID         string
	UserName       string
	Command        string
	Text           string
	ResponseURL    string
	TriggerID      string
}

type CommandResponse struct {
	ResponseType string       `json:"response_type,omitempty"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

func (cr *CommandResponse) AddAttachment(attachment Attachment) {
	cr.Attachments = append(cr.Attachments, attachment)
}

type Attachment struct {
	Fallback   string  `json:"fallback,omitempty"`
	Color      string  `json:"color,omitempty"`
	Pretext    string  `json:"pretext,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Title      string  `json:"title,omitempty"`
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short string `json:"short"`
}

func NewInChannelCommandResponse(text string) CommandResponse {
	return CommandResponse{
		ResponseType: "in_channel",
		Text:         text,
	}
}
