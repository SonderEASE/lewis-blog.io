#ifndef WHAT_IS_CALLBACK_HOTEL_HPP
#define WHAT_IS_CALLBACK_HOTEL_HPP

#include <string>
#include <vector>
#include <map>
#include <iostream>
#include <zconf.h>
#include <functional>

class Hotel {
public:
    explicit Hotel(const std::string& name) {m_name = name;};
    using WakeUpMode = std::function<void()>;

public:
    static void Knock() {std::cout<<"咚咚咚!"<<std::endl;}
    static void Call() {std::cout<<"滴铃铃!"<<std::endl;}
    void OrderWakeUpServer(int room_id, const WakeUpMode& mode) {
        wake_lists.insert(std::make_pair(room_id, mode));
    };
    void WakeUp(int room_id) {
        auto it = wake_lists.find(room_id);
        if(it != wake_lists.end()) it->second();
        else {
            std::cout<<"They dont order the server"<<std::endl;
        }
    };

public:
    std::string m_name;

private:
    std::map<int, WakeUpMode>wake_lists;
};

#endif //WHAT_IS_CALLBACK_HOTEL_HPP
