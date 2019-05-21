package db

import (
	"os"

	"github.com/prazd/nodes_mon_bot/db/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

var (
	host       = os.Getenv("HOST")
	database   = os.Getenv("DB")
	username   = os.Getenv("USER")
	password   = os.Getenv("PASS")
	user_collection = os.Getenv("USER_COLL")
	endpoints_collection =  os.Getenv("ENDPOINTS_COLL")
	apis_collection = os.Getenv("API_COLL")
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

	c := session.DB(database).C(user_collection)

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

	c := session.DB(database).C(user_collection)

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

	c := session.DB(database).C(user_collection)

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

	c := session.DB(database).C(user_collection)

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

	c := session.DB(database).C(user_collection)

	err = c.Find(bson.M{"subscription": true}).All(&users)
	if err != nil {
		return nil
	}

	for i := 0; i < len(users); i++ {
		usersId = append(usersId, users[i].Telegram_id)
	}

	return usersId
}

func GetAddresses(currency string)([]string, error){
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var entry schema.NodeInfo

	c := session.DB(database).C(endpoints_collection)

	err = c.Find(bson.M{"currency": currency}).One(&entry)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	for i:=0; i<len(entry.Addresses);i++{
			entry.Addresses[i] = re.FindString(entry.Addresses[i])
	}


	return entry.Addresses, nil
}

func GetPort(currency string)(int, error){
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return 0, err
	}
	defer session.Close()

	var entry schema.NodeInfo

	c := session.DB(database).C(endpoints_collection)

	err = c.Find(bson.M{"currency": currency}).One(&entry)
	if err != nil {
		return 0, err
	}

	return entry.Port, nil
}

func GetAllNodesEntrys()(map[string]schema.NodeInfo, error){
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	var entrys []schema.NodeInfo

	c := session.DB(database).C(endpoints_collection)

	err = c.Find(nil).All(&entrys)
	if err != nil {
		return nil, err
	}

	for _, keyEntry := range entrys {
		for i:=0;i<len(keyEntry.Addresses);i++{
			keyEntry.Addresses[i] = re.FindString(keyEntry.Addresses[i])
		}
	}

	result := map[string]schema.NodeInfo{
		"btc": entrys[0],
		"eth": entrys[1],
		"etc": entrys[2],
		"bch": entrys[3],
		"ltc": entrys[4],
		"xlm": entrys[5],
	}

	return result, nil
}

func GetApiEndpoint(currency string)(string, error){
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var entry schema.NodesApi

	c := session.DB(database).C(apis_collection)

	err = c.Find(bson.M{"currency": currency}).One(&entry)
	if err != nil {
		return "", err
	}

	return entry.Endpoint, nil
}


