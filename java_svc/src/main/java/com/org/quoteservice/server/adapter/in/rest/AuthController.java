package com.org.quoteservice.server.adapter.in.rest;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import com.org.quoteservice.server.adapter.in.rest.dto.LoginRequest;
import com.org.quoteservice.server.adapter.in.rest.dto.LoginResponse;
import com.org.quoteservice.server.adapter.in.rest.dto.RegisterRequest;
import com.org.quoteservice.server.adapter.in.rest.dto.RegisterResponse;
import com.org.quoteservice.server.application.service.AuthService;
import com.org.quoteservice.server.common.Result;

import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/auth")
@RequiredArgsConstructor
public class AuthController {

  private final AuthService authService;

  @PostMapping("/login")
  public Result<LoginResponse> login(@RequestBody LoginRequest request) {
    return Result.success(authService.authenticate(request));
  }

  @PostMapping("/register")
  public Result<RegisterResponse> register(@RequestBody RegisterRequest request) {
    return Result.success(authService.register(request));
  }

  @PostMapping("/logout")
  public ResponseEntity<Result<Void>> logout(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader) {
    authService.logout(authHeader);
    return ResponseEntity.ok(Result.success());
  }
}