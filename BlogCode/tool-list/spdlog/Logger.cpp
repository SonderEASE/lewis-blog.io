#include "Logger.hpp"

#if !defined(WX_DISABLE_LOGGER)
#include <spdlog/sinks/stdout_sinks.h>
#include <spdlog/sinks/rotating_file_sink.h>
#if defined(OS_ANDROID)
#include <android/asset_manager.h>
#include <android/asset_manager_jni.h>
#include "AndroidUtils.hpp"
#include <spdlog/spdlog.h>
#include <spdlog/sinks/android_sink.h>
#endif

#include <map>
#include <nlohmann/json.hpp>
#include <iostream>


std::shared_ptr<spdlog::logger> Logger::server;


MonitorManager::Response Logger::SetLevel(const MonitorManager::Request& req) {
    struct RequestJSON {
        std::string level;

        bool Parse(const char* json) {
            nlohmann::json root = nlohmann::json::parse(json, nullptr, false);
            if (root.is_discarded()) {
                return false;
            }

            nlohmann::json item = root["level"];
            if (!item.is_string()) return false;
            level = item;

            return true;
        }

        bool Parse(const std::string& json) {
            return Parse(json.c_str());
        }
    };
    struct ResponseJSON {
        static std::string Format(UpdateLoggerLevelResult code, const char* msg) {
            char response[1024];
            int n = std::sprintf(response, R"({"code":%d,"message":"%s"})", (int)code, msg);
            return std::string(response, n);
        }
        static std::string Format(UpdateLoggerLevelResult code, const std::string& msg) {
            return Format(code, msg.c_str());
        }
    };

    RequestJSON req_json;
    MonitorManager::Response res{};
    res.headers["Content-Type"] = "application/json";

    if (!req_json.Parse(req.body)) {
        res.body = ResponseJSON::Format(UpdateLoggerLevelResult::BadRequest, "Bad Request");
        return res;
    }

    using namespace std::string_literals;
    std::map<std::string, spdlog::level::level_enum> levels{
        {"debug"s,    spdlog::level::debug},
        {"info"s,     spdlog::level::info},
        {"warn"s,     spdlog::level::warn},
        {"error"s,    spdlog::level::err},
        {"critical"s, spdlog::level::critical},
    };
    auto l = levels.find(req_json.level);
    if (l == levels.end()) {
        res.body = ResponseJSON::Format(UpdateLoggerLevelResult::WrongLevel, "Wrong Level");
        return res;
    }

    try {
        LOG_INFO("[Logger] Set log level to: '{}'", req_json.level);
        server->flush();
        server->set_level(l->second);
        res.body = ResponseJSON::Format(UpdateLoggerLevelResult::OK, "OK");
    } catch (const std::exception& e) {
        res.body = ResponseJSON::Format(UpdateLoggerLevelResult::UpdateLevelFail, e.what());
    }
    return res;
}


void Logger::ParseConfig(const std::string& json) {
    _config.name = "platinum";
    _config.pattern = "[%x %X.%e] %L %v      %@";
    _config.file = "../logs/server.log";
    _config.size = 1048576;
    _config.count = 10;
    _config.level = spdlog::level::warn;

    if (json.empty()) {
        return;
    }

    nlohmann::json root = nlohmann::json::parse(json, nullptr, false);
    if (root.is_discarded()) {
        return;
    }

    nlohmann::json item = root["name"];
    if (item.is_string()) {
        _config.name = item;
    }

    item = root["pattern"];
    if (item.is_string()) {
        _config.pattern = item;
    }

    item = root["file"];
    if (item.is_string()) {
        _config.file = item;
    }

    item = root["size"];
    if (item.is_number_unsigned()) {
        _config.size = item;
    }

    item = root["count"];
    if (item.is_number_unsigned()) {
        _config.count = item;
    }

    using namespace std::string_literals;
    std::map<std::string, spdlog::level::level_enum> levels{
        {"debug"s,    spdlog::level::debug},
        {"info"s,     spdlog::level::info},
        {"warn"s,     spdlog::level::warn},
        {"error"s,    spdlog::level::err},
        {"critical"s, spdlog::level::critical},
    };
    item = root["level"];
    if (item.is_string()) {
        std::string value = item;
        auto l = levels.find(value);
        if (l != levels.end()) {
            _config.level = l->second;
        }
    }

    item = root["sink"];
    if (item.is_string()) {
        _config.sink = item;
    }
}


Logger::Logger(const std::string& config_json) {
    // 1. Parse config
//    std::cout<<config_json<<std::endl;
    ParseConfig(config_json);
    std::cout<<config_json<<std::endl;
    // 2. Create logger
    try {
#if defined(__APPLE__) && (defined(TARGET_IPHONE_SIMULATOR) || defined(TARGET_IOS_IPHONE))
        server = spdlog::stdout_logger_mt(_config.name);
#elif !defined(OS_ANDROID)
        if (config_json.empty()) {
            server = spdlog::stdout_logger_mt(_config.name);
        } else {
            server = spdlog::rotating_logger_mt(_config.name, _config.file, _config.size, _config.count);
        }
#else
        if (config_json.empty()) {
            server = spdlog::stdout_logger_mt(_config.name);
        } else if(_config.sink == "rotating") {
            DiskUtils::MkDir(DiskUtils::GetDirPath(_config.file));
            server = spdlog::rotating_logger_mt(_config.name, _config.file, _config.size, _config.count);
        } else {
            server = spdlog::android_logger_mt(_config.name, _config.name);
        }
#endif
        server->set_pattern(_config.pattern);
        server->flush_on(spdlog::level::info);
        server->set_level(_config.level);
    }
    catch (const std::exception& e) {
        std::fprintf(stderr, "create spdlog failed: %s\n", e.what());
        _error = Error::Fail;
    }
}


Logger::~Logger() {
    if (!server) return;

    try {
        spdlog::drop_all();
        spdlog::shutdown();
    }
    catch (const std::exception& e) {
        std::fprintf(stderr, "shutdown spdlog failed: %s\n", e.what());
    }
    server = nullptr;
}

#endif
