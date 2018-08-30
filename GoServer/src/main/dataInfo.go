package main

import (

)

type UserInfo struct {
	Username		string		`bson:"username"`
	Hasconn 		bool		`bson:"hasconn"`
}

type MqttUserInfo struct {
	Username		string		`bson:"username"`
	Password		string		`bson:"password"`
	Hasconn 		bool		`bson:"hasconn"`
}

type MqttUserDeviceInfo struct {
	Username		string		`bson:"username"`
	Deviceid		string		`bson:"deviceid"`
}

type Message struct {
	Id				string		`bson:"id"`
	Sender			string		`bson:"sender"`
	Receiver 		string		`bson:"receiver"`
	SendTime		string		`bson:"sendtime"`
	Body			string		`bson:"body"`
}

type MqttMessage struct {
	Sender			string		`bson:"sender"`
	Receiver 		string		`bson:"receiver"`
	Device			string		`bson:"device"`
	SendTime		string		`bson:"sendtime"`
	Body			string		`bson:"body"`
}

type UserLog struct {
	Action			string		`bson:"action"`
	Info 			string		`bson:"info"`
	Time			string		`bson:"time"`
	Result			string		`bson:"result"`
}

type SendLog struct {
	Action			string		`bson:"action"`
	MessageId		string		`bson:"id"`
	Sender 			string		`bson:"sender"`
	Receiver 		string		`bson:"receiver"`
	Body 			string		`bson:"body"`
	Time			string		`bson:"time"`
	Result			string		`bson:"result"`
}

type MqttSendLog struct {
	Action			string		`bson:"action"`
	Sender 			string		`bson:"sender"`
	Receiver 		string		`bson:"receiver"`
	Body 			string		`bson:"body"`
	Time			string		`bson:"time"`
	Result			string		`bson:"result"`
}