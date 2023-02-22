#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import csv
import re
import time


class LocationEngineDongleLogItem:
    def __init__(self, line):
        self.logtime = self.__parse_logtime__(line)
        self.driveID = self.__parse_driverID__(line)
        self.devicets= self.__parse_device_ts__(line)
        self.donglets= self.__parse_dongle_ts__(line)
        self.deviceLat, self.deviceLng = self.__parse_device_position__(line)
        self.dongleLat, self.dongleLng = self.__parse_dongle_position__(line)

    def __parse_logtime__(self, line):
        #2023-02-01T09:18:20.603666
        timestr = str(line[:26])
        timestr = timestr.replace("T", " ")
        timelocal = time.strptime(timestr, "%Y-%m-%d %H:%M:%S.%f")
        return timelocal

    def __parse_driverID__(self, line):
        parttern = "driver (\d+)"
        result = re.search(parttern, line)
        if result:
            return result.group(1)
        else:
            return None

    def __parse_device_position__(self, line):
        parttern = r'device lat: (\d+\.?\d+), device lng: (\d+\.?\d+)'
        result = re.search(parttern, line)
        if result:
            lat = float(result.group(1))
            lon = float(result.group(2))
            return lat, lon
        else:
            return None, None

    def __parse_device_ts__(self, line):
        parttern = r'device ts: (\d+)'
        result = re.search(parttern, line)
        if result:
            ts = int(result.group(1))
            return ts
        else:
            return None

    def __parse_dongle_position__(self, line):
        parttern = r'dongle lat: (\d+\.?\d+), dongle lng: (\d+\.?\d+)'
        result = re.search(parttern, line)
        if result:
            lat = float(result.group(1))
            lon = float(result.group(2))
            return lat, lon
        else:
            return None, None

    def __parse_dongle_ts__(self, line):
        parttern = r'dongle ts: (\d+)'
        result = re.search(parttern, line)
        if result:
            ts = int(result.group(1))
            return ts
        else:
            return None

class LocationEngineDongleLogParser:
    def __init__(self, logpath):
        self.logfile  = logpath
        self.rowlist  = []
        self.itemlist = []
        rowlist = []
        with open(self.logfile) as log:
            row = log.readline()
            while row:
                item = LocationEngineDongleLogItem(row)
                rowlist.append(item)
                row = log.readline()
        self.rowlist = rowlist

    def writeDonglePositionToCSVFile(self, filename):
        with open(filename, "w") as f:
            writer = csv.writer(f)
            keys = ["timestamp", "humantime", "lat", "lon"]
            writer.writerow(keys)
            for r in self.rowlist:
                values = [
                    r.donglets,
                    time.strftime("%Y-%m-%d %H:%M:%S", r.logtime),
                    r.dongleLat,
                    r.dongleLng
                ]
                writer.writerow(values)

    def writeDevicePositionToCSVFile(self, filename):
        with open(filename, "w") as f:
            writer = csv.writer(f)
            keys = ["timestamp", "humantime", "lat", "lon"]
            writer.writerow(keys)
            for r in self.rowlist:
                values = [
                    r.devicets,
                    time.strftime("%Y-%m-%d %H:%M:%S", r.logtime),
                    r.deviceLat,
                    r.deviceLng
                ]
                writer.writerow(values)

if __name__ == '__main__':
    logfile = "/Users/xiaolong.ji/Downloads/rtK/20230215/decodedongle-donglereplacement-driver13770990.log"
    parser = LocationEngineDongleLogParser(logfile)
    parser.writeDonglePositionToCSVFile("/Users/xiaolong.ji/Downloads/rtK/20230215/decodedongle-donglereplacement-driver13770990-dongle.csv")