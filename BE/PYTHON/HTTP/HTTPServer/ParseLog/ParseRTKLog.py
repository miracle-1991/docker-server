#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import csv
import json
import re
import time

from geopy.distance import geodesic

#解析RTK日志
class LocationEngineRTKFilterLogItem:
    def __init__(self, line):
        self.logtime                = self.__parse_logtime__(line)
        self.ip                     = self.__parse_ip__(line)
        self.driverID               = self.__parse_driverID__(line)
        self.adrlat, self.adrlon    = self.__parse_adr_ping__(line)
        self.ts                     = self.__parse_adr_ts__(line)
        self.staleDuration          = self.__parse_staleDuration__(line)
        self.accuracy               = self.__parse_accuracy__(line)
        self.gnssReq                = self.__parse_RTK_REQ_JSON__(line)
        self.rtkstate               = self.__parse_RTK_state__(line)
        self.rtkValidSatNum         = self.__parse_RTK_validSatNum__(line)
        self.rtkgpstime             = self.__parse_RTK_gpstime__(line)
        self.rtklat, self.rtklon    = self.__parse_RTK_result(line)
        self.disBetAdrRtk           = self.__get_dis_bet_adr_rtk__()

    def __parse_logtime__(self, line):
        #2023-02-01T09:18:20.603666
        timestr = str(line[:26])
        timestr = timestr.replace("T", " ")
        timelocal = time.strptime(timestr, "%Y-%m-%d %H:%M:%S.%f")
        return timelocal
    def __parse_ip__(self, line):
        pattern = "ip-(\d+)-(\d+)-(\d+)-(\d+)"
        result = re.search(pattern, line)
        if result:
            return result.group(1) + "." + result.group(2) + "." + result.group(3) + "." + result.group(4)
        else:
            return None

    def __parse_driverID__(self, line):
        parttern = "driverID:(\d+)"
        result = re.search(parttern, line)
        if result:
            return result.group(1)
        else:
            return None

    def __parse_adr_ping__(self, line):
        parttern = r'ADR ping: LatLng:{"lat":(\d+\.?\d+),"lng":(\d+\.?\d+)}'
        result = re.search(parttern, line)
        if result:
            lat = float(result.group(1))
            lon = float(result.group(2))
            return lat, lon
        else:
            return None, None

    def __parse_adr_ts__(self, line):
        parttern = r'ts:(\d+)'
        result = re.search(parttern, line)
        if result:
            ts = int(result.group(1))
            return ts
        else:
            return None

    def __parse_staleDuration__(self, line):
        parttern = r'staleDuration:(\d+)'
        result = re.search(parttern, line)
        if result:
            dur = int(result.group(1))
            return dur
        else:
            return None

    def __parse_accuracy__(self, line):
        parttern = r'accuracy:(\d+\.?\d+)'
        result = re.search(parttern, line)
        if result:
            accuracy = float(result.group(1))
            return accuracy
        else:
            return None

    def __parse_RTK_REQ_JSON__(self, line):
        try:
            startindex = line.index("RTK ping: rawGnss:") + 18
            endindex = line.index(", state")
            rawGnss = line[startindex:endindex]
            req = json.loads(rawGnss)
            return req
        except:
            print("parse failed:")
            print(line)
    def __parse_RTK_state__(self, line):
        parttern = r'state:(\d+)'
        result = re.search(parttern, line)
        if result:
            state = int(result.group(1))
            return state
        else:
            return None

    def __parse_RTK_validSatNum__(self, line):
        parttern = r'validSatNum:(\d+)'
        result = re.search(parttern, line)
        if result:
            validSatNum = int(result.group(1))
            return validSatNum
        else:
            return None

    def __parse_RTK_gpstime__(self, line):
        parttern = r'gpstime:(\d+)'
        result = re.search(parttern, line)
        if result:
            gpstime = int(result.group(1))
            return gpstime
        else:
            return None

    def __parse_RTK_result(self, line):
        parttern = r'rtkResult:{"lat":(\d+\.?\d+),"lng":(\d+\.?\d+)}'
        result = re.search(parttern, line)
        if result:
            lat = float(result.group(1))
            lon = float(result.group(2))
            return lat, lon
        else:
            return None, None

    def __getLatLonFromStr__(self, gpsstr):
        l = gpsstr.split(",")
        lat, lon = l[0], l[1]
        return lat, lon

    def __getDisBetGps__(self, gps1str, gps2str):
        lat1, lon1 = self.__getLatLonFromStr__(gps1str)
        lat2, lon2 = self.__getLatLonFromStr__(gps2str)
        return geodesic((lat1, lon1), (lat2, lon2)).m

    def __get_dis_bet_adr_rtk__(self):
        adrlat, adrlon = self.adrlat, self.adrlon
        rtklat, rtklon = self.rtklat, self.rtklon
        if adrlat is None or rtklat is None:
            return -1
        return self.__getDisBetGps__(str(adrlat)+","+str(adrlon), str(rtklat)+","+str(rtklon))

class LocationEngineRTKLogParser:
    def __init__(self, logpath):
        self.logfile = logpath
        self.rowlist = []
        rowlist = []
        with open(self.logfile) as log:
            row = log.readline()
            while row:
                item = LocationEngineRTKFilterLogItem(row)
                rowlist.append(item)
                row = log.readline()
        self.rowlist = rowlist

    def writeGNSSIntoJsonFile(self, filename):
        with open(filename, "w") as f:
            for r in self.rowlist:
                json.dump(r.gnssReq, f)
                f.write("\n")

    def writeADRPositionToCSVFile(self, filename):
        with open(filename, "w") as f:
            writer = csv.writer(f)
            keys = ["timestamp", "humantime", "lat", "lon"]
            writer.writerow(keys)
            for r in self.rowlist:
                values = [r.ts, time.strftime("%Y-%m-%d %H:%M:%S", r.logtime), r.adrlat, r.adrlon]
                writer.writerow(values)

    def writeRTKPositionToCSVFile(self, filename):
        with open(filename, "w") as f:
            writer = csv.writer(f)
            keys = ["timestamp", "humantime", "lat", "lon", "rtk-state", "rtk-valid-sat-num", "rtk-gpstime", "disToAdrPosition"]
            writer.writerow(keys)
            for r in self.rowlist:
                values = [r.ts,
                          time.strftime("%Y-%m-%d %H:%M:%S", r.logtime),
                          r.rtklat, r.rtklon,
                          r.rtkstate, r.rtkValidSatNum, r.rtkgpstime,
                          r.disBetAdrRtk]
                writer.writerow(values)

    def getSuccessRate(self):
        successItemCount = 0
        for r in self.rowlist:
            if r.rtkstate is not None and r.rtkstate != 0:
                successItemCount += 1
        return successItemCount/len(self.rowlist)

if __name__ == '__main__':
    logfile = "/Users/xiaolong.ji/Downloads/rtK/20230215/rtkFilter-driverID-13770990.log"
    parser = LocationEngineRTKLogParser(logfile)
    parser.writeGNSSIntoJsonFile("/Users/xiaolong.ji/Downloads/rtK/20230215/s3-rtk-gnss-13770990.json")
    parser.writeADRPositionToCSVFile("/Users/xiaolong.ji/Downloads/rtK/20230215/s3-adr-position-13770990.csv")
    parser.writeRTKPositionToCSVFile("/Users/xiaolong.ji/Downloads/rtK/20230215/s3-rtk-position-13770990.csv")
    print(parser.getSuccessRate())