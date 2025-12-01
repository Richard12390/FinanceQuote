package com.org.quoteservice.server.adapter.in.rest.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class RegisterRequest {
  @NotBlank
  private String account;

  @NotBlank
  private String password;

  private String displayName;
}