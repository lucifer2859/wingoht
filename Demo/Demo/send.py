#coding=utf-8

import json

from django.http import HttpResponse

from umessage.pushclient import PushClient
from umessage.iospush import *
from umessage.androidpush import *

from umessage.errorcodes import UMPushError, APIServerErrorCode

from igetui.igt_message import *
from igetui.igt_target import *
from igt_push import IGeTui
from igetui.template.igt_notification_template import *

#友盟
umengAppKey = '5b7e7ce98f4a9d109200004b'
umengAppMasterSecret = 't23anusbo9ofyqcjmamcoj9nrhwqmcct'
#deviceToken = 'AmnAtj3pt3WvR7dRYV9C01jxn1ktUlb7txA8tNlG-GvM'
#groupFilter = r'{"where":{"and":[{"or":[{"app_version":">1.0"}]}]}}'

#个推
getuiAppKey = "Ytou67dSem7jrUAedUVfq6"
getuiAppId = "2ESG1htby99WIDIE7id4H2"
getuiMasterSecret = "ZcFYtFR9k38TsiDTieCu37"
getuiHost = 'http://sdk.open.api.igexin.com/apiex.htm'
#CID = "1f5c167996c6c7a7fc46380d83f935d2"

# 友盟
def sendUmengUnicast(deviceToken, ticker, title, text):
    unicast = AndroidUnicast(umengAppKey, umengAppMasterSecret)
    unicast.setDeviceToken(deviceToken)
    unicast.setTicker(ticker)
    unicast.setTitle(title)
    unicast.setText(text)
    unicast.goAppAfterOpen()
    unicast.setDisplayType(AndroidNotification.DisplayType.notification)
    unicast.setTestMode()
    pushClient = PushClient()
    pushClient.send(unicast)

def sendUmengBroadcast(ticker, title, text):
    broadcast = AndroidBroadcast(umengAppKey, umengAppMasterSecret)
    broadcast.setTicker(ticker)
    broadcast.setTitle(title)
    broadcast.setText(text)
    broadcast.goAppAfterOpen()
    broadcast.setDisplayType(AndroidNotification.DisplayType.notification)
    broadcast.setTestMode()
    pushClient = PushClient()
    pushClient.send(broadcast)

# def sendUmengGroupcast(groupFilter, ticker, title, text):
#     groupcast = AndroidGroupcast(umengAppKey, umengAppMasterSecret)
#     groupcast.setFilter(groupFilter)
#     groupcast.setTicker(ticker)
#     groupcast.setTitle(title)
#     groupcast.setText(text)
#     groupcast.goAppAfterOpen()
#     groupcast.setDisplayType(AndroidNotification.DisplayType.notification)
#     groupcast.setTestMode()
#     pushClient = PushClient()
#     pushClient.send(groupcast)

def sendUmengListcast(deviceTokenList, ticker, title, text):
    listcast = AndroidListcast(umengAppKey, umengAppMasterSecret)
    listcast.setDeviceToken(deviceTokenList)
    listcast.setTicker(ticker)
    listcast.setTitle(title)
    listcast.setText(text)
    listcast.goAppAfterOpen()
    listcast.setDisplayType(AndroidNotification.DisplayType.notification)
    listcast.setTestMode()
    pushClient = PushClient()
    pushClient.send(listcast)

#个推
def pushGetuiMessageToApp(content, title, text):
    push = IGeTui(getuiHost, getuiAppKey, getuiMasterSecret)

    template = NotificationTemplateDemo(content, title, text)

    message = IGtAppMessage()
    message.data = template
    message.isOffline = True
    message.offlineExpireTime = 1000 * 3600 * 12
    message.appIdList.extend([getuiAppId])

    ret = push.pushMessageToApp(message, 'toApp')
    print ret

def pushGetuiMessageToSingle(cid, content, title, text):
    push = IGeTui(getuiHost, getuiAppKey, getuiMasterSecret)

    template = NotificationTemplateDemo(content, title, text)
    
    message = IGtSingleMessage()
    message.isOffline = True
    message.offlineExpireTime = 1000 * 3600 * 12
    message.data = template
    message.pushNetWorkType = 1
    
    target = Target()
    target.appId = getuiAppId
    target.clientId = cid

    ret = push.pushMessageToSingle(message, target)
    print ret

def pushGetuiMessageToList(cidList, content, title, text):
    push = IGeTui(getuiHost, getuiAppKey, getuiMasterSecret)

    template = NotificationTemplateDemo(content, title, text)

    message = IGtListMessage()
    message.data = template
    message.isOffline = True
    message.offlineExpireTime = 1000 * 3600 * 12
    message.pushNetWorkType = 0

    arr = []

    for cid in cidList.split(','):
        target = Target()
        target.appId = getuiAppId
        target.clientId = cid
        arr.append(target)
    
    contentId = push.getContentId(message, 'ToList')
    ret = push.pushMessageToList(contentId, arr)
    print ret

# 通知透传模板动作内容
def NotificationTemplateDemo(content, title, text):
    template = NotificationTemplate()
    template.appId = getuiAppId
    template.appKey = getuiAppKey
    template.transmissionType = 1
    template.transmissionContent = content
    template.title = title
    template.text = text
    template.logo = ""
    template.logoURL = ""
    template.isRing = True
    template.isVibrate = True
    template.isClearable = True
    return template

def sendUnicast(request):
    apiType = request.POST['ApiType']
    deviceToken = request.POST['DeviceToken']
    ticker = request.POST['Ticker']
    title = request.POST['Title']
    text = request.POST['Text']
    if apiType == "umeng":
        sendUmengUnicast(deviceToken, ticker, title, text)
    elif apiType == "getui":
        pushGetuiMessageToSingle(deviceToken, ticker, title, text)
    else:
        return HttpResponse("Error!")
    return HttpResponse("Unicast!")

def sendBroadcast(request):
    apiType = request.POST['ApiType']
    ticker = request.POST['Ticker']
    title = request.POST['Title']
    text = request.POST['Text']
    if apiType == "umeng":
        sendUmengBroadcast(ticker, title, text)
    elif apiType == "getui":
        pushGetuiMessageToApp(ticker, title, text)
    else:
        return HttpResponse("Error!")
    return HttpResponse("Broadcast!")

def sendListcast(request):
    apiType = request.POST['ApiType']
    deviceToken = request.POST['DeviceToken']
    ticker = request.POST['Ticker']
    title = request.POST['Title']
    text = request.POST['Text']
    if apiType == "umeng":
        sendUmengListcast(deviceToken, ticker, title, text)
    elif apiType == "getui":
        pushGetuiMessageToList(deviceToken, ticker, title, text)
    else:
        return HttpResponse("Error!")
    return HttpResponse("Listcast!")

# def sendGroupcast(request):
#     apiType = request.POST['ApiType']
#     groupFilter = request.POST['groupFilter']
#     ticker = request.POST['Ticker']
#     title = request.POST['Title']
#     text = request.POST['Text']
#     if apiType == "umeng":
#         sendUmengGroupcast(groupFilter, ticker, title, text)
#     elif apiType == "getui":
#         print ""
#     else:
#         print ""
#     return HttpResponse("Groupcast!")