#pragma once

#include <string>
#include <sstream>
#include <fstream>
#include <vector>
#include <algorithm>

namespace StringUtils {

static const char* kWhiteChars = "\t\n\v\f\r ";

static inline std::vector<std::string> Compact(const std::vector<std::string>& tokens) {
    std::vector<std::string> compacted;

    for (const auto& token : tokens) {
        if (!token.empty()) {
            compacted.push_back(token);
        }
    }

    return compacted;
}

static inline std::vector<std::string>
Split(const std::string& str, const std::string& delim, const bool trim_empty = false) {
    size_t pos, last_pos = 0, len;
    std::vector<std::string> tokens;

    while (true) {
        pos = str.find(delim, last_pos);
        if (pos == std::string::npos) {
            pos = str.size();
        }

        len = pos - last_pos;
        if (!trim_empty || len != 0) {
            tokens.push_back(str.substr(last_pos, len));
        }

        if (pos == str.size()) {
            break;
        } else {
            last_pos = pos + delim.size();
        }
    }

    return tokens;
}

static inline std::string
Join(const std::vector<std::string>& tokens, const std::string& delim, const bool trim_empty = false) {
    if (trim_empty) {
        return Join(Compact(tokens), delim, false);
    } else {
        std::stringstream ss;
        for (size_t i = 0; i < tokens.size() - 1; ++i) {
            ss << tokens[i] << delim;
        }
        ss << tokens[tokens.size() - 1];

        return ss.str();
    }
}

static inline std::string Trim(const std::string& str) {
    const auto b = str.find_first_not_of(kWhiteChars);
    if (b == std::string::npos) {
        return "";
    }
    const auto e = str.find_last_not_of(kWhiteChars);
    return str.substr(b, e - b + 1);
}

static inline std::string TrimLeft(const std::string& str) {
    const auto b = str.find_first_not_of(kWhiteChars);
    if (b == std::string::npos) {
        return "";
    }
    return str.substr(b);
}

static inline std::string TrimRight(const std::string& str) {
    const auto e = str.find_last_not_of(kWhiteChars);
    if (e == std::string::npos) {
        return "";
    }
    return str.substr(0, e + 1);
}

static inline std::string Repeat(const std::string& str, unsigned int times) {
    std::stringstream ss;
    for (unsigned int i = 0; i < times; ++i) {
        ss << str;
    }
    return ss.str();
}

static inline std::string
ReplaceAll(const std::string& source, const std::string& target, const std::string& replacement) {
    return Join(Split(source, target, false), replacement, false);
}

static inline std::string ToUpper(const std::string& str) {
    std::string s(str);
    std::transform(s.begin(), s.end(), s.begin(), toupper);
    return s;
}

static inline std::string ToLower(const std::string& str) {
    std::string s(str);
    std::transform(s.begin(), s.end(), s.begin(), tolower);
    return s;
}

static inline bool StartsWith(const std::string& source, const std::string& compare,
                              bool case_sensitive = true) {
    if (source.length() < compare.length()) { return false; }

    return case_sensitive ? source.substr(0, compare.length()) == compare :
           ToLower(source.substr(0, compare.length())) == ToLower(compare);
}

static inline bool EndsWith(const std::string& source, const std::string& compare,
                            bool case_sensitive = true) {
    if (source.length() < compare.length()) { return false; }

    size_t start_position = source.length() - compare.length();
    return case_sensitive ? source.substr(start_position, compare.length()) == compare :
           ToLower(source.substr(start_position, compare.length())) == ToLower(compare);
}

static inline std::string ReadFile(const std::string& filepath) {
    std::ifstream ifs(filepath.c_str());
    std::string content((std::istreambuf_iterator<char>(ifs)),
                        (std::istreambuf_iterator<char>()));
    ifs.close();
    return content;
}

static inline void WriteFile(const std::string& filepath, const std::string& content) {
    std::ofstream ofs(filepath.c_str());
    ofs << content;
    ofs.close();
}

}
