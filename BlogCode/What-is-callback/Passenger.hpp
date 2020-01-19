#ifndef WHAT_IS_CALLBACK_PASSENGER_HPP
#define WHAT_IS_CALLBACK_PASSENGER_HPP
#include <utility>

#include "Hotel.hpp"
class Hotel;
class Passenger {
public:
    Passenger(std::string  name, Hotel& hotel, int id)
        : m_name(std::move(name)), m_hotel(hotel),room_id(id) {}

public:
    Hotel& m_hotel;
    int room_id;

private:
    std::string m_name;
};

#endif //WHAT_IS_CALLBACK_PASSENGER_HPP
