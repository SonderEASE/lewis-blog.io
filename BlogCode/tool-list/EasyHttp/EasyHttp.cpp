#include "EasyHttp.hpp"


EasyHttp::EasyHttp(std::string url)
    : _handle(curl_easy_init()),
      _original_url(std::move(url)) {}


EasyHttp::~EasyHttp() {
    if (_handle) {
        curl_easy_cleanup(_handle);
        _handle = nullptr;
    }

    if (_header_list) {
        curl_slist_free_all(_header_list);
        _header_list = nullptr;
    }
}


void EasyHttp::Reset() {
    if (_handle) {
        curl_easy_reset(_handle);
    }

    if (_header_list) {
        curl_slist_free_all(_header_list);
        _header_list = nullptr;
    }

    _curl_code = CURLE_OK;
    _response_head.clear();
    _response_body.clear();
    _response_code = 0;

    _max_redirects = 5;
    _redirects = 0;

    _original_url.clear();
    _redirect_url.clear();
}


void EasyHttp::AddHeader(const std::string& field, const std::string& value) {
    _header_list = curl_slist_append(_header_list, (field + ": " + value).c_str());
}


void EasyHttp::SetOption(CURLoption option, long value) {
    curl_easy_setopt(_handle, option, value);
}


void EasyHttp::SetOption(CURLoption option, const char* value) {
    curl_easy_setopt(_handle, option, value);
}


void EasyHttp::SetOption(CURLoption option, void* value) {
    curl_easy_setopt(_handle, option, value);
}


void EasyHttp::Get() {
    curl_easy_setopt(_handle, CURLOPT_HTTPGET, 1);
    Perform();
}


void EasyHttp::Post(const char* data) {
    return Post(data, std::strlen(data));
}


void EasyHttp::Post(const std::string& buffer) {
    Post(buffer.data(), buffer.size());
}


void EasyHttp::Post(const char* data, size_t size) {
    curl_easy_setopt(_handle, CURLOPT_POST, 1);
    curl_easy_setopt(_handle, CURLOPT_POSTFIELDS, data);
    curl_easy_setopt(_handle, CURLOPT_POSTFIELDSIZE, size);

    Perform();
}


void EasyHttp::Perform() {
    curl_easy_setopt(_handle, CURLOPT_HEADERDATA, &_response_head);
    curl_easy_setopt(_handle, CURLOPT_HEADERFUNCTION, OnWriteHeader);
    curl_easy_setopt(_handle, CURLOPT_WRITEDATA, &_response_body);
    curl_easy_setopt(_handle, CURLOPT_WRITEFUNCTION, OnWriteBody);
    curl_easy_setopt(_handle, CURLOPT_HTTPHEADER, _header_list);
    curl_easy_setopt(_handle, CURLOPT_NOSIGNAL, 1);
    curl_easy_setopt(_handle, CURLOPT_USERAGENT, curl_version());
    curl_easy_setopt(_handle, CURLOPT_SSL_VERIFYPEER, 0);
    curl_easy_setopt(_handle, CURLOPT_SSL_VERIFYHOST, 0);

    const char sep = _original_url.find_last_of('?') == std::string::npos ? '?' : '&';
    auto url = _original_url + sep + "TPSecNotice&TPNotCheck";

    do {
        curl_easy_setopt(_handle, CURLOPT_URL, url.c_str());
        _curl_code = curl_easy_perform(_handle);

        char* ip = nullptr;
        curl_easy_getinfo(_handle, CURLINFO_PRIMARY_IP, &ip);
        if (ip) _remote_ip.assign(ip);

        if (_curl_code == CURLE_OK) {
            char* redirect_url = nullptr;
            curl_easy_getinfo(_handle, CURLINFO_REDIRECT_URL, &redirect_url);
            if (redirect_url == nullptr) {
                curl_easy_getinfo(_handle, CURLINFO_RESPONSE_CODE, &_response_code);
                break;
            }

            _response_code = 0;
            _response_head.clear();
            _response_body.clear();

            if (_redirects > _max_redirects) {
                _curl_code = CURLE_TOO_MANY_REDIRECTS;
            } else {
                ++_redirects;
                url = redirect_url;
                _redirect_url = redirect_url;
                //LOG_INFO("Redirect (count:{}) to: {}", _redirects, url);
            }
        } else break;
    } while (true);
}


size_t EasyHttp::OnWriteBody(char* buffer, size_t size, size_t nmemb, std::string* response) {
    if (response == nullptr || buffer == nullptr || response->size() > (10ul << 20u)) { // 10M
        return 0;
    }
    response->append(buffer, size * nmemb);
    return size * nmemb;
}


size_t EasyHttp::OnWriteHeader(char* buffer, size_t size, size_t nmemb, std::string* header) {
    if (header == nullptr || buffer == nullptr || header->size() > (1ul << 20ul)) { // 1M
        return 0;
    }
    header->append(buffer, size * nmemb);
    return size * nmemb;
}
