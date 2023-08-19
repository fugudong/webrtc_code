package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Student struct {
	Id bson.ObjectId 	`bson:"_id"	`
	Name  string		`bson:"name"`
	Phone string        `bson:"phone"`
	Email string        `bson:"email"`
	Sex   string        `bson:"sex"`
}

func ConnecToDB() *mgo.Collection {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("medex").C("student")
	return c
}

func InsertToMogo() {
	c := ConnecToDB()
	stu1 := Student{
		Name:"zhangsan",
		Phone: "13480989765",
		Email: "329832984@qq.com",
		Sex:   "F",
	}
	stu2 := Student{
		Name:  "liss",
		Phone: "13980989767",
		Email: "12832984@qq.com",
		Sex:   "M",
	}
	err := c.Insert(&stu1, &stu2)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDataViaSex() {
	c := ConnecToDB()
	result := Student{}
	err := c.Find(bson.M{"sex": "M"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("student", result)
	students := make([]Student, 20)
	err = c.Find(nil).All(&students)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(students)

}

func GetDataViaId() {
	id := bson.ObjectIdHex("5a66a96306d2a40a8b884049")
	c := ConnecToDB()
	stu := &Student{}
	err := c.FindId(id).One(stu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stu)
}

func UpdateDBViaId() {
	//id := bson.ObjectIdHex("5a66a96306d2a40a8b884049")
	c := ConnecToDB()
	err := c.Update(bson.M{"email": "12832984@qq.com"}, bson.M{"$set": bson.M{"name": "haha", "phone": "37848"}})
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveFromMgo() {
	c := ConnecToDB()
	_, err := c.RemoveAll(bson.M{"phone": "13480989765"})
	if err != nil {
		log.Fatal(err)
	}
}

func main()  {
	InsertToMogo()
	GetDataViaSex()
	GetDataViaId()
	UpdateDBViaId()
	RemoveFromMgo()
}