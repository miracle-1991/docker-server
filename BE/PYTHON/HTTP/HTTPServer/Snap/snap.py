#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import csv
import json

import requests
from geopy.distance import geodesic

class Snap:
    def __init__(self):
        self.defaultUrl = "https://snap-engine.stg-myteksi.com/v3/snap"
        self.defaultAltitide = 0
        self.defaultAccuracy = 1
        self.defaultBearing  = 90
        self.defaultSpeed    = 20
        self.snapResult      = {}

    def __getLatLonFromStr__(self, gpsstr):
        l = gpsstr.split(",")
        lat, lon = l[0], l[1]
        return lat, lon

    def __getDisBetGps__(self, gps1str, gps2str):
        lat1, lon1 = self.__getLatLonFromStr__(gps1str)
        lat2, lon2 = self.__getLatLonFromStr__(gps2str)
        dis = geodesic((lat1, lon1), (lat2, lon2)).m
        # 过滤掉无效的距离
        if dis < 0 or dis > 5000:
            return -1
        else:
            return dis

    def __send_request_to_snap_engine__(self, locationlist):
        if len(locationlist) > 100 or len(locationlist) == 0:
            print("Invalid input for snap engine")
        data = [
            ["source", "mustang_automation"],
            ["vecID", 69]
        ]
        for pos in locationlist:
            location = ",".join([
                str(pos["lat"]),
                str(pos["lon"]),
                str(self.defaultAltitide),
                str(self.defaultAccuracy),
                str(self.defaultBearing),
                str(self.defaultSpeed),
                str(pos["timestamp"])
            ]),
            data.append(["location", location])
        response = requests.get(url=self.defaultUrl, params=data)
        jsonresp = json.loads(response.text)
        resultlist = []
        if jsonresp["status"] == True:
            snappedPoints = jsonresp["snappedPoints"]
            for sp in snappedPoints:
                originalIndex = sp["originalIndex"]
                snappedFlag = sp["snappedFlag"]
                locationlat = sp["location"]["lat"]
                locationlon = sp["location"]["lon"]
                originlat = locationlist[originalIndex]["lat"]
                originlon = locationlist[originalIndex]["lon"]
                origintime = locationlist[originalIndex]["timestamp"]
                result = {
                    "snappedFlag": snappedFlag,
                    "locationlat": locationlat,
                    "locationlon": locationlon,
                    "originlat": originlat,
                    "originlon": originlon,
                    "origintime": origintime,
                    "distance": self.__getDisBetGps__(str(locationlat) + "," + str(locationlon),
                                                      str(originlat) + "," + str(originlon))
                }
                resultlist.append(result)
        return resultlist
    def getSnapResult(self, latlonlist):
        locationslicelist = []
        for i in range(0, len(latlonlist), 100):
            locationlist = latlonlist[i:i + 100]
            item = {
                "startIndex": i,
                "locationlist": locationlist,
                "size": len(locationlist)
            }
            locationslicelist.append(item)

        resultlist = []
        for slice in locationslicelist:
            resp = self.__send_request_to_snap_engine__(slice["locationlist"])
            resultlist.extend(resp)
        self.snapResult.clear()
        for r in resultlist:
            self.snapResult[int(r["origintime"])] = r
        return resultlist

    def findSnapResult(self, timestamp):
        return self.snapResult.get(int(timestamp))

class SNAPCSVReader:
    def __init__(self, filepath):
        self.path = filepath
        self.rowlist = []
        rowlist = []
        with open(self.path) as f:
            reader = csv.DictReader(f)
            for row in reader:
                rowlist.append(row)
        self.rowlist = rowlist

    def __writeToCSVFile__(self):
        with open(self.path, "w") as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(self.rowlist[0].keys())
            for r in self.rowlist:
                writer.writerow(r.values())
            print("over write ", self.path, " success")

    def OverWriteSnapResultToOriginFile(self, LatitudeColumnName = "lat", LongitudeColumnName = "lon", TimestampColumnName="timestamp"):
        latlonlist = []
        for r in self.rowlist:
            lat = float(r[LatitudeColumnName])
            lng = float(r[LongitudeColumnName])
            ts  = int(r[TimestampColumnName])
            latlonitem = {
                "lat": lat,
                "lon": lng,
                "timestamp": ts
            }
            latlonlist.append(latlonitem)
        s = Snap()
        s.getSnapResult(latlonlist)
        rowlist = []
        for r in self.rowlist:
            newrow = r
            result = s.findSnapResult(r[TimestampColumnName])
            if result is not None:
                newrow["snapedLat"] = result["locationlat"]
                newrow["snapedLon"] = result["locationlon"]
                newrow["distanceToSnapGPS"] = result["distance"]
                newrow["snappedFlag"] = result["snappedFlag"]
            else:
                print("can't get result for ", r[TimestampColumnName])
                newrow["snapedLat"] = -1
                newrow["snapedLon"] = -1
                newrow["distanceToSnapGPS"] = -1
                newrow["snappedFlag"] = False
            rowlist.append(newrow)
        if len(rowlist) == len(self.rowlist):
            self.rowlist = rowlist
            self.__writeToCSVFile__()
            return True
        else:
            return False

if __name__ == '__main__':
    csvFilePath = "/Users/xiaolong.ji/code/position-tool-python/road-test/BE/ADR raw.csv"
    reader = SNAPCSVReader(csvFilePath)
    reader.OverWriteSnapResultToOriginFile(LatitudeColumnName="lat", LongitudeColumnName="lon", TimestampColumnName="timestamp")

