//
// Created by Xiaolong Ji on 2023/1/12.
//

#ifndef RTKTEST_INTERFACE_DEBUG_H
#define RTKTEST_INTERFACE_DEBUG_H
#include <string>
#include "rtklib.h"

extern void startOBSStream(const char* stationUrl, int navsys, int rcv);

extern std::string getLatestRoverObsData();
extern std::string getLatestStationObsData();

extern prcopt_t getDefaultDemo5Prcopt(int pmode, int sys);

#endif //RTKTEST_INTERFACE_DEBUG_H
