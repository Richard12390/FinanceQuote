package com.org.quoteservice.server.application.port.out;

import java.util.List;

import com.org.quoteservice.server.domain.model.Favorite;

public interface FavoriteRepositoryPort {
  List<Favorite> getByUserId(long userId);

  void save(Favorite favorite);

  void delete(long userId, String symbol);
}
