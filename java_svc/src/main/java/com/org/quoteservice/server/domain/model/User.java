package com.org.quoteservice.server.domain.model;

import java.time.LocalDateTime;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class User {
  private Long id;
  private String account;
  private String password;
  private String displayName;
  private LocalDateTime createdAt;
}
