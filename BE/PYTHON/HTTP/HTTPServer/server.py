#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import json

from flask import Flask, request, jsonify
from flask_cors import CORS

from ParseLog.ParseRTKLog import LocationEngineRTKLogParser
from ParseLog.ParseDongle import LocationEngineDongleLogParser
from ParseRoadTest.CSVReader import ADRCSVReader, STGRTKCSVReader, OFFLINERTKCSVReader
from ParseRoadTest.DistanceDiff import DistanceDiffItem, DistanceDiff
from ParseRoadTest.GroundTruth import GroundTruth
from ParseRoadTest.Kepler import KeplerDraw
from Snap.snap import SNAPCSVReader
from S3Log.S3Log import AWSS3

app = Flask(__name__)
CORS(app)

user = "ssm-user"
bucket = "grabtaxi-logs-stg"
awss3 = AWSS3(user, bucket)

# hello-world 测试连通情况
@app.route('/', methods=['GET','POST'])
def helloworld():
    print(request.headers)
    print(request.json)
    resp = {
        "code": 0,
        "msg": "hello world",
    }
    return jsonify(resp)

#下载日志
@app.route('/downlog', methods=['POST'])
def downloadLogFromS3():
    print(request.headers)
    print(request.json)
    starttimestr = request.json["starttime"]
    endtimestr   = request.json["endtime"]
    targetpath   = request.json["outputpath"]
    filterlist   = request.json["filter"]
    appname      = request.json["prefix"]
    flist = []
    for filterstr in filterlist:
        f = lambda x : x.find(filterstr) != -1
        flist.append(f)
    outfilename = "-".join(map(lambda x : x.replace(":", "-"), filterlist)) + ".log"
    targetfile = targetpath + "/" + outfilename
    targetfile = targetfile.replace(" ", "").replace("\n", "")
    try:
        awss3.pull_log_into_file(appname + "/app", starttimestr, endtimestr, targetfile, filterlist)
    except Exception as e:
        print(e)
        resp = {
            "code": -1,
            "error": str(e)
        }
        return jsonify(resp)

    resp = {
        "code": 0,
        "targetfile": targetfile,
    }
    return jsonify(resp)

#下载日志进度
@app.route('/downlogProcessing', methods=['GET'])
def downlogProcessing():
    total, cur, state = awss3.getProcessingFilesCnt()
    if state == "end":
        processing = 1
    else:
        if total == 0:
            processing = 0
        else:
            processing = cur / total

    resp = {
        "code": 0,
        "total": total,
        "current": cur,
        "processing": processing
    }
    return jsonify(resp)

def parseRTKlog(filename, filepath):
    targetfile = filename.replace(".log", "")
    parser = LocationEngineRTKLogParser(filepath + "/" + filename)
    gnssjsonfile = filepath + "/" + targetfile + "-gnss.json"
    parser.writeGNSSIntoJsonFile(gnssjsonfile)
    adrcsvfile = filepath + "/" + targetfile + "-adr.csv"
    parser.writeADRPositionToCSVFile(adrcsvfile)
    rtkcsvfile = filepath + "/" + targetfile + "-rtk.csv"
    parser.writeRTKPositionToCSVFile(rtkcsvfile)
    rtksuccessrate = parser.getSuccessRate()
    resp = {
        "code": 0,
        "logfile": filename,
        "gnssjsonfile": gnssjsonfile,
        "adrcsvfile": adrcsvfile,
        "rtkcsvfile": rtkcsvfile,
        "rtksuccessrate": rtksuccessrate
    }
    return resp

def parseDongleLog(filename, filepath):
    targetfile = filename.replace(".log", "")
    parser  = LocationEngineDongleLogParser(filepath + "/" + filename)
    donglecsvfile = filepath + "/" + targetfile + "-dongle.csv"
    parser.writeDonglePositionToCSVFile(donglecsvfile)
    devicecsvfile = filepath + "/" + targetfile + "-device.csv"
    parser.writeDevicePositionToCSVFile(devicecsvfile)
    resp = {
        "code": 0,
        "logfile": filename,
        "donglecsvfile": donglecsvfile,
        "devicecsvfile": devicecsvfile
    }
    return resp

#解析日志
@app.route('/parselog', methods=['POST'])
def parseLog():
    print(request.headers)
    print(request.json)
    filepath = request.json["filepath"]
    loglist = request.json["loglist"]
    respList = []
    for item in loglist:
        logtype = item["logtype"]
        logfile = item["logfile"]
        if logtype == "rtk":
            respList.append(parseRTKlog(logfile, filepath))
        elif logtype == "dongle":
            respList.append(parseDongleLog(logfile, filepath))
    return respList

