//
// Created by Xiaolong Ji on 2023/2/1.
//

#include "offline.h"
#include "cmdline.h"
#include "rtklib.h"
#include "interface-debug.h"
#include "interface.h"
#include "json/json.h"
#include "common.h"
#include <string>
#include <iostream>
#include <algorithm>
#include <vector>

using std::string;
using std::ifstream;
using std::vector;
using std::cout;
using std::endl;
using std::pair;

std::map<::time_t, vector<obsd_t>> collectObsDataToMap(obs_t* obs, int rcv) {
    std::map<::time_t, vector<obsd_t>> obsMap;
    if (obs == nullptr || obs-> n <= 0) { return obsMap; }
    int size = obs->n;
    for (int i = 0; i < size; i++) {
        obsd_t rawOBSD = obs->data[i];
        rawOBSD.rcv = rcv;
        ::time_t key = rawOBSD.time.time;
        obsMap[key].push_back(rawOBSD);
    }
    return obsMap;
}

obs_t* Offline::getStationObs(std::string stationObsFile) {
    std::string fname = m_station_obs.first;
    if (fname != stationObsFile) {
        obs_t* obsptr = m_station_obs.second;
        deleteObs(obsptr);
        Reader r;
        obs_t* obs = r.ReadObsDataFromRinex(stationObsFile);
        m_station_obs.first  = stationObsFile;
        m_station_obs.second = obs;
    }
    return m_station_obs.second;
}

nav_t* Offline::getStationNav(std::string stationNavFile) {
    std::string fname = m_station_nav.first;
    if (fname != stationNavFile) {
        nav_t* navptr = m_station_nav.second;
        deleteNav(navptr);
        Reader r;
        nav_t* nav = r.ReadNavDataFromRinex(stationNavFile);
        m_station_nav.first = stationNavFile;
        m_station_nav.second = nav;
    }
    return m_station_nav.second;
}

void Offline::deleteObs(obs_t *ptr) {
    if (ptr == nullptr) { return; }
    freeobs(ptr);
}

void Offline::deleteNav(nav_t* ptr) {
    if (ptr == nullptr) { return; }
    freenav(ptr, 0xFF);
}

std::vector<CPPObs> Offline::collectObsDataToArray(const obs_t* obs, int rcv) {
    std::map<::time_t, vector<obsd_t>> obsMap;
    if (obs == nullptr || obs-> n <= 0) { return std::vector<CPPObs>{}; }
    int size = obs->n;
    vector<::time_t> tarr;
    for (int i = 0; i < size; i++) {
        obsd_t* rawOBSD = &obs->data[i];
        rawOBSD->rcv = rcv;
        ::time_t key = rawOBSD->time.time;
        obsMap[key].push_back(*rawOBSD);
        tarr.push_back(key);
    }
    std::sort(tarr.begin(), tarr.end());
    tarr.erase(unique(tarr.begin(), tarr.end()), tarr.end());
    vector<CPPObs> obsarr;
    for (::time_t t : tarr) {
        vector<obsd_t>& obsdvec = obsMap.at(t);
        std::sort(obsdvec.begin(), obsdvec.end(), [](const obsd_t& left, const obsd_t& right) ->bool {
            if (left.sat < right.sat) { return true; }
            if (left.time.sec < right.time.sec) { return true; }
            return false;
        });
        CPPObs cppobs{obsdvec};
        obsarr.push_back(cppobs);
    }
    return obsarr;
}

