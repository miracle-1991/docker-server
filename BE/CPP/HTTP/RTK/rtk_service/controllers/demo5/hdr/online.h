//
// Created by Xiaolong Ji on 2023/2/1.
//

#ifndef RTKTEST_ONLINE_H
#define RTKTEST_ONLINE_H
#include <string>

using std::string;

class Online {
public:
    int RTK(int argc, char *argv[]);
private:
    int RTK(string phoneobsurl, string stationobsurl, string stationnavurl, int pmode, int navsys);
};


#endif //RTKTEST_ONLINE_H
