#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import csv

import numpy as np
from geopy.distance import geodesic

class DistanceDiffItem:
    def __init__(self, timestamp, groundTruthLat, groundTruthLon, curentLat, curentLon):
        self.timestamp = timestamp
        self.glat = groundTruthLat
        self.glon = groundTruthLon
        self.clat = curentLat
        self.clon = curentLon
        self.dis  = self.__getDisBetGps__(lat1=groundTruthLat,lon1=groundTruthLon,lat2=curentLat,lon2=curentLon)

    def __getDisBetGps__(self, lat1, lon1, lat2, lon2):
        return geodesic((lat1, lon1), (lat2, lon2)).m

    def getDis(self):
        return self.dis

    def getKeys(self):
        return ["timestamp", "groundTruthLat","groundTruthLon","curentLat","curentLon","distanceToGroundTruth"]

    def getVals(self):
        return [self.timestamp,self.glat,self.glon,self.clat,self.clon,self.dis]

    def getGroundTruthLatLon(self):
        return self.glat, self.glon

class DistanceDiff:
    def __init__(self):
        self.diffList = []
        self.disList = []

    def PushBack(self, item):
        self.diffList.append(item)
        self.disList.append(item.getDis())
    def P50(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return np.quantile(arr, q = 0.5)

    def P95(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return np.quantile(arr, q = 0.95)

    def P99(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return np.quantile(arr, q = 0.99)

    def Min(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return arr.min()

    def Max(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return arr.max()

    #平均数
    def Average(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return arr.mean()

    #中位数
    def Median(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return np.median(arr)

    #标准差
    def Std(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return arr.std()

    #方差
    def Var(self):
        if (len(self.disList) == 0):
            return None
        arr = np.array(self.disList)
        return arr.var()

    def Empty(self):
        if (len(self.disList) == 0):
            return True
        return False

    def SummaryDistance(self):
        if (len(self.disList) == 0):
            return None
        return {
            "p50": self.P50(),
            "p95": self.P95(),
            "p99": self.P99(),
            "min": self.Min(),
            "max": self.Max(),
            "average": self.Average(),
            "median": self.Median(),
            "std": self.Std(),
            "var": self.Var()
        }

    def WriteToCSV(self, filename):
        if len(self.diffList) == 0:
            return False
        with open(filename, "w") as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(self.diffList[0].getKeys())
            for r in self.diffList:
                writer.writerow(r.getVals())
        return True

    def getCenterLatLon(self):
        if len(self.diffList) == 0:
            return None, None
        return self.diffList[0].getGroundTruthLatLon()




