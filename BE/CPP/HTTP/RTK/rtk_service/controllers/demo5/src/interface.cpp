#include <unistd.h>
#include <iostream>
#include <memory>
#include <ctime>
#include <mutex>
#include <algorithm>
#include <thread>
#include "rtklib.h"
#include "interface.h"
#include "cpp_obs.h"

void freeNavFunc(nav_t* nav) {
    if (nav == nullptr) { return; }
    if (nav->eph != nullptr) { delete[] nav->eph; nav->eph = nullptr; }
    if (nav->geph != nullptr) { delete[] nav->geph; nav->geph = nullptr; }
    if (nav->seph != nullptr) { delete[] nav->seph; nav->seph = nullptr; }
    nav->eph = 0;
    nav->geph = 0;
    nav->seph = 0;
    delete nav;
    nav = nullptr;
}

void freeRtkFunc(rtk_t* rtk) {
    if (rtk == nullptr) { return; }
    rtk->nx = 0;
    rtk->na = 0;
    if (rtk->x != nullptr) { delete[] rtk->x; rtk->x = nullptr; }
    if (rtk->P != nullptr) { delete[] rtk->P; rtk->P = nullptr; }
    if (rtk->xa != nullptr) { delete[] rtk->xa; rtk->xa = nullptr; }
    if (rtk->Pa != nullptr) { delete[] rtk->Pa; rtk->Pa = nullptr; }
}

// cache rover obs data
std::vector<std::shared_ptr<CPPObs>> g_rover_obs_list;
std::mutex g_rover_obs_list_lock;
const int g_rover_obs_list_max_size = 1200;

// cache base obs data
std::vector<std::shared_ptr<CPPObs>> g_base_obs_list;
std::mutex g_base_obs_list_lock;
const int g_base_obs_list_max_size = 1200;

// cache rover station pos
double* globalRoverStationpos = new double[3];
std::mutex globalRoverStationposLock;

// cache base station pos
double* globalBaseStationpos = new double[3];
std::mutex globalBaseStationposLock;

// cache nav data
std::shared_ptr<nav_t> globalnav(new nav_t, freeNavFunc);
std::mutex g_nav_lock;

void printObsDataInner(const obsd_t& data) {
    std::cout  << "{time:{" << data.time.time << ",sec:" << data.time.sec
         << "}sat:" << (int)data.sat << ",rcv:" << (int)data.rcv
         << ",snr:[" << (int)data.SNR[0] << "," << (int)data.SNR[1] << ","<< (int)data.SNR[2] << "],"
         << "LLI:[" << (int)data.LLI[0] << "," << (int)data.LLI[1] << "," << (int)data.LLI[2] << "],"
         << "code:[" << (int)data.code[0] << "," << (int)data.code[1] << "," << (int)data.code[2] << "],"
         << "L:[" << data.L[0] << "," << data.L[1] << "," << data.L[2] << "],"
         << "P:[" << data.P[0] << "," << data.P[1] << "," << data.P[2] << "],"
         << "D:[" << data.D[0] << "," << data.D[1] << "," << data.D[2] << "]}"
         << std::endl;
}

void printObsInner(obs_t* obs) {
    if (obs == nullptr || obs->n <= 0) { return; }
    for (int i = 0; i < obs->n; i++) {
        printObsDataInner(obs->data[i]);
    }
}

extern void openTrace(int minLevel) {
    traceopen("rtknavi_%Y%m%d%h%M.trace");
    tracelevel(minLevel);
}

