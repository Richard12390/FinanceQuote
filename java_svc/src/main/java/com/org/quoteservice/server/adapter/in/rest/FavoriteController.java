package com.org.quoteservice.server.adapter.in.rest;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

import com.org.quoteservice.server.application.port.in.FavoriteUseCase;
import com.org.quoteservice.server.common.Result;
import com.org.quoteservice.server.common.auth.JwtUserPrincipal;
import com.org.quoteservice.server.domain.model.Favorite;

import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/favorites")
@RequiredArgsConstructor
public class FavoriteController {

  private final FavoriteUseCase favoriteUseCase;

  @GetMapping
  public Result<List<Favorite>> listFavorites(@AuthenticationPrincipal JwtUserPrincipal principal) {
    long userId = requireUserId(principal);
    return Result.success(favoriteUseCase.getByUserId(userId));
  }

  @PutMapping("/{symbol}")
  public Result<Void> addFavorite(@PathVariable String symbol,
      @AuthenticationPrincipal JwtUserPrincipal principal) {
    long userId = requireUserId(principal);
    favoriteUseCase.add(userId, symbol);
    return Result.success();
  }

  @DeleteMapping("/{symbol}")
  public Result<Void> removeFavorite(@PathVariable String symbol,
      @AuthenticationPrincipal JwtUserPrincipal principal) {
    long userId = requireUserId(principal);
    favoriteUseCase.remove(userId, symbol);
    return Result.success();
  }

  private long requireUserId(JwtUserPrincipal principal) {
    if (principal == null || principal.userId() == null) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "User is not authenticated");
    }
    return principal.userId();
  }

}
