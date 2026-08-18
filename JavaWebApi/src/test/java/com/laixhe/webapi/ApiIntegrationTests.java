package com.laixhe.webapi;

import tools.jackson.databind.JsonNode;
import tools.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;
import org.springframework.transaction.annotation.Transactional;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

/**
 * 接口集成测试 (H2 内存库, 每个用例事务回滚)
 */
@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("h2")
@Transactional
class ApiIntegrationTests {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @BeforeEach
    void clearSecurity() {
        // 防止上个用例的认证上下文泄漏到本用例
        SecurityContextHolder.clearContext();
    }

    private String registerUser() throws Exception {
        return registerUser("test@laixhe.com");
    }

    private String registerUser(String email) throws Exception {
        String body = """
                {"email":"%s","password":"abc12345","nickname":"测试用户"}
                """.formatted(email);
        MvcResult result = mockMvc.perform(post("/api/v1/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.token").isNotEmpty())
                .andExpect(jsonPath("$.user.uid").isNumber())
                .andExpect(jsonPath("$.user.email").value(email))
                .andExpect(jsonPath("$.user.nickname").value("测试用户"))
                .andExpect(jsonPath("$.user.states").value(1))
                .andReturn();
        JsonNode node = objectMapper.readTree(result.getResponse().getContentAsString());
        return node.path("token").asText();
    }

    @Test
    void health_ok() throws Exception {
        mockMvc.perform(get("/api/v1/health"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").value("ok"))
                .andExpect(jsonPath("$.database").value("up"))
                .andExpect(jsonPath("$.version").value("1.0.0"))
                .andExpect(jsonPath("$.started_at").isNotEmpty())
                .andExpect(jsonPath("$.now").isNotEmpty());
    }

    @Test
    void register_then_login() throws Exception {
        registerUser();

        String body = """
                {"email":"test@laixhe.com","password":"abc12345"}
                """;
        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.token").isNotEmpty())
                .andExpect(jsonPath("$.user.email").value("test@laixhe.com"));
    }

    @Test
    void register_duplicateEmail() throws Exception {
        registerUser();
        String body = """
                {"email":"test@laixhe.com","password":"abc12345","nickname":"另一个用户"}
                """;
        mockMvc.perform(post("/api/v1/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.code").value(422))
                .andExpect(jsonPath("$.message").value("邮箱已存在"));
    }

    @Test
    void register_badEmail() throws Exception {
        String body = """
                {"email":"not-an-email","password":"abc12345","nickname":"测试用户"}
                """;
        mockMvc.perform(post("/api/v1/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.code").value(422))
                .andExpect(jsonPath("$.message").value("邮箱格式错误"));
    }

    @Test
    void register_badPassword() throws Exception {
        String body = """
                {"email":"test@laixhe.com","password":"short","nickname":"测试用户"}
                """;
        mockMvc.perform(post("/api/v1/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.message").value("密码格式错误，需 6~64 位，只能包含字母 数字 _ @ $"));
    }

    @Test
    void register_badNickname() throws Exception {
        String body = """
                {"email":"test@laixhe.com","password":"abc12345","nickname":"a"}
                """;
        mockMvc.perform(post("/api/v1/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.message").value("昵称长度不能小于2位"));
    }

    @Test
    void login_wrongPassword() throws Exception {
        registerUser();
        String body = """
                {"email":"test@laixhe.com","password":"wrong123"}
                """;
        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.message").value("邮箱或密码错误"));
    }

    @Test
    void refresh_withoutToken() throws Exception {
        mockMvc.perform(post("/api/v1/auth/refresh"))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.code").value(401))
                .andExpect(jsonPath("$.message").value("Unauthorized"));
    }

    @Test
    void refresh_withToken() throws Exception {
        String token = registerUser();
        mockMvc.perform(post("/api/v1/auth/refresh")
                        .header("Authorization", "Bearer " + token))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.token").isNotEmpty())
                .andExpect(jsonPath("$.user.email").value("test@laixhe.com"));
    }

    @Test
    void update_user() throws Exception {
        String token = registerUser();
        String body = """
                {"nickname":"新昵称","avatar_url":"https://example.com/a.png"}
                """;
        mockMvc.perform(post("/api/v1/user/update")
                        .header("Authorization", "Bearer " + token)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.nickname").value("新昵称"))
                .andExpect(jsonPath("$.avatar_url").value("https://example.com/a.png"));
    }

    @Test
    void update_withoutToken() throws Exception {
        String body = """
                {"nickname":"新昵称"}
                """;
        mockMvc.perform(post("/api/v1/user/update")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.code").value(401));
    }

    @Test
    void update_badAvatar() throws Exception {
        String token = registerUser();
        String body = """
                {"nickname":"新昵称","avatar_url":"ftp://example.com/a.png"}
                """;
        mockMvc.perform(post("/api/v1/user/update")
                        .header("Authorization", "Bearer " + token)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(body))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.message").value("头像地址必须以http或https开头"));
    }

    @Test
    void userInfo_notFound() throws Exception {
        mockMvc.perform(get("/api/v1/user/info").param("uid", "999999"))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.code").value(422))
                .andExpect(jsonPath("$.message").value("用户不存在"));
    }

    @Test
    void userInfo_invalidUid() throws Exception {
        mockMvc.perform(get("/api/v1/user/info").param("uid", "0"))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.message").value("无效的用户ID"));
    }

    @Test
    void userList_pagination() throws Exception {
        registerUser();
        mockMvc.perform(get("/api/v1/user/list"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.page").value(1))
                .andExpect(jsonPath("$.page_size").value(12))
                .andExpect(jsonPath("$.total").value(1))
                .andExpect(jsonPath("$.list[0].email").value("test@laixhe.com"));

        // page_size 超上限被钳制为 100
        mockMvc.perform(get("/api/v1/user/list").param("page", "1").param("page_size", "1000"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.page_size").value(100));
    }

    @Test
    void notFound_route() throws Exception {
        mockMvc.perform(get("/api/v1/unknown"))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.code").value(404))
                .andExpect(jsonPath("$.message").value("Not Found"));
    }

    @Test
    void swaggerDocs_available() throws Exception {
        mockMvc.perform(get("/api/v1/swagger.yaml"))
                .andExpect(status().isOk());
        mockMvc.perform(get("/api/v1/swagger.json"))
                .andExpect(status().isOk());
        mockMvc.perform(get("/api/v1/swagger"))
                .andExpect(status().isOk());
    }
}
