package br.com.rocket.projeto_integrador.modules.user.repositories;

import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;
import java.util.UUID;

public interface UserRepository extends JpaRepository<UserEntity, UUID> {
    Optional<UserEntity> findByUsernameOrEmail(String username, String email);
    Optional<UserEntity> findByUsername(String username);
    Optional<UserEntity> findById(UUID id);
}
