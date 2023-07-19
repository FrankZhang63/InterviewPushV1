package models

type ParamChat struct {
	//RobotCode   string   `json:"robotCode"`
	UserIds  []string `json:"userIds"`
	Msgkey   string   `json:"msgKey"`
	MsgParam string   `json:"msgParam"`
}
type DingParamReveiver struct {
	Header
	Body
}
type Header struct {
	Content_Type string `header:"Content-Type"`
	Timestamp    string `header:"Timestamp"`
	Sign         string `header:"Sign"`
}
type Body struct {
	SenderId                  string    `json:"senderId"`                  //加密的发送者ID。
	ConversationId            string    `json:"conversationId"`            //加密的会话ID。
	AtUsers                   []atUsers `json:"atUsers"`                   //被@人的信息。
	ChatbotCorpId             string    `json:"chatbotCorpId"`             //加密的机器人所在的企业corpId。
	ChatbotUserId             string    `json:"chatbotUserId"`             //加密的机器人id
	MsgId                     string    `json:"msgId"`                     //加密的消息ID。
	SenderNick                string    `json:"senderNick"`                //发送者昵称。
	IsAdmin                   bool      `json:"isAdmin"`                   //是否是管理员
	SenderStaffId             string    `json:"senderStaffId"`             //企业内部群中@该机器人的成员userid。
	SessionWebhookExpiredTime int64     `json:"sessionWebhookExpiredTime"` //当前会话的Webhook地址过期时间。
	CreateAt                  int64     `json:"createAt"`                  //消息的时间戳，单位ms。
	SenderCorpId              string    `json:"senderCorpId"`              //企业内部群有的发送者当前群的企业corpId。
	ConversationTitle         string    `json:"conversationTitle"`         //群聊时才有的会话标题。
	IsInAtList                bool      `json:"isInAtList"`                //是否在@列表中。
	SessionWebhook            string    `json:"sessionWebhook"`            //当前会话的Webhook地址。
	Text                      text      `json:"text"`                      //机器人收到的信息
	Msgtype                   string    `json:"msgtype"`                   //目前只支持text
	ConversationType          string    `json:"conversationType"`
}
type atUsers struct {
	DingtalkId string `json:"dingtalkId"`
	StaffId    string `json:"staffId"`
}
type text struct {
	Content string `json:"content"`
}
