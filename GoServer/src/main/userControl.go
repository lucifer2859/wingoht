package main

import (
    "net/http"
    "io/ioutil"
	"encoding/json" 
	"time"
	"strconv"
	"github.com/streadway/amqp"
)

func checkToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	r.ParseForm()
	token := r.Form["token"][0]
	
	t := time.Now()
	timestamp := strconv.FormatInt(t.UnixNano() + 8*60*60*1000000000, 10)
	timestamp = timestamp[:13]
	
	var userlog UserLog
	userlog.Action = "Token验证"
	userlog.Info = token
	userlog.Time = t.Format("2006-01-02 15:04:05")
	
	client := &http.Client{}
    url := "http://192.168.40.12:8081/com.wingoht.yiyun.core/api/User/ConfirmAccessToken/"
    
    reqest, err := http.NewRequest("GET", url, nil)

    reqest.Header.Set("Content-Type", "application/json")
    reqest.Header.Add("AccessToken", token)
    reqest.Header.Add("timestamp", timestamp)
    reqest.Header.Add("token", "00000000-0000-0000-0000-000000000000")

    fatalErrorLog(err, "http.NewRequest")

    res, err := client.Do(reqest)  
	
	defer res.Body.Close()
	
    jsonStr, err := ioutil.ReadAll(res.Body)
   
    fatalErrorLog(err, "ioutil.ReadAll")

	var msg TokenResponse
	err = json.Unmarshal(jsonStr, &msg)
	
	fatalErrorLog(err, "json.Unmarshal")
        
	var result CheckResponse
	
	if msg.VerLogin == 0 {
		w.WriteHeader(http.StatusOK)
		result.Data = msg.Data
		if hasUser(msg.Data) {		
			result.Message = "0"
			userlog.Result = "token验证成功，已注册用户！"
		} else {
			addUser(result.Data)
			result.Message = "1"
			userlog.Result = "token验证成功，新注册用户！"
		}
		
		conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/cd")
	    fatalErrorLog(err, "amqp.Dial")
	    defer conn.Close()

	    ch, err := conn.Channel()
	    fatalErrorLog(err, "conn.Channel")
	    defer ch.Close()

		_, err = ch.QueueDeclare(
			  msg.Data, // name
			  true,   // durable
			  true,   // delete when unused
			  false,   // exclusive
			  false,   // no-wait
			  nil,     // arguments
		)
		
		fatalErrorLog(err, "ch.QueueDeclare")
		
		userLogInfo(userlog)
		
	} else {
		userlog.Result = "token验证失败！"
		userLogInfo(userlog)
		
		w.WriteHeader(http.StatusNotFound)
		result.Message = "2"
	}
	
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}

func getUnreadMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	r.ParseForm()
	receiver := r.Form["user"][0]
	
	var userlog UserLog
	userlog.Action = "获取离线消息"
	userlog.Info = receiver
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")

	result := getMessage(receiver)
	
	if len(result) > 0 {
		removeMessage(receiver)	
	}
	
	userLogInfo(userlog)
	
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	username := jsonReq["Username"]
	password := jsonReq["Password"]
	
	var userlog UserLog
	userlog.Action = "mqtt用户登录"
	userlog.Info = "username=" + username + "&password=" + password
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")
	
	res := mqttConfirmUser(username, password)
	
	var msg string
	
	if res == 0 {
		mqttUpdateUser(username, true)	
		userlog.Result = "用户登录成功！"
		w.WriteHeader(http.StatusOK)
		msg = "0"	
	} else if res > 0 {
		userlog.Result = "该用户已登录！"
		w.WriteHeader(http.StatusNotFound)
		msg = "1"
	} else if mqttHasUser(username) {		
		userlog.Result = "用户名与密码不匹配！"
		w.WriteHeader(http.StatusNotFound)
		msg = "2"
	} else {
		var user MqttUserInfo
		user.Username = username
		user.Password = password
		user.Hasconn = true
		mqttAddUser(user)
		userlog.Result = "新用户注册登录成功！"
		w.WriteHeader(http.StatusOK)
		msg = "3"
	}	
	
	mqttUserLogInfo(userlog)
	w.Write([]byte(msg))
}

func userLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	r.ParseForm()
	username := r.Form["user"][0]
	
	var userlog UserLog
	userlog.Action = "mqtt用户登出"
	userlog.Info = username
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")
	
	mqttUpdateUser(username, false)
	
	mqttUserLogInfo(userlog)
	
	w.WriteHeader(http.StatusOK)
}

func addDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	username := jsonReq["Username"]
	deviceid := jsonReq["DeviceId"]
	
	var userlog UserLog
	userlog.Action = "mqtt添加设备"
	userlog.Info = "username=" + username + "&deviceid=" + deviceid
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")	
	
	var msg string
	
	if mqttHasUser(username) {
		if mqttHasDevice(username, deviceid) {
			userlog.Result = "该设备已被添加！"
			w.WriteHeader(http.StatusNotFound)
			msg = "0"
		} else {
			mqttAddUserDevice(username, deviceid)
			userlog.Result = "设备添加成功！"
			w.WriteHeader(http.StatusOK)
			msg = "1"
		}
	} else {
		userlog.Result = "该用户不存在！"
		w.WriteHeader(http.StatusNotFound)
		msg = "2"
	}
	
	mqttUserLogInfo(userlog)
	
	w.Write([]byte(msg))
}

func removeDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	username := jsonReq["Username"]
	deviceid := jsonReq["DeviceId"]
	
	var userlog UserLog
	userlog.Action = "mqtt删除设备"
	userlog.Info = "username=" + username + "&deviceid=" + deviceid
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")	
	
	var msg string
	
	if mqttHasUser(username) {
		if mqttHasDevice(username, deviceid) {
			mqttRemoveUserDevice(username, deviceid)
			userlog.Result = "设备删除成功！"
			w.WriteHeader(http.StatusOK)
			msg = "0"
		} else {
			userlog.Result = "该设备未被添加！"
			w.WriteHeader(http.StatusNotFound)
			msg = "1"
		}
	} else {
		userlog.Result = "该用户不存在！"
		w.WriteHeader(http.StatusNotFound)
		msg = "2"
	}
	
	mqttUserLogInfo(userlog)
	
	w.Write([]byte(msg))
}

func showDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	r.ParseForm()
	username := r.Form["user"][0]
	
	var userlog UserLog
	userlog.Action = "mqtt设备显示"
	userlog.Info = username
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")
	
	result := mqttGetDevice(username)
	
	mqttUserLogInfo(userlog)
	
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}

func getUnreadMqttMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	r.ParseForm()
	receiver := r.Form["user"][0]
	
	var userlog UserLog
	userlog.Action = "mqtt获取离线消息"
	userlog.Info = receiver
	userlog.Time = time.Now().Format("2006-01-02 15:04:05")

	result := mqttGetMessage(receiver)
	
	if len(result) > 0 {
		mqttRemoveMessage(receiver)	
	}
	
	mqttUserLogInfo(userlog)
	
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}