int pushGObsListWithLock(std::vector<obsd_t>& obsdbuf, int rcv) {
    if (obsdbuf.size() == 0) { return 0; }
    //1、sort
    std::sort(obsdbuf.begin(), obsdbuf.end(), [](const obsd_t& left, const obsd_t& right) -> bool {
        if (left.rcv        < right.rcv)       { return true; }
        if (left.rcv        > right.rcv)       { return false;}

        if (left.sat        < right.sat)       { return true; }
        if (left.sat        > right.sat)       { return false;}

        if (left.time.time  < right.time.time) { return true; }
        if (left.time.time  > right.time.time) { return false;}

        return false;
    });

    //2、remove duplicated data
    std::vector<obsd_t> uniqueObsBuf;
    uniqueObsBuf.push_back(obsdbuf.front());
    for (int i = 1; i < obsdbuf.size(); i++) {
        const obsd_t& lastUniqueObs = uniqueObsBuf.back();
        const obsd_t& curObs = obsdbuf[i];
        if (curObs.sat != lastUniqueObs.sat ||
            curObs.rcv != lastUniqueObs.rcv ||
            timediff(curObs.time,lastUniqueObs.time)!=0.0) {
            uniqueObsBuf.push_back(curObs);
        }
    }
    //3、trans to obs
    std::shared_ptr<CPPObs> obs = std::make_shared<CPPObs>(uniqueObsBuf);
//    time_t rawtime;
//    time(&rawtime);
//    struct tm *ptminfo;
//    ptminfo = localtime(&rawtime);
//    printf("[current][rcv%d][utc: %d] %02d-%02d-%02d %02d:%02d:%02d\n",
//           rcv,rawtime, ptminfo->tm_year + 1900, ptminfo->tm_mon + 1, ptminfo->tm_mday,
//           ptminfo->tm_hour, ptminfo->tm_min, ptminfo->tm_sec);

    //4、save
    if (rcv == 1) {
        std::lock_guard<std::mutex> autolock(g_rover_obs_list_lock);
        if (g_rover_obs_list.size() > g_rover_obs_list_max_size) {
            auto startIter = g_rover_obs_list.begin();
            auto stopIter = g_rover_obs_list.begin() + g_rover_obs_list_max_size/2;
            g_rover_obs_list.erase(startIter, stopIter);
        }
        g_rover_obs_list.push_back(obs);
        trace(3, "[globalRoverOBSBuffer] buffer size:%d\n", g_rover_obs_list.size());
    } else if (rcv == 2) {
        std::lock_guard<std::mutex> autolock(g_base_obs_list_lock);
        if (g_base_obs_list.size() > g_base_obs_list_max_size) {
            auto startIter = g_base_obs_list.begin();
            auto stopIter = g_base_obs_list.begin() + g_base_obs_list_max_size/2;
            g_base_obs_list.erase(startIter, stopIter);
        }
        g_base_obs_list.push_back(obs);
        trace(3, "[globalBaseOBSBuffer] buffer size:%d\n", g_base_obs_list.size());
    }
}

std::shared_ptr<CPPObs> popGObsListWithLock(int rcv) {
    if (rcv == 1) {
        std::lock_guard<std::mutex> autolock(g_rover_obs_list_lock);
        if (g_rover_obs_list.empty()) { return nullptr; }
        std::shared_ptr<CPPObs> back = g_rover_obs_list.back();
        g_rover_obs_list.pop_back();
        return back;
    }else if (rcv == 2) {
        std::lock_guard<std::mutex> autolock(g_base_obs_list_lock);
        if (g_base_obs_list.empty()) { return nullptr; }
        std::shared_ptr<CPPObs> back = g_base_obs_list.back();
        g_base_obs_list.pop_back();
        return back;
    }
}

std::shared_ptr<CPPObs> topGObsListWithLock(int rcv) {
    if (rcv == 1) {
        std::lock_guard<std::mutex> autolock(g_rover_obs_list_lock);
        if (g_rover_obs_list.empty()) { return nullptr; }
        std::shared_ptr<CPPObs> back = g_rover_obs_list.back();
        return back;
    }else if (rcv == 2) {
        std::lock_guard<std::mutex> autolock(g_base_obs_list_lock);
        if (g_base_obs_list.empty()) { return nullptr; }
        std::shared_ptr<CPPObs> back = g_base_obs_list.back();
        return back;
    }
}

extern std::string getLatestRoverObsData() {
    std::shared_ptr<CPPObs> obs = topGObsListWithLock(1);
    return obs == nullptr ? "" : obs->ToJsonString();
}

extern std::string getLatestStationObsData() {
    std::shared_ptr<CPPObs> obs = topGObsListWithLock(2);
    return obs == nullptr ? "" : obs->ToJsonString();
}

