#include <memory>
#include "Hotel.hpp"
#include "Passenger.hpp"
int main() {
    auto hotel = std::make_shared<Hotel>("wanda");
    auto Mia = std::make_shared<Passenger>("Mia", *hotel, 601);

    // Mia选择的叫醒方式是大喊三声mia
    Mia->m_hotel.OrderWakeUpServer(Mia->room_id, []()->void {
        std::cout<<"Mia!!!\nMia!!!\nMia!!!"<<std::endl;
    });

    // 当然Mia也可以选择酒店提供的叫醒方式, 比如敲门
    Mia->m_hotel.OrderWakeUpServer(Mia->room_id, Hotel::Knock);

    //这个wanda酒店非常的low, 他只提供8点整的叫醒服务.
    //到8点了~ 比如住在601的mia.
    hotel->WakeUp(601);

    return 0;
}
