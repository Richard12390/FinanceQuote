package com.org.quoteservice.server.adapter.out;

import java.util.Optional;

import org.springframework.stereotype.Repository;

import com.org.quoteservice.server.adapter.out.persistence.mapper.UserMapper;
import com.org.quoteservice.server.application.port.out.UserRepositoryPort;
import com.org.quoteservice.server.domain.model.User;

import lombok.RequiredArgsConstructor;

@Repository
@RequiredArgsConstructor
public class UserRepositoryAdapter implements UserRepositoryPort {

  private final UserMapper userMapper;

  @Override
  public Optional<User> findByAccount(String account) {
    return Optional.ofNullable(userMapper.findByAccount(account));
  }

  @Override
  public boolean existsAccount(String account) {
    return userMapper.existsAccount(account);
  }

  @Override
  public User save(User user) {
    userMapper.insert(user);
    return user;
  }
}