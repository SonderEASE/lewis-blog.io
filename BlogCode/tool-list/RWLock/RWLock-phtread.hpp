#ifndef __MUTEX_H__
#define __MUTEX_H__

#include <pthread.h>

class Mutex
{
    public:
        Mutex()
        {
            ::pthread_mutex_init(&m_mutex, nullptr);
        }

        virtual ~Mutex()
        {
            ::pthread_mutex_destroy(&m_mutex);
        }

        void Lock()
        {
            ::pthread_mutex_lock(&m_mutex);
        }

        void Unlock()
        {
            ::pthread_mutex_unlock(&m_mutex);
        }

    private:
        pthread_mutex_t m_mutex{};
};

class RWMutex
{
public:
    RWMutex()
    {
        ::pthread_rwlock_init(&m_rwlock, nullptr);
    }
    virtual ~RWMutex()
    {
        ::pthread_rwlock_destroy(&m_rwlock);
    }
    void RLock()
    {
        ::pthread_rwlock_rdlock(&m_rwlock);
    }
    void RUnlock()
    {
        ::pthread_rwlock_unlock(&m_rwlock);
    }
    void Lock()
    {
        ::pthread_rwlock_wrlock(&m_rwlock);
    }
    void Unlock()
    {
        ::pthread_rwlock_unlock(&m_rwlock);
    }

private:
    pthread_rwlock_t m_rwlock {};
};

class ScopedLocker
{
    public:
        explicit ScopedLocker(Mutex& mutex)
            :m_mutex(mutex)
        {
            m_mutex.Lock();
        }

        virtual ~ScopedLocker()
        {
            m_mutex.Unlock();
        }
    private:
        Mutex& m_mutex;
};


class ScopedRLocker
{
public:
    explicit ScopedRLocker(RWMutex& mutex)
            :m_mutex(mutex)
    {
        m_mutex.RLock();
    }

    virtual ~ScopedRLocker()
    {
        m_mutex.RUnlock();
    }
private:
    RWMutex& m_mutex;
};


class ScopedWLocker
{
public:
    explicit ScopedWLocker(RWMutex& mutex)
            :m_mutex(mutex)
    {
        m_mutex.RLock();
    }

    virtual ~ScopedWLocker()
    {
        m_mutex.RUnlock();
    }
private:
    RWMutex& m_mutex;
};
#endif //__MUTEX_H__