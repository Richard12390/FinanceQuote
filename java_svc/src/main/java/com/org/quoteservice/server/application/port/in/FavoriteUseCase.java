package com.org.quoteservice.server.application.port.in;

import java.util.List;

import com.org.quoteservice.server.domain.model.Favorite;

public interface FavoriteUseCase {

  List<Favorite> getByUserId(long userId);

  void add(long userId, String symbol);

  void remove(long userId, String symbol);
}
