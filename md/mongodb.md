## MongoDB

1. 启动数据库

   ```shell
   mongod -dbpath e:\mongodb\data\db --bind_ip_all
   ```

2. 连接数据库

   ```shell
   mongosh mongodb://user:password@127.0.0.1:27017
   ```

3. 查询所有数据库

   ```shell
   show dbs # 如果当前所在的数据库没有集合，是不会显示在查询的结果中
   ```

4. 查询当前数据库

   ```shell
   db
   ```

5. 新建数据库

   ```shell
   use newdatabase
   ```

6. 删除数据库

   ```shell
   db.dropDatabase() # 删除当前数据库,貌似没有删除目标数据库的命令
   ```

7. 查询集合

   ```shell
   show collections
   ```

8. 创建集合

   ```shell
   # db.createCollection(name, option)
   db.createCollection("mycollection", {capped : true, autoIndexId : true, size : 6142800, max : 10000})
   # db.mycol.insert({"name": "菜鸟教程"})
   ```

9. 删除集合

   ```
   db.mycol.drop()
   ```

10. 插入文档

   ```shell
   db.mycol.insertOne({name:"myname"})
   db.mycol.insertMany({name:"myname"},{name:"myname1"},{name:"myname2"})
   ```

11. 更新文档

    ```shell
    # 更新一个文档
    db.mycol.updateOne({name:'myname'},{$set:{name:'newname'}})
    # 更新全部匹配文档
    db.mycol.updateMany({name:'myname'},{$set:{name:'newname'}})
    # 更新全部文档
    db.mycol.updateMany({},{$set:{name:'newname'}})
    # 移除集合中的键值对
    db.mycol.updateOne({name:'xxname'},{$unset:{name:'xxname'}})
    # 数据存在是不进行操作
    db.mycol.update({title:'mongodb 教程'}, {$setOnInsert:{title:'MongoDB'}})
    ```

12. 删除文档

    ```shell
    db.mycol.deleteOne({name:"myname"}) # 删除name等于myname的一个文档
    db.mycol.deleteMany({name:"myname"}) # 删除name等于myname的全部文档
    ```

13. 查询文档

    ```shell
    db.mycol.find() # 查找所有结果
    db.mycol.findOne() # 查找一个结果
    db.mycol.find({name:'myname'})
    
    db.mycol.find({num:{$gt:10}}) # 查询num>10
    db.mycol.find({num:{$gte:10}}) # 查询num>=10
    db.mycol.find({num:10}) # 查询num=10
    db.mycol.find({num:{$lt:10}}) # 查询num<10
    db.mycol.find({num:{$ne:10}}) # 查询num!=10
    db.mycol.find({num:{$gte:6,$lt:10}}) # 查询6<=num<10
    ```

14. 替换文档

    ```shell
    # 前面的参数为表Id，后面为新的表
    db.mycol.replaceOne({_id:ObjectId("64c32be12f4957dfc5a58215")},{num:7})
    ```

15. 分页和排序

    ```shell
    # limit：限制3个  skip：跳过前面n个
    db.mycol.find().limit(3).skip(0) 
    # 排序，根据key值来进行排序，key:1升序，key:-1降序
    db.mycol.find().sort({num:1})
    ```

16. 索引

    ```shell
    # 查看集合索引
    db.mycol.getIndexes()
    # 查看集合索引大小
    db.mycol.totalIndexSize()
    # 删除集合的所有索引
    db.mycol.dropIndexes()
    # 删除集合指定索引
    db.mycol.dropIndex('索引名称')
    ```

17. 利用TTL集合对存储的数据进行失效时间设置：经过指定的时间段后或在指定的时间点过期，MongoDB独立线程去清除数据。

    ```
    
    ```

