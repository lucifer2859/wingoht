package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbUrl = "mongodb://192.168.200.3:27017"
var dbName = "mq_userlist"

func hasUser(username string) bool {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("user_info")

	var user UserInfo
	err = userList.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return false
	}
	return true
}

func IsHasConnection(username string) int {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("user_info")

	var user UserInfo
	err = userList.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return -1
	}
	if user.Hasconn {
		return 1
	} else {
		return 0
	}
}

func addUser(username string) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()
	
	var user UserInfo
	user.Username = username
	user.Hasconn = false
	
	userList := session.DB(dbName).C("user_info")

	err = userList.Insert(user)
	fatalErrorLog(err, "userList.Insert")
}

func updateUser(username string, change bool) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("user_info")

	var user UserInfo
	user.Username = username
	user.Hasconn = change
	
	err = userList.Update(bson.M{"username":username}, &user)
	fatalErrorLog(err, "userList.Update")
}

func getMessage(receiver string) []Message {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("message_info")

	var messages []Message
	err = messageList.Find(bson.M{"receiver": receiver}).All(&messages)
	fatalErrorLog(err, "messageList.Find")
	return messages
}

func addMessage(msg Message) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("message_info")

	err = messageList.Insert(msg)
	fatalErrorLog(err, "messageList.Insert")
}

func removeMessage(receiver string) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("message_info")

	_ , err = messageList.RemoveAll(bson.M{"receiver": receiver})
	fatalErrorLog(err, "messageList.RemoveAll")
}

func userLogInfo(userlog UserLog) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userLogList := session.DB(dbName).C("user_log_info")

	err = userLogList.Insert(userlog)
	fatalErrorLog(err, "userLogList.Insert")
}

func sendLogInfo(sendlog SendLog) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	sendLogList := session.DB(dbName).C("send_log_info")

	err = sendLogList.Insert(sendlog)
	fatalErrorLog(err, "sendLogList.Insert")
}

func mqttHasUser(username string) bool {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("mqtt_user_info")

	var user UserInfo
	err = userList.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return false
	}
	return true
}

func mqttConfirmUser(username, password string) int {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("mqtt_user_info")

	var user MqttUserInfo
	err = userList.Find(bson.M{"username": username, "password": password}).One(&user)
	if err != nil {
		return -1
	}
	
	if user.Hasconn {
		return 1
	} else {
		return 0
	}
}

func mqttHasDevice(username, deviceid string) bool {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("mqtt_user_device_info")

	var user MqttUserDeviceInfo
	err = userList.Find(bson.M{"username": username, "deviceid": deviceid}).One(&user)
	if err != nil {
		return false
	}
	return true
} 

func mqttAddUser(user MqttUserInfo) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()
	
	userList := session.DB(dbName).C("mqtt_user_info")

	err = userList.Insert(user)
	fatalErrorLog(err, "mqttUserList.Insert")
}

func mqttUpdateUser(username string, change bool) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("mqtt_user_info")
	
	var user MqttUserInfo
	err = userList.Find(bson.M{"username": username}).One(&user)
	fatalErrorLog(err, "MqttUserList.Find")
	
	user.Hasconn = change
	
	err = userList.Update(bson.M{"username":username}, &user)
	fatalErrorLog(err, "MqttUserList.Update")
}

func mqttAddUserDevice(username, deviceid string) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()
	
	userList := session.DB(dbName).C("mqtt_user_device_info")

	var user MqttUserDeviceInfo
	user.Username = username
	user.Deviceid = deviceid
	err = userList.Insert(user)
	fatalErrorLog(err, "mqttUserDeviceList.Insert")
}

func mqttRemoveUserDevice(username, deviceid string) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()
	
	userList := session.DB(dbName).C("mqtt_user_device_info")
	
	_ , err = userList.RemoveAll(bson.M{"username": username, "deviceid": deviceid})
	fatalErrorLog(err, "mqttUserDeviceList.RemoveAll")
}

func mqttGetDevice(username string) []string {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("mqtt_user_device_info")

	var userDevices []MqttUserDeviceInfo
	err = messageList.Find(bson.M{"username": username}).All(&userDevices)
	fatalErrorLog(err, "mqttMessageList.Find")
	
	var devices []string
	
	for _, userDevice := range userDevices {
		devices = append(devices, userDevice.Deviceid)
	}
	 
	return devices
}

func mqttIsHasConnection(username string) int {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userList := session.DB(dbName).C("mqtt_user_info")

	var user MqttUserInfo
	err = userList.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return -1
	}
	
	if user.Hasconn {
		return 1
	} else {
		return 0
	}
}

func mqttAddMessage(msg MqttMessage) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("mqtt_message_info")

	err = messageList.Insert(msg)
	fatalErrorLog(err, "mqttMessageList.Insert")
}

func mqttGetMessage(receiver string) []MqttMessage {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("mqtt_message_info")

	var messages []MqttMessage
	err = messageList.Find(bson.M{"receiver": receiver}).All(&messages)
	fatalErrorLog(err, "mqttMessageList.Find")
	return messages
}

func mqttRemoveMessage(receiver string) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	messageList := session.DB(dbName).C("mqtt_message_info")

	_ , err = messageList.RemoveAll(bson.M{"receiver": receiver})
	fatalErrorLog(err, "mqttMessageList.RemoveAll")
}

func mqttUserLogInfo(userlog UserLog) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	userLogList := session.DB(dbName).C("mqtt_user_log_info")

	err = userLogList.Insert(userlog)
	fatalErrorLog(err, "mqttUserLogList.Insert")
}

func mqttSendLogInfo(sendlog MqttSendLog) {
	session, err := mgo.Dial(dbUrl)
	fatalErrorLog(err, "mgo.Dial")
	defer session.Close()

	sendLogList := session.DB(dbName).C("mqtt_send_log_info")

	err = sendLogList.Insert(sendlog)
	fatalErrorLog(err, "mqttSendLogList.Insert")
}