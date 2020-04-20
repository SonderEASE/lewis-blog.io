#pragma once

#if defined(WX_DISABLE_LOGGER)
#define LOG_DEBUG(...)
#define LOG_INFO(...)
#define LOG_WARN(...)
#define LOG_ERROR(...)
#define LOG_CRITICAL(...)
#else

#include "Monitor/MonitorManager.hpp"
#include <spdlog/spdlog.h>


enum class UpdateLoggerLevelResult {
    BadRequest = -1,
    WrongLevel = -2,
    UpdateLevelFail = -3,
    OK = 0,
};


class Logger {
public:
    static std::shared_ptr<spdlog::logger> server;

public:
    explicit Logger(const std::string& config_json);
    ~Logger();

    int error() const { return int(_error); }
    explicit operator bool() {
        return (_error == Error::OK) && server;
    }

public:
    static MonitorManager::Response SetLevel(const MonitorManager::Request& req);

private:
    struct Config {
        std::string name;
        std::string pattern;
        std::string file;
        std::size_t size;
        std::size_t count;
        spdlog::level::level_enum level;
        std::string sink;
    };

    void ParseConfig(const std::string& json);

private:
    Config _config;
    enum class Error {
        OK,
        Fail
    };
    Error _error = Error::OK;
};


#define LOG_DEBUG(...) SPDLOG_LOGGER_DEBUG(Logger::server, __VA_ARGS__)
#define LOG_INFO(...) SPDLOG_LOGGER_INFO(Logger::server, __VA_ARGS__)
#define LOG_WARN(...) SPDLOG_LOGGER_WARN(Logger::server, __VA_ARGS__)
#define LOG_ERROR(...) SPDLOG_LOGGER_ERROR(Logger::server, __VA_ARGS__)
#define LOG_CRITICAL(...) SPDLOG_LOGGER_CRITICAL(Logger::server, __VA_ARGS__)

#endif

