package main

import (
	"log"
	"time"
	"net/http"
	"github.com/streadway/amqp"
	"github.com/eclipse/paho.mqtt.golang"
)

func main() {	
	conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/")
    fatalErrorLog(err, "amqp.Dial")
    defer conn.Close()

    ch, err := conn.Channel()
    fatalErrorLog(err, "conn.Channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
	                "queue.event",    // name
	                true, 				// durable
	                false,				// delete when usused
	                false,  			// exclusive
	                false, 				// no-wait
	                nil,   				// arguments
                )
    
    fatalErrorLog(err, "ch.QueueDeclare")

	err = ch.QueueBind(
                    q.Name,       			// queue name
                    "queue.*",    	// routing key
                    "amq.rabbitmq.event", 	// exchange
                    false,
                    nil)
    
    fatalErrorLog(err, "ch.QueueBind")

	msgs, err := ch.Consume(
	                q.Name, // queue
	                "",     // consumer
	                true,   // auto ack
	                false,  // exclusive
	                false,  // no local
	                false,  // no wait
	                nil,    // args
		        )
        
    fatalErrorLog(err, "Failed to register a consumer")

    go func() {
	    for d := range msgs {
		    routingKey := d.RoutingKey
		    queueName := d.Headers["name"].(string)
		    timestamp := d.Timestamp.Format("2006-01-02 15:04:05")
		    vhost := d.Headers["vhost"].(string)
		    
		    var userlog UserLog
			
			userlog.Info = queueName
			userlog.Time = timestamp
		    
		    if vhost == "cd" {
			    if routingKey == "queue.deleted" {
			    	if hasUser(queueName) {	
				    	updateUser(queueName, false)
				    	userlog.Action = "队列删除成功"
			    	} else {
				    	userlog.Action = "非法队列删除"
			    	}
			    } else if routingKey == "queue.created"{
			    	if hasUser(queueName) {	
					    updateUser(queueName, true)
					    userlog.Action = "队列创建成功"
			    	} else {
			    		conn, err := amqp.Dial("amqp://ishowfun:123456@dev.corp.wingoht.com:5672/cd")
					    fatalErrorLog(err, "amqp.Dial")
					    defer conn.Close()
	
					    ch, err := conn.Channel()
					    fatalErrorLog(err, "conn.Channel")
					    defer ch.Close()
	
						_, err = ch.QueueDelete(queueName, false, false, false)
									 
						fatalErrorLog(err, "ch.QueueDeclare")
						
				    	userlog.Action = "非法队列创建"
			    	}			    
			    }
			    
			    userLogInfo(userlog)
		    }
	    }
    }()

	opts := mqtt.NewClientOptions().AddBroker("tcp://dev.corp.wingoht.com:1883").SetClientID("will_topic_listener") 
  
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
	    panic(token.Error())
	}
	
	msgRcvd := func(client mqtt.Client, message mqtt.Message) {
	    username := string(message.Payload())
	    
	    var userlog UserLog
		userlog.Action = "mqtt用户异常登出"
		userlog.Info = username
		userlog.Time = time.Now().Format("2006-01-02 15:04:05")
	
		mqttUpdateUser(username, false)
	
		mqttUserLogInfo(userlog)
	}
  
	if token := c.Subscribe("will/topic", 0, msgRcvd); token.Wait() && token.Error() != nil {
	    log.Fatal(token.Error())
	}
	
	http.HandleFunc("/user/check", checkToken)
	http.HandleFunc("/user/unread", getUnreadMessage)
	
	http.HandleFunc("/user/mqtt/login", userLogin)
	http.HandleFunc("/user/mqtt/logout", userLogout)
	http.HandleFunc("/user/mqtt/device/add", addDevice)
	http.HandleFunc("/user/mqtt/device/remove", removeDevice)
	http.HandleFunc("/user/mqtt/device/show", showDevice)
	http.HandleFunc("/user/mqtt/unread", getUnreadMqttMessage)
	
	http.HandleFunc("/send/basicsend", basicSend)
	http.HandleFunc("/send/multisend", multiSend)
	http.HandleFunc("/send/hasread", hasRead)
	
	http.HandleFunc("/send/mqtt/basicsend", basicMqttSend)
	http.HandleFunc("/send/mqtt/multisend", multiMqttSend)
	
	err = http.ListenAndServe("192.168.40.47:9090", nil)
	fatalErrorLog(err, "http.ListenAndServe")
}

func fatalErrorLog(err error, f string) {
	if err != nil {
		log.Fatal(f + ": " + err.Error())
	}
}