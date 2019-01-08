package main

import (
	"log"
	"fmt"
        s "strings"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CCSDB struct {
	Server   string
	Database string
	Username string
	Password string
}

var db *mgo.Database
var dbprimary *mgo.Database

const (
	COLLECTION = "ccscatalog"
)

const (
	MongoDBHosts = "internal-LoadBalancer-GP-2070204421.us-west-1.elb.amazonaws.com:27017"
	AuthDatabase = "admin"
	AuthUserName = "admin"
	AuthPassword = "cmpe281"
	MongoServer1 = "10.0.1.123"
	MongoServer2 = "10.0.1.115"
	MongoServer3 = "10.0.1.54"
	MongoServer4 = "10.0.3.168"
	MongoServer5 = "10.0.3.254"
)


var mgoSession   *mgo.Session
//Initiate connection to database
func (m *CCSDB) Connect() {
	fmt.Printf("In Connect1")
	mongoDBDialInfo := &mgo.DialInfo{
        Addrs:    []string{MongoDBHosts},
        Database: AuthDatabase,
        Username: AuthUserName,
        Password: AuthPassword,
    }
    session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
        fmt.Printf("In Connect Error1")
		log.Fatal(err)
    }
	db = session.DB("testing")
}

func GetMongoSession() *mgo.Session {
        mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoDBHosts},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        if mgoSession == nil {
            var err error
            mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
            if err != nil {
                log.Fatal("Failed to start the Mongo session")
            }
        }
        return mgoSession.Clone()
 }

// Make sure the write happens to Master of the Mongo Database
func (m *CCSDB) ConnecttoPrimary() {
	result := make(map[string]interface{})
	fmt.Printf("Came Here")
	err := db.Run(bson.M{"isMaster": 1}, result)
	 if err != nil {
                log.Fatal(err)
        }
	value := result["primary"].(string)
	if s.Contains(value,"Primary"){
		fmt.Printf("In Primary")
            mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoServer1},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        sessionone, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                log.Fatal(err)
        }
        dbprimary = sessionone.DB("testing")
	}
	if s.Contains(value,"Secondary1"){
                fmt.Printf("In Secondary1")
            mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoServer2},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        sessionone, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                log.Fatal(err)
        }
        dbprimary = sessionone.DB("testing")
        }
	if s.Contains(value,"Secondary2"){
                fmt.Printf("In Secondary2")
		  mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoServer3},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        sessionone, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                log.Fatal(err)
        }
        dbprimary = sessionone.DB("testing")
        }
	if s.Contains(value,"Secondary3"){
                fmt.Printf("In Secondary3")
		  mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoServer4},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        sessionone, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                log.Fatal(err)
        }
        dbprimary = sessionone.DB("testing")
        }
	if s.Contains(value,"Secondary4"){
                fmt.Printf("In Secondary4")
		  mongoDBDialInfo := &mgo.DialInfo{
                Addrs:    []string{MongoServer5},
                Database: AuthDatabase,
                Username: AuthUserName,
                Password: AuthPassword,
        }
        sessionone, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                log.Fatal(err)
        }
        dbprimary = sessionone.DB("testing")
        }
}

// Find list of All Drinks in the Catalog
func (m *CCSDB) FindAll() ([]Order, error) {
        session := GetMongoSession()
	db = session.DB("testing")
	var orders []Order
        err := db.C(COLLECTION).Find(bson.M{}).All(&orders)
        defer session.Close()
	return orders, err
}

// Find a Drink by its id
func (m *CCSDB) FindById(id string) (Order, error) {
	session := GetMongoSession()
	db = session.DB("testing")
	var order Order
	fmt.Printf("finding name = ",id)
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&order)
	fmt.Println("%v",order)
	defer session.Close()
	return order, err
}

func (m *CCSDB) FindByUserId(id1 string) (Order, error) {
	session := GetMongoSession()
	db = session.DB("testing")
        var order Order
        fmt.Printf(id1)
        err := db.C(COLLECTION).Find(bson.M{ "name" : id1}).One(&order)
	fmt.Printf("%v",order)
	defer session.Close()
	return order, err
}eturn order, err
}

func (m *CCSDB) Insert(order Order) error {
        session := GetMongoSession()
	db = session.DB("testing")
        err := dbprimary.C(COLLECTION).Insert(&order)
        defer session.Close()
	return err
}

func (m *CCSDB) Delete(order Order) error {
	session := GetMongoSession()
        db = session.DB("testing")
	err := db.C(COLLECTION).Remove(&order)
	defer session.Close()
	return err
}

