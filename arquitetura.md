Princípios de Arquitetura e Boas Práticas: signet-go
Status: Documento Vivo
Última Atualização: 1 de julho de 2025

Introdução
Este documento serve como a fonte canónica da verdade para as decisões de arquitetura, padrões de código e boas práticas adotadas no desenvolvimento da biblioteca signet-go e de todo o ecossistema Signet. O seu objetivo é garantir a consistência, qualidade, manutenibilidade e segurança do nosso código.

Todos os desenvolvedores que contribuem para o projeto DEVEM aderir aos princípios aqui descritos. Este é um documento vivo, que será enriquecido à medida que novas decisões forem tomadas e novas lições forem aprendidas.

1. Arquitetura em Camadas (Clean Architecture)
O signet-go adota uma arquitetura em camadas rigorosa para garantir a separação de responsabilidades e um fluxo de dependências unidirecional.

Camada de Núcleo (internal/core): O coração criptográfico. Contém lógica pura, sem conhecimento do mundo externo ou das regras de negócio do Signet. Depende apenas da biblioteca padrão do Go.

Camada de Aplicação (Pacote signet): A API pública. Orquestra a lógica, conhece as regras da especificação Signet (validação de claims, perfis) e depende da Camada de Núcleo. Não conhece protocolos de transporte.

Camada de Adaptadores (Pacotes como grpcinterceptor): Adapta a lógica do Signet para casos de uso específicos (ex: gRPC). Depende da Camada de Aplicação. A Camada de Aplicação NUNCA deve depender de um Adaptador.

Regra Fundamental: As dependências fluem sempre de fora para dentro (Adaptadores -> Aplicação -> Núcleo).

2. Padrões de Código e Boas Práticas
2.1. Tratamento de Erros: Sentinel Errors
Para garantir um tratamento de erros robusto e programaticamente verificável, o projeto DEVE utilizar o padrão de Sentinel Errors.

Definição: Erros comuns e reutilizáveis devem ser declarados como variáveis exportadas no nível do pacote onde ocorrem (ex: var ErrInvalidPrivateKey = errors.New(...)).

Motivação: Isso permite que o código que chama uma função verifique o tipo de erro de forma confiável usando errors.Is(), em vez de recorrer à frágil prática de comparar strings de erro.

Aplicação: Todas as funções que podem retornar um conjunto previsível de erros devem usar este padrão.

2.2. Idioma e Convenções de Comentários
Para garantir a clareza e a acessibilidade internacional do código, ao mesmo tempo que se mantém uma documentação explicativa rica, adotamos a seguinte convenção:

Idioma do Código: Todo o código — nomes de variáveis, funções, tipos, etc. — DEVE ser escrito em Inglês Americano (en-US).

Idioma dos Comentários: Comentários explicativos, que detalham o "porquê" de uma implementação ou adicionam contexto, DEVEM ser escritos em Português do Brasil (pt-BR). A documentação GoDoc, que é extraída dos comentários e serve como documentação da API, seguirá a mesma regra.