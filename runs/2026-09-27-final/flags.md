# Sinalizações dos pareceres da rodada definitiva

Esta é a classificação manual das inconsistências que os 30 pareceres desta
rodada apontaram e que a política não avalia. Ela sustenta as contagens do
artigo (Resultados e discussões: 37 sinalizações em 20 pareceres) e da tese. Cada sinalização
foi confrontada com `generate_cenarios.py` e `generate_cadastro.py`.

## Regras de contagem

- Conta-se cada inconsistência distinta entre campos ou entre fontes que o
  parecer aponta, uma vez por parecer. Subitens de um mesmo tópico contam
  separadamente.
- Ficam de fora:
  - as observações de risco, como pedidos de crédito, comprometimento alto,
    fluxo de caixa em queda, risco de evasão e salário que deixou de entrar;
  - os pontos que a política já trata: listas vazias, vínculo sem
    contrapartida no eSocial, consentimento, pré-condições e regras
    eliminatórias, mesmo quando não acionadas;
  - as menções de que uma regra não se aplica.
- **Planejado.** Conflito que o cenário foi desenhado para conter.
- **Artefato.** Inconsistência que o gerador produz ao criar campos de uma
  mesma pessoa por regras independentes.

## Classificação

| Cenário | Sinalização | Origem |
|---|---|---|
| 1 | Sexo incompatível com o nome | Artefato |
| 2 | Aposentadoria no INSS iniciada aos 15 anos | Artefato |
| 2 | Sexo incompatível com o nome | Artefato |
| 2 | Consignado ativo sem parcelas internas registradas | Artefato |
| 2 | E-mail e telefone verificados no birô, mas não na validação cadastral | Artefato |
| 3 | Segundo vínculo (PJ, R$ 2.200) ausente das receitas recorrentes do Open Finance | Artefato |
| 3 | Sexo incompatível com o nome | Artefato |
| 6 | Titular falecido com a conta em uso | Planejado |
| 6 | Sexo incompatível com o nome | Artefato |
| 6 | E-mail e telefone verificados no birô, mas não na validação cadastral | Artefato |
| 7 | Score interno de 372 contra 812 no birô | Planejado |
| 7 | Telefone do birô diferente das outras fontes | Planejado |
| 7 | Endereço do birô diferente das outras fontes | Planejado |
| 7 | Sexo incompatível com o nome | Artefato |
| 7 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |
| 9 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |
| 10 | Dias de atraso incompatíveis com as datas de vencimento | Artefato |
| 10 | Inadimplência interna com o birô em dia | Planejado |
| 11 | Declaração de renda do tipo autônomo com origem em folha de pagamento | Artefato |
| 11 | Sexo incompatível com o nome | Artefato |
| 12 | Empréstimo pessoal sem parcelas registradas | Artefato |
| 12 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |
| 13 | Sexo incompatível com o nome | Artefato |
| 15 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |
| 15 | Sexo incompatível com o nome | Artefato |
| 16 | Comprometimento de renda de 0,44 com parcelas mensais de R$ 0 no birô | Artefato |
| 16 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |
| 16 | Sexo incompatível com o nome | Artefato |
| 17 | Avaliação de risco do birô com nota 0 e nível alto, mas fatores positivos | Artefato |
| 18 | Empréstimo interno sem registro no birô | Artefato |
| 23 | Validação na Receita declarada válida para um número de documento em branco | Artefato |
| 26 | CPF cancelado em validação anterior | Planejado |
| 27 | Sexo incompatível com o nome | Artefato |
| 28 | Três pré-aprovações simultâneas de valores diferentes | Planejado |
| 29 | Parcelas internas anteriores ao início do relacionamento | Artefato |
| 30 | Empréstimo, consignado e cartão ativos sem parcelas internas registradas | Artefato |
| 30 | Parcela do empréstimo diferente entre o registro interno e o Open Finance | Artefato |

## Totais

| | Sinalizações | Pareceres |
|---|---|---|
| Planejadas | 7 | 5 (6, 7, 10, 26, 28) |
| Artefatos | 30 | 18 |
| Total | 37 | 20 |

Dez pareceres sem sinalização: 4, 5, 8, 14, 19, 20, 21, 22, 24, 25.

O sexo incompatível com o nome está nos 30 cenários: o gerador alterna o sexo
pelo índice, e os nomes alternam entre feminino e masculino na ordem inversa.
Os pareceres o apontaram em 10 cenários (1, 2, 3, 6, 7, 11, 13, 15, 16, 27). Nas
150 sessões de `2026-10-18-repeat`, o apontaram 51. Essa contagem vem de uma
busca por "sexo", "gênero", "feminino" e "masculino" fora do bloco estruturado,
conferida manualmente.

## Casos de fronteira

- Cenário 16: o salário que deixou de constar das receitas recorrentes
  ativas foi contado como observação de risco, porque indica perda de renda,
  e não como inconsistência entre campos.
- Cenário 23: o CPF vazio é a pré-condição M3, que a política trata. Conta-se
  apenas a validação declarada válida sobre um documento em branco, que a
  política não avalia.
- Cenário 12: o empréstimo sem parcelas e a parcela diferente no Open Finance
  aparecem no mesmo tópico, mas são duas inconsistências.
