//
// Created by Xiaolong Ji on 2023/1/16.
//

#ifndef RTKTEST_CPP_OBS_H
#define RTKTEST_CPP_OBS_H
#include <vector>
#include <memory>
#include <string>
#include "rtklib.h"

class CPPObs {
public:
    explicit CPPObs(const std::vector<obsd_t>& obs);
    explicit CPPObs(const std::string& obsInJson);
    explicit CPPObs(const CPPObs& rover, const CPPObs& base);
    std::string ToJsonString() const;
    std::shared_ptr<obs_t> GetCopy();
    bool isValid();
private:
    void initFromObsdVec(const std::vector<obsd_t>& obs);
private:
    int n;
    int nmax;
    int flag;
    int rcvcount;
    int tmcount;
    std::vector<obsd_t> data;
};


#endif //RTKTEST_CPP_OBS_H
