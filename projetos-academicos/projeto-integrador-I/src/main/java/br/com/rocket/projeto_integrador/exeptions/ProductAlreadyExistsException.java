package br.com.rocket.projeto_integrador.exeptions;

public class ProductAlreadyExistsException extends RuntimeException {
    public ProductAlreadyExistsException() {
        super("Produto já existente");
    }
}
