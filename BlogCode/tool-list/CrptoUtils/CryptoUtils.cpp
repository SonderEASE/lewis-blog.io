#include "CryptoUtils.hpp"
#include "SHA0.h"
#include <openssl/md5.h>
#include <openssl/sha.h>


std::string CryptoUtils::GetHexMD5(const std::string& str) {
    return GetHexMD5(str.data(), str.size());
}


std::string CryptoUtils::GetHexMD5(const char* data, size_t size) {
    unsigned char digests[MD5_DIGEST_LENGTH];
    MD5((const unsigned char*) data, size, digests);

    std::string hex(sizeof(digests)*2, 0);
    for (size_t i = 0; i < sizeof(digests); ++i) {
        std::sprintf(&hex[i * 2], "%02x", digests[i]);
    }
    return hex;
}


std::string CryptoUtils::GetHexSHA0(const std::string& str) {
    return GetHexSHA0(str.data(), str.size());
}


std::string CryptoUtils::GetHexSHA0(const char* data, size_t size) {
    unsigned char digests[SHA_DIGEST_LENGTH];
    SHA0((const unsigned char*) data, size, digests);

    std::string hex(sizeof(digests)*2, 0);
    for (size_t i = 0; i < sizeof(digests); ++i) {
        std::sprintf(&hex[i*2], "%02x", digests[i]);
    }
    return hex;
}


std::string CryptoUtils::GetHexSHA1(const std::string& str) {
    return GetHexSHA1(str.data(), str.size());
}


std::string CryptoUtils::GetHexSHA1(const char* data, size_t size) {
    unsigned char digests[SHA_DIGEST_LENGTH];
    SHA1((const unsigned char*) data, size, digests);

    std::string hex(sizeof(digests)*2, 0);
    for (size_t i = 0; i < sizeof(digests); ++i) {
        std::sprintf(&hex[i*2], "%02x", digests[i]);
    }
    return hex;
}


std::string CryptoUtils::GetHexSHA256(const std::string& str) {
    return GetHexSHA256(str.data(), str.size());
}


std::string CryptoUtils::GetHexSHA256(const char* data, size_t size) {
    unsigned char digests[SHA256_DIGEST_LENGTH];
    SHA256((const unsigned char*) data, size, digests);

    std::string hex(sizeof(digests)*2, 0);
    for (size_t i = 0; i < sizeof(digests); ++i) {
        std::sprintf(&hex[i*2], "%02x", digests[i]);
    }
    return hex;
}


std::string CryptoUtils::GetHex(const void* data, size_t size) {
    std::string hex(size*2, 0);
    for (size_t i = 0; i < size; ++i) {
        std::sprintf(&hex[i * 2], "%02x", ((unsigned char*)data)[i]);
    }
    return hex;
}
