#pragma once


#include <uv.h>


namespace uv {

    template <typename Object>
    class Timer {
        using Callback = void (Object::*)();
        struct Context {
            uv_timer_t timer;
            Object*    object;
            Callback   callback;
        };

    public:
        Timer(const Timer&) = delete;
        Timer(Timer&&) = delete;
        Timer& operator=(const Timer&) = delete;
        Timer& operator=(Timer&&) = delete;
        Timer() = default;

        ~Timer() {
            Close();
        }

        void Close() {
            if (_context) {
                if (uv_is_active((uv_handle_t*)&_context->timer)) {
                    uv_timer_stop(&_context->timer);
                }
                if (!uv_is_closing((uv_handle_t*)&_context->timer)) {
                    uv_close((uv_handle_t*)&_context->timer, [](uv_handle_t* h) {
                        delete static_cast<Context*>(h->data);
                    });
                }
                _context = nullptr;
            }
        }

        int Init(uv_loop_t* loop, Object* object, Callback callback) {
            if (_context) {
                return UV_EINVAL;
            }

            _context = new Context{};
            int ret = uv_timer_init(loop, &_context->timer);
            if (ret < 0) {
                delete _context;
                _context = nullptr;
            } else {
                _context->timer.data = _context;
                _context->object = object;
                _context->callback = callback;
            }

            return ret;
        }

        int Start(uint64_t timeout, uint64_t repeat = 0) {
            if (_context) {
                return uv_timer_start(&_context->timer, [](uv_timer_t* h) {
                    auto ctx = static_cast<Context*>(h->data);
                    (ctx->object->*ctx->callback)();
                }, timeout, repeat);
            }
            return UV_EINVAL;
        }

        int Stop() {
            if (_context == nullptr) {
                return UV_EINVAL;
            }
            if (uv_is_active((uv_handle_t*)&_context->timer)) {
                return uv_timer_stop(&_context->timer);
            }
            return 0;
        }

        bool IsActive() {
            return _context && uv_is_active((uv_handle_t*)&_context->timer);
        }

    private:
        Context* _context = nullptr;
    };

}
