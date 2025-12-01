package com.org.quoteservice.server.adapter.in.rest.dto;

import lombok.Data;

@Data
public class LoginRequest {
  private String account;
  private Long userId;
  private String password;
}
