//
// Copyright(c) 2015 Gabi Melman; Lewis
// Distributed under the MIT License (http://opensource.org/licenses/MIT)
//
#pragma once

#ifndef SPDLOG_H
#include "spdlog/spdlog.h"
#endif

#include "spdlog/details/file_helper.h"
#include "spdlog/details/null_mutex.h"
#include "spdlog/fmt/fmt.h"
#include "spdlog/sinks/base_sink.h"

#include <chrono>
#include <cstdio>
#include <ctime>
#include <mutex>
#include <string>

namespace spdlog {
namespace sinks {

/*
 * Generator of daily log file names in format basename.YYYY-MM-DD-HH.ext
 */
struct hour_filename_calculator
{
    // Create filename for the form basename.YYYY-MM-DD-HH
    static filename_t calc_filename(const filename_t &filename, const tm &now_tm)
    {
        filename_t basename, ext;
        std::tie(basename, ext) = details::file_helper::split_by_extension(filename);
        std::conditional<std::is_same<filename_t::value_type, char>::value, fmt::memory_buffer, fmt::wmemory_buffer>::type w;
        fmt::format_to(
                w, SPDLOG_FILENAME_T("{}-{:04d}-{:02d}-{:02d}-{:02d}{}"),
                basename, now_tm.tm_year + 1900, now_tm.tm_mon + 1, now_tm.tm_mday, now_tm.tm_hour, ext);
        return fmt::to_string(w);
    }
};


/*
 * Rotating file sink based on date. rotates at midnight
 */
    template<typename Mutex, typename FileNameCalc = hour_filename_calculator>
    class hour_file_sink final : public base_sink<Mutex>
    {
    public:
        // create daily file sink which rotates on given time
        hour_file_sink(filename_t base_filename, int rotation_hour, bool truncate = false)
                : base_filename_(std::move(base_filename))
                , truncate_(truncate)
        {
            auto now = log_clock::now();
            file_helper_.open(FileNameCalc::calc_filename(base_filename_, now_tm(now)), truncate_);
            rotation_tp_ = next_rotation_tp_();
        }

    protected:
        void sink_it_(const details::log_msg &msg) override
        {

            if (msg.time >= rotation_tp_)
            {
                file_helper_.open(FileNameCalc::calc_filename(base_filename_, now_tm(msg.time)), truncate_);
                rotation_tp_ = next_rotation_tp_();
            }
            fmt::memory_buffer formatted;
            sink::formatter_->format(msg, formatted);
            file_helper_.write(formatted);
        }

        void flush_() override
        {
            file_helper_.flush();
        }

    private:
        tm now_tm(log_clock::time_point tp)
        {
            time_t tnow = log_clock::to_time_t(tp);
            return spdlog::details::os::localtime(tnow);
        }

        log_clock::time_point next_rotation_tp_()
        {
            auto now = log_clock::now();
            tm date = now_tm(now);
            return now + std::chrono::seconds(60-date.tm_sec) + std::chrono::minutes(59-date.tm_min);
        }

        filename_t base_filename_;
        log_clock::time_point rotation_tp_;
        details::file_helper file_helper_;
        bool truncate_;
    };

    using hour_file_sink_mt = hour_file_sink<std::mutex>;
    using hour_file_sink_st = hour_file_sink<details::null_mutex>;

} // namespace sinks

//
// factory functions
//
    template<typename Factory = default_factory>
    inline std::shared_ptr<logger> hour_logger_mt(
            const std::string &logger_name, const filename_t &filename, bool truncate = false)
    {
        return Factory::template create<sinks::hour_file_sink_mt>(logger_name, filename, truncate);
    }

    template<typename Factory = default_factory>
    inline std::shared_ptr<logger> hour_logger_st(
            const std::string &logger_name, const filename_t &filename, bool truncate = false)
    {
        return Factory::template create<sinks::hour_file_sink_st>(logger_name, filename, truncate);
    }
}
