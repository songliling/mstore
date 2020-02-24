### 编译

```go build -o test main.go```

#### 生成数据

levelDB:

* 10W量级： ```./test --case 1```
* 100W量级： ```./test --case 2```

iavlDB:

* 10W量级： ```./test --case 3```
* 100W量级： ```./test --case 4```

#### get key平均时间

levelDB:

* 10W量级： ```./test --case 1 --time```
* 100W量级： ```./test --case 2 --time```

iavlDB:

* 10W量级： ```./test --case 3 --time```
* 100W量级： ```./test --case 4 --time```

#### set key平均时间

levelDB:

* 10W量级： ```./test --case 1 --set```
* 100W量级： ```./test --case 2 --set```

iavlDB:

* 10W量级： ```./test --case 3 --set```
* 100W量级： ```./test --case 4 --set```