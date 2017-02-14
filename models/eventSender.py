import requests, time, random
a = time.time()
while 1:
   event = int( raw_input( "Enter type of event: " ))
   if event == 0:
      break
   if event == 1:
      stationId = int( raw_input( "Enter the station ID: " ))
      data = { "stationId" : stationId , "eventId" : 1 }
      r = requests.post( "http://localhost:9999/handleStationEvent", data )
      print r.status_code, r.reason
   elif event == 2:
      name = raw_input( "Enter Name to add to queue: " )
      data = { "licensePlate" : "6XDS849", "name" : name }
      r = requests.post( "http://localhost:9999/", data )
      print r.status_code, r.reason
   else:
      print "Please enter a valid event"
      
   
   
   
