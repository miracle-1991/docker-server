#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import easyocr

ocrreader = easyocr.Reader(['en'])

class JPGReader:
    def __init__(self, jpgfile):
        self.filename = jpgfile
        self.lat, self.lon, self.timestr = self.__getLatLngFromJPG__(jpgfile)

    def __getLatLngFromJPG__(self, jpgfile):
        resultlist = ocrreader.readtext(jpgfile)
        lat, lng, timestr = None, None, None
        for index in range(len(resultlist)):
            piclocation, text, prob = resultlist[index]
            if text.find("Latitude") != -1:
                lat = float(resultlist[index+1][1].replace(",",""))
            if text.find("Longitude") != -1:
                lng = float(resultlist[index+1][1].replace(",",""))
            if text.find("UTC") != -1:
                curindex = index - 1
                timestr = resultlist[curindex][1]
                if resultlist[curindex-1][1].find("MAP VIEW") == -1:
                    timestr = resultlist[curindex-1][1] + " : " + timestr
                    if resultlist[curindex-2][1].find("MAP VIEW") == -1:
                        timestr = resultlist[curindex - 2][1] + " : " + timestr
        return lat, lng, timestr

    def getLatLon(self):
        return self.lat, self.lon
