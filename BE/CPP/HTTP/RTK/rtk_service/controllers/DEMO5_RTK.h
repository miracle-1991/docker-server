#pragma once

#include <drogon/HttpController.h>
#include <shared_mutex>

using namespace drogon;

namespace DEMO5
{
class RTK : public drogon::HttpController<RTK>
{
  public:
    METHOD_LIST_BEGIN
    METHOD_ADD(RTK::CallOfflineRTK, "/demo5", Post, Options);
    METHOD_ADD(RTK::GetProcessingState, "/rtkProcessing", Get, Options);
    METHOD_ADD(RTK::HelloWorld, "/", Get, Options);
    METHOD_LIST_END

    void CallOfflineRTK(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback);
    void GetProcessingState(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback) const;
    void HelloWorld(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback);
private:
    void setRunningStart();
    void setRunningEnd();
    void setFileCntToProcess(int cnt);
    void incFileCntInProcessing();
private:
    mutable std::shared_mutex m_running_state_rwlock;
    int m_total_cnt = 0;
    int m_cur_cnt = 0;
    std::string m_runing_state = "stop";
};
}