void updateBaseStationPos(rtcm_t* rtcm, int rcv) {
    if (rtcm != nullptr) { return; }
    sta_t* sta = &rtcm->sta;
    double pos[3]   = {0,0,0};
    double rb[3]    = {0,0,0};
    double del[3]   = {0,0,0};
    double dr[3]    = {0,0,0};
    int i;
    /* update base station position */
    for (i=0;i<3;i++) {
        rb[i]=sta->pos[i];
    }
    /* antenna delta */
    ecef2pos(rb,pos);
    if (sta->deltype) { /* xyz */
        del[2]=sta->hgt;
        enu2ecef(pos,del,dr);
        for (i=0;i<3;i++) {
            rb[i]+=sta->del[i]+dr[i];
        }
    }
    else { /* enu */
        enu2ecef(pos,sta->del,dr);
        for (i=0;i<3;i++) {
            rb[i]+=dr[i];
        }
    }

    if (1 == rcv) {
        trace(4, "[updateRoverStationPos] update rover station pos:%f,%f,%f\n", rb[0],rb[1],rb[2]);
        std::lock_guard<std::mutex> autolock(globalRoverStationposLock);
        for (int i = 0; i < 3; i++) {
            globalRoverStationpos[i] = rb[i];
        }
    }else if (2 == rcv) {
        trace(4, "[updateBaseStationPos] update base station pos:%f,%f,%f\n", rb[0],rb[1],rb[2]);
        std::lock_guard<std::mutex> autolock(globalBaseStationposLock);
        for (int i = 0; i < 3; i++) {
            globalBaseStationpos[i] = rb[i];
        }
    }
}

void updateSSR(rtcm_t* rtcm) {
    int i,sys,prn,iode;

    std::lock_guard<std::mutex> autolock(g_nav_lock);
    for (i=0;i<MAXSAT;i++) {
        if (!rtcm->ssr[i].update) continue;

        /* check consistency between iods of orbit and clock */
        if (rtcm->ssr[i].iod[0] != rtcm->ssr[i].iod[1]) {
            continue;
        }
        rtcm->ssr[i].update=0;

        iode= rtcm->ssr[i].iode;
        sys=satsys(i+1,&prn);

        /* check corresponding ephemeris exists */
        if (sys==SYS_GPS||sys==SYS_GAL||sys==SYS_QZS) {
            if (globalnav->eph[i       ].iode!=iode&&
                globalnav->eph[i+MAXSAT].iode!=iode) {
                continue;
            }
        }
        else if (sys==SYS_GLO) {
            if (globalnav->geph[prn-1          ].iode!=iode&&
                globalnav->geph[prn-1+MAXPRNGLO].iode!=iode) {
                continue;
            }
        }
        globalnav->ssr[i]=rtcm->ssr[i];
    }
}

void getRoverStationPos(double* x, double* y, double* z) {
    std::lock_guard<std::mutex> autolock(globalRoverStationposLock);
    *x = globalRoverStationpos[0];
    *y = globalRoverStationpos[1];
    *z = globalRoverStationpos[2];
}

void getBaseStationPos(double* x, double* y, double* z) {
    std::lock_guard<std::mutex> autolock(globalBaseStationposLock);
    *x = globalBaseStationpos[0];
    *y = globalBaseStationpos[1];
    *z = globalBaseStationpos[2];
}

