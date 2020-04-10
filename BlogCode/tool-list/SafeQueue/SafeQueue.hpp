#ifndef __SAFEQUEUE_H__
#define __SAFEQUEUE_H__

#include "Mutex.h"
#include <pthread.h>
#include <list>
#include <mutex>

template<class T>
class SafeQueue
{
public:
	SafeQueue() = default;

	~SafeQueue() = default;

	void Add(T item)
	{
        std::lock_guard<std::mutex> scoped_lock(m_mutex);
		m_container.push_back(item);
	}

	void PushFront(T item)
	{
        std::lock_guard<std::mutex> scoped_lock(m_mutex);
		m_container.push_front(item);		
	}

	bool Get(T* item)
	{
        std::lock_guard<std::mutex> scoped_lock(m_mutex);
		bool ret;
		if(m_container.empty() )
		{
			ret = false;
		}
		else
		{
			(*item) = m_container.front();
			m_container.pop_front();
			ret = true;
		}
		return ret;
	}

	bool RefGet(T& item)
	{
        std::lock_guard<std::mutex> scoped_lock(m_mutex);
		bool ret;
		if (m_container.empty())
		{
			ret = false;
		}
		else
		{
			item = m_container.front();
			m_container.pop_front();
			ret = true;
		}
		return ret;
	}

	int GetCnt()
	{
        std::lock_guard<std::mutex> scoped_lock(m_mutex);
		int size = 0;
		size = m_container.size();
		return size;
	}


private:
	std::mutex m_mutex;
	std::list<T> m_container;
};

#endif //__SAFEQUEUE_H__