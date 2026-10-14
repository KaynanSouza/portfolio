package br.com.rocket.projeto_integrador.modules.products.useCases;

import br.com.rocket.projeto_integrador.modules.products.entity.ProductsEntity;
import br.com.rocket.projeto_integrador.modules.products.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class ListProductUseCase {

    @Autowired
    private ProductRepository productRepository;

    public List<ProductsEntity> listAllProducts(UUID userId) {
        return productRepository.findByUser_Id(userId);
    }
}
