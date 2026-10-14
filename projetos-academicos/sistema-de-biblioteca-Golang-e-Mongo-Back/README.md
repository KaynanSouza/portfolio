# Biblioteca Virtual API

Bem-vindo à API da Biblioteca Virtual! Com esta API, é possível gerenciar e consultar informações detalhadas sobre livros, usuários e empréstimos em nossa biblioteca digital.

## Documentação Completa

Para acessar a documentação completa e interativa da API, visite o link abaixo:

[Documentação da API da Biblioteca Virtual no Postman](https://documenter.getpostman.com/view/35028840/2sA3QniE4T)

## Funcionalidades Principais

- Administração de Livros: O administrador pode criar, atualizar, deletar e consultar livros disponíveis na biblioteca.
- Administração de Usuários: O administrador pode criar, atualizar, deletar e consultar usuários da biblioteca.
- Empréstimos: Os usuários podem alugar e devolver livros. O administrador pode visualizar todos os empréstimos.


## Como Começar

1. Registro: Registre-se na nossa plataforma para obter uma chave de API.
2. Autenticação: Utilize a chave de API para autenticar suas requisições.
3. Exploração: Use a documentação para explorar e testar os endpoints disponíveis.


## Exemplos de Endpoints

### Cokie
- Adicionar cokie: GET /protected

### Usuarios
- Logar Usuário: POST /user/
- Detalhes do Usuário: GET /user/{username}
- Adicionar Usuário: POST /user/signup
- Adicionar livro ao Usuário: PUT /user/{username}/addBooks
- Remover livro ao Usuário: PUT /user/{username}/removeBooks
- Modificar Usuário: PUT /user/{username}

### Livros
- Listar Livros: GET /books
- Detalhes do Livro: GET /books/{id}
- Adicionar Livro: POST /books

### Admin
- Adicionar admin: POST /admin/signup
- Listar Usuários: GET /admin/users
- Deletar Usuários: DELETE /admin/user/{username}
- Atualizar Livro: PUT /admin/book/{title}
- Deletar Livro: DELETE /admin/books/{title}


## Suporte

Se você encontrar algum problema ou tiver dúvidas, entre em contato com nossa equipe de suporte pelo e-mail: [suporte@biblioteca-virtual.com](https://mailto:suporte@biblioteca-virtual.com).

## Códigos de Status de Erro

- 400 Bad Request: A requisição foi malformada ou contém dados inválidos.
- 401 Unauthorized: A autenticação falhou ou o usuário não tem permissão para acessar o recurso.
- 404 Not Found: O recurso solicitado não foi encontrado.
- 500 Internal Server Error: Ocorreu um erro no servidor que impediu o processamento da requisição.


Aproveite para explorar e integrar a API da Biblioteca Virtual nos seus projetos!