# 从csv文件中解析groundtruth,文件必须包含"lat", "lon"字段
def parseGroundTruth(gttype, content):
    gt = GroundTruth(gttype, content)
    return gt

def comparePositionWithGroundTruth(reader, truth):
    tslist = reader.getTimestampList()
    ds = DistanceDiff()
    for ts in tslist:
        glat, glon = truth.getGroundTruthByTimeStamp(ts)
        if glat is None or glon is None:
            continue
        clat,clon = reader.getLatLonByTimeStamp(ts)
        if clat is None or clon is None:
            continue
        item = DistanceDiffItem(timestamp=ts, groundTruthLat=glat,groundTruthLon=glon,curentLat=clat, curentLon=clon)
        ds.PushBack(item)
    return ds

# 从图片中解析groundtruth
def compareADRPositionWithGroundTruth(filename, truth):
    reader = ADRCSVReader(filename)
    return comparePositionWithGroundTruth(reader, truth)

def compareSTGPositionWithGroundTruth(filename, truth):
    reader = STGRTKCSVReader(filename)
    return comparePositionWithGroundTruth(reader, truth)

def compareOfflineRTKPositionWithGroundTruth(filename, truth):
    reader = OFFLINERTKCSVReader(filename)
    return comparePositionWithGroundTruth(reader, truth)

@app.route('/parsertkroadtest', methods=['POST'])
def parsertkroadtest():
    print(request.headers)
    print(request.json)
    filepath = request.json["filepath"]
    outputpath = request.json["outputpath"]

    #解析ground truth
    groundtruth         = request.json["groundtruth"]
    groundtruthtype     = groundtruth["type"]
    groundtruthcontent  = groundtruth["content"]
    if groundtruthtype not in ["text", "dongle", "img"]:
        resp = {
            "code": -1,
            "error": "only support type: text, dongle or img"
        }
        return jsonify(resp)
    gt = parseGroundTruth(gttype=groundtruthtype, content=filepath + "/" + groundtruthcontent)
    resp = {
        "code": 0
    }

    #解析安卓的定位结果
    adrcsv = request.json["adrcsv"]
    r = compareADRPositionWithGroundTruth(filepath + "/" + adrcsv, gt)
    if r.Empty() == False:
        centerLat, centerLon = r.getCenterLatLon()
        h = KeplerDraw(filepath, groundtruthcontent, adrcsv,centerLat, centerLon, outputpath)
        r.WriteToCSV(outputpath + "/" + adrcsv.replace(".csv", "-diff.csv"))
        resp[adrcsv.replace(".csv","")] = r.SummaryDistance()
        resp[adrcsv.replace(".csv","")]["html"] = h.WriteToHTML()

    #解析在线rtk的定位结果
    rtkcsv = request.json["rtkcsv"]
    r = compareSTGPositionWithGroundTruth(filepath + "/" + rtkcsv, gt)
    if r.Empty() == False:
        centerLat, centerLon = r.getCenterLatLon()
        h = KeplerDraw(filepath, groundtruthcontent, rtkcsv, centerLat, centerLon, outputpath)
        resp[rtkcsv.replace(".csv","")] = r.SummaryDistance()
        r.WriteToCSV(outputpath + "/" + rtkcsv.replace(".csv", "-diff.csv"))
        resp[rtkcsv.replace(".csv", "")]["html"] = h.WriteToHTML()

    #解析所有的离线定位结果
    offlinertkcsv = request.json["offlinertkcsv"]
    for c in offlinertkcsv:
        f = filepath + "/" + c
        r = compareOfflineRTKPositionWithGroundTruth(f, gt)
        if r.Empty() == False:
            centerLat, centerLon = r.getCenterLatLon()
            h = KeplerDraw(filepath, groundtruthcontent, c, centerLat, centerLon, outputpath)
            resp[c.replace(".csv", "")] = r.SummaryDistance()
            r.WriteToCSV(outputpath + "/" + c.replace(".csv", "-diff.csv"))
            resp[c.replace(".csv", "")]["html"] = h.WriteToHTML()
    return resp

@app.route('/snap', methods=['POST'])
def SnapProxy():
    print(request.headers)
    print(request.json)
    filepath = request.json["filepath"]
    filename = request.json["filename"]
    latitudeColumnName  = request.json["latitude_column_name"]
    longitudeColumnName = request.json["longitude_column_name"]
    timestampColumnName = request.json["timestamp_column_name"]
    file = filepath + "/" + filename
    s = SNAPCSVReader(file)
    if s.OverWriteSnapResultToOriginFile(latitudeColumnName, longitudeColumnName, timestampColumnName) == True:
        resp = {
            "code": 0,
            "message": "overwrite origin file success"
        }
    else:
        resp = {
            "code": 0,
            "message": "overwrite origin file failed"
        }
    return resp

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8000, debug=True)
