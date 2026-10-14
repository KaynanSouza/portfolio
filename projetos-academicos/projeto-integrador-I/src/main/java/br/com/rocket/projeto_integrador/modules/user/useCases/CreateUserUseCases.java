package br.com.rocket.projeto_integrador.modules.user.useCases;

import br.com.rocket.projeto_integrador.exeptions.UserFoundException;
import br.com.rocket.projeto_integrador.modules.user.entity.UserEntity;
import br.com.rocket.projeto_integrador.modules.user.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class CreateUserUseCases {

    @Autowired
    private UserRepository userRepository;

    public UserEntity execute(UserEntity userEntity) {
        this.userRepository.findByUsernameOrEmail(userEntity.getUsername(), userEntity.getEmail())
                .ifPresent((user) ->{
                    throw new UserFoundException();
                });

        return this.userRepository.save(userEntity);
    }
}