void obsStreamLoop(std::string path, int navsys, int rcv) {
    stream_t *stream = new stream_t;
    //1、init
    strinit(stream);

    //2、open stream
    int rw = STR_MODE_RW;
    int strcli = STR_NTRIPCLI;
    int ok = stropen(stream, strcli, rw, path.c_str());
    if (!ok) {
        trace(1, "[obsStreamLoop] failed to stropen %s\n", path.c_str());
        return;
    }

    trace(4, "[obsStreamLoop] stropen success\n");
    const int buffsize = 32768;
    uint8_t* buffer = new uint8_t[buffsize];
    uint32_t tick = 0;
    while (true) {
        tick=tickget();
        memset(buffer, 0, sizeof(uint8_t) * buffsize);
        int bytesnum = 0;
        uint8_t* startPos = buffer + bytesnum;
        uint8_t* endPos = buffer + buffsize;
        int n = strread(stream, startPos, endPos-startPos);
        if (n <= 0) {
            continue;
        }

        bytesnum += n;
        rtcm_t* rtcm = new rtcm_t;
        memset(rtcm,0,sizeof(rtcm_t));
        init_rtcm(rtcm);
        gtime_t time=utc2gpst(timeget());
        rtcm->time = time;
        int lastRTCMLen = rtcm->len;
        int curRTCMLen = rtcm->len;

        std::vector<obsd_t> obsdBuf;
        int missedObsCnt = 0;
        auto updateObsBuf = [&obsdBuf, &missedObsCnt, navsys, rcv](rtcm_t* rtcm) {
            int bufSize = obsdBuf.size();
            if (bufSize >= MAXOBSBUF) {
                missedObsCnt++;
                trace(5, "[obsStreamLoop] bigger than buffer size, pass cnt:%d\n", missedObsCnt);
                return;
            }
            const obs_t& obs = rtcm->obs;
            if (obs.n <= 0) {
                trace(5, "[obsStreamLoop] no obs data, cnt:%d\n", obs.n);
                return;
            }

            for (int i = 0; i < obs.n; i++) {
                int sat = obs.data[i].sat;
                int sys = satsys(sat, NULL);
                if ((sys & navsys) == 0) { continue; }
                obs.data[i].rcv = rcv;
                obsdBuf.push_back(obs.data[i]);
            }
            int size = obsdBuf.size();
            trace(4, "[obsStreamLoop] obs buffer size:%d\n", size);
            return;
        };

        for (int i = 0; i < bytesnum; i++) {
            curRTCMLen = rtcm->len;
            if (curRTCMLen != lastRTCMLen) {
                lastRTCMLen = curRTCMLen;
            }
            int ret = input_rtcm3(rtcm, buffer[i]);
            if (1 == ret) {
                /* observation data */
                updateObsBuf(rtcm);
            } else if (5 == ret) { /* antenna postion */
                updateBaseStationPos(rtcm, rcv);
            } else{
                missedObsCnt++;
            }
        }

        if (obsdBuf.size() > 0) {
            pushGObsListWithLock(obsdBuf, rcv);
        }
        int cputime= tickget()-tick;
        sleepms(1000 - cputime);
        free_rtcm(rtcm);
        delete rtcm;
    }
}

extern void startOBSStream(const char* stationUrl, int navsys, int rcv) {
    trace(2, "startStream\n");
    std::string path{stationUrl};
    std::thread loop(obsStreamLoop, path, navsys, rcv);
    loop.detach();
}

void initNav(std::shared_ptr<nav_t> nav) {
    if (nav == nullptr) {
        return;
    }
    memset(nav.get(), 0,sizeof(nav_t));
    if (!(nav->eph =(eph_t  *)malloc(sizeof(eph_t )*MAXSAT*4 ))||
        !(nav->geph=(geph_t *)malloc(sizeof(geph_t)*NSATGLO*2))||
        !(nav->seph=(seph_t *)malloc(sizeof(seph_t)*NSATSBS*2))) {
        tracet(1,"initGlobalNav: malloc error\n");
        return;
    }

    for (int i=0;i<MAXSAT*4 ;i++) nav->eph [i]={0,-1,-1};
    for (int i=0;i<NSATGLO*2;i++) nav->geph[i]={0,-1};
    for (int i=0;i<NSATSBS*2;i++) nav->seph[i]={0};
    nav->n =MAXSAT *2;
    nav->ng=NSATGLO*2;
    nav->ns=NSATSBS*2;

    for (int i=0;i<MAXSAT*4 ;i++) nav->eph [i].ttr={0};
    for (int i=0;i<NSATGLO*2;i++) nav->geph[i].tof={0};
    for (int i=0;i<NSATSBS*2;i++) nav->seph[i].tof={0};
}

