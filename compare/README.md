### 编译

```go build -o 12 main.go types.go```

#### 生成数据
iavldb

* 10W量级： ```./12 --case 3```
* 100W量级： ```./12 --case 4```

#### get key平均时间

```$xslt
get
./12.sh

set
./12.sh set
```