import requests, time, random
from pymongo import MongoClient

class RequestSender:
   def __init__( self ):
      self.client = MongoClient( "localhost" )
      self.stationTable = self.client.test.stations
      self.queue = self.client.test.queue
      self.chat = self.client.test.chat
      self.user = self.client.test.user

   def addToQueue( self, user, status="Active" ):
      print status
      if status == "Active":
         data = { "licensePlate" : "6XDS849", "name" : user, "status" : status, "time": 0 }
      elif status == "Suspended":
         timeStamp = self.stationTable.find_one( { "user.name" : user } )[ 'user' ][ 'tstamp' ]
         data = { "licensePlate" : "6XDS849", "name" : user, "status" : status, "time": timeStamp }
      r = requests.post( "http://localhost:9999/", data )
      print r.status_code, r.reason
      
   def freeStation( self, station ):
      data = { "stationId" : station , "eventId" : 1 }
      r = requests.post( "http://localhost:9999/handleStationEvent", data )
      print r.status_code, r.reason

   def findStationForUser( self, user ):
      record = { "user.name" : user }
      if self.stationTable.find( record ).count() > 0:
         return self.stationTable.find_one( record )[ 'id' ]
      else:
         return 0

   def allocateFreeStation( self ):
      data = {}
      r = requests.post( "http://localhost:9999/findFreeStation", data )
      print r.status_code, r.reason

   def findUserInQueue( self, user ):
      record = { "name" : user }
      return self.queue.find( record ).count()

   def removeUser( self, user ):
      record = { "name" : user }
      return self.queue.remove( record )

   def updateStatusForUserInQueue( self, user, status ):
      record = { "name" : user }
      result = self.queue.update_one( record, { "$set": { "status" : status } } )
      return result.matched_count

   def findUsername( self, user ):
      record = { "username" : user }
      if self.chat.find( record ).count() > 0:
         return self.chat.find_one( record )[ 'name' ]
      else:
         return 0

   def addUserName( self, user, licensePlate ):
      record = { "username" : user }
      if self.chat.find( record ).count():
         return 1
      record = { "licensePlate" : licensePlate  }
      name = self.user.find_one( record )[ 'name' ]
      record = { "username" : user, "name" : name }
      self.chat.insert( record )
      return 0

if __name__ == "__main__":
   requestSender = RequestSender()
   while 1:
      event = int( raw_input( "Enter type of event: " ))
      if event == 0:
         break
      if event == 1:
         stationId = int( raw_input( "Enter the station ID: " ))
         requestSender.freeStation( stationId )
      elif event == 2:
         name = raw_input( "Enter Name to add to queue: " )
         requestSender.addToQueue( name, status="Active" )
      else:
         print "Please enter a valid event"
   
   
