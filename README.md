## ddz-api

### 部署方法
```bash
# 服务部署在localhost:3005
# 一共有三个接口
# /post 单步智能出牌调的接口
# /deck 读取后端预设牌组 牌组格式为json 路径在根目录decks下
# /gametable 根据传入的参数显示ai走的每一步

cd  template/go
go build .

./go
``` 