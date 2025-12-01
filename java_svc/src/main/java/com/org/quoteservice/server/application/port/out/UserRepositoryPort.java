package com.org.quoteservice.server.application.port.out;

import java.util.Optional;

import com.org.quoteservice.server.domain.model.User;

public interface UserRepositoryPort {
  Optional<User> findByAccount(String account);

  boolean existsAccount(String account);

  User save(User user);
}