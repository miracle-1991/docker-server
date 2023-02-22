//
// Created by Xiaolong Ji on 2023/2/16.
//

#ifndef RTK_SERVICE_RTKRUNARG_H
#define RTK_SERVICE_RTKRUNARG_H
#include <string_view>
#include <json/json.h>

class RTKRunArg {
public:
    std::string Init(std::shared_ptr<Json::Value> jsonObj);
    std::vector<int> GetPModeList() const;
    std::vector<int> GetNavSysList() const;
    std::string GetRoverObs() const;
    std::string GetStationObs() const;
    std::string GetStationNav() const;
    std::string GetOutputPath() const;
private:
    int transPmode(std::string pmodestr);
    int transNavSys(std::string navsysstr);
private:
    std::vector<int> m_pmode;
    std::vector<int> m_navsys;
    std::string m_rover_obs;
    std::string m_station_obs;
    std::string m_station_nav;
    std::string m_outputpath;
};


#endif //RTK_SERVICE_RTKRUNARG_H
