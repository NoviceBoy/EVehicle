import xmpp, time, re
from eventSender import RequestSender

def userName( user ):
    return str( user ).split( "@" )[ 0 ]

def handleAddCommand( client, user, text=None ):
    print "Recieved Add Command :", user
    userId = userName( user ) 
    stationForUser = requestSender.findStationForUser( userId )
    userInQueue = requestSender.findUserInQueue( userId )
    if not stationForUser and not userInQueue :
        requestSender.addToQueue( userId, status="Active"  )
        stationForUser = requestSender.findStationForUser( userId )
        if stationForUser:
            client.sendMessage( user, "Station free, Assigned to use station: %d " % stationForUser )
        else:
            client.sendMessage( user, "Successfully added to Queue" )
    else:
        if userInQueue:
            client.sendMessage( user, "You are already in the queue" )
        else:
            client.sendMessage( user, "You are already using station %d" % stationForUser )

def handleTimeCommand( client, user, text=None ):
    print "Recieved Time Command :", user
    client.sendMessage( user, time.asctime() )

def handleSuspendCommand( client, user, text=None ):
    print "Recieved Suspend Command :", user
    userId = userName( user ) 
    stationForUser = requestSender.findStationForUser( userId )
    userInQueue = requestSender.findUserInQueue( userId )
    # Checking for 2 conditions, user can be using a station/ in queue
    if stationForUser:
        requestSender.addToQueue( userId, status="Suspended"  )
        requestSender.freeStation( int(stationForUser) )
        client.sendMessage( user, "Station %d released, Successfully added back to Queue" % stationForUser )
    elif userInQueue:
        requestSender.updateStatusForUserInQueue( userId, status="Suspended" )
        client.sendMessage( user, "Status changed to suspended")
    else:
        client.sendMessage( user, "User not using any Stations or Queue, Add to Queue first" )

def handleActivateCommand( client, user, text=None ):
    print "Recieved activate Command :", user
    userId = userName( user ) 
    userInQueue = requestSender.findUserInQueue( userId )
    if userInQueue:
        requestSender.updateStatusForUserInQueue( userId, status="Active" )
        client.sendMessage( user, "Status changed to active" )
        requestSender.allocateFreeStation()
    else:
        client.sendMessage( user, "User not in Queue, Add to Queue first" )

def handleStationFreeCommand( client, user, text=None ):
    print "Recieved Free Command :", user
    match = re.match( "ev free", text )
    userId = userName( user ) 
    stationForUser = requestSender.findStationForUser( userId )
    userInQueue = requestSender.findUserInQueue( userId )
    if stationForUser:
        requestSender.freeStation( int(stationForUser) )
        client.sendMessage( user, "Station %d Free-ed up" % stationForUser )
    elif userInQueue:
        requestSender.removeUser( userId )
        client.sendMessage( user, "Deleted user from the queue, use 'ev add' to add again" )        
    else:
        client.sendMessage( user, "User not using any Stations or Queue, Add to Queue first" )

commands = {
    # A dictionary based on Regexs and methods
    "ev add" : handleAddCommand,
    "time" : handleTimeCommand,
    "ev suspend" : handleSuspendCommand,
    "ev activate" : handleActivateCommand,
    "ev free" : handleStationFreeCommand,
}
class AdminBot:
    # Client for XMPP
    client = None 

    def __init__( self, username, password, host='talk.google.com', port=5223 ):  
        self.client = xmpp.Client('gmail.com', debug=[] )
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