vector<pair<::time_t, vector<obsd_t>>> mergeObsDataFromUserAndStation(obs_t* userData, obs_t* stationData) {
    //rcv=1 rover, rcv=2 reference sorted by receiver and satellte
    std::map<::time_t, vector<obsd_t>> userDataMap = collectObsDataToMap(userData, 1);
    std::map<::time_t, vector<obsd_t>> stationDataMap = collectObsDataToMap(stationData, 2);

    auto sortFuncForEveryMinute = [](std::map<::time_t, vector<obsd_t>>& obs) {
        for (auto iter = obs.begin(); iter != obs.end(); iter++) {
            vector<obsd_t>& val = iter->second;
            std::sort(val.begin(), val.end(), [](const obsd_t& left, const obsd_t& right) ->bool {
                if (left.sat < right.sat) { return true; }
                if (left.time.sec < right.time.sec) { return true; }
                return false;
            });
        }
    };
    sortFuncForEveryMinute(userDataMap);
    sortFuncForEveryMinute(stationDataMap);

    vector<pair<::time_t, vector<obsd_t>>> obsList;
    for (auto userIter = userDataMap.begin(); userIter != userDataMap.end(); userIter++) {
        ::time_t key = userIter->first;
        for (auto stationIter = stationDataMap.begin(); stationIter != stationDataMap.end(); stationIter++) {
            ::time_t beginTime = key - 15;
            ::time_t endTime = key + 15;
            ::time_t stationTime = stationIter->first;
            if (stationTime > beginTime && stationTime <= endTime) {
                vector<obsd_t>& userDataVec = userIter->second;
                vector<obsd_t>& stationDataVec = stationIter->second;
                userDataVec.insert(userDataVec.end(), stationDataVec.begin(), stationDataVec.end());
                obsList.push_back(std::make_pair(key, userDataVec));
            }
        }
    }
    std::sort(obsList.begin(), obsList.end(), [](const pair<::time_t, vector<obsd_t>>& left, const pair<::time_t, vector<obsd_t>>& right) ->bool {
        if (left.first < right.first) { return true; }
        return false;
    });
    return obsList;
}

void printGPS(rtk_t* rtk, int timestamp) {
    if (rtk == nullptr) { return; }

    ::uint8_t  stat = rtk->sol.stat;
    switch (stat) {
        case SOLQ_NONE:
            ::printf("solution status: no solution\t");
            break;
        case SOLQ_FIX:
            ::printf("solution status: fix\t");
            break;
        case SOLQ_FLOAT:
            ::printf("solution status: float\t");
            break;
        case SOLQ_SBAS:
            ::printf("solution status: SBAS\t");
            break;
        case SOLQ_DGPS:
            ::printf("solution status: DGPS/DGNSS\t");
            break;
        case SOLQ_SINGLE:
            ::printf("solution status: single\t");
            break;
        case SOLQ_PPP:
            ::printf("solution status: PPP\t");
            break;
        case SOLQ_DR:
            ::printf("solution status: dead reconing\t");
            break;
        default:
            ::printf("no solution\t");
            break;
    }

    double rr[3];
    double pos[3];
    matcpy(rr,rtk->sol.rr,3,1);
    ecef2pos(rr,pos);
    pos[0] *= R2D;
    pos[1] *= R2D;
    double lat = pos[0];
    double lng = pos[1];
    double high = pos[2];
    int ns = rtk->sol.ns;
    ::printf("timestamp:%d, time:%s, pos: {%lf,%lf,%lf}, valid sat num: %d\n", timestamp, time_str(rtk->sol.time,0), lat, lng, high, ns);
}

void getRTKResult(const rtk_t* rtk, RTKResult& r) {
    double rr[3];
    double pos[3];
    matcpy(rr,rtk->sol.rr,3,1);
    ecef2pos(rr,pos);
    pos[0] *= R2D;
    pos[1] *= R2D;
    r.lat = pos[0];
    r.lng = pos[1];
    r.high = pos[2];
    r.validSatNum = rtk->sol.ns;
    r.state = rtk->sol.stat;
    r.gpsTime = rtk->sol.time.time;
    r.humanTime = time_str(rtk->sol.time,0);
    return;
}

string getFileNameAccordOption(const prcopt_t& p) {
    string filename = "";
    string mode = "";
    switch (p.mode) {
        case PMODE_SINGLE:
            mode = "PMODE_SINGLE";
            break;
        case PMODE_DGPS:
            mode = "PMODE_DGPS";
            break;
        case PMODE_KINEMA:
            mode = "PMODE_KINEMA";
            break;
        case PMODE_STATIC:
            mode = "PMODE_STATIC";
            break;
        case PMODE_MOVEB:
            mode = "PMODE_MOVEB";
            break;
        case PMODE_FIXED:
            mode = "PMODE_FIXED";
            break;
        case PMODE_PPP_KINEMA:
            mode = "PMODE_PPP_KINEMA";
            break;
        case PMODE_PPP_STATIC:
            mode = "PMODE_PPP_STATIC";
            break;
        case PMODE_STATIC_START:
            mode = "PMODE_STATIC_START";
            break;
        case PMODE_PPP_FIXED:
            mode = "PMODE_PPP_FIXED";
            break;
        default:
            mode = "PMODE_UNKNOWN";
            break;
    }
    string sys = "";
    sys += p.navsys & SYS_GPS ? "_GPS" : "";
    sys += p.navsys & SYS_SBS ? "_SBS" : "";
    sys += p.navsys & SYS_GLO ? "_GLO" : "";
    sys += p.navsys & SYS_GAL ? "_GAL" : "";
    sys += p.navsys & SYS_QZS ? "_QZS" : "";
    sys += p.navsys & SYS_CMP ? "_CMP" : "";
    sys += p.navsys & SYS_IRN ? "_IRN" : "";
    sys += p.navsys & SYS_LEO ? "_LEO" : "";


    filename += mode + sys;
    return filename;
}

