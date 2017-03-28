package users

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"time"
    	"strings"
)

func openDB( coll string )( m *mgo.Session, c *mgo.Collection ){
    sess, err := mgo.Dial("mongodb://localhost")
    if err != nil {
            fmt.Printf("Can't connect to mongo, go error %v\n", err)
            os.Exit(1)
    }
    sess.SetSafe(&mgo.Safe{})
    collection := sess.DB("test").C( coll )
    return sess, collection
}
type Person struct {
        LicensePlate string `bson:"licensePlate"`
        Email string
        Name string
}

type Common struct{
    Id int
    User QueueObj
}

type QueueObj struct{
    Name string
    Tstamp int64 `bson:"tstamp"`
    Status string `bson:"status"`
}

type Dashboard struct {
    Stations []Common
    Queue []QueueObj
}

func GetUser( r string ) ( s Person ){
    sess, collection := openDB( "user")
    defer sess.Close()
    result := Person{}
    collection.Find(bson.M{"licensePlate": strings.ToUpper( r ) }).One(&result)
    return result
}

func InsertUser( licensePlate string, name string , email string )( e error, code int ){
    sess, collection := openDB( "user" )
    defer sess.Close()
    //check for duplicates
    result := Person{}
    collection.Find(bson.M{"licensePlate": strings.ToUpper( licensePlate ) }).One(&result)
    if result.Email != "" {
        return nil, -1
    }
    record := &Person{ LicensePlate: strings.ToUpper( licensePlate ), Email: email, Name: name }
    error := collection.Insert(record)
    return error, 0
}

func GetAllUsers() ( r []Person ){
    sess, collection := openDB( "user" )
    defer sess.Close()
    results := []Person{}
    collection.Find(bson.M{"_id": bson.M{"$exists": 1}}).All(&results)
    return results
}

func GetDashBoard() ( d Dashboard , r int ){
    sess, collection := openDB( "stations" )
    defer sess.Close()
    results := Dashboard{}
    collection.Find(bson.M{"_id": bson.M{"$exists": 1}}).All(&results.Stations )
    sess, collection = openDB( "queue" )
    defer sess.Close()
    collection.Find(bson.M{"_id": bson.M{"$exists": 1}}).Sort( "tstamp" ).All(&results.Queue )
    //fmt.Printf( "%s" , queue.Users )
    return results, 1
}

func FindFreeStation(){
    sStation, cStations := openDB( "stations" )
    defer sStation.Close()
    station := Common{}
    cStations.Find(bson.M{"user": "" }).One(&station)
    if station.Id != 0 {
        fmt.Println( "Free station, calling station free event" )
	// raising event 1 for now, change it later on
	HandleEventStation( station.Id, 1 )
    }
}

func AddToQueue( name string, status string, timeStamp int64 ){
    fmt.Println( "Adding  to Queue: " + name + " Status: " + status )
    sQueue, cQueue := openDB( "queue" )
    defer sQueue.Close()
    record := &QueueObj{ Name : name , Tstamp : time.Now().UnixNano() , Status : status }
    if status == "Suspended" {
       record.Tstamp = timeStamp
    }
    cQueue.Insert(record)
    // calling free stations
    FindFreeStation()
}

func HandleEventStation( stationId int , eventId int ){
    sStation, cStations := openDB( "stations" )
    defer sStation.Close()
    sQueue, cQueue := openDB( "queue" )
    defer sQueue.Close()
    fmt.Println( "Freeing Station Id: ", stationId )
    cStations.Update(bson.M{ "id" : stationId }, bson.M{ "$set" : bson.M{ "user" : "" }} )
    queue := []QueueObj{}
    cQueue.Find(bson.M{"status": "Active" }).Sort( "tstamp" ).All(&queue )
    if len( queue ) != 0 {
        fmt.Println( "Free station, adding the user:", queue[0] )
        cStations.Update( bson.M{ "id" : stationId }, bson.M{ "$set" : bson.M{ "user" : queue[0] }} )
	cQueue.Remove( bson.M{ "name" : queue[0].Name } )
    }
}
