#!/usr/bin/python3
# -*- coding: UTF-8 -*-
from ParseRoadTest.CSVReader import DongleCSVReader
# from ParseRoadTest.OCRReader import JPGReader
class GroundTruth:
    def __init__(self, gttype, content):
        self.gttype       = gttype
        self.content    = content
        self.recordMap  =  {}
        if gttype == "text":
            self.__parseText__(content)
        elif gttype == "dongle":
            self.__parseDongle__(content)
        elif gttype == "img":
            self.__parseIMG__(content)

    #解析手动输入的经纬度，比如 1.2910918,103.7929637
    def __parseText__(self, text):
        text = text.replace(" ", "").replace("\n", "")
        latlon = text.split(",")
        self.recordMap = {
            0: {
                "lat": float(latlon[0]),
                "lon": float(latlon[1])
            }
        }

    #解析dongle采集到的经纬度
    def __parseDongle__(self, filename):
        reader = DongleCSVReader(filename)
        self.recordMap = reader.getResultMap()

    # #解析截图中的经纬度
    # def __parseIMG__(self, imgfile):
    #     reader = JPGReader(imgfile)
    #     lat, lon = reader.getLatLon()
    #     self.recordMap = {
    #         0 : {
    #             "lat": float(lat),
    #             "lon": float(lon)
    #         }
    #     }

    def getGroundTruthByTimeStamp(self, timestamp):
        if self.gttype == "text":
            return self.recordMap[0]["lat"], self.recordMap[0]["lon"]
        elif self.gttype == "dongle":
            ts = int(timestamp)
            item = self.recordMap.get(ts)
            if item is None:
                return None, None
            else:
                return item["lat"], item["lon"]
        elif self.gttype == "img":
            return self.recordMap[0]["lat"], self.recordMap[0]["lon"]