#ifndef  _RNX2RTCM_H_
#define _RNX2RTCM_H_

#include <vector>
#include "rtklib.h"

using namespace std;
/* obs2rtcm3 将obs数据转化为rtcm3格式的数据
 * */
extern int obs2rtcm3(obs_t* obs, vector<unsigned char>& buff);

#endif