void updateGlobalEPH(rtcm_t* rtcm) {
    if (rtcm == nullptr) { return; }
    std::lock_guard<std::mutex> autolock(g_nav_lock);
    int ephsat = rtcm->ephsat;
    int ephset = rtcm->ephset;
    nav_t* nav = &rtcm->nav;
    int prn = 0;

    if (satsys(ephsat,&prn)!=SYS_GLO) {
        eph_t *eph1 = nav->eph + ephsat - 1 + MAXSAT * ephset;      /* received */
        eph_t *eph2 = globalnav->eph + ephsat - 1 + MAXSAT * ephset;     /* current */
        eph_t *eph3 = globalnav->eph + ephsat - 1 + MAXSAT * (2 + ephset); /* previous */
        if (eph2->ttr.time == 0 ||
            (eph1->iode != eph3->iode && eph1->iode != eph2->iode) ||
            (timediff(eph1->toe, eph3->toe) != 0.0 &&
             timediff(eph1->toe, eph2->toe) != 0.0) ||
            (timediff(eph1->toc, eph3->toc) != 0.0 &&
             timediff(eph1->toc, eph2->toc) != 0.0)) {
            *eph3 = *eph2; /* current ->previous */
            *eph2 = *eph1; /* received->current */
        }
    } else {
        geph_t* geph1=nav->geph+prn-1;
        geph_t* geph2=globalnav->geph+prn-1;
        geph_t* geph3=globalnav->geph+prn-1+MAXPRNGLO;
        if (geph2->tof.time==0||
            (geph1->iode!=geph3->iode&&geph1->iode!=geph2->iode)) {
            *geph3=*geph2;
            *geph2=*geph1;
        }
    }
    return;
}

void updateGlobalIONUTC(rtcm_t* rtcm) {
    if (rtcm == nullptr) { return; }
    std::lock_guard<std::mutex> autolock(g_nav_lock);
    nav_t* nav = &rtcm->nav;
    matcpy(globalnav->utc_gps,nav->utc_gps,8,1);
    matcpy(globalnav->utc_glo,nav->utc_glo,8,1);
    matcpy(globalnav->utc_gal,nav->utc_gal,8,1);
    matcpy(globalnav->utc_qzs,nav->utc_qzs,8,1);
    matcpy(globalnav->utc_cmp,nav->utc_cmp,8,1);
    matcpy(globalnav->utc_irn,nav->utc_irn,9,1);
    matcpy(globalnav->utc_sbs,nav->utc_sbs,4,1);
    matcpy(globalnav->ion_gps,nav->ion_gps,8,1);
    matcpy(globalnav->ion_gal,nav->ion_gal,4,1);
    matcpy(globalnav->ion_qzs,nav->ion_qzs,8,1);
    matcpy(globalnav->ion_cmp,nav->ion_cmp,8,1);
    matcpy(globalnav->ion_irn,nav->ion_irn,8,1);
}

std::shared_ptr<nav_t> copyGlobalNav() {
    std::lock_guard<std::mutex> autolock(g_nav_lock);
    std::shared_ptr<nav_t> tempNav(new nav_t, freeNavFunc);
    initNav(tempNav);
    for (int i = 0; i < globalnav->n; i++) {
        tempNav->eph[i] = globalnav->eph[i];
    }
    for (int i = 0; i < globalnav->ng; i++) {
        tempNav->geph[i] = globalnav->geph[i];
    }
    matcpy(tempNav->utc_gps,globalnav->utc_gps,8,1);
    matcpy(tempNav->utc_glo,globalnav->utc_glo,8,1);
    matcpy(tempNav->utc_gal,globalnav->utc_gal,8,1);
    matcpy(tempNav->utc_qzs,globalnav->utc_qzs,8,1);
    matcpy(tempNav->utc_cmp,globalnav->utc_cmp,8,1);
    matcpy(tempNav->utc_irn,globalnav->utc_irn,9,1);
    matcpy(tempNav->utc_sbs,globalnav->utc_sbs,4,1);
    matcpy(tempNav->ion_gps,globalnav->ion_gps,8,1);
    matcpy(tempNav->ion_gal,globalnav->ion_gal,4,1);
    matcpy(tempNav->ion_qzs,globalnav->ion_qzs,8,1);
    matcpy(tempNav->ion_cmp,globalnav->ion_cmp,8,1);
    matcpy(tempNav->ion_irn,globalnav->ion_irn,8,1);
    return tempNav;
}

