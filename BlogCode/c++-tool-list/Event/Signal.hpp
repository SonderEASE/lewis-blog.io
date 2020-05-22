#pragma once


#include <uv.h>


namespace uv {

    template <typename Object>
    class Signal {
        using Callback = void (Object::*)(int);
        struct Context {
            uv_signal_t signal;
            Object*     object;
            Callback    callback;
        };

    public:
        Signal(const Signal&) = delete;
        Signal(Signal&&) = delete;
        Signal& operator=(const Signal&) = delete;
        Signal& operator=(Signal&&) = delete;
        Signal() = default;

        ~Signal() {
            Close();
        }

        void Close() {
            if (_context) {
                if (uv_is_active((uv_handle_t*)&_context->signal)) {
                    uv_signal_stop(&_context->signal);
                }
                if (!uv_is_closing((uv_handle_t*)&_context->signal)) {
                    uv_close((uv_handle_t*)&_context->signal, [](uv_handle_t* h) {
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
            int ret = uv_signal_init(loop, &_context->signal);
            if (ret < 0) {
                delete _context;
                _context = nullptr;
            } else {
                _context->signal.data = _context;
                _context->object = object;
                _context->callback = callback;
            }

            return ret;
        }

        int Start(int signum) {
            if (!_context) {
                return UV_EINVAL;
            }
            return uv_signal_start(&_context->signal, [](uv_signal_t* h, int signum) {
                auto ctx = static_cast<Context*>(h->data);
                (ctx->object->*ctx->callback)(signum);
            }, signum);
        }

        int Stop() {
            if (!_context) {
                return UV_EINVAL;
            }
            return uv_signal_stop(&_context->signal);
        }

    private:
        Context* _context = nullptr;
    };

} // namespace uv
