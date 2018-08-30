class Notify():
    def __init__(self):
        self.title = None
        self.content = None
        self.payload = None

    def setTitle(self, title):
        self.title = title
    def getTitle(self):
        return self.title
    def setContent(self, content):
        self.content = content
    def getContent(self):
        return self.content
    def setPayload(self, payload):
        self.payload = payload
    def getPayload(self):
        return self.payload;