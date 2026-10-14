package br.com.rocket.projeto_integrador.modules.products.useCases;

import br.com.rocket.projeto_integrador.exeptions.ProductAlreadyExistsException;
import br.com.rocket.projeto_integrador.modules.products.entity.ProductsEntity;
import br.com.rocket.projeto_integrador.modules.products.repository.ProductRepository;
import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import br.com.rocket.projeto_integrador.modules.user.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class CreateProductUseCase {

    @Autowired
    private ProductRepository productRepository;

    public ProductsEntity create(UUID userId, ProductsEntity productsEntity) {

        productRepository.findByNameAndUser_Id(productsEntity.getName(), userId)
                .ifPresent(p -> {
                    throw new ProductAlreadyExistsException();
                });

        UserEntity user = new UserEntity();
        user.setId(userId);
        productsEntity.setUser(user);; // <- Agora setamos o relacionamento

        return productRepository.save(productsEntity);
    }

}
