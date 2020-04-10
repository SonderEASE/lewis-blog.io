#pragma once

#include <mutex>
#include <condition_variable>

class RWLock {
public:
    RWLock() = default;

    ~RWLock() = default;

public:
    void read_lock() {
        std::unique_lock<std::mutex> ul(mutex);
        cond_r.wait(ul, [this] { return write_count == 0; });
        ++read_count;
    }

    void read_unlock() {
        std::unique_lock<std::mutex> ul(mutex);
        if (--read_count == 0 && write_count > 0) {
            cond_w.notify_one();
        }
    }

    void write_lock() {
        std::unique_lock<std::mutex> ul(mutex);
        ++write_count;
        cond_w.wait(ul, [this] { return read_count == 0 && can_write; });
        can_write = false;
    }

    void write_unlock() {
        std::unique_lock<std::mutex> ul(mutex);
        if (--write_count == 0) {
            cond_r.notify_all();
        } else {
            cond_w.notify_one();
        }
        can_write = true;
    }

private:
    volatile size_t read_count{0};
    volatile size_t write_count{0};
    volatile bool can_write{true};
    std::mutex mutex;
    std::condition_variable cond_w;
    std::condition_variable cond_r;
};

class WriteGuard {
public:
    explicit WriteGuard(RWLock& rw_lockable)
        : rw_lockable_(rw_lockable) {
        rw_lockable_.write_lock();
    }

    ~WriteGuard() {
        rw_lockable_.write_unlock();
    }

public:
    WriteGuard() = delete;

    WriteGuard(const WriteGuard&) = delete;

    WriteGuard& operator=(const WriteGuard&) = delete;

private:
    RWLock& rw_lockable_;
};

class ReadGuard {
public:
    explicit ReadGuard(RWLock& rw_lockable)
        : rw_lockable_(rw_lockable) {
        rw_lockable_.read_lock();
    }

    ~ReadGuard() {
        rw_lockable_.read_unlock();
    }

public:
    ReadGuard() = delete;

    ReadGuard(const ReadGuard&) = delete;

    ReadGuard& operator=(const ReadGuard&) = delete;

private:
    RWLock& rw_lockable_;
};
