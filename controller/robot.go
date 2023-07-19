package controller

import (
	"InterviewPush/dao/mysql"
	"InterviewPush/logic"
	"InterviewPush/models"
	reqm "InterviewPush/models/request"
	"InterviewPush/pkg"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var groupIds []string

func init() {
	// 1 三期校招   2 四期
	groupIds = append(groupIds, "cid7i28j4j/4R0EtyezLtSMtg==", "cidoXnlPCUoRU5W60UFGyHT8g==")
}

// PushInformation 填写表单后发送给大学姐或者李老师
func PushInformation(p *reqm.InterviewMsg) {
	//面试类型
	var interviewType string
	switch p.InterviewType {
	case 1:
		interviewType = "实习"
	case 2:
		interviewType = "校招"
	case 3:
		interviewType = "社招"
	}
	//约面试途径
	var interviewApproach string
	switch p.InterviewApproach {
	case 1:
		interviewApproach = "网申"
	case 2:
		interviewApproach = "boss"
	case 3:
		interviewApproach = "脉脉"
	case 4:
		interviewApproach = "内推"
	}
	//发消息
	userIDs := make([]string, 3)
	userIDs[0] = "011433126224-1070024747" //张腾飞
	userIDs[1] = "manager6736"             //李老师
	userIDs[2] = "19036662541994219728"    //大学姐
	//群聊推送三期四期
	groupchatS := fmt.Sprintf("%s%s在%s有一场面试，欢迎各位同学前来聆听", p.InterviewUsername, p.InterviewTime.Format("2006-01-02 15:04"), p.InterviewLocation)
	for i := 0; i < 2; i++ {
		GroupChatRobot(groupchatS, groupIds[i])
	}
	//单聊推送
	s := fmt.Sprintf("面试信息推送：\n面试人:%s\n面试公司:%s\n面试岗位:%s\n面试地点:%s\n面试类型:%s\n约面试途径:%s\n面试困惑:%s\n面试时间:%s\n", p.InterviewUsername, p.InterviewCompany, p.InterviewPosition, p.InterviewLocation, interviewType, interviewApproach, p.InterviewConfusion, pkg.UTCTransLocal(p.InterviewTime.Format("2006-01-02T15:04:05.000Z")))
	SingleChatRobot(userIDs, s)
}

// ReceptionHandler 机器人接收信息
func ReceptionHandler(c *gin.Context) {
	var p models.DingParamReveiver
	err := c.ShouldBindJSON(&p)
	err = c.ShouldBindHeader(&p)
	if err != nil {
		return
	}
	//单聊群聊 p.ConversationType	@成员 p.SenderStaffId	被@人的信息 p.AtUsers[0].DingtalkId p.AtUsers[0].StaffId		是否在@列表中 p.IsInAtList
	p.Text.Content = strings.Trim(p.Text.Content, " ")
	userIDs := make([]string, 1)
	userIDs[0] = p.SenderStaffId
	var s string
	if p.Text.Content == "你好" {
		s = "你好，小面很高兴为你服务，如果有需要请发送help"
	} else if p.Text.Content == "help" {
		s = fmt.Sprintf("可以对小面说：\n1.填写面试表单\n2.查看个人面试记录\n3.查看面试记录-xxx")
	} else if p.Text.Content == "填写面试表单" {
		s = fmt.Sprintf("面试信息模板如下，复制下面的内容编辑完发送即可\n*面试人：%s\n*面试公司：百度\n*面试岗位：后端研发工程师\n*面试地点：0#906\n*面试类型：(实习、校招、社招)\n*约面试途径：(网申、boss、脉脉、内推)\n*面试困惑：目前没有\n*面试时间：格式2006-01-02 15:04:05\n", p.SenderNick)
	} else if strings.Contains(p.Text.Content, "查看") {
		if p.Text.Content == "查看个人面试记录" {
			s = ViewInterviewTranscripts(p.SenderNick)
		} else {
			arr := strings.Split(p.Text.Content, "-")
			if len(arr) == 2 {
				s = ViewInterviewTranscripts(arr[1])
			} else {
				s = "请按照第三条预定格式发送"
			}
		}
	} else if strings.ContainsAny(p.Text.Content, "面试人") {
		Msg := strings.FieldsFunc(p.Text.Content, func(c rune) bool {
			return c == '：' || c == '*' || c == '\n'
		})
		s = ChatValidated(Msg, p.SenderNick)
	} else {
		s = "不好意思，小面还看不懂呢，如果想问关于面试的信息，请发送help"
	}
	//单聊
	if p.ConversationType == "1" {
		SingleChatRobot(userIDs, s)
	} else {
		GroupChatRobot(s, p.ConversationId)
	}
}

// 对面试信息进行参数校验
func ChatValidated(Msg []string, Name string) (s string) {
	interviewType := "实习校招社招"
	interviewApproach := "网申boss脉脉内推"
	// Msg[1]面试人 	Msg[3]面试公司 	Msg[5]面试岗位 	Msg[7]面试地点 	Msg[9]面试类型 	Msg[11]约面试途径 	Msg[13]面试困惑 	Msg[15]面试时间
	// 完整性检验
	if len(Msg) != 16 {
		return "填写面试格式不正确，请按照下文复制后编辑后再试\n*面试人：" + Name + "\n*面试公司：百度\n*面试岗位：后端研发工程师\n*面试地点：0#906\n*面试类型：(实习、校招、社招)\n*约面试途径：(网申、boss、脉脉、内推)\n*面试困惑：目前没有\n*面试时间：格式2006-01-02 15:04:05\n"
	}
	//非空检验
	for i := 0; i < len(Msg); i++ {
		Msg[i] = strings.Trim(Msg[i], "\n")
		if Msg[i] == "" {
			return "未填完，请补充后再提交"
		}
		if len(Msg[i]) < 2 && i != 13 {
			return "无效数据，请修改后再试"
		}
	}
	interviewtype := strings.Index(interviewType, Msg[9])          //匹配面试类型
	interviewapproach := strings.Index(interviewApproach, Msg[11]) //匹配约面试途径
	// 中文校验
	isChinese := regexp.MustCompile("^[\u4e00-\u9fa5]") //中文匹配原则
	if !isChinese.MatchString(Msg[1]) {
		return "名字填写错误"
	}
	if interviewtype == -1 {
		return "面试类型填写错误"
	}
	if interviewapproach == -1 {
		return "约面试途径填写错误"
	}
	// 时间校验
	interviewtime, err := time.ParseInLocation("2006-01-02 15:04:05", Msg[15], time.Local)
	if interviewtime.Before(time.Now()) {
		return "时间不正确"
	} else if err != nil {
		return "时间格式不正确"
	}
	switch Msg[9] {
	case "实习":
		interviewtype = 1
	case "校招":
		interviewtype = 2
	case "社招":
		interviewtype = 3
	}
	switch Msg[11] {
	case "网申":
		interviewapproach = 1
	case "boss":
		interviewapproach = 2
	case "脉脉":
		interviewapproach = 3
	case "内推":
		interviewapproach = 4
	}
	interviewMsg := reqm.InterviewMsg{
		InterviewUsername:  Msg[1],
		InterviewCompany:   Msg[3],
		InterviewPosition:  Msg[5],
		InterviewLocation:  Msg[7],
		InterviewType:      interviewtype,
		InterviewApproach:  interviewapproach,
		InterviewConfusion: Msg[13],
		InterviewTime:      interviewtime,
	}
	//加入数据库中
	err = logic.AddMessage(interviewMsg)
	if err != nil {
		return err.Error()
	}
	//单聊推送
	singlechatS := fmt.Sprintf("面试信息推送：\n面试人:%s\n面试公司:%s\n面试岗位:%s\n面试地点:%s\n面试类型:%s\n约面试途径:%s\n面试困惑:%s\n面试时间:%s\n", Msg[1], Msg[3], Msg[5], Msg[7], Msg[9], Msg[11], Msg[13], interviewtime.Format("2006-01-02 15:04:05"))
	userIDs := make([]string, 3)
	userIDs[0] = "011433126224-1070024747"
	userIDs[1] = "manager6736"
	userIDs[2] = "19036662541994219728"
	SingleChatRobot(userIDs, singlechatS)
	//群聊推送
	groupchatS := fmt.Sprintf("%s%s在%s有一场面试，欢迎各位同学前来聆听", Msg[1], interviewtime.Format("2006-01-02 15:04"), Msg[7])
	for i := 0; i < 2; i++ {
		GroupChatRobot(groupchatS, groupIds[i])
	}
	return "添加面试记录成功"
}

// 单聊机器人回复信息
func SingleChatRobot(UserIDs []string, s string) {
	accesstoken, err := GetAccessToken()
	//发消息
	var client *http.Client
	var request *http.Request
	var resp *http.Response
	URL := "https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend"
	client = &http.Client{Transport: &http.Transport{ //对客户端进行一些配置
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, Timeout: time.Duration(time.Second * 5)}
	//此处是post请求的请求体，我们先初始化一个对象
	b := struct {
		MsgParam  string   `json:"msgParam"`
		MsgKey    string   `json:"msgKey"`
		RobotCode string   `json:"robotCode"`
		UserIds   []string `json:"userIds"`
	}{MsgParam: fmt.Sprintf("{       \"content\": \"%s\"   }", s),
		MsgKey:    "sampleText",
		RobotCode: "dingjfxu75nqiffmdbb8",
		UserIds:   UserIDs,
	}
	//然后把结构体对象序列化一下
	bodymarshal, err := json.Marshal(&b)
	if err != nil {
		return
	}
	//再处理一下
	reqBody := strings.NewReader(string(bodymarshal))
	//然后就可以放入具体的request中的
	request, _ = http.NewRequest(http.MethodPost, URL, reqBody)
	request.Header.Set("x-acs-dingtalk-access-token", accesstoken) //拿到access-token
	request.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(request)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) //把请求到的body转化成byte[]
	h := struct {
		Code                      string   `json:"code"`
		Message                   string   `json:"message"`
		ProcessQueryKey           string   `json:"processQueryKey"`
		InvalidStaffIdList        []string `json:"invalidStaffIdList"`
		FlowControlledStaffIdList []string `json:"flowControlledStaffIdList"`
	}{}
	//把请求到的结构反序列化到专门接受返回值的对象上面
	err = json.Unmarshal(body, &h)
}

// openConversationId ===== cidvJVIH/fZbHC3bbyQfg8CWg==   面试推送系统测试群id
// 群聊推送信息
func GroupChatRobot(s string, openConversationId string) {
	accesstoken, err := GetAccessToken()
	//发消息
	var client *http.Client
	var request *http.Request
	var resp *http.Response
	URL := "https://api.dingtalk.com/v1.0/robot/groupMessages/send"
	client = &http.Client{Transport: &http.Transport{ //对客户端进行一些配置
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, Timeout: time.Duration(time.Second * 5)}
	//此处是post请求的请求体，我们先初始化一个对象
	b := struct {
		MsgParam           string `json:"msgParam"`
		MsgKey             string `json:"msgKey"`
		OpenConversationId string `json:"openConversationId"`
		RobotCode          string `json:"robotCode"`
	}{MsgParam: fmt.Sprintf("{       \"content\": \"%s\"   }", s),
		MsgKey:             "sampleText",
		OpenConversationId: openConversationId,
		RobotCode:          "dingjfxu75nqiffmdbb8",
	}
	//然后把结构体对象序列化一下
	bodymarshal, err := json.Marshal(&b)
	if err != nil {
		return
	}
	//再处理一下
	reqBody := strings.NewReader(string(bodymarshal))
	//然后就可以放入具体的request
	request, _ = http.NewRequest(http.MethodPost, URL, reqBody)
	request.Header.Set("x-acs-dingtalk-access-token", accesstoken) //拿到access-token
	request.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(request)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) //把请求到的body转化成byte[]
	h := struct {
		Code                      string   `json:"code"`
		Message                   string   `json:"message"`
		ProcessQueryKey           string   `json:"processQueryKey"`
		InvalidStaffIdList        []string `json:"invalidStaffIdList"`
		FlowControlledStaffIdList []string `json:"flowControlledStaffIdList"`
	}{}
	//把请求到的结构反序列化到专门接受返回值的对象上面
	err = json.Unmarshal(body, &h)
}

// 查看面试记录
func ViewInterviewTranscripts(sendName string) string {
	var s string
	interviewMessage, err := mysql.SelectMessageByName(sendName)
	if err != nil {
		zap.L().Error("查看面试记录错误", zap.Error(err))
		return "服务器繁忙，请稍后再试"
	}
	if len(interviewMessage) == 0 {
		s = "此人最近没有面试记录"
	} else {
		var interviewType string
		var interviewApproach string
		s = fmt.Sprintf("面试信息推送：\n")
		for _, msg := range interviewMessage {
			//面试类型
			switch msg.InterviewType {
			case 1:
				interviewType = "实习"
			case 2:
				interviewType = "校招"
			case 3:
				interviewType = "社招"
			}
			//约面试途径
			switch msg.InterviewApproach {
			case 1:
				interviewApproach = "网申"
			case 2:
				interviewApproach = "boss"
			case 3:
				interviewApproach = "脉脉"
			case 4:
				interviewApproach = "内推"
			}
			ss := fmt.Sprintf("面试人:%s\n面试公司:%s\n面试岗位:%s\n面试地点:%s\n面试类型:%s\n约面试途径:%s\n面试困惑:%s\n面试时间:%s\n\n", msg.InterviewUsername, msg.InterviewCompany, msg.InterviewPosition, msg.InterviewLocation, interviewType, interviewApproach, msg.InterviewConfusion, pkg.UTCTransLocal(msg.InterviewTime.Format("2006-01-02T15:04:05.000Z")))
			s += ss
		}
	}
	return s
}