int Offline::RTK(string adr_file, string station_obs_file, string station_nav_file, string outputpath, int pmode, int sys, string& outfile) {
    traceopen("rtknavi_%Y%m%d%h%M.trace");
    tracelevel(4);
    prcopt_t prcopt = getDefaultDemo5Prcopt(pmode, sys);
    prcopt.refpos = POSOPT_POS;
    prcopt.rb[0] = -1507971.0500;
    prcopt.rb[1] = 6195614.1100;
    prcopt.rb[2] = 148487.8700;

    rtk_t* rtk = new rtk_t;
    for (int i = 0; i < 6; i++) {
        rtk->rb[i] = 0;
    }
    rtkinit(rtk, &prcopt);

    obs_t* userData = nullptr;
    Reader r;
    if (adr_file.find("json") != string::npos) {
        userData = r.ReadObsDataFromJson(adr_file);
    }else {
        userData = r.ReadObsDataFromRinex(adr_file);
        auto obsarr = collectObsDataToArray(userData, 1);
        Writer w;
        string filepath = outputpath;
        string filename = "rover_obs";
        w.WriteOBSToJSONFile(filepath, filename, obsarr);
    }

    obs_t* stationObsData = getStationObs(station_obs_file);
    nav_t* nav = getStationNav(station_nav_file);
    vector<pair<::time_t, vector<obsd_t>>> obslist = mergeObsDataFromUserAndStation(userData, stationObsData);
    deleteObs(userData);

    vector<RTKResult> rArr;
    for (int i = 0; i < obslist.size(); i++) {
        pair<::time_t, vector<obsd_t>>& item = obslist[i];
        ::time_t timestamp = item.first;
        vector<obsd_t>& itemData = item.second;
        obs_t* obs = new obs_t;
        obs->n = itemData.size();
        obs->data = itemData.data();
        rtkpos(rtk, obs->data, obs->n, nav);
        printGPS(rtk, timestamp);
        RTKResult r;
        getRTKResult(rtk, r);
        rArr.push_back(r);
        delete obs;
    }

    string filePath = outputpath;
    string fileName = getFileNameAccordOption(prcopt);
    Writer w;
    w.WriteGPSResultToFile(filePath, fileName, rArr);
    outfile=fileName;
    return 0;
}

int Offline::RTK(int argc, char *argv[]) {
    cmdline::parser cmd;
    cmd.add<string>("stationnav", 'n', "nav rinex file from station", true, ""); //基站星历数据
    cmd.add<string>("stationobs", 'o', "obs rinex file from station", true, ""); //基站观测数据
    cmd.add<string>("phoneobs", 'p', "obs rinex file from phone", true, ""); //手机观测数据
    cmd.add<string>("logpath",'l',"path to write output", true, ""); //结果输出目录
    cmd.add<int>("mode",'m', "rtk mode, PMODE_SINGLE:0 PMODE_DGPS:1 PMODE_KINEMA:2 PMODE_STATIC:3", true, 1);
    cmd.add<int>("sys",'s',"navigation system, SYS_GPS:0x01, SYS_GLO:0x04, SYS_GAL:0x08, SYS_GPS|SYS_GLOSYS_GAL=13",true, 1);
    cmd.parse_check(argc, argv);

    int pmode = cmd.get<int>("mode");
    int sys = cmd.get<int>("sys");
    string phone_obs_file = cmd.get<string>("phoneobs");
    string station_obs_file = cmd.get<string>("stationobs");
    string station_nav_file = cmd.get<string>("stationnav");
    string output_path = cmd.get<string>("logpath");
    string outfile;
    return RTK(phone_obs_file, station_obs_file, station_nav_file, output_path, pmode, sys, outfile);
}