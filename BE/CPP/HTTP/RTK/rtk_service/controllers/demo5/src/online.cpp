//
// Created by Xiaolong Ji on 2023/2/1.
//

#include "online.h"
#include "rtklib.h"
#include "interface.h"
#include "interface-debug.h"
#include "cmdline.h"
#include <string>
#include <fstream>
#include <sstream>
#include <iostream>

using std::string;
using std::ifstream;
using std::stringstream;
using std::cerr;
using std::endl;
using std::cout;

string parseObsDataFromFile(string jsonFilePath) {
    ifstream t(jsonFilePath);
    if (!t.is_open()) { return ""; }
    stringstream buffer;
    buffer << t.rdbuf();
    return buffer.str();
}

int Online::RTK(string roverObsStationPath, string baseObsStationPath, string navStationPath, int pmode, int navsys) {
    openTrace(5);
    int ok = init(baseObsStationPath.c_str(), navStationPath.c_str(), pmode, navsys);
    if (ok != 0) {
        cerr << "init rtk failed" << endl;
        return -1;
    }

    // recv rover obs
    startOBSStream(roverObsStationPath.c_str(), navsys, 1);
    while (true) {
        sleepms(1000);
        string obsdata = getLatestStationObsData();
        if (obsdata.size() == 0) { continue; }
        double x = -1507971.0500;
        double y = 6195614.1100;
        double z = 148487.8700;
        int solstate = 0;
        int validSatNum = 0;
        long gpstime = 0;
        double lat = 0;
        double lng = 0;
        double high = 0;
        ::time_t timep;
        struct tm *p;
        ::time(&timep);
        int ok = rtk(obsdata.c_str(), pmode, navsys,
                     x, y, z,
                     &solstate, &validSatNum,
                     &gpstime, &lat, &lng, &high);
        if (ok != 0) {
            cerr << "call rtk failed" << endl;
        }else {
            cout << "state: " << solstate
                 << ", validSatNum: " << validSatNum
                 << ", utc: " << timep
                 << ", gpstime: " << gpstime
                 << ", lat: " << lat
                 << ", lon: " << lng
                 << ", high: "<< high
                 << endl;
        }
    }
}

int Online::RTK(int argc, char **argv) {
    cmdline::parser cmd;
    // -n选项是离线独有，online不能占用
    cmd.add<string>("phoneobsfile", 'f', "obs rinex/json file from phone", false, ""); //手机观测数据文件
    cmd.add<string>("phoneobsurl", 'x', "obs rinex/json file from phone", false, "huxiao1224:huxiao1224@Igs-ip.net:2101/SIN100SGP0"); //手机观测数据文件
    cmd.add<string>("stationobsurl", 'y', "obs rinex/json file from phone", false, "huxiao1224:huxiao1224@Igs-ip.net:2101/SIN100SGP0"); //基站观测数据文件
    cmd.add<string>("stationnavurl", 'z', "obs rinex/json file from phone", false, "huxiao1224:huxiao1224@products.igs-ip.net:2101/BCEP00BKG0"); //基站星历数据文件

    cmd.add<int>("mode",'m', "rtk mode, PMODE_SINGLE:0 PMODE_DGPS:1 PMODE_KINEMA:2 PMODE_STATIC:3", true, 1);
    cmd.add<int>("sys",'s',"navigation system, SYS_GPS:0x01, SYS_GLO:0x04, SYS_GAL:0x08, SYS_GPS|SYS_GLOSYS_GAL=13",true, 1);
    cmd.parse_check(argc, argv);

    int pmode = cmd.get<int>("mode");
    int navsys = cmd.get<int>("sys");
    string phone_obs_file = cmd.get<string>("phoneobsfile");
    string phone_obs_url = cmd.get<string>("phoneobsurl");
    string station_obs_url = cmd.get<string>("stationobsurl");
    string station_nav_url = cmd.get<string>("stationnavurl");

    string roverObsStationPath, baseObsStationPath, navStationPath;
    if (!phone_obs_file.empty()) {
        //从文件解析 TODO
    } else if (!phone_obs_url.empty()) {
        //从url中解析
        return RTK(phone_obs_url, station_obs_url, station_nav_url, pmode, navsys);
    }
}