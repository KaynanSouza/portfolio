package br.com.rocket.projeto_integrador.modules.user.useCases;

import br.com.rocket.projeto_integrador.exeptions.UserNotFoundException;
import br.com.rocket.projeto_integrador.modules.user.dto.LoginUserDTO;
import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import br.com.rocket.projeto_integrador.modules.user.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class LoginUserUseCase {

    @Autowired
    private UserRepository userRepository;

    public UserEntity execute(LoginUserDTO loginUserDTO) {
        var user = this.userRepository.findByUsername(loginUserDTO.username())
                .orElseThrow(() -> {
                    throw new UserNotFoundException();
                });

        if (!user.getPassword().equals(loginUserDTO.password())) {
            throw new UserNotFoundException();
        }

        return user;
    }
}
