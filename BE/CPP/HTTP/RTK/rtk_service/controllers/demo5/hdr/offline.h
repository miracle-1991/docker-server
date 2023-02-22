//
// Created by Xiaolong Ji on 2023/2/1.
//

#ifndef RTKTEST_OFFLINE_H
#define RTKTEST_OFFLINE_H
#include <string>
#include "rtklib.h"
#include "cpp_obs.h"

using std::string;
using std::vector;

class Offline {
public:
    int RTK(int argc, char *argv[]);
    int RTK(string adr_file, string station_obs_file, string station_nav_file, string outputpath, int pmode, int sys, string& outfile);
private:
    vector<CPPObs> collectObsDataToArray(const obs_t* obs, int rcv);
    void deleteObs(obs_t* ptr);
    void deleteNav(nav_t* ptr);
    obs_t* getStationObs(std::string stationObsFile);
    nav_t* getStationNav(std::string stationNavFile);
private:
    std::pair<std::string, obs_t*> m_station_obs = std::make_pair("", nullptr);
    std::pair<std::string, nav_t*> m_station_nav = std::make_pair("", nullptr);
};


#endif //RTKTEST_OFFLINE_H
