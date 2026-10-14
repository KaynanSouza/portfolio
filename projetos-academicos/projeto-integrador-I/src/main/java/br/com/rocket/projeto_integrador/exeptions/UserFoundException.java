package br.com.rocket.projeto_integrador.exeptions;

public class UserFoundException extends RuntimeException {
    public UserFoundException() {
        super("Usuário já existente");
    }
}
