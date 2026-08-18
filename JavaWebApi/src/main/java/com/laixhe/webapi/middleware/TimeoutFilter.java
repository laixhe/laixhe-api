package com.laixhe.webapi.middleware;

import com.laixhe.webapi.common.Error;
import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletResponse;
import tools.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * 请求超时过滤器 (与 Go 版 timeout 中间件对齐): 超过 http.timeout 秒未完成返回 408 统一 JSON。
 *
 * 实现说明: 将过滤链提交到独立线程执行并等待超时, 超时后返回 408;
 * 注意超时后下游请求线程无法被强制中断, 会继续运行到结束 (教学规模下仅用于兜底异常慢请求)。
 */
public class TimeoutFilter implements Filter {

    private final int timeoutSeconds;
    private final ObjectMapper objectMapper;
    private final ExecutorService executor = Executors.newCachedThreadPool(r -> {
        Thread t = new Thread(r, "request-timeout");
        t.setDaemon(true);
        return t;
    });

    public TimeoutFilter(int timeoutSeconds, ObjectMapper objectMapper) {
        this.timeoutSeconds = Math.max(timeoutSeconds, 1);
        this.objectMapper = objectMapper;
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        Future<?> future = executor.submit(() -> {
            try {
                chain.doFilter(request, response);
            } catch (Throwable e) {
                throw new WrappedServletException(e);
            }
        });
        try {
            future.get(timeoutSeconds, TimeUnit.SECONDS);
        } catch (TimeoutException e) {
            // 超时: 返回 408 统一 JSON (与 Go/Python 端一致)
            if (response instanceof HttpServletResponse httpResponse) {
                httpResponse.setStatus(408);
                httpResponse.setContentType("application/json");
                httpResponse.setCharacterEncoding(StandardCharsets.UTF_8.name());
                httpResponse.getWriter().write(objectMapper.writeValueAsString(new Error(408, "Request Timeout")));
            }
        } catch (ExecutionException e) {
            // 下游异常经 WrappedServletException 跨线程传递后还原为 ServletException
            Throwable cause = e.getCause();
            if (cause instanceof WrappedServletException wse) {
                throw wse.toServletException();
            }
            throw new ServletException(cause);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new ServletException(e);
        }
    }

    /** 包装受检异常以跨线程传递 */
    private static final class WrappedServletException extends RuntimeException {

        WrappedServletException(Throwable cause) {
            super(cause);
        }

        ServletException toServletException() {
            return new ServletException(getCause());
        }
    }
}
