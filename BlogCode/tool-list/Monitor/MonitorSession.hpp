#pragma once


#include "MonitorManager.hpp"
#include <uv.h>


class TcpStream;


class MonitorSession {
public:
    enum class Status { New, WaitRequest, Processing, Response, Stopped };

public:
    static const int E_OK               = 0;
    static const int E_ParseRequestFail = 1;
    static const int E_PathNotRegister  = 2;
    static const int E_SendResponseFail = 3;

public:
    MonitorSession(MonitorManager& server, TcpStream* stream);
    ~MonitorSession();

    bool Start();
    void Stop();

private:
    int OnStreamData(const char* data, size_t size);
    void OnStreamError(int error);
    void OnFinishResponse();

private:
    void CloseSession();

    int ParseRequest(MonitorManager::Request& request);
    static void ParseQuery(MonitorManager::Request& request);

    int HandleRequest(MonitorManager::Request& request, MonitorManager::Response& response);
    void SendResponse(const MonitorManager::Response& response);

private:
    Status _status = Status::New;

    MonitorManager& _manager;
    TcpStream* _stream{};

    http_parser _request_parser{};
    std::string _request_buffer{};
};