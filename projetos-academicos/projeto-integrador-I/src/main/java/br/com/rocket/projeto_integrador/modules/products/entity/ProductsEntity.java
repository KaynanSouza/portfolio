package br.com.rocket.projeto_integrador.modules.products.entity;

import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import jakarta.persistence.*;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.Data;
import org.springframework.format.annotation.NumberFormat;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@Entity
@Table(name = "produtos")
public class ProductsEntity {

    @Id
    @GeneratedValue
    private UUID id;

    @NotBlank
    private String name;

    private String description;

    @NotNull
    @Positive
    private Float price;

    @NotNull
    @Min(0)
    private Integer stock;

    @ManyToOne
    @JoinColumn(name = "user_id", nullable = false)
    private UserEntity user;
}
