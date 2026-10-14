package br.com.rocket.projeto_integrador.modules.products.repository;

import br.com.rocket.projeto_integrador.modules.products.entity.ProductsEntity;
import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface ProductRepository extends JpaRepository<ProductsEntity, UUID> {
    Optional<ProductsEntity> findByNameAndUser_Id(String name, UUID userId);
    List<ProductsEntity> findByUser_Id(UUID userId);

}
