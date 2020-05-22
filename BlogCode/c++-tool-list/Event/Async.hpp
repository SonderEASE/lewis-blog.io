#pragma once


#include <uv.h>


namespace uv {

    template <typename Object>
    class Async {
        using Callback = void (Object::*)();
        struct Context {
            uv_async_t async;
            Object*    object;
            Callback   callback;
        };

    public:
        Async(const Async&) = delete;
        Async(Async&&) = delete;
        Async& operator=(const Async&) = delete;
        Async& operator=(Async&&) = delete;
        Async() = default;

        ~Async() {
            Close();
        }

        void Close() {
            if (_context) {
                uv_close((uv_handle_t*)&_context->async, [](uv_handle_t* h) {
                    delete static_cast<Context*>(h->data);
                });
                _context = nullptr;
            }
        }

        int Init(uv_loop_t* loop, Object* object, Callback callback) {
            if (_context) {
                return UV_EINVAL;
            }

            _context = new Context{};
            int ret = uv_async_init(loop, &_context->async, [](uv_async_t* h) {
                auto ctx = static_cast<Context*>(h->data);
                (ctx->object->*ctx->callback)();
            });
            if (ret < 0) {
                delete _context;
                _context = nullptr;
            } else {
                _context->async.data = _context;
                _context->object = object;
                _context->callback = callback;
            }
            return ret;
        }

        int Send() {
            if (_context) {
                return uv_async_send(&_context->async);
            }
            return UV_EINVAL;
        }

    private:
        Context* _context = nullptr;
    };

} // namespace uv