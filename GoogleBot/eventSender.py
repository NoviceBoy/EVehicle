import requests, time, random
from pymongo import MongoClient

class RequestSender:
   def __init__( self ):
      self.client = MongoClient( "localhost" )
      self.stationTable = self.client.test.stations
      self.queue = self.client.test.queue

   def addToQueue( self, user ):
      data = { "licensePlate" : "6XDS849", "name" : user }
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

   def findUserInQueue( self, user ):
      record = { "name" : user }
      return self.queue.find( record ).count()

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
         requestSender.addToQueue( name )
      else:
         print "Please enter a valid event"
   
   
