package kbot

type Bot struct {
	Port    string
	Host    string
	Token   string
	Handler func(Update, chan<- OutMessage, chan<- OutQuery)
}

type UserResponse struct {
	Ok          bool
	Error_code  int
	Description string
	Result      User
}

type UpdateResponse struct {
	Ok          bool
	Error_code  int
	Description string
	Result      []Update
}

type OutQuery struct {
	Inline_query_id string       `json:"inline_query_id"`
	Results         QueryResults `json:"results"`
}

type QueryResults []QueryResultArticle

// TODO: QueryResults: Audio,Contact,Document,Gif,Location,Photo,Venue,Video,Voice
type QueryResultArticle struct {
	Type                  string                  `json:"type"`
	Id                    string                  `json:"id"`
	Title                 string                  `json:"title"`
	Input_message_content InputTextMessageContent `json:"input_message_content"`
	Url                   string                  `json:"url"`
}

// TODO: Input MessageContent: Location,Venue,Contact
type InputTextMessageContent struct {
	Message_text string `json:"message_text"`
	Parse_mode   string `json:"parse_mode"`
}

type OutMessage struct {
	Chat_id    int    `json:"chat_id"`
	Text       string `json:"text"`
	Parse_mode string `json:"parse_mode"`
}

type Update struct {
	Update_id    int
	Message      Message
	Inline_query InlineQuery
	//chosen_inline_result ChosenInlineResult // TODO: inline result
	//callback_query       CallbackQuery // TODO: callback query
}

type User struct {
	Id         int
	First_name string
	Last_name  string
	Username   string
}

type Chat struct {
	Id         int
	Type       string
	Title      string
	First_name string
	Last_name  string
	Username   string
}

type Message struct {
	Message_id       int
	From             User
	Date             int
	Chat             Chat
	Text             string
	Entities         MessageEntities
	New_chat_member  User
	Left_chat_member User
	Location         Location
	// TODO: audio,document,photo,venue...
}

type MessageEntity struct {
	EntityType string `json:"type"`
	Offset     int
	Length     int
	Url        string
}

type MessageEntities []MessageEntity

type InlineQuery struct {
	Id       string
	From     User
	Query    string
	Offset   string
	Location Location
}

type Location struct {
	Longitude float64
	Latitude  float64
}
