package com.org.quoteservice.server.application.service;

import java.util.Map;
import java.util.Objects;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import com.org.quoteservice.server.adapter.in.rest.dto.LoginRequest;
import com.org.quoteservice.server.adapter.in.rest.dto.LoginResponse;
import com.org.quoteservice.server.adapter.in.rest.dto.RegisterRequest;
import com.org.quoteservice.server.adapter.in.rest.dto.RegisterResponse;
import com.org.quoteservice.server.application.port.out.UserRepositoryPort;
import com.org.quoteservice.server.common.auth.JwtUtil;
import com.org.quoteservice.server.domain.model.User;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class AuthService {

  private final UserRepositoryPort userRepository;
  private final PasswordEncoder passwordEncoder;

  public LoginResponse authenticate(LoginRequest request) {
    User user = userRepository.findByAccount(request.getAccount())
        .orElseThrow(() -> new ResponseStatusException(HttpStatus.UNAUTHORIZED, "帳號或密碼錯誤"));

    if (request.getUserId() != null && !Objects.equals(user.getId(), request.getUserId())) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "帳號與使用者 ID 不符");
    }

    if (!passwordEncoder.matches(request.getPassword(), user.getPassword())) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "帳號或密碼錯誤");
    }

    Map<String, Object> claims = Map.of(
        "userId", user.getId(),
        "account", user.getAccount());

    String token = JwtUtil.generateToken(user.getAccount(), claims);
    String displayName = Optional.ofNullable(user.getDisplayName()).orElse(user.getAccount());

    return new LoginResponse(token, user.getId(), displayName);
  }

  public RegisterResponse register(RegisterRequest request) {
    if (userRepository.existsAccount(request.getAccount())) {
      throw new ResponseStatusException(HttpStatus.CONFLICT, "帳號已存在");
    }

    User toSave = new User();
    toSave.setAccount(request.getAccount());
    toSave.setPassword(passwordEncoder.encode(request.getPassword()));
    toSave.setDisplayName(request.getDisplayName());

    User saved = userRepository.save(toSave);
    return new RegisterResponse(saved.getId(), saved.getAccount(), saved.getDisplayName());
  }

  public void logout(String authHeader) {
    String token = JwtUtil.extractToken(authHeader);
    JwtUtil.parseToken(token);
  }
}
