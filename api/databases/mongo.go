package dbs

import (
	"log"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type MgoSession struct {
	Session *mgo.Session
}

func newMgoSession(s *mgo.Session) *MgoSession {
	return &MgoSession{s}
}

func StartMongoDB(msg string) *MgoSession {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{os.Getenv("MONGODB_URL")},
		Timeout: 60 * time.Second,
		//Database: "",
		//Username: "",
		//Password: "",
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("[MongoDB] CreateSession: %s\n", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	log.Printf("[MongoDB] connected! %s", msg)
	return newMgoSession(mongoSession)
}
