#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import gzip
import os
import time

#开始此脚本之前，请确保已经配置好了awscli
#配置awscli请参考: https://wiki.grab.com/display/NET/AWS+CLI%2C+Boto3+and+Web+Console+Access+with+AWS+SSO+and+Cisco+Duo+MFA#expand-Screenshot
import boto3
import botocore
from readerwriterlock import rwlock
locker = rwlock.RWLockFair()

class AWSS3:
    #初始化，需要提供awscli配置后~/.aws/config中的用户名
    def __init__(self, user, bucket):
        boto3.setup_default_session(profile_name=user)
        self.bucket = bucket
        self.clientlow = boto3.client("s3")
        self.clienthigh = boto3.resource("s3")
        self.totalCntOfFiles = 0
        self.curCntOfFiles = 0
        self.processingstate=""
        self.writeLock = locker.gen_wlock()
        self.readerLock = locker.gen_rlock()

    #获取某个前缀下的文件名，前缀prefix必须精确到小时，比如location-engine/app/2023/01/18/10
    def get_filenames(self, prefix):
        """Get a list of keys in an S3 bucket."""
        if prefix.endswith("/") == False:
            prefix += "/"
        files = []
        resp = self.clientlow.list_objects(Bucket=self.bucket, Prefix=prefix)
        if 'Contents' in resp:
            for obj in resp['Contents']:
                if 'Key' in obj:
                    files.append(os.path.basename(obj['Key']))
        return files

    def __binary_search__(self, filelist, target):
        left = 0
        right = len(filelist) - 1
        while left <= right:
            mid = int((left + right)/2)
            if filelist[mid]["time"] <= target:
                left = mid + 1
            else:
                right = mid - 1
        return left

    def __get_log_files_by_time__(self, prefix, starttime, endtime):
        startyear, startmonth, startday, starthour, startminu, startsec= time.strftime("%Y:%m:%d:%H:%M:%S", starttime).split(":")
        endyear, endmonth, endday, endhour, endminu, endsec= time.strftime("%Y:%m:%d:%H:%M:%S", endtime).split(":")
        if startyear != endyear or startmonth != endmonth or startday != endday:
            return None
        prefix += "/" + startyear + "/" + startmonth + "/" + startday
        hourlist = self.get_subdir(prefix)
        fhmslist = []
        for h in hourlist:
            fileprefix = prefix + "/" + h
            flist = self.get_filenames(fileprefix)
            for f in flist:
                if f.endswith("err.gz"):
                    continue
                ftime = f.split("_")[1].split("-")
                fhms = ftime[len(ftime) - 1]
                fhmslist.append({ "time": fhms, "prefix": fileprefix, "name": f })
        fhmslist.sort(key=lambda x : x["time"])
        s = starthour + startminu + startsec
        e = endhour + endminu + endsec
        startindex = self.__binary_search__(fhmslist, s)
        endindex = self.__binary_search__(fhmslist, e)
        return fhmslist[startindex: endindex+1]

    def __run_filter_list__(self, filterlist, line):
        for fstr in filterlist:
            if line.find(fstr) == -1:
                return False
        return True

    def __parse_logtime__(self, line):
        #2023-02-01T09:18:20.603666
        timestr = str(line[:26])
        timestr = timestr.replace("T", " ")
        timelocal = time.strptime(timestr, "%Y-%m-%d %H:%M:%S.%f")
        return timelocal

    #下载某个时间范围内所有的日志文件，并按照过滤器filter过滤出感兴趣的日志，统一输出到targetfile中,过滤器会在日志中的每一行上运行
    #前缀prefix不能包含日期，比如location-engine/app
    #开始时间starttime和结束时间endtime必须是字符串，形式为"%Y/%m/%d %H:%M:%S"
    def pull_log_into_file(self, prefix, starttimestr, endtimestr, targetfile, filterlist):
        self.__setProcessingStart__()
        starttime = time.strptime(starttimestr, "%Y-%m-%d %H:%M:%S")
        endtime = time.strptime(endtimestr, "%Y-%m-%d %H:%M:%S")
        flist = self.__get_log_files_by_time__(prefix, starttime, endtime)
        self.__setProcessingFileTotalCnt__(len(flist))
        with open(targetfile, "w") as outfile:
            tslist = []
            for f in flist:
                print(f)
                print("download file: " + f["name"])
                self.download_file(f["prefix"], f["name"])
                print("unzip file: " + f["name"])
                gf = gzip.GzipFile(f["name"])
                print("parse file: " + f["name"])
                bf = gf.read()
                s = bytes.decode(bf)
                sl = s.split("\n")

                if filterlist is None or len(filterlist) == 0:
                    for ts in sl:
                        tslist.append({
                            "time": self.__parse_logtime__(ts),
                            "data": ts
                        })
                else:
                    for ts in sl:
                        if self.__run_filter_list__(filterlist, ts):
                            tslist.append({
                                "time": self.__parse_logtime__(ts),
                                "data": ts
                            })
                print("rm file: " + f["name"])
                os.remove(f["name"])
                self.__addProcessingFileCnt__()
            #sort
            tslist = sorted(tslist, key=lambda x : x["time"])
            for ts in tslist:
                outfile.write(ts["data"] + "\n")
        self.__setProcessingEnd__()

    def __setProcessingEnd__(self):
        self.writeLock.acquire()
        self.processingstate = "end"
        self.writeLock.release()

    def __setProcessingStart__(self):
        self.writeLock.acquire()
        self.curCntOfFiles = 0
        self.totalCntOfFiles = 0
        self.processingstate = "runing"
        self.writeLock.release()

    def __setProcessingFileTotalCnt__(self, cnt):
        self.writeLock.acquire()
        self.totalCntOfFiles = cnt
        self.writeLock.release()

    def __addProcessingFileCnt__(self):
        self.writeLock.acquire()
        self.curCntOfFiles += 1
        self.writeLock.release()

    def getProcessingFilesCnt(self):
        self.readerLock.acquire()
        t, c, s = self.totalCntOfFiles, self.curCntOfFiles, self.processingstate
        self.readerLock.release()
        return t, c, s

    #获取某个前缀下所有的子文件夹名称,前缀比如location-engine/app/2023/01/18
    def get_subdir(self, prefix):
        if prefix.endswith("/") == False:
            prefix += "/"
        paths = []
        resp = self.clientlow.list_objects(Bucket=self.bucket, Prefix=prefix)
        if 'Contents' in resp:
            for obj in resp['Contents']:
                if 'Key' in obj:
                    paths.append(os.path.dirname(obj['Key']))
        prefixSplit = prefix.split("/")
        subdirList = []
        for p in paths:
            sonSplit = p.split("/")
            son = sonSplit[len(prefixSplit)-1]
            subdirList.append(son)
        uniqueSubdirList = []
        for item in subdirList:
            if item not in uniqueSubdirList:
                uniqueSubdirList.append(item)
        return uniqueSubdirList

    #判断某个路径是否存在
    def is_exist(self, path) -> bool:
        try:
            customPath = '/'.join(path.split("s3://")[1].split("/")[1:])
            customBucket = path.split("s3://")[1].split("/")[0]
            self.clienthigh.Object(customBucket, customPath).load()
        except botocore.exceptions.ClientError as e:
            if e.response['Error']['Code'] == "404":
                return False
            else:
                raise Exception(e)
        else:
            return True

    #下载某个文件，前缀要精确到目标文件的最后一级目录
    def download_file(self, prefix, file_name):
        path = "s3://" + self.bucket + "/" + prefix + "/" + file_name
        if not self.is_exist(path):
            print("The object does not exist: ", path)
            return

        try:
            key = prefix + "/" + file_name
            self.clienthigh.Bucket(self.bucket).download_file(key, file_name)
        except botocore.exceptions.ClientError as e:
            if e.response['Error']['Code'] == "404":
                print("The object does not exist when download.")
            else:
                raise Exception(e)

if __name__ == '__main__':
    # 使用方式:
    # 1、首先配置awscli : https://wiki.grab.com/display/NET/AWS+CLI%2C+Boto3+and+Web+Console+Access+with+AWS+SSO+and+Cisco+Duo+MFA#expand-Screenshot
    # 2、从~/.aws/config复制用户名，作为AWSS3的第一个参数，如果出现有关SSO登录的错误提示，请执行aws sso login --profile ssm-user重新登录
    # 3、指定时间范围，和输出文件的名称(带路径), 指定filer，只选取想要的日志
    user = "ssm-user"
    bucket = "grabtaxi-logs-stg"
    awss3 = AWSS3(user, bucket)
    starttimestr = "2023/02/15 07:54:10"
    endtimestr = "2023/02/15 08:21:19"
    targetfile = "/Users/xiaolong.ji/Downloads/rtK/20230215/s3-rtkFilter-13770990.log"
    filterlist = ["rtkFilter", "driverID:13770990"]
    awss3.pull_log_into_file("location-vendor/app", starttimestr, endtimestr, targetfile, filterlist)
