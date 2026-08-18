package com.laixhe.webapi.service;

import com.laixhe.webapi.common.ApiException;
import com.laixhe.webapi.dto.UserListResponse;
import com.laixhe.webapi.dto.UserResponse;
import com.laixhe.webapi.dto.UserUpdateRequest;
import com.laixhe.webapi.entity.User;
import com.laixhe.webapi.entity.UserState;
import com.laixhe.webapi.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 用户业务 (对应 Go 版 services/user.go)
 */
@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;

    /**
     * 更新用户信息 (Uid 由 JWT 提供, 返回更新后的预期值而非 DB 回读值, 与 Go 版一致)
     */
    public UserResponse update(int uid, UserUpdateRequest req) {
        User user = userRepository.findById(uid)
                .orElseThrow(() -> ApiException.paramError("用户不存在"));
        if (user.getStates() != UserState.NORMAL) {
            throw ApiException.unauthorized();
        }
        UserResponse resp = UserResponse.from(user, req.nickname(), req.avatarUrl());
        user.setNickname(req.nickname());
        // 与 Go 端 UpdateUser 非零字段更新语义一致: avatar_url 为空时不覆盖原值
        if (req.avatarUrl() != null && !req.avatarUrl().isEmpty()) {
            user.setAvatarUrl(req.avatarUrl());
        }
        user.setUpdatedAt(LocalDateTime.now());
        userRepository.save(user);
        return resp;
    }

    /**
     * 获取用户信息
     */
    public UserResponse info(int uid) {
        if (uid <= 0) {
            throw ApiException.paramError("无效的用户ID");
        }
        User user = userRepository.findById(uid)
                .orElseThrow(() -> ApiException.paramError("用户不存在"));
        return UserResponse.from(user);
    }

    /**
     * 获取用户列表 (按 ID 降序分页)
     * 分页归一化: page<=0→1, page_size<=0→12, page_size>100→100 (与 Go 版 normalizePagination 对齐)
     */
    public UserListResponse list(int page, int pageSize) {
        page = Math.max(page, 1);
        pageSize = pageSize <= 0 ? 12 : Math.min(pageSize, 100);
        Page<User> result = userRepository.findAll(
                PageRequest.of(page - 1, pageSize, Sort.by(Sort.Direction.DESC, "id")));
        List<UserResponse> list = result.getContent().stream().map(UserResponse::from).toList();
        return new UserListResponse((int) result.getTotalElements(), page, pageSize, list);
    }
}
