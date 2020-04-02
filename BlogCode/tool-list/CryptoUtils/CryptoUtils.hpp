#pragma once


#include <string>


struct CryptoUtils {
    static std::string GetHexMD5(const std::string& str);
    static std::string GetHexMD5(const char* data, size_t size);

    static std::string GetHexSHA0(const std::string& str);
    static std::string GetHexSHA0(const char* data, size_t size);

    static std::string GetHexSHA1(const std::string& str);
    static std::string GetHexSHA1(const char* data, size_t size);

    static std::string GetHexSHA256(const std::string& str);
    static std::string GetHexSHA256(const char* data, size_t size);

    static std::string GetHex(const void* hex, size_t size);
};
