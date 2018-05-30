package slack

type CommandRequest struct {
	Token          string
	TeamId         string
	TeamDomain     string
	EnterpriseId   string
	EnterpriseName string
	ChannelId      string
	ChannelName    string
	UserId         string
	UserName       string
	Command        string
	Text           string
	ResponseUrl    string
	TriggerId      string
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
	ImageUrl   string  `json:"image_url,omitempty"`
	ThumbUrl   string  `json:"thumb_url,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Ts         int     `json:"ts,omitempty"`
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
