package com.laixhe.mapper;

import org.apache.ibatis.annotations.Mapper;
import com.mybatisflex.core.BaseMapper;

import com.laixhe.entity.User;

/**
 * @author laixhe
 */
@Mapper
public interface UserMapper extends BaseMapper<User> {
}
