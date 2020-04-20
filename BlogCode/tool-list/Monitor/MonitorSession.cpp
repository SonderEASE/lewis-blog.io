#include "MonitorSession.hpp"
#include "Common/Logger.hpp"
#include "Common/TcpStream.hpp"
#include "MonitorUtils.hpp"
#include <cassert>
#include <sstream>


MonitorSession::MonitorSession(MonitorManager& server, TcpStream* stream)
    : _manager(server), _stream(stream) {
    assert(_stream != nullptr);

    _stream->SetOnDataCallback([](const char* data, size_t size, void* self) -> int {
        return static_cast<MonitorSession*>(self)->OnStreamData(data, size);
    }, this);
    _stream->SetOnErrorCallback([](int code, void* self) {
        static_cast<MonitorSession*>(self)->OnStreamError(code);
    }, this);
    _stream->SetOnFinishCallback([](void* self) {
        static_cast<MonitorSession*>(self)->OnFinishResponse();
    }, this);
}


MonitorSession::~MonitorSession() {
    delete _stream;
    _stream = nullptr;
}


bool MonitorSession::Start() {
    _status = Status::WaitRequest;
    return true;
}


void MonitorSession::Stop() {
    _status = Status::Stopped;

    delete _stream;
    _stream = nullptr;

    CloseSession();
}


int MonitorSession::OnStreamData(const char* data, size_t size) {
    _status = Status::Processing;
    MonitorManager::Request request{};
    MonitorManager::Response response{};

    _request_buffer.append(data, size);

    // TODO: request may not complete

    int ret = ParseRequest(request);
    if (ret != E_OK) {
        response.status_code = 400;
        response.headers["Content-Type"] = "application/json";
        response.body = R"({"error_message":"parse request fail"})";
        SendResponse(response);
        return E_OK;
    }

    ret = HandleRequest(request, response);
    if (ret != E_OK) {
        response.status_code = 404;
        response.headers["Content-Type"] = "application/json";
        response.body = R"({"error_message":"path is not register"})";
    }

    SendResponse(response);

    return E_OK;
}


void MonitorSession::OnStreamError(int error) {
    LOG_ERROR("[MonitorSession] Error: {}", error);

    if (_status == Status::Stopped) return;
    _status = Status::Stopped;

    Stop();
}


void MonitorSession::OnFinishResponse() {
    Stop();
}


void MonitorSession::CloseSession() {
    _manager.OnSessionStop(this);
}


int MonitorSession::ParseRequest(MonitorManager::Request& request) {
    http_parser_init(&_request_parser, HTTP_REQUEST);
    _request_parser.data = &request;

    http_parser_settings settings{};

    settings.on_url = [](http_parser* parser, const char* at, size_t length) -> int {
        static_cast<MonitorManager::Request*>(parser->data)->url.assign(at, length);
        return E_OK;
    };

    static std::string filed{};
    settings.on_header_field = [](http_parser*, const char* at, size_t length) -> int {
        filed.assign(at, length);
        return E_OK;
    };
    settings.on_header_value = [](http_parser* parser, const char* at, size_t length) -> int {
        std::string value(at, length);
        static_cast<MonitorManager::Request*>(parser->data)->headers.insert(std::make_pair(filed, value));
        return E_OK;
    };

    settings.on_body = [](http_parser* parser, const char* at, size_t length) -> int {
        static_cast<MonitorManager::Request*>(parser->data)->body.append(at, length);
        return E_OK;
    };

    size_t n = http_parser_execute(&_request_parser, &settings, _request_buffer.data(), _request_buffer.size());
    if (n != _request_buffer.size()) {
        LOG_ERROR("[MonitorSession] Http Parse error, expect size: {} , actual size: {}", _request_buffer.size(), n);
        return E_ParseRequestFail;
    }

    request.method = (http_method)_request_parser.method;

    ParseQuery(request);

    return E_OK;
}


void MonitorSession::ParseQuery(MonitorManager::Request& request) {
    auto question_pos = request.url.find('?');
    if (question_pos == std::string::npos) { // has no queries
        request.path = request.url;
        return;
    }

    request.path = request.url.substr(0, question_pos);

    auto query_string = request.url.substr(question_pos + 1);
    auto queries = MonitorUtils::SplitQuery(query_string);
    for (const auto& query : queries) {
        std::string field, value;
        auto ok = MonitorUtils::SplitQueryPair(query, field, value);
        if (!ok) continue;
        request.queries[field] = value;
    }
}


int MonitorSession::HandleRequest(MonitorManager::Request& request, MonitorManager::Response& response) {
    MonitorManager::Handler handler = _manager.GetMonitorHandler(request.path);
    if (handler == nullptr) {
        LOG_ERROR("[MonitorSession] path is not register, path: {}", request.path);
        return E_PathNotRegister;
    }

    response = handler(request);
    return E_OK;
}


void MonitorSession::SendResponse(const MonitorManager::Response& response) {
    std::stringstream ss{};

    const char* status_message = http_status_str((enum http_status)response.status_code);
    ss << "HTTP/1.1 " << response.status_code << ' ' << status_message << "\r\n";

    for (const auto& header : response.headers) {
        if (!MonitorUtils::IsCaseInsensitiveEqual("Connection", header.first) &&
            !MonitorUtils::IsCaseInsensitiveEqual("Content-Length", header.first)) {
            ss << header.first << ": " << header.second << "\r\n";
        }
    }
    ss << "Content-Length: " << response.body.size() << "\r\n";
    ss << "Connection: close\r\n";
    ss << "\r\n";

    ss << response.body;

    auto response_message = ss.str();
    int ret = _stream->Send(response_message.data(), response_message.size());
    if (ret) {
        LOG_ERROR("[MonitorSession] send response fail, error: {}", ret);
        OnStreamError(E_SendResponseFail);
    }

    _status = Status::Response;
}
