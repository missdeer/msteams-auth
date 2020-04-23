package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	code           string
	botAccessToken string
)

type AccessTokenRequestBody struct {
	GrantType   string `json:"grant_type"`
	ClientID    string `json:"client_id"`
	Scope       string `json:"scope"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

type BotMessageFromRecipient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BotMessageConversation struct {
	ID string `json:"id"`
}

type BotEndpointMessage struct {
	Type         string                  `json:"type"`
	ID           string                  `json:"id"`
	TimeStamp    time.Time               `json:"timestamp"`
	ServiceURL   string                  `json:"serviceUrl"`
	ChannelID    string                  `json:"channelId"`
	From         BotMessageFromRecipient `json:"from"`
	Conversation BotMessageConversation  `json:"conversation"`
	Recipient    BotMessageFromRecipient `json:"recipient"`
	TextFormat   string                  `json:"textFormat"`
	Locale       string                  `json:"locale"`
	Text         string                  `json:"text"`
	ChannelData  struct {
		ClientActivityID string    `json:"clientActivityID"`
		ClientTimeStamp  time.Time `json:"clientTimestamp"`
		PostBack         bool      `json:"postBack,omitempty"`
	} `json:"channelData"`
	Value struct {
		PreferCallingTool string `json:"preferCallingTool,omitempty"`
	} `json:"value,omitempty"`
}

type BotReplyAttachmentTapAction struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BotReplyAttachmentImage struct {
	URL string                      `json:"url"`
	Alt string                      `json:"alt"`
	Tap BotReplyAttachmentTapAction `json:"tap"`
}

type BotReplyAttachmentButton struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Image string `json:"image,omitempty"`
	Value string `json:"value"`
}

type AdaptiveChoice struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type AdaptiveElement struct {
	Type          string           `json:"type"`
	Text          string           `json:"text"`
	Size          string           `json:"size,omitempty"`
	Separation    string           `json:"separation,omitempty"`
	ID            string           `json:"id,omitempty"`
	Style         string           `json:"style,omitempty"`
	IsMultiSelect bool             `json:"isMultiSelect,omitempty"`
	Value         string           `json:"value,omitempty"`
	Choices       []AdaptiveChoice `json:"choices,omitempty"`
}

type AdaptiveAction struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type BotReplyAttachmentContent struct {
	Type    string            `json:"type"`
	Version string            `json:"version"`
	Body    []AdaptiveElement `json:"body"`
	Actions []AdaptiveAction  `json:"actions,omitempty"`
}

type BotReplyAttachment struct {
	ContentType string                    `json:"contentType"`
	Content     BotReplyAttachmentContent `json:"content"`
}

type BotReplyMessage struct {
	Type         string                  `json:"type"`
	From         BotMessageFromRecipient `json:"from"`
	Conversation BotMessageConversation  `json:"conversation"`
	Recipient    BotMessageFromRecipient `json:"recipient"`
	Text         string                  `json:"text"`
	TextFormat   string                  `json:"textFormat"`
	ReplyToId    string                  `json:"replyToId"`
	Attachments  []BotReplyAttachment    `json:"attachments,omitempty"`
}

type BotAccessTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func getBotAccessToken() error {
	body := `grant_type=client_credentials&client_id=24dd1e52-103c-4214-94e8-2b9e5d3085a4&client_secret=0G8b0?l_ZXH-Tf/srNgp3JLS?Q:1Pooa&scope=https%3A%2F%2Fapi.botframework.com%2F.default`
	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token", strings.NewReader(body))
	if err != nil {
		log.Println("generate request failed:", err)
		return err
	}
	req.Header.Set("Content-Type", `application/x-www-form-urlencoded`)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("post request failed:", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response body failed", err)
		return err
	}
	var botAccessTokenResponse BotAccessTokenResponse
	err = json.Unmarshal(b, &botAccessTokenResponse)
	if err != nil {
		log.Println("unmarshalling bot access token response failed", err)
		return err
	}
	botAccessToken = botAccessTokenResponse.AccessToken
	log.Println("bot access token:", botAccessTokenResponse)
	return nil
}

func replyBotMessage(msg BotEndpointMessage) error {
	id := strings.Split(msg.ID, "|")[0]
	reply := &BotReplyMessage{
		Type: `message`,
		//Text:         `I have received **` + msg.Text + `**`,
		TextFormat:   `markdown`,
		ReplyToId:    id,
		From:         msg.Recipient,
		Recipient:    msg.From,
		Conversation: msg.Conversation,
		Attachments: []BotReplyAttachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: BotReplyAttachmentContent{
					Type:    "AdaptiveCard",
					Version: "1.0",
					Body: []AdaptiveElement{
						{
							Type: "TextBlock",
							Text: "Choose Calling Tool",
						},
						{
							Type:          "Input.ChoiceSet",
							ID:            "preferCallingTool",
							Style:         "expanded",
							IsMultiSelect: false,
							Value:         "Cisco Jabber",
							Choices: []AdaptiveChoice{
								{
									Title: "Cisco Jabber",
									Value: "Cisco Jabber",
								},
								{
									Title: "Cisco WebEx Teams",
									Value: "Cisco WebEx Teams",
								},
							},
						},
					},
					Actions: []AdaptiveAction{
						{
							Type:  "Action.Submit",
							Title: "OK",
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(reply)
	if err != nil {
		log.Println("marshalling reply message failed:", err)
		return err
	}
	u := fmt.Sprintf(`%sv3/conversations/%s/activities/%s`, msg.ServiceURL, msg.Conversation.ID, id)
	req, err := http.NewRequest("POST", u, bytes.NewReader(b))
	if err != nil {
		log.Println("generate request failed:", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+botAccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("post request failed:", err)
		return err
	}
	defer resp.Body.Close()
	log.Println("post reply message:", u, string(b))
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response body failed", err)
		return err
	}
	log.Println("got reply message response:", string(b))
	return nil
}

func botReplySetPreferCallingTool(msg BotEndpointMessage) error {
	id := strings.Split(msg.ID, "|")[0]
	reply := &BotReplyMessage{
		Type:         `message`,
		Text:         `You are preferring **` + msg.Value.PreferCallingTool + `**`,
		TextFormat:   `markdown`,
		ReplyToId:    id,
		From:         msg.Recipient,
		Recipient:    msg.From,
		Conversation: msg.Conversation,
	}

	b, err := json.Marshal(reply)
	if err != nil {
		log.Println("marshalling reply message failed:", err)
		return err
	}
	u := fmt.Sprintf(`%sv3/conversations/%s/activities/%s`, msg.ServiceURL, msg.Conversation.ID, id)
	req, err := http.NewRequest("POST", u, bytes.NewReader(b))
	if err != nil {
		log.Println("generate request failed:", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+botAccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("post request failed:", err)
		return err
	}
	defer resp.Body.Close()
	log.Println("post reply message:", u, string(b))
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response body failed", err)
		return err
	}
	log.Println("got reply message response:", string(b))
	return nil
}

func botEndpoint(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	var msg BotEndpointMessage
	err = json.Unmarshal(rawData, &msg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	log.Println("botEndpoint raw data:", string(rawData), msg)
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
	if msg.Value.PreferCallingTool != "" {
		log.Println("botEndpoint got value:", msg.Value.PreferCallingTool)
		botReplySetPreferCallingTool(msg)
		return
	}
	replyBotMessage(msg)
}

func main() {
	if err := getBotAccessToken(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("bot-endpoint", botEndpoint)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://login.microsoftonline.com/afae2f63-1bcb-4d1f-b8c3-252a4cd3dd07/oauth2/v2.0/authorize?client_id=46442420-1b26-4bd7-a997-183e1880bbd5&response_type=code&redirect_uri=http://localhost:8765/individual_user_consent/&response_mode=query&scope=offline_access%20user.read.all%20chat.read&state=12345")
	})

	r.GET("/individual_user_consent/", func(c *gin.Context) {
		responseCode := c.Query("code")
		log.Println("code", responseCode)
		if responseCode != "" {
			code = responseCode
			c.Redirect(http.StatusFound, "https://login.microsoftonline.com/afae2f63-1bcb-4d1f-b8c3-252a4cd3dd07/v2.0/adminconsent?client_id=46442420-1b26-4bd7-a997-183e1880bbd5&state=12345&redirect_uri=http://localhost:8765/individual_user_consent/&scope=offline_access%20user.read.all%20chat.read")
			return
		}
		admin_consent := c.Query("admin_consent")
		log.Println("admin_consent:", admin_consent)
		state := c.Query("state")
		log.Println("state:", state)
		scope := c.Query("scope")
		log.Println("scope:", scope)

		req, err := http.NewRequest("POST", "https://login.microsoftonline.com/afae2f63-1bcb-4d1f-b8c3-252a4cd3dd07/oauth2/v2.0/token",
			strings.NewReader(`client_secret=msHRpSOTQLP24lCk9afnSTejW%3DlV%3F8%3D%40&grant_type=authorization_code&client_id=46442420-1b26-4bd7-a997-183e1880bbd5&scope=offline_access%20user.read.all%20chat.read&redirect_uri=http://localhost:8765/individual_user_consent/&code=`+code))
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var response gin.H
		if err = json.Unmarshal(respBody, &response); err != nil {
			log.Fatal(err)
		}

		accessToken, ok := response["access_token"]
		if ok {
			log.Println("access_token", accessToken)
		}

		refreshToken, ok := response["refresh_token"]
		if ok {
			log.Println("refresh_token", refreshToken)
		}

		c.JSON(http.StatusOK, &response)
	})

	bind := os.Getenv("BINDADDR")
	if bind == "" {
		bind = "127.0.0.1:8765"
	}
	log.Fatal(r.Run(bind))
}
