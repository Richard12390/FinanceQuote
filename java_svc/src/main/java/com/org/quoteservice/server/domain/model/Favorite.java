package com.org.quoteservice.server.domain.model;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class Favorite {
  private Long userId;
  private String symbol;
}
