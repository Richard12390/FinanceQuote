package com.org.quoteservice.server.adapter.in.rest.dto;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class LoginResponse {
  private String token;
  private long userId;
  private String displayName;
}
