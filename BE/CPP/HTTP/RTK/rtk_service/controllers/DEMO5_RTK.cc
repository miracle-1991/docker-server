#include "DEMO5_RTK.h"

using namespace DEMO5;

// Add definition of your processing function here
#include "RTKRunArg.h"
#include "offline.h"

void RTK::CallOfflineRTK(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback)
{
    setRunningStart();
    auto reqJson = req->getJsonObject();
    std::cout << reqJson->toStyledString() << std::endl;

    RTKRunArg arg;
    std::string ok = arg.Init(reqJson);
    if (!ok.empty()) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(drogon::k400BadRequest);
        resp->setContentTypeCode(drogon::CT_TEXT_HTML);
        resp->addHeader("Access-Control-Allow-Origin", "*");
        resp->setBody(ok);
        callback(resp);
        setRunningEnd();
        return;
    }

    auto pmodelist = arg.GetPModeList();
    auto navsyslist = arg.GetNavSysList();
    setFileCntToProcess(pmodelist.size() * navsyslist.size());
    vector<string> outfilelist;
    Offline o;
    for (auto pmode : pmodelist) {
        for (auto navsys: navsyslist) {
            string outfile;
            int r = o.RTK(arg.GetRoverObs(), arg.GetStationObs(), arg.GetStationNav(), arg.GetOutputPath(), pmode, navsys, outfile);
            if (r != 0) {
                auto resp = HttpResponse::newHttpResponse();
                resp->setStatusCode(drogon::k500InternalServerError);
                resp->setContentTypeCode(drogon::CT_TEXT_HTML);
                resp->addHeader("Access-Control-Allow-Origin", "*");
                std::string errmsg = string("run error for ") + outfile;
                resp->setBody(errmsg);
                callback(resp);
                setRunningEnd();
                return;
            }else {
                outfilelist.push_back(arg.GetOutputPath() + "/" + outfile);
                incFileCntInProcessing();
            }
        }
    }

    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(drogon::k200OK);
    resp->setContentTypeCode(drogon::CT_TEXT_HTML);
    resp->addHeader("Access-Control-Allow-Origin", "*");
    Json::Value rspJson;
    Json::Value filelistJson;
    for (string s : outfilelist) {
        filelistJson.append(s);
    }
    rspJson["outputfile"] = filelistJson;
    resp->setBody(rspJson.toStyledString());
    callback(resp);
    setRunningEnd();
}


void RTK::GetProcessingState(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback) const {
    std::shared_lock<std::shared_mutex> autolock(m_running_state_rwlock);
    double state = 0;
    if (m_runing_state == "stop") {
        state = 1;
    } else{
        if (m_total_cnt == 0) {
            state = 0;
        } else{
            state = m_cur_cnt / m_total_cnt;
        }
    }
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(drogon::k200OK);
    resp->setContentTypeCode(drogon::CT_TEXT_HTML);
    resp->addHeader("Access-Control-Allow-Origin", "*");
    Json::Value rspJson;
    rspJson["code"] =  0;
    rspJson["total"] =  m_total_cnt;
    rspJson["current"] =  m_cur_cnt;
    rspJson["processing"] =  state;
    resp->setBody(rspJson.toStyledString());
    LOG_DEBUG << rspJson.toStyledString();
    callback(resp);
}

void RTK::setRunningStart() {
    std::unique_lock<std::shared_mutex> autolock(m_running_state_rwlock);
    m_total_cnt = 0;
    m_cur_cnt = 0;
    m_runing_state = "running";
}

void RTK::setRunningEnd() {
    std::unique_lock<std::shared_mutex> autolock(m_running_state_rwlock);
    m_runing_state = "stop";
}

void RTK::setFileCntToProcess(int cnt) {
    std::unique_lock<std::shared_mutex> autolock(m_running_state_rwlock);
    m_total_cnt = cnt;
}
void RTK::incFileCntInProcessing() {
    std::unique_lock<std::shared_mutex> autolock(m_running_state_rwlock);
    m_cur_cnt++;
}

void RTK::HelloWorld(const HttpRequestPtr &req, std::function<void(const HttpResponsePtr &)> &&callback) {
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(drogon::k200OK);
    resp->setContentTypeCode(drogon::CT_TEXT_HTML);
    resp->addHeader("Access-Control-Allow-Origin", "*");
    Json::Value rspJson;
    rspJson["code"] =  "helloworld";
    resp->setBody(rspJson.toStyledString());
    callback(resp);
}
