# Política de Segurança

## Versões com suporte

Recebe correções de segurança a versão mais recente publicada na branch `main` (última tag `v*`).

## Como reportar uma vulnerabilidade

Não abra uma issue pública para falhas de segurança.

1. Use o **relato privado de vulnerabilidades** do GitHub: aba **Security → Report a vulnerability**.
2. Se essa opção não estiver disponível, envie uma mensagem privada pelo [LinkedIn](https://www.linkedin.com/in/kaynansouza) descrevendo o problema, os passos para reproduzir e o impacto.

Respondo em até 5 dias úteis com a avaliação inicial.

## O que é tratado como segurança

- Credenciais, tokens ou chaves expostos em qualquer arquivo ou no histórico do Git.
- Dependências com vulnerabilidades conhecidas.
- Falhas nos projetos armazenados que permitam acesso indevido a dados.

## Boas práticas adotadas neste repositório

- Autenticação em dois fatores na conta responsável.
- Branch `main` protegida, com Pull Request obrigatório e force push bloqueado.
- Secret scanning e push protection ativos.
- Dependabot para alertas e atualizações.
- `.gitignore` cobrindo `.env`, chaves e artefatos de build.

## Se você encontrar um segredo exposto

Avise pelo canal privado acima. O segredo será revogado e substituído; apenas remover o arquivo do repositório não é suficiente, pois ele continua no histórico.
