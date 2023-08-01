package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func TestInsertOne(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "runoob", Coll: "users"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	res, err := cli.InsertOne(ctx, &Person{Name: "王波", Phone: "123123123"})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestInsertMany(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "runoob", Coll: "users"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	many := []Person{
		{Name: "xxxxx", Phone: "123123"},
		{Name: "xxxxa", Phone: "123124"},
		{Name: "xxxxb", Phone: "123125"},
		{Name: "xxxxc", Phone: "123126"},
		{Name: "xxxxd", Phone: "123127"},
		{Name: "xxxxe", Phone: "123128"},
	}
	res, err := cli.InsertMany(ctx, many)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestFindOne(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.17.16.1:27017", Database: "runoob", Coll: "users"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	one := &Person{}
	err = cli.Find(ctx, bson.M{"name": "王波"}).One(one)
	if err != nil {
		panic(err)
	}
	fmt.Println(one)
}

func TestFindMany(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.17.16.1:27017", Database: "runoob", Coll: "users"})
	if err != nil {
		panic(err)
	}
	many := make([]Person, 10)
	err = cli.Find(ctx, bson.M{"name": "*"}).All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
}
