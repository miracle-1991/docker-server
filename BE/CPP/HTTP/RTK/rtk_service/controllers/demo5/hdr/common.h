//
// Created by Xiaolong Ji on 2023/2/1.
//

#ifndef RTKTEST_COMMON_H
#define RTKTEST_COMMON_H
#include <string>
#include <vector>
#include "rtklib.h"
#include "cpp_obs.h"

using std::string;
using std::vector;

struct RTKResult {
    int     gpsTime;
    int     state;
    int     validSatNum;
    double  lat;
    double  lng;
    double  high;
    string  humanTime;
};

class Writer {
public:
    int WriteGPSResultToFile(const string filepath, const string filename, const vector<RTKResult>& rArr);
    int WriteOBSToJSONFile(const string filepath, const string filename, const vector<CPPObs>& obsarr);
private:
};

class Reader {
public:
    obs_t* ReadObsDataFromRinex(string rinexfilepath);
    obs_t* ReadObsDataFromJson(string jsonFilePath);
    nav_t* ReadNavDataFromRinex(string filepath);
};

class Printer {
public:
    void PrintObs(obs_t* obs);
private:
    void PrintObsData(const obsd_t& data);
};

#endif //RTKTEST_COMMON_H
