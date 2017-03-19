import xmpp, time, re
from eventSender import RequestSender

def handleAddCommand( client, user, text=None ):
    print "Recieved Add Command"
    requestSender.addToQueue( str( user ) )
    client.sendMessage( user, "Successfully added to queue" )

def handleTimeCommand( client, user, text=None ):
    print "Recieved Time Command"
    client.sendMessage( user, time.asctime() )

def handleSuspendCommand( client, user, text=None ):
    print "Recieved Suspend Command"
    client.sendMessage( user, "Suspended, retaining position in queue" )

def handleDeleteCommand( client, user, text=None ):
    print "Recieved Delete Command"
    client.sendMessage( user, "Deleted from queue" )

def handleStationFreeCommand( client, user, text=None ):
    match = re.match( "free station (\d+)", text )
    requestSender.freeStation( int( match.group(1) ) )
    client.sendMessage( user, "Station is free" )

commands = {
    # A dictionary based on Regexs and methods
    "ev add" : handleAddCommand,
    "time" : handleTimeCommand,
    "ev suspend" : handleSuspendCommand,
    "ev delete" : handleDeleteCommand,
    "free station (\d+)" : handleStationFreeCommand,

}
class AdminBot:
    # Client for XMPP
    client = None 

    def __init__( self, username, password, host='talk.google.com', port=5223 ):  
        self.client = xmpp.Client('gmail.com', debug=[])
        self.client.connect(server=( host , port ))
        self.client.auth( username, password, 'Charge Bot')
        self.registerHandlers()
        self.client.sendInitPresence()

    def messageHandler( self, client, message ):
        text = message.getBody()
        user = message.getFrom()
        if text:
            text = text.encode( "utf-8", 'ignore' )
            for ( command, function ) in commands.items():
                match = re.match( command, text )
                if match:
                    function( self, user, text )
                    break
            else:
                self.sendMessage( user, "Invalid command, please check with Admin" )
                    
    def presenceHandler( self, client, presence ):
        pass

    def registerHandlers( self ):
        self.client.RegisterHandler( "message", self.messageHandler )
        self.client.RegisterHandler( "presence", self.presenceHandler )

    def sendMessage( self, to, msg ):
        message = xmpp.Message(to, msg)
        message.setAttr('type', 'chat')
        self.client.send(message)

    def process( self ):
        try:
            self.client.Process( 1 )
        except KeyboardInterrupt:
            return 0
        return 1

    def start( self ):
        while self.process():
            pass


if __name__ == "__main__":
    bot = AdminBot( "evchargepoint", "" )
    requestSender = RequestSender()
    bot.start()


