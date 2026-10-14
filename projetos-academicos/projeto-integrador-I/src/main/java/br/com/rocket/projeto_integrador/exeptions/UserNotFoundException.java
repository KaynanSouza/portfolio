package br.com.rocket.projeto_integrador.exeptions;

public class UserNotFoundException extends RuntimeException {
    public UserNotFoundException() {
        super("Username/password incorreto");
    }
}
