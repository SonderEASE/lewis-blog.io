#pragma once


#include <curl/curl.h>
#include <string>


class EasyHttp {
public:
    explicit EasyHttp(std::string url);
    ~EasyHttp();

    void Reset();

    void AddHeader(const std::string& field, const std::string& value);

    void SetOption(CURLoption option, long value);
    void SetOption(CURLoption option, const char* value);
    void SetOption(CURLoption option, void* value);

    void SetMaxRedirects(unsigned int count) { _max_redirects = count; }

    void Get();
    void Post(const char* data);
    void Post(const std::string& buffer);
    void Post(const char* data, size_t size);
    void Perform();

    CURL* GetHandle() const { return _handle; }
    CURLcode GetCurlCode() const { return _curl_code; }
    long GetResponseCode() const { return _response_code; }

    // Return non-const reference, so we can move the result.
    std::string& GetResponseHead() { return _response_head; }
    std::string& GetResponseBody() { return _response_body; }

    const std::string& GetRemoteIP() const { return _remote_ip; }

    const std::string& GetOriginalURL() const { return _original_url; }
    const std::string& GetRedirectURL() const { return _redirect_url; }

public:
    static size_t OnWriteBody(char* buffer, size_t size, size_t nmemb, std::string* response);
    static size_t OnWriteHeader(char* buffer, size_t size, size_t nmemb, std::string* header);

private:
    CURL* _handle = nullptr;
    std::string _original_url;
    std::string _redirect_url;
    curl_slist* _header_list = nullptr;

    CURLcode     _curl_code = CURLE_OK;
    long         _response_code = 0;
    std::string  _response_head;
    std::string  _response_body;
    unsigned int _max_redirects = 5;
    unsigned int _redirects = 0;

    std::string  _remote_ip;
};