void navStreamLoop(std::string path) {
    //1、init
    initNav(globalnav);
    stream_t *stream = new stream_t;
    strinit(stream);

    trace(4, "navStreamLoop\n");

    //2、open stream
    int rw = STR_MODE_RW;
    int strcli = STR_NTRIPCLI;
    int ok = stropen(stream, strcli, rw, path.c_str());
    if (!ok) {
        trace(1, "[navStreamLoop] failed to stropen %s\n", path.c_str());
        return;
    }

    trace(4, "[navStreamLoop] stropen success\n");
    const int buffsize = 32768;
    uint8_t* buffer = new uint8_t[buffsize];
    uint32_t tick = 0;
    while (true) {
        tick=tickget();
        memset(buffer, 0, sizeof(uint8_t) * buffsize);
        int n = strread(stream, buffer, buffsize);
        if (n <= 0) {
            trace(5, "[navStreamLoop] strread failed: %d\n", n);
            continue;
        }
        trace(5, "[navStreamLoop] strread success: %d\n", n);

        rtcm_t* rtcm = new rtcm_t;
        memset(rtcm,0,sizeof(rtcm_t));
        init_rtcm(rtcm);
        rtcm->dgps = globalnav->dgps;
        trace(5, "[navStreamLoop] init_rtcm success\n");

        int missedNavCnt = 0;
        for (int i = 0; i < n; i++) {
            int ret = input_rtcm3(rtcm, buffer[i]);
            if (2 == ret) { /* ephemeris */
                trace(5, "[navStreamLoop] get nav frame\n");
                /* nav data */
                updateGlobalEPH(rtcm);
            } else if (9 == ret) { /* ion/utc parameters */
                trace(4, "[navStreamLoop] get ion/utc frame\n");
                updateGlobalIONUTC(rtcm);
            } else if (10 == ret) {
                updateSSR(rtcm);
            } else {
                missedNavCnt++;
                trace(5, "[navStreamLoop] not nav frame, pass cnt:%d\n", missedNavCnt);
            }
        }

        int cputime= tickget()-tick;
        sleepms(1000 - cputime);
        free_rtcm(rtcm);
        delete rtcm;
    }
}

extern void startNAVStream(const char* stationUrl) {
    trace(2, "startNAVStream\n");
    std::string path{stationUrl};
    std::thread loop(navStreamLoop, path);
    loop.detach();
}

extern int init(const char *baseStationUrl, const char *navStationUrl, int module, int navsys) {
    trace(2, "init");
    if (baseStationUrl == NULL || navStationUrl == NULL) {
        trace(1, "path can't be null\n");
        return -1;
    }

    //1、start thread to receive obs from base station
    startOBSStream(baseStationUrl, navsys, 2);

    //2、start thread to receive nav from station
    startNAVStream(navStationUrl);

    sleepms(100); //wait thread create
    return 0;
}

