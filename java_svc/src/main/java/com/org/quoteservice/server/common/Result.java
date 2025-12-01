package com.org.quoteservice.server.common;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class Result<T> {

  private boolean success;
  private String message;
  private T data;

  public static <T> Result<T> success(T object) {
    return new Result<>(true, "success", object);
  }

  public static Result<Void> success() {
    return success(null);
  }

  public static <T> Result<T> error(String message, T object) {
    return new Result<>(false, message, object);
  }
}
