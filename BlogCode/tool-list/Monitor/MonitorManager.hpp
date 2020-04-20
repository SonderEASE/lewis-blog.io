#pragma once


#include <functional>
#include <string>
#include <unordered_map>
#include <vector>
#include <set>
#include <uv.h>
#include <http_parser.h>


class MonitorSession;
class TcpServer;
class TcpStream;


class MonitorManager {
public:
    enum class Status { New, Listening, Stopped };

    struct Request {
        http_method method;
        std::string url;
        std::string path;
        std::unordered_map<std::string, std::string> headers;
        std::unordered_map<std::string, std::string> queries;
        std::string body;
    };
    struct Response {
        int status_code = 200;
        std::unordered_map<std::string, std::string> headers;
        std::string body;
    };
    using Handler = std::function<Response(const Request&)>;
    using Routers = std::unordered_map<std::string, Handler>;

public:
    static const int E_OK = 0;

public:
    explicit MonitorManager(uv_loop_t* loop);
    ~MonitorManager();

    bool Start(const std::string& ip, uint16_t port);
    int Stop();

    uint16_t GetListenPort() const { return _listen_port; }

    void RegisterRouter(const std::string& path, Handler handler);

    Handler GetMonitorHandler(std::string& path);
    void OnSessionStop(MonitorSession* session);

private:
    void OnError(int error);
    void OnStream(TcpStream* stream);

private:
    Status _status = Status::New;

    uv_loop_t*  _loop{};
    TcpServer* _listener{};
    uint16_t _listen_port{};

    Routers _routers{};
    std::set<MonitorSession*> _streams{};
};