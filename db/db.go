package db

import (
	"os"

	"github.com/prazd/nodes_mon_bot/db/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	host       = os.Getenv("HOST")
	database   = os.Getenv("DB")
	username   = os.Getenv("USER")
	password   = os.Getenv("PASS")
	collection = os.Getenv("COLLECTION")
)

var info = mgo.DialInfo{
	Addrs:    []string{host},
	Database: database,
	Username: username,
	Password: password,
}

func IsInDb(id int) (bool, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return false, err
	}
	defer session.Close()

	var user schema.User

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"telegram_id": id}).One(&user)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func CreateUser(id int) error {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database).C(collection)

	err = c.Insert(&schema.User{Telegram_id: id, Subscription: false})
	if err != nil {
		return err
	}

	return nil
}

func SubscribeOrUnSubscribe(id int, subscription bool) error {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database).C(collection)

	err = c.Update(bson.M{"telegram_id": id}, bson.M{"$set": bson.M{"subscription": subscription}})
	if err != nil {
		return err
	}

	return nil
}

func GetSubStatus(id int) (string, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var user schema.User

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"telegram_id": id}).One(&user)
	if err != nil {
		return "", err
	}

	message := "Subscription"

	switch user.Subscription {
	case true:
		message += ": ✔"
	case false:
		message += ": ✖"
	}

	return message, nil
}

func GetAllSubscribers() []int {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil
	}
	defer session.Close()

	var users []schema.User
	var usersId []int

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"subscription": true}).All(&users)
	if err != nil {
		return nil
	}

	for i := 0; i < len(users); i++ {
		usersId = append(usersId, users[i].Telegram_id)
	}

	return usersId
}
