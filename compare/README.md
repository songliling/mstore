### 编译

```
flags="-X 'main.PruningType=sync'"
go build -ldflags "$flags" -x -o 13_sync main.go types.go

flags="-X 'main.PruningType=every'"
go build -ldflags "$flags" -x -o 13_every main.go types.go

flags="-X 'main.PruningType=nothing'"
go build -ldflags "$flags" -x -o 13_nothing main.go types.go
```

#### 生成数据

sync pruning   
* 10W量级： ```./13_sync --case 3```
* 100W量级： ```./13_sync --case 4```

every pruning
* 10W量级： ```./13_every --case 3```
* 100W量级： ```./13_every --case 4```

nothing pruning
* 10W量级： ```./13_nothing --case 3```
* 100W量级： ```./13_nothing --case 4```

#### 观测数据变化

```$xslt
get 
./sync.sh
set
./sync.sh set

get 
./every.sh
set
./every.sh set

get 
./nothing.sh
set
./nothing.sh set
```