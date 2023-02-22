/*------------------------------------------------------------------------------
* rnx2rtcm.c : rinex to rtcm converter
*
*          Copyright (C) 2012 by T.TAKASU, All rights reserved.
*
* version : $Revision: 1.1 $ $Date: 2008/07/17 21:55:16 $
* history : 2012/12/12  1.0 new
*-----------------------------------------------------------------------------*/
#include "rnx2rtcm.h"

static const char rcsid[]="$Id:$";
#define PROGNAME    "RNX2RTCM"           /* program name */
#define TRACEFILE   "rnx2rtcm.trace"     /* debug trace file */

/* print usage ---------------------------------------------------------------*/
static const char *help[]={
    "",
    "usage: rnx2rtcm [options] [infile ...]",
    "",
    "options:",
    "  -ts  y/m/d h:m:s    start time (gpst)",
    "  -te  y/m/d h:m:s    end time (gpst)",
    "  -ti  tint           time interval (s)",
    "  -sta staid          station id",
    "  -out outfile        output rtcm file",
    "  -typ type[,type...] rtcm message types",
    "  -x   level          debug trace level",
    ""
};
static void print_help(void)
{
    int i;
    for (i=0;i<sizeof(help)/sizeof(*help);i++) fprintf(stderr,"%s\n",help[i]);
    exit(0);
}
/* test rtcm nav data --------------------------------------------------------*/
static int is_nav(int type)
{
    return type==1019||type==1044||type==1045||type==1046;
}
/* test rtcm gnav data -------------------------------------------------------*/
static int is_gnav(int type)
{
    return type==1020;
}
/* test rtcm ant info --------------------------------------------------------*/
static int is_ant(int type)
{
    return type==1005||type==1006||type==1007||type==1008||type==1033;
}
/* generate rtcm obs data messages -------------------------------------------*/
static void gen_rtcm_obs(rtcm_t *rtcm, const int *type, int n, vector<unsigned char >& buff)
{
    int i,j=0;
    
    for (i=0;i<n;i++) {
        if (is_nav(type[i])||is_gnav(type[i])||is_ant(type[i])) continue;
        j=i; /* index of last message */
    }
    for (i=0;i<n;i++) {
        if (is_nav(type[i])||is_gnav(type[i])||is_ant(type[i])) continue;
        
        if (!gen_rtcm3(rtcm,type[i],0, i!=j)) continue;
        for (int k = 0; k < rtcm->nbyte; k++) {
            buff.push_back(rtcm->buff[k]);
        }
    }
}
/* generate rtcm nav data messages -------------------------------------------*/
static void gen_rtcm_nav(gtime_t time, rtcm_t *rtcm, const nav_t *nav,
                         int *index, const int *type, int n, vector<unsigned char >& buff)
{
    int i,j,sat,prn;
    
    for (i=index[0];i<nav->n;i++) {
        
        if (time.time&&timediff(nav->eph[i].ttr,time)>-0.1) continue;
        sat=nav->eph[i].sat;
        rtcm->time=nav->eph[i].ttr;
        rtcm->nav.eph[sat-1]=nav->eph[i];
        rtcm->ephsat=sat;
        
        for (j=0;j<n;j++) {
            if (!is_nav(type[j])) continue;
            
            if (!gen_rtcm3(rtcm,type[j],0,0)) continue;
            for (int k = 0; k < rtcm->nbyte; k++) {
                buff.push_back(rtcm->buff[k]);
            }
        }
        index[0]=i+1;
    }
    for (i=index[1];i<nav->ng;i++) {
        
        if (time.time&&timediff(nav->geph[i].tof,time)>-0.1) continue;
        sat=nav->geph[i].sat;
        if (satsys(sat,&prn)!=SYS_GLO) continue;
        rtcm->time=nav->geph[i].tof;
        rtcm->nav.geph[prn-1]=nav->geph[i];
        rtcm->ephsat=sat;
        
        for (j=0;j<n;j++) {
            if (!is_gnav(type[j])) continue;
            
            if (!gen_rtcm3(rtcm,type[j],0,0)) continue;
            for (int k = 0; k < rtcm->nbyte; k++) {
                buff.push_back(rtcm->buff[k]);
            }
        }
        index[1]=i+1;
    }
}
/* generate rtcm antenna info messages ---------------------------------------*/
static void gen_rtcm_ant(rtcm_t *rtcm, const int *type, int n, vector<unsigned char >& buff)
{
    int i;
    
    for (i=0;i<n;i++) {
        if (!is_ant(type[i])) continue;
        
        if (!gen_rtcm3(rtcm,type[i],0,0)) continue;
        for (int k = 0; k < rtcm->nbyte; k++) {
            buff.push_back(rtcm->buff[k]);
        }
    }
}
/* convert to rtcm messages --------------------------------------------------*/
static int conv_rtcm(const int *type, int n,
                     vector<unsigned char>& buff,
                     const obs_t *obs,
                     const nav_t *nav,
                     const sta_t *sta,
                     int staid)
{
    gtime_t time0={0};
    rtcm_t rtcm={0};
    eph_t eph0={0};
    geph_t geph0={0};
    int i,j,prn,index[2]={0};
    
    if (!(rtcm.nav.eph =(eph_t  *)malloc(sizeof(eph_t )*MAXSAT   ))||
        !(rtcm.nav.geph=(geph_t *)malloc(sizeof(geph_t)*MAXPRNGLO))) return 0;
    
    rtcm.staid=staid;
    rtcm.sta=*sta;
    
    for (i=0;i<MAXSAT   ;i++) rtcm.nav.eph [i]=eph0;
    for (i=0;i<MAXPRNGLO;i++) rtcm.nav.geph[i]=geph0;
    
    /* update glonass freq channel number */
    for (i=0;i<nav->ng;i++) {
        if (satsys(nav->geph[i].sat,&prn)!=SYS_GLO) continue;
        rtcm.nav.geph[prn-1]=nav->geph[i];
    }
    /* gerate rtcm antenna info messages */
    gen_rtcm_ant(&rtcm,type,n,buff);
    
    for (i=0;i<obs->n;i=j) {
        
        /* extract epoch obs data */
        for (j=i+1;j<obs->n;j++) {
            if (timediff(obs->data[j].time,obs->data[i].time)>DTTOL) break;
        }
        rtcm.time=obs->data[i].time;
        rtcm.seqno++;
        rtcm.obs.data=obs->data+i;
        rtcm.obs.n=j-i;
        
        /* generate rtcm obs data messages */
        gen_rtcm_obs(&rtcm,type,n,buff);
        
        /* generate rtcm nav data messages */
        gen_rtcm_nav(rtcm.time,&rtcm,nav,index,type,n,buff);
    }
    /* gerate rtcm nav data messages */
    gen_rtcm_nav(time0,&rtcm,nav,index,type,n,buff);

    free(rtcm.nav.eph);
    free(rtcm.nav.geph);
    return 1;
}
/* obs2rtcm3 translate obs_t data to rtcm3----------------------------------------------------------------------*/
int obs2rtcm3(obs_t* obs, vector<unsigned char>& buff)
{
    nav_t nav={0};
    sta_t sta={{0}};
    int type[16] = {1077};

    sortobs(obs);
    uniqnav(&nav);
    
    /* convert to rtcm messages */
    if (!conv_rtcm(type,1,buff,obs,&nav,&sta,0)) {
        return -1;
    }
    
    freenav(&nav,0xFF);
    return 0;
}