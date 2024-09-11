package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// FCM单条推送消息中的通知内容
type FCMNotification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Icon      string `json:"icon,omitempty"`
	EventTime string `json:"event_time,omitempty"`
}

// FCM单条发送消息
type FCMPushMsg struct {
	Notification FCMNotification   `json:"notification"`
	Token        string            `json:"token,omitempty"`
	Topic        string            `json:"topic,omitempty"`
	Data         map[string]string `json:"data"`
}

// FCM单条发送请求
type FCMPushRequest struct {
	Message FCMPushMsg `json:"message"`
}

// FCM单条发送结果
type FCMPushResult struct {
	Name  string       `json:"name"`
	Error FCMPushError `json:"error"`
}

type FCMPushError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type PushMsg struct {
	Token        string `protobuf:"bytes,1,opt,name=token" json:"token"`
	Title        string `protobuf:"bytes,2,opt,name=title" json:"title"`
	Content      string `protobuf:"bytes,3,opt,name=content" json:"content"`
	PushSeq      string `protobuf:"bytes,4,opt,name=pushSeq" json:"pushSeq"`
	PushType     string `protobuf:"bytes,5,opt,name=pushType" json:"pushType"`
	PushConfigId string `protobuf:"bytes,6,opt,name=pushConfigId" json:"pushConfigId"`
	Icon         string `protobuf:"bytes,7,opt,name=icon" json:"icon"`
	EventTime    int64  `protobuf:"varint,8,opt,name=eventTime" json:"eventTime"`
	Topic        string `protobuf:"bytes,9,opt,name=topic" json:"topic"`
}

func main() {
	fmt.Println(send("http://127.0.0.1:9090/", "ssssssssssssssssssssssssssssssssssssssssssssss"))
}

func send(sendApi string, accessToken string) error {
	msg := &PushMsg{
		Title:     "title",
		Content:   "dddddddddddddddddddddddd",
		Icon:      "icon",
		EventTime: 123123123,
		Token:     "token",
	}
	fcmReq := FCMPushRequest{
		Message: FCMPushMsg{
			Notification: FCMNotification{
				Title:     msg.Title,
				Body:      msg.Content,
				Icon:      msg.Icon,
				EventTime: time.UnixMilli(msg.EventTime).Format(time.RFC3339),
			},
			Token: msg.Token,
			Topic: msg.Topic,
			Data: map[string]string{
				"pushSeq":      msg.PushSeq,
				"pushType":     msg.PushType,
				"pushConfigId": msg.PushConfigId,
			},
		},
	}
	data, err := json.Marshal(fcmReq)
	if err != nil {
		return fmt.Errorf("json.Marshal(%#v) error: %s", fcmReq, err.Error())
	}
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", sendApi, buff)
	if err != nil {
		return fmt.Errorf("http.NewRequest() error: %s", err.Error())
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	//生成请求对象
	cli := &http.Client{Timeout: time.Second * 30}
	//发送请求
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("http POST %s error: %s, data: %s", sendApi, err.Error(), string(data))
	}
	//判断状态
	if resp.StatusCode != 200 {
		return fmt.Errorf("http POST %s status: %d, data: %s", sendApi, resp.StatusCode, string(data))
	}
	//读取返回数据
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("http POST %s failed, ioutil.ReadAll() error: %s|data: %s", sendApi, err.Error(), string(data))
	}
	//解析结果
	ret := FCMPushResult{}
	if err = json.Unmarshal(content, &ret); err != nil {
		return fmt.Errorf("json.Unmarshal(%s) error: %s", string(content), err.Error())
	}
	//判断结果
	if ret.Error.Code != 0 {
		return fmt.Errorf("push %s error: %s", string(data), string(content))
	}

	return nil
}
