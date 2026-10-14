package br.com.rocket.projeto_integrador.exeptions;

public class AuthenticationException extends RuntimeException {
    public AuthenticationException() {
        super("password incorreto");
    }
}
