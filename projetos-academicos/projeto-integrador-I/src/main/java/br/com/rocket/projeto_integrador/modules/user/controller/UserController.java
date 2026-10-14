package br.com.rocket.projeto_integrador.modules.user.controller;

import br.com.rocket.projeto_integrador.modules.user.dto.LoginUserDTO;
import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import br.com.rocket.projeto_integrador.modules.user.useCases.CreateUserUseCases;
import br.com.rocket.projeto_integrador.modules.user.useCases.LoginUserUseCase;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private CreateUserUseCases createUserUseCases;
    @Autowired
    private LoginUserUseCase loginUserUseCase;

    @PostMapping("/")
    public ResponseEntity<UserEntity> createUser(@Valid @RequestBody UserEntity userEntity) {
        try {
            var result = this.createUserUseCases.execute(userEntity);
            return ResponseEntity.ok(result);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.badRequest().build();
        }
    }

    @PostMapping("/login")
    public ResponseEntity<UserEntity> updateUser(@Valid @RequestBody LoginUserDTO loginUserDTO) {
        try{
            var result = this.loginUserUseCase.execute(loginUserDTO);
            return ResponseEntity.ok(result);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        }
    }
}
