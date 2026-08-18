package com.laixhe.webapi;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

/**
 * 接口限流测试: 同一 IP 在窗口内超过 max 次请求返回 429 (app.limit.max=3),
 * 健康检查路径豁免限流。
 * 独立测试上下文, 保证令牌桶从零开始。
 */
@SpringBootTest(properties = "app.limit.max=3")
@AutoConfigureMockMvc
@ActiveProfiles("h2")
class RateLimitTests {

    @Autowired
    private MockMvc mockMvc;

    @Test
    void overLimit_returns429_andHealthExempt() throws Exception {
        // 窗口内前 3 次正常通过
        mockMvc.perform(get("/api/v1/user/list")).andExpect(status().isOk());
        mockMvc.perform(get("/api/v1/user/list")).andExpect(status().isOk());
        mockMvc.perform(get("/api/v1/user/list")).andExpect(status().isOk());

        // 第 4 次触发限流 → 429 统一 JSON
        mockMvc.perform(get("/api/v1/user/list"))
                .andExpect(status().isTooManyRequests())
                .andExpect(jsonPath("$.code").value(429))
                .andExpect(jsonPath("$.message").value("请求过于频繁，请稍后再试"));

        // 健康检查豁免限流, 即使已超限仍返回 200
        mockMvc.perform(get("/api/v1/health")).andExpect(status().isOk());
    }
}
