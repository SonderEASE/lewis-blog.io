#include "Common/TcpStream.hpp"
#include "MonitorManager.hpp"
#include "MonitorSession.hpp"
#include "Common/Logger.hpp"
#include "Common/TcpServer.hpp"


MonitorManager::MonitorManager(uv_loop_t* loop) : _loop(loop) {
    _listener = new TcpServer(_loop);
    _listener->SetOnStreamCallback([](TcpStream* stream, void* self) {
        static_cast<MonitorManager*>(self)->OnStream(stream);
    }, this);
    _listener->SetOnErrorCallback([](int code, void* self) {
        static_cast<MonitorManager*>(self)->OnError(code);
    }, this);
}


MonitorManager::~MonitorManager() {
    Stop();
}


bool MonitorManager::Start(const std::string& ip, uint16_t port) {
    int result = _listener->Start(ip, port, 128);
    if (result != TcpServer::E_OK) {
        LOG_ERROR("[MonitorManager] Listen {}:{} failed: {}", ip, port, result);
        return false;
    }

    _listen_port = _listener->GetListenPort();

    _status = Status::Listening;
    return true;
}


int MonitorManager::Stop() {
    delete _listener;
    _listener = nullptr;

    _status = Status::Stopped;

    return E_OK;
}


void MonitorManager::RegisterRouter(const std::string& path, Handler handler) {
    _routers.insert(std::make_pair(path, std::move(handler)));
}


MonitorManager::Handler MonitorManager::GetMonitorHandler(std::string& path) {
    auto handler = _routers.find(path);
    if (handler != _routers.end()) {
        return handler->second;
    }
    return nullptr;
}


void MonitorManager::OnSessionStop(MonitorSession* session) {
    _streams.erase(session);
    delete session;
}


void MonitorManager::OnError(int error) {
    LOG_ERROR("[MonitorManager] Error: {}", error);
    if (_status == Status::Stopped) return;
    _status = Status::Stopped;

    Stop();
}


void MonitorManager::OnStream(TcpStream* stream) {
    auto s = new MonitorSession(*this, stream);
    LOG_INFO("[MonitorManager] Get new http connect, {}:{}", stream->GetRemoteIP(), stream->GetRemotePort());
    if (s->Start()) {
        _streams.insert(s);
    } else {
        delete s;
    }
}