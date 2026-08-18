package com.laixhe.webapi.service;

import com.laixhe.webapi.common.ApiException;
import com.laixhe.webapi.dto.AuthLoginRequest;
import com.laixhe.webapi.dto.AuthLoginResponse;
import com.laixhe.webapi.dto.AuthRefreshResponse;
import com.laixhe.webapi.dto.AuthRegisterRequest;
import com.laixhe.webapi.dto.AuthRegisterResponse;
import com.laixhe.webapi.dto.UserResponse;
import com.laixhe.webapi.entity.User;
import com.laixhe.webapi.entity.UserExtend;
import com.laixhe.webapi.entity.UserSex;
import com.laixhe.webapi.entity.UserState;
import com.laixhe.webapi.entity.UserThirdParty;
import com.laixhe.webapi.entity.UserType;
import com.laixhe.webapi.repository.UserExtendRepository;
import com.laixhe.webapi.repository.UserRepository;
import com.laixhe.webapi.repository.UserThirdPartyRepository;
import com.laixhe.webapi.security.JwtService;
import lombok.RequiredArgsConstructor;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

/**
 * 鉴权业务 (对应 Go 版 services/auth.go)
 */
@Service
@RequiredArgsConstructor
public class AuthService {

    private final UserRepository userRepository;
    private final UserExtendRepository userExtendRepository;
    private final UserThirdPartyRepository userThirdPartyRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtService jwtService;

    /**
     * 注册: 事务内创建用户 + 扩展信息 + 第三方关联 (与 Go 版 CreateUser 对齐)
     */
    @Transactional
    public AuthRegisterResponse register(AuthRegisterRequest req) {
        // 先查邮箱是否已注册, 避免无效的 bcrypt 计算; email 唯一索引兜底并发防重
        if (userRepository.existsByEmail(req.email())) {
            throw ApiException.paramError("邮箱已存在");
        }
        LocalDateTime now = LocalDateTime.now();
        User user = new User();
        user.setTypeId(UserType.ORDINARY);
        user.setAccount(UUID.randomUUID().toString().replace("-", ""));
        user.setMobile("");
        user.setNickname(req.nickname());
        user.setEmail(req.email());
        user.setPassword(passwordEncoder.encode(req.password()));
        user.setAvatarUrl("");
        user.setSex(UserSex.UNKNOWN);
        user.setStates(UserState.NORMAL);
        user.setCreatedAt(now);
        user.setUpdatedAt(now);
        try {
            user = userRepository.save(user);
        } catch (DataIntegrityViolationException e) {
            // 唯一键冲突仅出现在并发注册同邮箱等极端情况
            throw ApiException.paramError("注册失败，请稍后再试");
        }
        UserExtend extend = new UserExtend();
        extend.setUid(user.getId());
        userExtendRepository.save(extend);
        UserThirdParty thirdParty = new UserThirdParty();
        thirdParty.setUid(user.getId());
        userThirdPartyRepository.save(thirdParty);

        String token = jwtService.generateToken(user.getId());
        return new AuthRegisterResponse(token, UserResponse.from(user));
    }

    /**
     * 登录: 账号封禁与密码错误返回同一提示, 避免暴露账号状态 (与 Go 版一致)
     */
    public AuthLoginResponse login(AuthLoginRequest req) {
        User user = userRepository.findByEmail(req.email())
                .orElseThrow(() -> ApiException.paramError("邮箱或密码错误"));
        if (user.getStates() != UserState.NORMAL) {
            throw ApiException.paramError("邮箱或密码错误");
        }
        if (!passwordEncoder.matches(req.password(), user.getPassword())) {
            throw ApiException.paramError("邮箱或密码错误");
        }
        String token = jwtService.generateToken(user.getId());
        return new AuthLoginResponse(token, UserResponse.from(user));
    }

    /**
     * 刷新 JWT: 用户不存在或非正常状态统一返回 401
     */
    public AuthRefreshResponse refresh(int uid) {
        User user = userRepository.findById(uid)
                .orElseThrow(ApiException::unauthorized);
        if (user.getStates() != UserState.NORMAL) {
            throw ApiException.unauthorized();
        }
        String token = jwtService.generateToken(uid);
        return new AuthRefreshResponse(token, UserResponse.from(user));
    }
}
