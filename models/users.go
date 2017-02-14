package users

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
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
    User string
}

type Queue struct{
    Name string
    Users []string 
}
type Dashboard struct {
    Stations []Common
    Queue []string
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
    queue := Queue{}
    collection.Find(bson.M{"name": "evqueue" }).One(&queue)
    //fmt.Printf( "%s" , queue.Users )
    results.Queue = queue.Users 
    return results, 1
}

func AddToQueue( name string ){
    sStation, cStations := openDB( "stations" )
    defer sStation.Close()
    station := Common{}
    cStations.Find(bson.M{"user": "" }).One(&station)
    if station.Id != 0 {
        fmt.Println( "Free station, adding the user: ", name )
        cStations.Update( bson.M{ "id" : station.Id }, bson.M{ "$set" : bson.M{ "user" : name }} )
    } else {
        fmt.Println( "No Free station, Adding  to Queue: ", name )
        sQueue, cQueue := openDB( "queue" )
        defer sQueue.Close()
        cQueue.Update(bson.M{"name": "evqueue" }, bson.M{ "$push" : bson.M{ "users" : name  }})
    }
}

func HandleEventStation( stationId int , eventId int ){
    sStation, cStations := openDB( "stations" )
    defer sStation.Close()
    sQueue, cQueue := openDB( "queue" )
    defer sQueue.Close()
    fmt.Println( "Freeing Station Id: ", stationId )
    cStations.Update(bson.M{ "id" : stationId }, bson.M{ "$set" : bson.M{ "user" : "" }} )
    queue := Queue{}
    cQueue.Find(bson.M{"name": "evqueue" }).One(&queue)
    if len( queue.Users ) != 0 {
        cQueue.Update(bson.M{"name": "evqueue" }, bson.M{ "$pop" : bson.M{ "users" : -1 }})
        fmt.Println( "Free station, adding the user:", queue.Users[0] )
        cStations.Update( bson.M{ "id" : stationId }, bson.M{ "$set" : bson.M{ "user" : queue.Users[0] }} )
    }
}
