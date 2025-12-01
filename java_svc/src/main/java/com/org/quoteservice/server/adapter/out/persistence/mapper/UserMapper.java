package com.org.quoteservice.server.adapter.out.persistence.mapper;

import org.apache.ibatis.annotations.Param;

import com.org.quoteservice.server.domain.model.User;

public interface UserMapper {

  User findByAccount(@Param("account") String account);

  boolean existsAccount(@Param("account") String account);

  void insert(User user);
}