18. 聚合

    ```shell
    # 计算总和
    db.user.aggregate([{$group:{_id: '$phone', num:{$sum:'$age'}}}])
    # 计算平均值
    db.user.aggregate([{$group:{_id: '$phone', num:{$avg:'$age'}}}])
    # 获取集合中所有文档对应值的最小值
    db.user.aggregate([{$group:{_id: '$phone', num:{$min:'$age'}}}])
    # 获取集合中所有文档对应值的最大值
    db.user.aggregate([{$group:{_id: '$phone', num:{$max:'$age'}}}])
    # 将值加入一个数组中，不会判断是否有重复的值
    db.user.aggregate([{$group:{_id: '$phone', ages:{$push:'$age'}}}])
    # 将值加入一个数组中，会判断是否有重复的值，若相同的值在数组中已经存在了，则不加入
    db.user.aggregate([{$group:{_id: '$phone', ages:{$addToSet:'$age'}}}])
    # 根据资源文档的排序获取第一个文档数据
    db.user.aggregate([{$group:{_id: '$phone', ages:{$last:'$age'}}}])
    # 根据资源文档的排序获取最后一个文档数据
    db.user.aggregate([{$group:{_id: '$phone', ages:{$first:'$age'}}}])
    ```

    管道在Unix和Linux中一般用于将当前命令的输出结果作为下一个命令的参数。

    MongoDB的聚合管道将MongoDB文档在一个管道处理完毕后将结果传递给下一个管道处理。管道操作是可以重复的。

    这里我们介绍一下聚合框架中常用的几个操作：

    - $project：修改输入文档的结构。可以用来重命名、增加或删除域，也可以用于创建计算结果以及嵌套文档。

      ```shell
      # 显示_id,tilte和author三个字段
      db.article.aggregate({$project:{title:1,author:1}})
      # 不显示_id字段
      db.article.aggregate({$project:{_id:0,title:1,author:1}})
      ```

    - $match：用于过滤数据，只输出符合条件的文档。$match使用MongoDB的标准查询操作。

    - $limit：用来限制MongoDB聚合管道返回的文档数。

    - $skip：在聚合管道中跳过指定数量的文档，并返回余下的文档。

    - $unwind：将文档中的某一个数组类型字段拆分成多条，每条包含数组中的一个值。

    - $group：将集合中的文档分组，可用于统计结果。

      ```shell
      db.user.aggregate([{$project:{_id:0,phone:1,age:1}},{$group:{_id:'$phone',count:{$sum:2}}}])
      ```

    - $sort：将输入文档排序后输出。

    - $geoNear：输出接近某一地理位置的有序文档。

19. 副本集

    副本集：同一份数据被保存在N台机器上，每台机器上都有一份数据。

    ```shell
    # --replSet 副本集合只有四同一个副本集合的
    mongod --port 27017 --dbpath e:\mongodb\data\db --bind_ip_all --replSet rs0
    mongod --port 27018 --dbpath e:\mongodb\data\db27018 --bind_ip_all --replSet rs0
    mongod --port 27019 --dbpath e:\mongodb\data\db27019 --bind_ip_all --replSet rs0
    # 连接server
    mongosh --port 27017
    # 初始化副本集
    rs.initiate({_id:'rs0',members:[{_id:0,host: 'YLMF-2020ZSMPQF:27020'},{_id:1,host: 'YLMF-2020ZSMPQF:27021'}]})
    # 新增加副本
    rs.reconfig({_id:'rs0',members:[{_id:0,host: 'YLMF-2020ZSMPQF:27020'},{_id:1,host: 'YLMF-2020ZSMPQF:27021'}]})
    # 添加目标机器的host地址(initiate方法生成的host)为副节点
    rs.add('YLMF-2020ZSMPQF:27019')
    # 查看状态
    rs.status()
    # 查看配置信息
    rs.config()
    ```

