package com.org.quoteservice.server.adapter.in.rest.dto;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class RegisterResponse {
  private long userId;
  private String account;
  private String displayName;
}