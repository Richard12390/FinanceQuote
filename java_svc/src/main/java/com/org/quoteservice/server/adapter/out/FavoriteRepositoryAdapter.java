package com.org.quoteservice.server.adapter.out;

import java.util.List;

import org.springframework.stereotype.Repository;

import com.org.quoteservice.server.adapter.out.persistence.mapper.FavoriteMapper;
import com.org.quoteservice.server.application.port.out.FavoriteRepositoryPort;
import com.org.quoteservice.server.domain.model.Favorite;

import lombok.RequiredArgsConstructor;

@Repository
@RequiredArgsConstructor
public class FavoriteRepositoryAdapter implements FavoriteRepositoryPort {

  private final FavoriteMapper favoriteMapper;

  @Override
  public List<Favorite> getByUserId(long userId) {
    return favoriteMapper.getByUserId(userId);
  }

  @Override
  public void save(Favorite favorite) {
    favoriteMapper.upsert(favorite);
  }

  @Override
  public void delete(long userId, String symbol) {
    favoriteMapper.delete(userId, symbol);
  }
}

