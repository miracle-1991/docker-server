//
// Created by Xiaolong Ji on 2023/2/1.
//

#include "common.h"
#include <fstream>
#include <iostream>
#include "json/json.h"

using std::ofstream;
using std::ifstream;
using std::cerr;
using std::cout;
using std::endl;
using std::ios;
using std::setprecision;

int Writer::WriteGPSResultToFile(const string filepath, const string filename, const vector<RTKResult>& rArr) {
    string path = filepath;
    if (path.back() != '/') {
        path += "/";
    }
    path += filename;
    path += ".csv";

    ofstream fd;
    fd.open(path, ios::out | ios::trunc);
    if (!fd.is_open()) {
        cerr << "failed to open file " << path << endl;
        return -1;
    }

    fd << "gpstime" << "," << "humantime" << "," << "state" << "," << "validsatnum" << "," <<"lat" << "," << "lng" << endl;
    for (int i = 0; i < rArr.size(); i++) {
        const RTKResult& r = rArr[i];
        fd << setprecision(20)
           << r.gpsTime << ","
           << r.humanTime << ","
           << r.state << ","
           << r.validSatNum << ","
           << r.lat << ","
           << r.lng
           << endl;
    }

    fd.close();
}

int Writer::WriteOBSToJSONFile(const string filepath, const string filename, const vector<CPPObs>& obsarr) {
    string path = filepath;
    if (path.back() != '/') {
        path += "/";
    }
    path += filename;
    path += ".json";
    ofstream fd;

    fd.open(path, ios::out | ios::trunc);
    if (!fd.is_open()) {
        cerr << "failed to open file " << path << endl;
        return -1;
    }

    for (const CPPObs o : obsarr) {
        fd << o.ToJsonString() << endl;
    }
    fd.close();
}

obs_t* Reader::ReadObsDataFromRinex(string filepath) {
    sta_t* sta = new sta_t ;
    nav_t* nav = new nav_t ;
    obs_t* obs = new obs_t ;
    obs->n = 0;
    obs->nmax = 0;
    obs->data = nullptr;
    char* opt = "";
    int rcv = 0;
    int stat = ::readrnx(filepath.c_str(), rcv, opt, obs, nav, sta);
    if(stat!=1){
        ::printf("failed to read rinex file: %d", stat);
        return nullptr;
    }
    return obs;
}

obs_t* Reader::ReadObsDataFromJson(std::string jsonFilePath) {
    ifstream t(jsonFilePath);
    if (!t.is_open()) { return nullptr; }
    std::string obsInJson;
    vector<obsd_t> obsDataList;
    while (std::getline(t, obsInJson)) {
        Json::Reader reader;
        Json::Value root;
        if (reader.parse(obsInJson, root)) {
            Json::Value data = root["data"];
            int n = root.get("n", "").asInt();
            for (int i = 0; i < n; i++) {
                obsd_t tempObs;
                Json::Value item = data[i];
                Json::Value t = item["time"];
                tempObs.time.time = t.get("time", "").asInt();
                tempObs.time.sec = t.get("sec","").asDouble();
                tempObs.sat = item.get("sat", "").asInt();
                if (tempObs.sat == 0) {
                    continue; // sat = 0 is invalid
                }
                tempObs.rcv = item.get("rcv", "").asInt();
                Json::Value snrArr = item["snr"];
                int maxArrSize = NFREQ+NEXOBS;
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.SNR[k] = snrArr[k].asInt();
                }
                Json::Value lliArr = item["lli"];
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.LLI[k] = lliArr[k].asInt();
                }
                Json::Value codeArr = item["code"];
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.code[k] = codeArr[k].asInt();
                }
                Json::Value lArr = item["l"];
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.L[k] = lArr[k].asDouble();
                }
                Json::Value pArr = item["p"];
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.P[k] = pArr[k].asDouble();
                }
                Json::Value dArr = item["d"];
                for (int k = 0; k < maxArrSize; k++) {
                    tempObs.D[k] = dArr[k].asDouble();
                }
                obsDataList.push_back(tempObs);
            }
        }
    }

    obs_t* dst = new obs_t;
    dst->n = obsDataList.size();
    dst->nmax = 0;
    dst->flag = 0;
    dst->rcvcount = 0;
    dst->tmcount = 0;
    dst->data = new obsd_t[dst->n]{0};
    for (int i = 0; i < dst->n; i++) {
        dst->data[i] = obsDataList[i];
    }

    return dst;
}

nav_t* Reader::ReadNavDataFromRinex(string filepath) {
    sta_t* sta = new sta_t ;
    nav_t* nav = new nav_t ;
    obs_t* obs = new obs_t ;
    char* opt = "";
    int rcv = 0;
    int stat = ::readrnx(filepath.c_str(), rcv, opt, obs, nav, sta);
    if(stat!=1){
        ::printf("failed to read rinex file: %d", stat);
        return nullptr;
    }
    return nav;
}


void Printer::PrintObsData(const obsd_t& data) {
    cout << setprecision(20) << "{time:{" << data.time.time << ",sec:" << data.time.sec
         << "}sat:" << (int)data.sat << ",rcv:" << (int)data.rcv
         << ",snr:[" << (int)data.SNR[0] << "," << (int)data.SNR[1] << ","<< (int)data.SNR[2] << "],"
         << "LLI:[" << (int)data.LLI[0] << "," << (int)data.LLI[1] << "," << (int)data.LLI[2] << "],"
         << "code:[" << (int)data.code[0] << "," << (int)data.code[1] << "," << (int)data.code[2] << "],"
         << "L:[" << data.L[0] << "," << data.L[1] << "," << data.L[2] << "],"
         << "P:[" << data.P[0] << "," << data.P[1] << "," << data.P[2] << "],"
         << "D:[" << data.D[0] << "," << data.D[1] << "," << data.D[2] << "]}"
         << endl;
}

void Printer::PrintObs(obs_t* obs) {
    if (obs == nullptr || obs->n <= 0) { return; }
    for (int i = 0; i < obs->n; i++) {
        PrintObsData(obs->data[i]);
    }
}