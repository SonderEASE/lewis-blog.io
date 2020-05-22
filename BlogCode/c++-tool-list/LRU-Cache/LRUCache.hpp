#pragma once


#include <unordered_map>
#include <list>


template <typename key_t, typename value_t>
class LRUCache {
public:
    typedef typename std::pair<key_t, value_t> key_item_pair_t;
    typedef typename std::list<key_item_pair_t>::iterator list_iterator_t;

    void Init(uint32_t max_size) {
        _max_size = max_size;
    }

    void Put(key_t key, value_t value) {
        auto it = _cache_items_map.find(key);
        if (it != _cache_items_map.end()) {
            _cache_items_list.erase(it->second);
            _cache_items_map.erase(it);
        }

        _cache_items_list.push_front(key_item_pair_t{key, std::move(value)});
        _cache_items_map.insert(std::make_pair(std::move(key), _cache_items_list.begin()));

        while (_cache_items_map.size() > _max_size) {
            auto& last = _cache_items_list.back();
            _cache_items_map.erase(last.first);
            _cache_items_list.pop_back();
        }
    }

    value_t* Get(const key_t& key, bool query = true) {
        if (query) {
            _query_count++;
        }
        auto it = _cache_items_map.find(key);
        if (it == _cache_items_map.end()) {
            return nullptr;
        }

        if (query) {
            _hit_count++;
        }
        _cache_items_list.splice(_cache_items_list.begin(), _cache_items_list, it->second);
        return &it->second->second;
    }

    bool Exist(const key_t& key) {
        auto it = _cache_items_map.find(key);
        return it != _cache_items_map.end();
    }

    bool Remove(const key_t& key) {
        auto it = _cache_items_map.find(key);
        if (it == _cache_items_map.end()) {
            return false;
        }
        _cache_items_list.erase(it->second);
        _cache_items_map.erase(it);
        return true;
    }

    size_t Size() const {
        return _cache_items_map.size();
    }

    uint32_t HitCount() const {
        return _hit_count;
    }

    uint32_t QueryCount() const {
        return _query_count;
    }

    double GetHitRate() {
        return _query_count == 0 ? 0.0 : (_hit_count * 1.0 / _query_count);
    }

    void ResetHitRate() {
        _hit_count = 0;
        _query_count = 0;
    }

private:
    std::list<key_item_pair_t> _cache_items_list;
    std::unordered_map<key_t, list_iterator_t> _cache_items_map;
    uint32_t _max_size{};
    uint32_t _hit_count{};
    uint32_t _query_count{};
};
