#pragma once


#include <vector>
#include <string>


namespace MonitorUtils {
    // The query string is composed of a series of field-value pairs.
    // Within each pair, the field name and value are separated by an equals sign, "=".
    // The series of pairs is separated by the ampersand, "&"
    // https://en.wikipedia.org/wiki/Query_string

    static std::vector<std::string> SplitQuery(const std::string& query_string) {
        size_t last_pos = 0;
        std::vector<std::string> pairs;

        while (true) {
            auto pos = query_string.find('&', last_pos);
            if (pos == std::string::npos) {
                pos = query_string.size();
            }

            auto len = pos - last_pos;
            if (len != 0) {
                pairs.push_back(query_string.substr(last_pos, len));
            }

            if (pos == query_string.size()) {
                break;
            } else {
                last_pos = pos + 1;
            }
        }

        return pairs;
    }

    static bool SplitQueryPair(const std::string& pair, std::string& field, std::string& value) {
        auto pos = pair.find('=');
        if (pos == std::string::npos) {
            return false;
        }

        field = pair.substr(0, pos);
        value = pair.substr(pos + 1);

        return true;
    }

    static bool IsCaseInsensitiveEqual(const std::string& first, const std::string& second) {
        if (first.size() != second.size()) return false;

        for (size_t i = 0; i < first.size(); ++i) {
            if (std::tolower(first[i]) != std::tolower(second[i])) {
                return false;
            }
        }

        return true;
    }
}