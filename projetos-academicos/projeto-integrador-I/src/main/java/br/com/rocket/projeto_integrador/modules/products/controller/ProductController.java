package br.com.rocket.projeto_integrador.modules.products.controller;

import br.com.rocket.projeto_integrador.modules.products.entity.ProductsEntity;
import br.com.rocket.projeto_integrador.modules.products.useCases.CreateProductUseCase;
import br.com.rocket.projeto_integrador.modules.products.useCases.ListProductUseCase;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/{id}/product")
public class ProductController {

    @Autowired
    private CreateProductUseCase createProductUseCase;

    @Autowired
    private ListProductUseCase listProductUseCase;

    @PostMapping("/create")
    public ResponseEntity<ProductsEntity> create(@PathVariable("id") UUID userId, @Valid @RequestBody ProductsEntity productsEntity) {
        try{
            var result = this.createProductUseCase.create(userId, productsEntity);
            return ResponseEntity.ok(result);
        } catch (Exception e){
            e.printStackTrace();
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping("/list")
    public ResponseEntity<List<ProductsEntity>> list(@PathVariable("id") UUID userId) {
        try{
            var result = this.listProductUseCase.listAllProducts(userId);
            return ResponseEntity.ok(result);
        } catch ( Exception e){
            e.printStackTrace();
            return ResponseEntity.badRequest().build();
        }
    }
}
