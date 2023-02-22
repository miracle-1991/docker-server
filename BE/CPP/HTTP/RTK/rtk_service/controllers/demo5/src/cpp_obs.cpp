//
// Created by Xiaolong Ji on 2023/1/16.
//

#include "cpp_obs.h"
#include "json/json.h"
#include <algorithm>
#include <regex>

std::string removeSpace(std::string input) {
    std::regex newlines_re(" ");
    return std::regex_replace(input, newlines_re, "");
}

std::string removeNewLineFlag(std::string input) {
    std::regex newlines_re("\n+");
    return std::regex_replace(input, newlines_re, "");
}

void CPPObs::initFromObsdVec(const std::vector<obsd_t>& obs) {
    int size = obs.size();
    n = size;
    nmax = 0;
    flag = 0;
    rcvcount = 0;
    tmcount = 0;
    data.insert(data.end(), obs.begin(), obs.end());
}

CPPObs::CPPObs(const std::vector<obsd_t>& obs) {
    initFromObsdVec(obs);
}

Json::Value gtimeToJson(gtime_t t) {
    Json::Value root;
    root["time"] = (unsigned int )t.time;
    root["sec"] = t.sec;
    return root;
}

Json::Value uint16ArrToJson(const uint16_t* arr, int n) {
    Json::Value datalist;
    for (int i = 0; i < n; i++) {
        Json::Value item = arr[i];
        datalist.append(item);
    }
    return datalist;
}

Json::Value uint8ArrToJson(const uint8_t* arr, int n) {
    Json::Value datalist;
    for (int i = 0; i < n; i++) {
        Json::Value item = arr[i];
        datalist.append(item);
    }
    return datalist;
}

Json::Value doubleArrToJson(const double * arr, int n) {
    Json::Value datalist;
    for (int i = 0; i < n; i++) {
        Json::Value item = arr[i];
        datalist.append(item);
    }
    return datalist;
}

Json::Value floatArrToJson(const float * arr, int n) {
    Json::Value datalist;
    for (int i = 0; i < n; i++) {
        Json::Value item = arr[i];
        datalist.append(item);
    }
    return datalist;
}

Json::Value obsdToJson(const obsd_t& obsd) {
    Json::Value root;
    root["time"] = gtimeToJson(obsd.time);
    root["sat"] = obsd.sat;
    root["rcv"] = obsd.rcv;
    int size = NFREQ+NEXOBS;
    root["snr"] = uint16ArrToJson(obsd.SNR, size);
    root["lli"] = uint8ArrToJson(obsd.LLI, size);
    root["code"] = uint8ArrToJson(obsd.code, size);
    root["l"] = doubleArrToJson(obsd.L, size);
    root["p"] = doubleArrToJson(obsd.P, size);
    root["d"] = floatArrToJson(obsd.D, size);
    return root;
}

std::string CPPObs::ToJsonString() const {
    if (n <= 0 ) { return ""; }

    Json::Value root;
    Json::Value dataList;
    root["n"]       = n;
    root["nmax"]    = nmax;
    for (int i = 0; i < n; i++) {
        const obsd_t& tmpdata   = data[i];
        Json::Value dataItem    = obsdToJson(tmpdata);
        dataList.append(dataItem);
    }
    root["data"] = dataList;
    std::string structedstr = root.toStyledString();
    std::string ret = removeNewLineFlag(removeSpace(structedstr));
    return ret;
}

bool CPPObs::isValid() {
    if (n <= 0) { return false; }
    if (data.empty()) { return false; }
    return true;
}

void initObsd(obsd_t& obsdata) {
    obsdata.time.time   = 0;
    obsdata.time.sec    = 0;
    obsdata.sat         = 0;
    obsdata.rcv         = 0;
    obsdata.timevalid   = 0;
    obsdata.eventime.time=0;
    obsdata.eventime.sec= 0;
    obsdata.freq        = 0;
    for (int j = 0; j < NFREQ+NEXOBS; j++) {
        obsdata.SNR[j] = 0;
        obsdata.LLI[j] = 0;
        obsdata.code[j] = 0;
        obsdata.L[j] = 0;
        obsdata.P[j] = 0;
        obsdata.D[j] = 0;
        obsdata.Lstd[j] = 0;
        obsdata.Pstd[j] = 0;
    }
}

CPPObs::CPPObs(const std::string& obsInJson) {
    Json::Reader reader;
    Json::Value root;
    if (reader.parse(obsInJson, root)) {
        n = root["n"].asInt();
        nmax = root["nmax"].asInt();
        flag = 0;
        rcvcount = 0;
        tmcount = 0;
        Json::Value rawdata = root["data"];
        for (int i = 0; i < n; i++) {
            obsd_t tmpdata;
            initObsd(tmpdata);
            Json::Value item       = rawdata[i];
            Json::Value t          = item["time"];
            tmpdata.time.time      = t["time"].asUInt();
            tmpdata.time.sec       = t["sec"].asDouble();

            tmpdata.sat            = item["sat"].asUInt();
            tmpdata.rcv            = item["rcv"].asUInt();
            Json::Value snrArr     = item["snr"];
            int maxArrSize = NFREQ+NEXOBS;
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.SNR[k] = snrArr[k].asInt();
            }
            Json::Value lliArr = item["lli"];
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.LLI[k] = lliArr[k].asInt();
            }
            Json::Value codeArr = item["code"];
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.code[k] = codeArr[k].asUInt();
            }
            Json::Value lArr = item["l"];
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.L[k] = lArr[k].asDouble();
            }
            Json::Value pArr = item["p"];
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.P[k] = pArr[k].asDouble();
            }
            Json::Value dArr = item["d"];
            for (int k = 0; k < maxArrSize; k++) {
                tmpdata.D[k] = dArr[k].asDouble();
            }
            data.push_back(tmpdata);
        }
    }
}

bool sortObsd(const obsd_t& left, const obsd_t& right) {
    if (left.rcv        < right.rcv)       { return true; }
    if (left.rcv        > right.rcv)       { return false;}

    if (left.sat        < right.sat)       { return true; }
    if (left.sat        > right.sat)       { return false;}

    if (left.time.time  < right.time.time) { return true; }
    if (left.time.time  > right.time.time) { return false;}

    return false;
}

CPPObs::CPPObs(const CPPObs& rover, const CPPObs& base) {
    std::vector<obsd_t> obslist;
    for (int i = 0; i < rover.n; i++) {
        if (rover.data[i].time.time == 0 || rover.data[i].sat == 0) { continue; }
        obslist.push_back(rover.data[i]);
        obslist.back().rcv = 1;
    }
    for (int j = 0; j < base.n; j++) {
        if (base.data[j].time.time == 0 || base.data[j].sat == 0) { continue; }
        obslist.push_back(base.data[j]);
        obslist.back().rcv = 2;
    }

    std::sort(obslist.begin(), obslist.end(), sortObsd);
    initFromObsdVec(obslist);
}

std::shared_ptr<obs_t> CPPObs::GetCopy() {
    std::shared_ptr<obs_t> copy(new obs_t, [](obs_t* obs) {
        if (obs == nullptr) { return; }
        if (obs->data != nullptr) { delete[] obs->data; obs->data = nullptr; }
        delete obs;
        obs = nullptr;
    });

    copy->n     = n;
    copy->nmax  = nmax;
    copy->flag  = flag;
    copy->rcvcount  = rcvcount;
    copy->tmcount   = tmcount;
    copy->data = new obsd_t[n];
    for (int i = 0; i < n; i++) {
        copy->data[i] = data[i];
    }
    return copy;
}