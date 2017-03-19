import time, random
from pymongo import MongoClient

client = MongoClient( "localhost" )
station = client.test.stations
queue = client.test.queue

class Stations:
   nofityQueue = None

   def printStations( self ):
      for s in station.find() :
         print s[ "id" ], s[ "user" ]

   def assignStation( self, stationId, name ):
      print "Assigning Station %d to %s " %( stationId, name )
      time.sleep( 2 )
      station.update( { "id" : stationId }, { "$set" : { "user" : name } } )

   def addToNotify( self , queue ):
      self.notifyQueue = queue

   def freeupStation( self, stationId ):
      print "Freeing up station: %d" % stationId
      station.update( { "id" : stationId }, { "$set" : { "user" : "" } } )
      time.sleep( 3 )
      if not self.notifyQueue.isEmpty():
         a = self.notifyQueue.getName()
         self.assignStation( stationId, a )

   def getFreeStations( self ):
      freeStations = []
      for item in station.find({ "user" : "" }):
         freeStations.append( item[ "id" ] )
      return freeStations



class Queue:
   notifyStations = None

   def isEmpty( self ):
      return len( queue.find_one( { "name": "evqueue" } )[ "users" ] ) == 0

   def getName( self ):
      if not self.isEmpty():
         name = queue.find_one( { "name": "evqueue" } )[ "users" ][0]
         queue.update( { "name": "evqueue" } , { "$pop" : { "users" : -1 } } )
         return name
      else:
         return ""

   def addName( self, name ):
      free = self.notifyStations.getFreeStations()
      if free:
         stationId = random.choice( free )
         self.notifyStations.assignStation( stationId, name )
      else:
         queue.update( {"name" :"evqueue"}, { "$push" : { "users" : name } } )

   def display( self ):
      for i in queue.find( { "name": "evqueue" } ):
         print i[ "users" ]

   def addToNotify( self, station ):
      self.notifyStations = station

a = Stations()
b = Queue()
a.addToNotify( b )
b.addToNotify( a )

while 1:
   c = int( raw_input( "Enter option : " ) )
   if c == 0:
      break
   elif c == 99:
      name = raw_input( "Enter the name : " )
      b.addName( name )
   elif c == 98:
      s = int( raw_input( "Enter Station: " ) )
      a.freeupStation( s )
   else:
      a.printStations()
      b.display()
