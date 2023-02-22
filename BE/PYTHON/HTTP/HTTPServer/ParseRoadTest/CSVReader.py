#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import csv

#在线RTK定位结果解析
class STGRTKCSVReader:
    def __init__(self, filepath):
        self.path = filepath
        self.rowlist = []
        self.rowMap = {}
        with open(self.path) as f:
            reader = csv.DictReader(f)
            for row in reader:
                self.rowlist.append(row)
                self.rowMap[int(row["timestamp"])] = row
    def getTimestampList(self):
        tsl = []
        for row in self.rowlist:
            tsl.append(int(row["timestamp"]))
        return tsl

    #按照时间戳查找
    def getLatLonByTimeStamp(self, timestamp):
        row = self.rowMap.get(int(timestamp))
        if row is None:
            return None, None
        return float(row["lat"]), float(row["lon"])
    def getSuccessCount(self):
        count = 0
        for r in self.rowlist:
            if r["rtk-state"] != "0":
                count += 1
        return count
    def getSuccessRate(self):
        c = self.getSuccessCount()
        a = len(self.rowlist)
        return c/a
    def printSuccessRate(self):
        r = self.getSuccessRate()
        print("%.2f%%"%r)

#离线定位结果解析
class OFFLINERTKCSVReader:
    def __init__(self, filepath):
        self.path = filepath
        self.rowlist = []
        self.rowMap = {}
        with open(self.path) as f:
            reader = csv.DictReader(f)
            for row in reader:
                self.rowMap[int(row["gpstime"])] = row
                self.rowlist.append(row)

    def getTimestampList(self):
        tsl = []
        for row in self.rowlist:
            tsl.append(int(row["gpstime"]))
        return tsl

    def getLatLonByTimeStamp(self, timestamp):
        row = self.rowMap.get(int(timestamp))
        if row is None:
            return None, None
        if row["state"] == "0" or row["state"] == 0:
            return None, None
        return float(row["lat"]), float(row["lng"])

    def getSuccessCount(self):
        count = 0
        for r in self.rowlist:
            if r["state"] != "0":
                count += 1
        return count

    def getSuccessRate(self):
        c = self.getSuccessCount()
        a = len(self.rowlist)
        return c/a
    def printSuccessRate(self):
        r = self.getSuccessRate()
        print("%.2f%%"%r)

#在线安卓定位结果解析
class ADRCSVReader:
    def __init__(self, filepath):
        self.path = filepath
        self.rowlist = []
        self.rowMap = {}
        with open(self.path) as f:
            reader = csv.DictReader(f)
            for row in reader:
                self.rowMap[int(row["timestamp"])] = row
                self.rowlist.append(row)
    def getLatLonByTimeStamp(self, timestamp):
        row = self.rowMap.get(int(timestamp))
        if row is None:
            return None, None
        return float(row["lat"]), float(row["lon"])

    def getTimestampList(self):
        tsl = []
        for row in self.rowlist:
            tsl.append(int(row["timestamp"]))
        return tsl

#在线dongle定位结果解析
class DongleCSVReader:
    def __init__(self, filepath):
        self.path = filepath
        self.rowlist = []
        self.rowMap = {}
        with open(self.path) as f:
            reader = csv.DictReader(f)
            for row in reader:
                self.rowMap[int(row["timestamp"])] = row
                self.rowlist.append(row)

    def getTimestampList(self):
        tsl = []
        for row in self.rowlist:
            tsl.append(int(row["timestamp"]))
        return tsl

    def getLatLonByTimeStamp(self, timestamp):
        row = self.rowMap.get(int(timestamp))
        if row is None:
            return None, None
        return float(row["lat"]), float(row["lon"])

    def getResultMap(self):
        resultMap = {}
        for row in self.rowlist:
            ts  = int(row["timestamp"])
            lat = row["lat"]
            lon = row["lon"]
            resultMap[ts] = {
                "lat": float(lat),
                "lon": float(lon)
            }
        return resultMap
