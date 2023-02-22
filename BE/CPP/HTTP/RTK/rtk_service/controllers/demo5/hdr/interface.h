#ifndef _INTERFACE_H_
#define _INTERFACE_H_

#ifdef __cplusplus
extern "C" {
#endif

/* init 初始化rtk
 * input:
 *  module: 模式类型, single:0，dgps:1 kinematic:2
 *  navsys：系统坐标类型
 *  inpstrPath: 基站路径，不可以为空,至少是"", 例如: username:password@Igs-ip.net:2101/SIN100SGP0
 * output:
 *  -1 初始化失败, 0 初始化成功
 * */
extern int init(const char *baseStationUrl, const char *navStationUrl, int module, int navsys);

/* openTrace 打开trace，默认关闭，打开后会在可执行文件所在的目录下生成以.trace结尾的文件
 * input:
 *  最小的trace level，level范围1-5,5是最低，1是最高
 * */
extern void openTrace(int minLevel);

/* inputUserObsData 输入用户的obs数据
 * input:
 * obsInJson 用户obs数据,以json字符串的方式进行组织
 * output:
 * -2 入参错误 -1 内部调用失败 0 内部调用成功 1 server线程未启动
 * */
extern int rtk(const char* obsInJson, int pmode,const int sys,
               double stationX, double stationY, double stationZ,
               int* solState, int* validSatNum, long* gpstime, double *lat, double *lng, double *high);

#ifdef __cplusplus
}
#endif

#endif // HELLO_H_