extern prcopt_t getDefaultDemo5Prcopt(int pmode, int sys) {
    prcopt_t prcopt = prcopt_default;
    prcopt.mode = pmode; //PMODE_KINEMA
    prcopt.navsys = sys; //SYS_GPS
    prcopt.nf = 3;
    prcopt.soltype = 3;
    prcopt.elmin = 15*D2R;
    prcopt.snrmask.ena[0] = 1;
    for (int i =0; i < 9; i++) {
        prcopt.snrmask.mask[0][i] = 24;
    }
    prcopt.snrmask.ena[1] = 1;
    for (int i = 0; i < 9; i++) {
        prcopt.snrmask.mask[1][i] = 34;
    }
    for (int i = 0; i < 9; i++) {
        prcopt.snrmask.mask[2][i] = 24;
    }
    prcopt.dynamics = 1;
    prcopt.tidecorr = 0;
    prcopt.ionoopt = IONOOPT_BRDC;
    prcopt.tropopt = TROPOPT_SAAS;
    prcopt.sateph = EPHOPT_BRDC;
    for (int i = 0; i < 6; i++) {
        prcopt.posopt[i] = 0;
    }
    prcopt.modear = 0;
    prcopt.glomodear = 0;
    prcopt.bdsmodear = 1;
    prcopt.arfilter = 1;
    prcopt.thresar[0] = 3;
    prcopt.thresar[5] = 3;
    prcopt.thresar[6] = 3;
    prcopt.thresar[1] = 0.05;
    prcopt.thresar[2] = 0;
    prcopt.thresar[3] = 1e-9;
    prcopt.thresar[4] = 1e-5;
    prcopt.varholdamb = 0.1;
    prcopt.gainholdamb = 0.01;
    prcopt.minlock = 5;
    prcopt.minfixsats = 4;
    prcopt.minholdsats = 5;
    prcopt.mindropsats = 10;
    prcopt.minfix = 10;
    prcopt.armaxiter = 1;
    prcopt.maxout = 4;
    prcopt.maxtdiff = 30;
    prcopt.syncsol = 0;
    prcopt.thresslip = 0.1;
    prcopt.thresdop = 5;
    prcopt.maxinno = 1;
    prcopt.maxgdop = 30;
    prcopt.niter = 1;
    prcopt.baseline[0] = 0;
    prcopt.baseline[1] = 0;
    prcopt.outsingle = 0;
    prcopt.eratio[0] = 300;
    prcopt.eratio[1] = 300;
    prcopt.eratio[2] = 100;

    prcopt.err[1] = 0.003;
    prcopt.err[2] = 0.003;
    prcopt.err[3] = 0;
    prcopt.err[4] = 1;
    prcopt.err[5] = 52;
    prcopt.err[6] = 0.00;
    prcopt.err[7] = 0;

    prcopt.std[0] = 30;
    prcopt.std[1] = 0.03;
    prcopt.std[2] = 0.3;
    prcopt.prn[3] = 3;
    prcopt.prn[4] = 1;
    prcopt.prn[0] = 0.01;
    prcopt.prn[1] = 0.001;
    prcopt.prn[2] = 0.0001;
    prcopt.prn[5] = 0;
    prcopt.sclkstab = 5e-12;
    prcopt.initrst = 1;
    prcopt.intpref = 1;
    prcopt.sbassatsel = 0;
    return prcopt;
}

extern int rtk(const char* obsInJson, int pmode, int sys,
               double stationX, double stationY, double stationZ,
               int* solState, int* validSatNum, long* gpstime, double *lat, double *lng, double *high) {
    if (obsInJson == nullptr ) { return -1; }
    std::shared_ptr<CPPObs> latestBaseObsData = topGObsListWithLock(2);
    if (latestBaseObsData == nullptr || !latestBaseObsData->isValid()) {
        return -1;
    }

    std::shared_ptr<nav_t> latestNavData = copyGlobalNav();
    if (latestNavData == nullptr) {
        return -1;
    }
    prcopt_t prcopt = getDefaultDemo5Prcopt(pmode, sys);
    prcopt.rb[0] = stationX;
    prcopt.rb[1] = stationY;
    prcopt.rb[2] = stationZ;

    std::shared_ptr<rtk_t> rtk(new rtk_t, freeRtkFunc);
    rtkinit(rtk.get(), &prcopt);

    std::string obsStr{obsInJson};
    std::shared_ptr<CPPObs> usrObs = std::make_shared<CPPObs>(obsStr);
    std::shared_ptr<CPPObs> obs = std::make_shared<CPPObs>(*usrObs, *latestBaseObsData);
    std::shared_ptr<obs_t> copy = obs->GetCopy();
    rtkpos(rtk.get(), copy.get()->data, copy.get()->n, latestNavData.get());

    double rr[3];double pos[3];
    matcpy(rr,rtk->sol.rr,3,1);
    ecef2pos(rr,pos);
    pos[0] *= R2D;pos[1] *= R2D;
    *lat = pos[0];*lng = pos[1];*high = pos[2];
    *validSatNum = rtk->sol.ns;
    *solState = rtk->sol.stat;
    *gpstime = rtk->sol.time.time;
    return 0;
}