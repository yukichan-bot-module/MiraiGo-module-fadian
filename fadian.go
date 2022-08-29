package fadian

import (
	_ "embed" // go embed
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

//go:embed assets/post.json
var postJSONData []byte

var instance *fadian
var logger = utils.GetModuleLogger("com.aimerneige.fadian")

type fadian struct {
}

// PostJSON post json object
type PostJSON struct {
	Post []string `json:"post"`
}

func init() {
	instance = &fadian{}
	bot.RegisterModule(instance)
}

func (f *fadian) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.fadian",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (f *fadian) Init() {
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (f *fadian) PostInit() {
}

// Serve 注册服务函数部分
func (f *fadian) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		textMsg := ""
		isText := false
		for _, element := range msg.Elements {
			switch e := element.(type) {
			case *message.TextElement:
				if !isText {
					textMsg = e.Content
					isText = true
				}
			}
		}
		if !isText {
			return
		}
		if strings.HasPrefix(textMsg, "每日发癫") {
			name := "小乌贼"
			nameString := strings.TrimSpace(textMsg[12:])
			if nameString != "" {
				name = nameString
			}
			replyText := getFadianText(name)
			c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText(replyText)))
		}
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (f *fadian) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (f *fadian) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func getFadianText(name string) string {
	var postJSON PostJSON
	err := json.Unmarshal(postJSONData, &postJSON)
	if err != nil {
		logger.WithError(err).Error("Fail to unmarshal post.json")
		return "解析 JSON 失败，请查阅后台日志。"
	}
	rand.Seed(time.Now().Unix())
	postString := postJSON.Post[rand.Intn(len(postJSON.Post))]
	postString = strings.ReplaceAll(postString, "阿咪", name)
	return postString
}
