//
// Created by Xiaolong Ji on 2023/2/16.
//

#include "RTKRunArg.h"
#include <iostream>
#include "rtklib.h"

std::string RTKRunArg::Init(std::shared_ptr<Json::Value> jsonObj) {
    if (jsonObj == nullptr) { return "empty input"; }

    auto rtk = jsonObj->get("rtk", "");
    auto pmode = rtk.get("pmode", "");
    auto navsys = rtk.get("navsys", "");
    auto obsnav = rtk.get("obsnav", "");

    //解析所有的模式
    std::vector<int> pmodelist;
    for (int i = 0; i < pmode.size(); i++) {
        auto m = pmode.get(i, "").asString();
        int im = transPmode(m);
        if (im == -1) { return "Parse pmode failed";}
        pmodelist.push_back(im);
    }
    m_pmode = pmodelist;

    //解析所有的卫星系统
    std::vector<int> navsyslist;
    for (int i = 0; i < navsys.size(); i++) {
        auto s = navsys.get(i, "").asString();
        int is = transNavSys(s);
        if (is == -1) { return "Parse navsys failed"; }
        navsyslist.push_back(is);
    }
    m_navsys = navsyslist;

    //用户的obs文件
    auto roverobs = obsnav.get("roverobs", "").asString();
    if (roverobs.empty()) { return "Parse roverobs failed"; }
    m_rover_obs = roverobs;

    //基站的obs文件
    auto stationobs = obsnav.get("stationobs", "").asString();
    if (stationobs.empty()) { return "Parse stationobs failed"; }
    m_station_obs = stationobs;

    //基站的星历文件
    auto stationnav = obsnav.get("stationnav","").asString();
    if (stationnav.empty()) { return "Parse stationnav failed"; }
    m_station_nav = stationnav;

    auto outputpath = jsonObj->get("outputpath", "").asString();
    if (outputpath.empty()) { return "Parse outputpath failed"; }
    m_outputpath = outputpath;
    return "";
}

int RTKRunArg::transPmode(std::string pmodestr) {
    if (pmodestr.empty()) { return -1; }
    if (pmodestr == "spp")    { return PMODE_SINGLE;}
    if (pmodestr == "dgps")   { return PMODE_DGPS;}
    if (pmodestr == "kinema") { return PMODE_KINEMA;}
    if (pmodestr == "static") { return PMODE_STATIC;}
    return -1;
}

int RTKRunArg::transNavSys(std::string navsysstr) {
    if (navsysstr.empty()) { return -1; }
    int navsys = 0;
    if (navsysstr.find("gps") != std::string::npos) { navsys |= SYS_GPS; }
    if (navsysstr.find("gal") != std::string::npos) { navsys |= SYS_GAL; }
    if (navsysstr.find("glo") != std::string::npos) { navsys |= SYS_GLO; }
    if (navsysstr.find("cmp") != std::string::npos) { navsys |= SYS_CMP; }
    return navsys == 0 ? -1 : navsys;
}

std::vector<int> RTKRunArg::GetPModeList() const { return m_pmode; }
std::vector<int> RTKRunArg::GetNavSysList() const { return m_navsys; }
std::string RTKRunArg::GetRoverObs() const { return m_rover_obs; }
std::string RTKRunArg::GetStationObs() const { return m_station_obs; }
std::string RTKRunArg::GetStationNav() const { return m_station_nav; }
std::string RTKRunArg::GetOutputPath() const { return m_outputpath; }