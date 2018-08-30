package main

import (
	"net/http"
	"github.com/streadway/amqp"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/eclipse/paho.mqtt.golang"
)

func messageCheck(content string) bool {
	return false
}

func hasRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	messageId := jsonReq["MessageId"]
	receiver := jsonReq["Receiver"]
	
	c, err := redis.Dial("tcp", "192.168.200.3:6379")
	fatalErrorLog(err, "redis.Dial")
	defer c.Close()
	
	unread, err := redis.Int((c.Do("HGET", messageId, "unread")))
	if err == nil {
		sender, _ := redis.String((c.Do("HGET", messageId, "sender")))
		
		conns := IsHasConnection(sender)
		
		if conns < 0 {	
			w.WriteHeader(http.StatusNotFound)
		} else {
			var message Message
			message.Id = messageId
			message.Sender = "admin"
			message.Receiver = sender
			message.SendTime = time.Now().Format("2006-01-02 15:04:05")
			message.Body = "Read:" + receiver
			
			if conns > 0 {
				conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/cd")
				fatalErrorLog(err, "amqp.Dial")
				defer conn.Close()

				ch, err := conn.Channel()
				fatalErrorLog(err, "conn.Channel")
				defer ch.Close()
			
				content, _ := json.Marshal(message)
			
				err = ch.Publish(
						"",     // exchange
						sender, // routing key
						false,  // mandatory
						false,  // immediate
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType: "apllication/json",
							Body:         content,
					})
				
				fatalErrorLog(err, "ch.Publish")	
				
			} else {
				addMessage(message)
			}
			w.WriteHeader(http.StatusOK)
		}			
		
		if unread <= 1 {
			_, err = c.Do("DEL", messageId)
			fatalErrorLog(err, "redis.DEL")
		} else { 
			_, err = c.Do("HSET", messageId, "unread", unread - 1)
			fatalErrorLog(err, "redis.HSET")
			_, err = c.Do("HSET", messageId, receiver, 0)
			fatalErrorLog(err, "redis.HSET")
		}
		
	} else {
		w.WriteHeader(http.StatusNotFound)
	} 
}

func basicSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	messageId := jsonReq["MessageId"]
	sender := jsonReq["Sender"]
	receiver := jsonReq["Receiver"]
	body := jsonReq["Body"]
	sendtime := time.Now().Format("2006-01-02 15:04:05")
	
	var sendlog SendLog
	sendlog.Action = "单播"
	sendlog.MessageId = messageId
	sendlog.Sender = sender
	sendlog.Receiver = receiver
	sendlog.Body = body
	sendlog.Time = sendtime
	
	var msg string
	
	if messageCheck(body) {
		sendlog.Result = "发送内容审核未通过！"
		sendLogInfo(sendlog)
		
		w.WriteHeader(http.StatusNotFound)
		msg = "3"
	} else {
		conns := IsHasConnection(receiver)
		if conns < 0 {
			sendlog.Result = "不存在该接收用户！"
			sendLogInfo(sendlog)
			
			w.WriteHeader(http.StatusNotFound)
			msg = "2"
		} else {
			var message Message
			message.Id = messageId
			message.Sender = sender
			message.Receiver = receiver
			message.SendTime = sendtime
			message.Body = body
			
			if conns > 0 {
				conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/cd")
				fatalErrorLog(err, "amqp.Dial")
				defer conn.Close()

				ch, err := conn.Channel()
				fatalErrorLog(err, "conn.Channel")
				defer ch.Close()
			
				content, _ := json.Marshal(message)
				
				err = ch.Publish(
						"",     // exchange
						receiver, // routing key
						false,  // mandatory
						false,  // immediate
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType: "apllication/json",
							Body:         content,
				})
				
				fatalErrorLog(err, "ch.Publish")
				
				msg = "0"
			} else {
				addMessage(message)
				msg = "1"
			}			
			
			w.WriteHeader(http.StatusOK)
			
			c, err := redis.Dial("tcp", "192.168.200.3:6379")
			fatalErrorLog(err, "redis.Dial")
			defer c.Close()
	
			_, err = c.Do("HSET", messageId, "unread", 1)
			fatalErrorLog(err, "redis.HSET")
			
			_, err = c.Do("HSET", messageId, "sender", sender)
			fatalErrorLog(err, "redis.HSET")
			
			sendLogInfo(sendlog)
		}
		w.Write([]byte(msg))
	}
}

func multiSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	messageId := jsonReq["MessageId"]
	sender := jsonReq["Sender"]
	receiverList := strings.Split(jsonReq["Receiver"], ",")
	body := jsonReq["Body"]
	sendtime := time.Now().Format("2006-01-02 15:04:05")
	
	var sendlog SendLog
	sendlog.Action = "多播"
	sendlog.MessageId = messageId
	sendlog.Sender = sender
	sendlog.Receiver = jsonReq["receiver"]
	sendlog.Body = body
	sendlog.Time = sendtime
	
	msg := ""
	
	if messageCheck(body) {
		sendlog.Result = "发送内容审核未通过！"
		sendLogInfo(sendlog)
		
		w.WriteHeader(http.StatusNotFound)
		msg = "3"
	} else {
		conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/cd")
		fatalErrorLog(err, "amqp.Dial")
		defer conn.Close()
	
		ch, err := conn.Channel()
		fatalErrorLog(err, "conn.Channel")
		defer ch.Close()
		
		isSend := true  
		
		var message Message
		message.Sender = sender
		message.SendTime = sendtime
		message.Body = body
		message.Id = messageId
		
		var sendList []string
				
		for _ , receiver := range receiverList {
			conns := IsHasConnection(receiver)
			message.Receiver = receiver
			sendlog.Receiver = receiver
			
			if conns < 0 {
				isSend = false
				msg += "2"
				sendlog.Result = "不存在该接收用户！"
				sendLogInfo(sendlog)
			} else {
				sendlog.Result = message.Id
							
				if conns > 0 {
					content, _ := json.Marshal(message)	
						
					
					err = ch.Publish(
							"",     // exchange
							receiver, // routing key
							false,  // mandatory
							false,  // immediate
							amqp.Publishing{
								DeliveryMode: amqp.Persistent,
								ContentType: "application/json",
								Body:        content,
					})
					
					fatalErrorLog(err, "ch.Publish")				
					
					msg += "0"
				} else {
					addMessage(message)
					msg += "1"
				}
				
				sendList = append(sendList, receiver)
				sendLogInfo(sendlog)
			}
		}
		
		if isSend {
			w.WriteHeader(http.StatusOK)
		} else {			
			w.WriteHeader(http.StatusNotFound)
		}
		
		if len(sendList) > 0 {
			c, err := redis.Dial("tcp", "192.168.200.3:6379")
			fatalErrorLog(err, "redis.Dial")
			defer c.Close()
	
			_, err = c.Do("HSET", messageId, "unread", len(sendList))
			fatalErrorLog(err, "redis.HSET")
			
			_, err = c.Do("HSET", messageId, "sender", sender)
			fatalErrorLog(err, "redis.HSET")
			
			for _, member := range sendList {
				_, err = c.Do("HSET", messageId, member, 1)
				fatalErrorLog(err, "redis.HSET")
			}
		}			
	}
	
	w.Write([]byte(msg))
}

func basicMqttSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	sender := jsonReq["Sender"]
	receiver := jsonReq["Receiver"]
	body := jsonReq["Body"]
	sendtime := time.Now().Format("2006-01-02 15:04:05")
	
	var sendlog MqttSendLog
	sendlog.Action = "mqtt单播"
	sendlog.Sender = sender
	sendlog.Receiver = receiver
	sendlog.Body = body
	sendlog.Time = sendtime
	
	var msg string
	
	if messageCheck(body) {
		sendlog.Result = "发送内容审核未通过！"
		mqttSendLogInfo(sendlog)
		
		w.WriteHeader(http.StatusNotFound)
		msg = "3"
	} else { 
		userDevice := strings.Split(receiver, "/")
		
		var username, device string
		
		if len(userDevice) == 2	{
			username = userDevice[0]
			device = userDevice[1]
		} 
		
		res := mqttIsHasConnection(username)
		
		if res > 0 {		
			if mqttHasDevice(username, device) {			
				content := sender + "," + sendtime + "," + body
		    
				opts := mqtt.NewClientOptions().AddBroker("tcp://dev.corp.wingoht.com:1883").SetClientID(sender + "_publisher") 
				
				c := mqtt.NewClient(opts)
				if token := c.Connect(); token.Wait() && token.Error() != nil {  
				    panic(token.Error())  
				}  
					
				token := c.Publish(receiver, 0, false, content)  
				token.Wait()  	
			
				c.Disconnect(250) 
				mqttSendLogInfo(sendlog)
			
				w.WriteHeader(http.StatusOK)
				msg = "0"
			} else {
				w.WriteHeader(http.StatusNotFound)
				msg = "4"
				sendlog.Result = "接受用户未注册该设备！"
				mqttSendLogInfo(sendlog)			
			}
		} else if res == 0 {
			var message MqttMessage
			message.Sender = sender
			message.Receiver = username
			message.Device = device
			message.SendTime = sendtime
			message.Body = body
			mqttAddMessage(message)
			
			mqttSendLogInfo(sendlog)
			
			w.WriteHeader(http.StatusOK)
			msg = "1"
		} else {
			sendlog.Result = "不存在该接收用户！"
			mqttSendLogInfo(sendlog)
			w.WriteHeader(http.StatusNotFound)
			msg = "2"
		}
	}
	w.Write([]byte(msg))
}

func multiMqttSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	
	var jsonReq map[string]string
	jsonBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(jsonBytes, &jsonReq)
	sender := jsonReq["Sender"]
	receiverList := strings.Split(jsonReq["Receiver"], ",")
	body := jsonReq["Body"]
	sendtime := time.Now().Format("2006-01-02 15:04:05")
	
	var sendlog MqttSendLog
	sendlog.Action = "mqtt多播"
	sendlog.Sender = sender
	sendlog.Receiver = jsonReq["receiver"]
	sendlog.Body = body
	sendlog.Time = sendtime
	
	msg := ""
	
	if messageCheck(body) {
		sendlog.Result = "发送内容审核未通过！"
		mqttSendLogInfo(sendlog)
		
		w.WriteHeader(http.StatusNotFound)
		msg = "3"
	} else {
		content := sender + "," + sendtime + "," + body
		isSend := true  
		var message MqttMessage
		message.Sender = sender
		message.SendTime = sendtime
		message.Body = body
		
		for _ , receiver := range receiverList {
			userDevice := strings.Split(receiver, "/")
			
			var username, device string
			
			if len(userDevice) == 2	{
				username = userDevice[0]
				device = userDevice[1]
			} 
			
			res := mqttIsHasConnection(username)
			
			if res > 0 {
				if mqttHasDevice(username, device) {
					opts := mqtt.NewClientOptions().AddBroker("tcp://dev.corp.wingoht.com:1883").SetClientID(sender + "_publisher") 
		
					c := mqtt.NewClient(opts)
					if token := c.Connect(); token.Wait() && token.Error() != nil {  
				        panic(token.Error())  
				    }  
				
				    token := c.Publish(receiver, 0, false, content)  
				    token.Wait()  	
					
					c.Disconnect(250)
					
					msg += "0"
				} else {
					isSend = false
					msg += "4"	
				}
			} else if res == 0 {
				message.Receiver = username
				message.Device = device
				mqttAddMessage(message)
				msg += "1"
			} else {
				isSend = false
				msg += "2"
			}
		}
		
		if isSend {
			mqttSendLogInfo(sendlog)
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(msg))
		} else {
			sendlog.Result = "存在未注册用户设备！"
			mqttSendLogInfo(sendlog)
			
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(msg))
		}		
	}
}