20. 分片

    一份数据被分开保存在N台机器上，N个机器上的数据组合起来是一份数据。

    ```shell
    # 副本集rs0
    mongod --port 27020 --dbpath=e:/mongodb/shard/s0 --logpath=e:/mongodb/shard/log/s0.log --logappend --shardsvr --replSet=rs0 --bind_ip_all
    mongod --port 27021 --dbpath=e:/mongodb/shard/s1 --logpath=e:/mongodb/shard/log/s1.log --logappend --shardsvr --replSet=rs0 --bind_ip_all
    
    mongosh --port 27020
    rs.initiate({_id:'rs0',members:[{_id:0,host: 'YLMF-2020ZSMPQF:27020'},{_id:1,host: 'YLMF-2020ZSMPQF:27021'}]})
    
    # 副本集rs1
    mongod --port 27022 --dbpath=e:/mongodb/shard/s2 --logpath=e:/mongodb/shard/log/s2.log --logappend --shardsvr --replSet=rs1 --bind_ip_all
    mongod --port 27023 --dbpath=e:/mongodb/shard/s3 --logpath=e:/mongodb/shard/log/s3.log --logappend --shardsvr --replSet=rs1 --bind_ip_all
    
    mongosh --port 27022
    rs.initiate({_id:'rs0',members:[{_id:0,host: 'YLMF-2020ZSMPQF:27020'},{_id:1,host: 'YLMF-2020ZSMPQF:27021'}]})
    
    # 副本集conf configsvr
    mongod --port 27100 --dbpath=e:/mongodb/shard/config0 --logpath=e:/mongodb/shard/log/config0.log --logappend --configsvr --replSet=conf --bind_ip_all
    mongod --port 27101 --dbpath=e:/mongodb/shard/config1 --logpath=e:/mongodb/shard/log/config1.log --logappend --configsvr --replSet=conf --bind_ip_all
    mongosh --port 27100
    rs.initiate({_id:'rs0',members:[{_id:0,host: 'YLMF-2020ZSMPQF:27100'},{_id:1,host: 'YLMF-2020ZSMPQF:27101'}]})
    
    
    mongos --port 28000 --configdb conf/localhost:27100,localhost:27101 --logpath=e:/mongodb/shard/log/route.log --bind_ip_all
    
    # 增加
    # 增加分片
    db.runCommand({addshard:'rs0/YLMF-2020ZSMPQF:27020,YLMF-2020ZSMPQF:27021'})
    sh.addShard('rs1/YLMF-2020ZSMPQF:27022,YLMF-2020ZSMPQF:27022')
    
    # 查看分片列表
    db.runCommand({listshards:1})
    # 查看分片状态
    sh.status()
    # 对集合所在的数据库启用分片功能
    db.runCommand({enablesharding: "testdb"})
    sh.enableSharding("testdb")
    # 对集合设置数据分片
    db.runCommand({shardcollection:"testdb.users",key:{name:1,age:1,phone:1}})
    sh.shardCollection("testdb.users",{name:1,age:1,phone:1})
    ```

    

21. Mongodb报错信息

    * 开启副本集后，连接次节点，查询时出错

      ```shell
      not primary and secondaryOk=false - consider using db.getMongo().setReadPref() or readPreference in the connection string
      # 只能从主节点读取数据
      ```

    * 进入mongodb后的提示

      ```shell
      Access control is not enabled for the database. Read and write access to data and configuration is unrestricted
      # 数据库的权限控制未启用，数据和配置的读写权限是没有限制的
      ad.createUser({user:"admin",pwd:"123123",roles:[{role:"root",db:"admin"}]})
      ```

    * 

22. 添加用户

    ```shell
    # 退出前
    db.addUser('sa','密码')
    # 登录后
    use admin
    db.auth('sa','密码')
    ```

23. 数据库恢复

    ```shell
    # 备份数据
    mongodump -h <hostname><:port> -d <dbname> -o <path>
    # 恢复数据
    mongorestore -h <hostname><:port> -d <newdbname> <path>
    ```

24. Mongodb监控

```shell
mongostat 
# mongostat是mongodb自带的状态检测工具，在命令行下使用。它会间隔固定时间获取mongodb的当前运行状态，并输出。如果你发现数据库突然变慢或者有其他问题的话，你第一手的操作就考虑采用mongostat来查看mongo的状态
mongostop
# mongotop也是mongodb下的一个内置工具，mongotop提供了一个方法，用来跟踪一个MongoDB的实例，查看哪些大量的时间花费在读取和写入数据。 mongotop提供每个集合的水平的统计数据。默认情况下，mongotop返回值的每一秒。

```

## Redis
1. 安装
```
# 下载最新的发行版
wget https://download.redis.io/redis-stable.tar.gz
# 解压
tar -xvf redis-stable.tar.gz
# 编译
make
# 安装
make install
# 清除
make clean
```

## ssh配置
~/.ssh/config
```
Host name
    HostName ip
    Port 6666
    User cyf

Host name
    HostName ip # 目标主机地址
    Port 6666   # 端口
    User cyf    # linux用户名
    ProxyJump ip # 跳板机地址
``` 








