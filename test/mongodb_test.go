package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name  string
	Phone string
	Age   int32
}

func TestInsertOne(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:28000", Database: "testdb", Coll: "users"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	res, err := cli.InsertOne(ctx, &User{Name: "王波"})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestInsertMany(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:28000", Database: "test", Coll: "users"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	many := make([]User, 50000-10000)
	for i := 0; i < 40000; i++ {
		many[i] = User{Name: "xxxxx", Phone: "123123", Age: int32(i + 10)}
	}
	res, err := cli.InsertMany(ctx, many)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestFindOne(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27018", Database: "runoob", Coll: "user"})
	if err != nil {
		panic(err)
	}
	defer cli.Close(ctx)
	one := &User{}
	err = cli.Find(ctx, bson.M{"name": "王波"}).One(one)
	if err != nil {
		panic(err)
	}
	fmt.Println(one)
}

func TestFindMany(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "runoob", Coll: "user"})
	if err != nil {
		panic(err)
	}
	many := make([]User, 10)
	err = cli.Find(ctx, bson.M{}).All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
}

// 删除集合
func TestDeleteCollection(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:28000", Database: "test", Coll: "users"})
	if err != nil {
		panic(err)
	}
	err = cli.DropCollection(ctx)
	if err != nil {
		panic(err)
	}
}

// 删除数据库
func TestDeleteDB(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "test", Coll: "users"})
	if err != nil {
		panic(err)
	}
	cli.DropDatabase(ctx)
}

// 删除键
func TestDeleteKey(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017/?connect=direct", Database: "runoob", Coll: "col"})
	if err != nil {
		panic(err)
	}
	cli.InsertOne(ctx, &User{Name: "达达", Phone: "123123", Age: 22})
	// cli.UpdateOne(ctx, bson.M{"name": "王波"}, bson.M{"$set": {"age": 20}})
}

func TestSort(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017/?connect=direct", Database: "runoob", Coll: "user"})
	if err != nil {
		panic(err)
	}
	many := make([]User, 10)
	// descend
	err = cli.Find(ctx, bson.M{}).Sort("-age").All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
	// ascend
	err = cli.Find(ctx, bson.M{}).Sort("+age").All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
}

func TestLimit(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "runoob", Coll: "user"})
	if err != nil {
		panic(err)
	}
	many := make([]User, 10)
	err = cli.Find(ctx, bson.M{}).Limit(3).Sort("+age").All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
}

func TestSkip(t *testing.T) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://172.20.11.80:27017", Database: "runoob", Coll: "user"})
	if err != nil {
		panic(err)
	}
	many := make([]User, 10)
	err = cli.Find(ctx, bson.M{}).Limit(3).Skip(0).Sort("+age").All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
	err = cli.Find(ctx, bson.M{}).Limit(3).Skip(3).Sort("+age").All(&many)
	if err != nil {
		panic(err)
	}
	for i := range many {
		fmt.Println(many[i])
	}
	err = cli.Find(ctx, bson.M{}).Limit(3).Skip(6).Sort("+age").All(&many)
	if err != nil {
		panic(err)
	}

	for i := range many {
		fmt.Println(many[i])
	}
}

// 副本集 查询主节点
func TestReplSet(t *testing.T) {
	ctx := context.Background()
	clis := make([]*qmgo.QmgoClient, 3)
	for i, port := range []int{27017, 27018, 27019} {
		cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: fmt.Sprintf("mongodb://172.20.11.80:%d/?connect=direct", port), Database: "runoob", Coll: "user"})
		if err != nil {
			continue
		}
		clis[i] = cli
	}

	var maindb *qmgo.QmgoClient

	for _, cli := range clis {
		if cli != nil {
			res := cli.RunCommand(ctx, bson.M{"isMaster": 1})
			xx := make(map[string]interface{})
			err := res.Decode(xx)
			isMain := xx["ismaster"].(bool)
			fmt.Println(err, xx, isMain)
			if isMain {
				maindb = cli
			}
		}
	}
	maindb.InsertOne(ctx, bson.M{"name": "巫医", "age": 44, "phone": "123123"})
}
