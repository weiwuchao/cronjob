package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type logRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type findByName struct {
	JobName string `bson:"jobName"`
}

func main() {

	var (
		client          *mongo.Client
		err             error
		clientOptions   *options.ClientOptions
		database        *mongo.Database
		connection      *mongo.Collection
		insertOneResult *mongo.InsertOneResult
		cursor          *mongo.Cursor
		delResult *mongo.DeleteResult
		updateResult *mongo.UpdateResult
	)

	/**
	建立mongo连接
	*/
	clientOptions = options.Client()
	// 设置最大连接的数量
	clientOptions = clientOptions.SetMaxPoolSize(uint64(10))
	// 设置连接超时时间 5000 毫秒
	clientOptions = clientOptions.SetConnectTimeout(2 * time.Minute)
	// 设置连接的空闲时间 毫秒
	clientOptions = clientOptions.SetMaxConnIdleTime(2 * time.Minute)
	clientOptions.ApplyURI("mongodb://192.168.92.129:27017")
	if client, err = mongo.Connect(context.Background(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}

	/**
	选择数据库
	*/
	database = client.Database("cron")

	/**
	选择表
	*/
	connection = database.Collection("log")

	//插入数据
	record := &logRecord{
		JobName: "job1",
		Command: "ls =a",
		Err:     "",
		Content: "test11",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	if insertOneResult, err = connection.InsertOne(context.Background(), record); err != nil {
		fmt.Println(err)
		return
	}
	docId := insertOneResult.InsertedID
	fmt.Println("插入ID:", docId)

	/**
		更新数据
	*/
	updateFilter:=bson.M{"jobName":"job1"}
	updateData:=bson.M{"$set":bson.M{"content":"oouwyw"}}
	if updateResult,err=connection.UpdateOne(context.Background(),updateFilter,updateData);err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("更新数据数目",updateResult.ModifiedCount)

	/**
	查询数据
	*/
	var limit int64 = 1
	var skip int64 = 0
	opts := &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}
	/*cond := &findByName{
		JobName: "job1",
	}
	if cursor, err = connection.Find(context.Background(), cond, opts); err != nil {
		fmt.Println(err)
		return
	}*/
	if cursor, err = connection.Find(context.Background(), bson.D{{"jobName","job1"}}, opts); err != nil {
		fmt.Println(err)
		return
	}
	//延迟关闭游标
	defer cursor.Close(context.Background())
	//遍历游标
	for cursor.Next(context.Background()){
		record=&logRecord{}
		//将bson反序列化到对象
		if err=cursor.Decode(record);err!=nil{
			fmt.Println(err)
			return
		}
		//打印对象
		fmt.Println(*record)
	}

	/**
	  	删除数据
		删除小于当前时间执行的{"timePoint.startTime":{"$lt":now()}
	 */
	/*filter:=bson.D{}
	filter=append(filter,bson.E{"timePoint.startTime",bson.D{{"$lt",time.Now().Unix()}}})*/
	filter:=bson.D{{"timePoint.startTime",bson.D{{"$lt",time.Now().Unix()}}}}
	if delResult,err=connection.DeleteMany(context.Background(),filter);err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("删除数目",delResult.DeletedCount)

}
