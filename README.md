# RTK离线测试工具

## 修改配置

### GO
修改[config.yaml](./BE/GOLANG/HTTP/SERVER/positioning_db_service/config.yaml)中presto的用户名和密码，用于访问presto

### 在下载目录下创建文件夹 positioning-data,将离线测试文件放到其中

## 启动:
```
docker-compose up
```
### 下载Log (postman):
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

### 解析Log
POST http://localhost:8000/parselog
```
{
    "filepath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra",
    "loglist": [
        {
            "logtype": "rtk",
            "logfile": "rtkFilter-driverID-13770990.log"
        },
        {
            "logtype": "rtk",
            "logfile": "rtkFilter-driverID-13773457.log"
        },
        {
            "logtype": "dongle",
            "logfile": "decodedongle-donglereplacement-driver13770990.log"
        }
    ],
    "outputpath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra"
}
```
### 运行RTK
POST http://localhost:8001/DEMO5/RTK/demo5
```
{
    "rtk": {
        "pmode": [
            "dgps"
        ],
        "navsys": [
            "gps-gal-glo"
        ],
        "obsnav": {
            "roverobs": "/Users/xiaolong.ji/Downloads/rtk/20230215/RoadTesting20230215/forestrouteloop1/Note20Ultra/gnss_log_2023_02_15_14_17_28.23o",
            "stationobs": "/Users/xiaolong.ji/Downloads/rtK/20230215/RoadTesting20230215/SIN100SGP_S_20230460000_01D_30S_MO.rnx",
            "stationnav": "/Users/xiaolong.ji/Downloads/rtK/20230215/RoadTesting20230215/BRDC00IGS_R_20230460000_01D_MN.rnx"
        } 
    },
    "outputpath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra"
}
```

### 自动解析结果
POST http://localhost:8000/parsertkroadtest
```
{
    "filepath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra",
    "groundtruth": {
        "type": "dongle",
        "content": "decodedongle-donglereplacement-driver13770990-device.csv"
    },
    "adrcsv": "rtkFilter-driverID-13770990-adr.csv",
    "rtkcsv": "rtkFilter-driverID-13770990-rtk.csv",
    "offlinertkcsv": [
        "PMODE_DGPS_GPS_GLO_GAL.csv",
        "PMODE_DGPS_GPS_GLO_GAL_CMP.csv",
        "PMODE_KINEMA_GPS_GLO_GAL.csv",
        "PMODE_KINEMA_GPS_GLO_GAL_CMP.csv",
        "PMODE_STATIC_GPS_GLO_GAL.csv",
        "PMODE_STATIC_GPS_GLO_GAL_CMP.csv"
    ],
    "outputpath": "/Users/xiaolong.ji/Downloads/rtk/20230215/result/forestrouteloop1/Note20Ultra/distanceSummary"
}
```

