布隆过滤器实现和测试性能
![测试数据图](https://cdn.jsdelivr.net/gh/fx-zpy/PictureBed@main/img/20230301180647.png)
使用
```
go get -u github.com/fx-zpy/BloomFilter
```
在代码中首先建立一个过滤器
```
filter := BloomFilter.New(100, 5, false)
filter.AddBatch([][]byte{[]byte("华科"), []byte("内核"), []byte("结束"), []byte("监控室的"), []byte("静安寺"), []byte("骄傲生死看淡"), []byte("2178126"),
[]byte("卡视角的"), []byte("怕视频的"), []byte("清华"), []byte("阿斯倒数第"), []byte("阿卡十多年"), []byte("人数好的"), []byte("啥计算机啊"), []byte("卡技术都能"), []byte("奥斯卡的年"), []byte("阿康"),
[]byte("那看到你"), []byte("爬山的"), []byte("明年"), []byte("美女好看"), []byte("啥的你可能"), []byte("snake"), []byte("阿克苏"), []byte("阿萨斯")})

fmt.Printf("snake exist：%t\n", filter.Test([]byte("snake"))) //true
fmt.Printf("aksdjaisdnka exist：%t\n", filter.Test([]byte("aksdjaisdnka"))) //false
fmt.Printf("啥的你可能 exist：%t", filter.Test([]byte("啥的你可能"))) //true
```
然后批量插入数据，然后可以查询数据是否在过滤器中，性能好。可以准确查出不在其中的数据，但是对于判断为true的数据存在误报的可能性。
