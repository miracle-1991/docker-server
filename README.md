# RTK离线测试工具

## 修改配置

### GO
修改[config.yaml](./BE/GOLANG/HTTP/SERVER/positioning_db_service/config.yaml)中presto的用户名和密码，用于访问presto

### 在下载目录下创建文件夹 positioning-data,将离线测试文件放到其中

## 启动:
```
docker-compose up
```
## 下载Log (postman):
POST http://localhost:8000/downlog
```
{
    "starttime": "2023-02-15 06:17:46",
    "endtime": "2023-02-15 06:38:33",
    "filter": [
        "rtkFilter",
        "driverID:13770990"
    ],
    "outputpath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra"
}
```