package com.org.quoteservice.server.application.service;

import java.util.List;

import org.springframework.stereotype.Service;

import com.org.quoteservice.server.application.port.in.FavoriteUseCase;
import com.org.quoteservice.server.application.port.out.FavoriteRepositoryPort;
import com.org.quoteservice.server.domain.model.Favorite;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class FavoriteService implements FavoriteUseCase {

  private final FavoriteRepositoryPort favoriteRepositoryPort;

  @Override
  public List<Favorite> getByUserId(long userId) {
    return favoriteRepositoryPort.getByUserId(userId);
  }

  @Override
  public void add(long userId, String symbol) {
    favoriteRepositoryPort.save(new Favorite(userId, symbol));
  }

  @Override
  public void remove(long userId, String symbol) {
    favoriteRepositoryPort.delete(userId, symbol);
  }
}
