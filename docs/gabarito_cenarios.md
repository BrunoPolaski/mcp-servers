# Gabarito dos cenarios sinteticos

> Gerado por `generate_cenarios.py` aplicando a politica v1.3 de
> `docs/instructions.md`. Regenerado junto com as fixtures.

## Distribuicao

| Faixa | Cenarios | %  |
|---|---|---|
| bom | 8 | 10% |
| complexo | 64 | 80% |
| ruim | 8 | 10% |

| Decisao | Cenarios |
|---|---|
| `APROVADO` | 29 |
| `APROVADO_COM_RESSALVAS` | 22 |
| `ANALISE_MANUAL` | 16 |
| `REPROVADO` | 13 |

## Por cenario

| # | Faixa | O que exercita | Decisao | Pontos | Regra | Limite | Criterios sem dado |
|---|---|---|---|---|---|---|---|
| 1 | bom | assalariado de carreira, quatro fontes concordantes | `APROVADO` | 100 | - | R$ 40.000 | - |
| 2 | bom | CLT estavel, relacionamento de sete anos | `APROVADO` | 100 | - | R$ 22.000 | - |
| 3 | bom | alta renda, sem passivo relevante | `APROVADO` | 100 | - | R$ 90.000 | - |
| 4 | bom | servidor publico, renda comprovada | `APROVADO` | 95 | - | R$ 14.000 | - |
| 5 | bom | MEI consolidado, faturamento estavel | `APROVADO` | 100 | - | R$ 30.000 | - |
| 6 | bom | CLT com poupanca, sem restricoes | `APROVADO` | 97 | - | R$ 18.000 | - |
| 7 | ruim | alerta de fraude grave e ativo | `REPROVADO` | - | K1 | - | - |
| 8 | ruim | score critico, abaixo do piso da politica | `REPROVADO` | - | K4 | - | - |
| 9 | ruim | tres negativacoes ativas | `REPROVADO` | - | K5 | - | - |
| 10 | ruim | CPF cancelado por obito na Receita Federal | `REPROVADO` | - | K10 | - | - |
| 11 | ruim | homonimo confirmado em lista restritiva | `REPROVADO` | - | K12 | - | - |
| 12 | ruim | falencia ativa somada a dividas em cobranca acima do teto | `REPROVADO` | - | K3 | - | - |
| 13 | complexo | biro otimo e comportamento interno pessimo: fontes discordam | `APROVADO` | 78 | - | R$ 15.000 | - |
| 14 | complexo | biro ruim e relacionamento interno impecavel de oito anos | `APROVADO_COM_RESSALVAS` | 64 | - | R$ 3.800 | - |
| 15 | complexo | vinculo declarado no biro sem contrapartida no eSocial | `APROVADO_COM_RESSALVAS` | 69 | - | R$ 4.800 | - |
| 16 | complexo | eSocial mostra o vinculo encerrado; o biro ainda o da como vigente | `APROVADO_COM_RESSALVAS` | 68 | - | R$ 5.000 | - |
| 17 | complexo | nome diverge da Receita sem alerta de fraude correspondente no biro | `APROVADO` | 80 | - | R$ 16.000 | - |
| 18 | complexo | renda declarada muito acima do que o open finance observa | `ANALISE_MANUAL` | - | K8 | - | - |
| 19 | complexo | fluxo de caixa do open finance desmente a renda declarada | `ANALISE_MANUAL` | - | K8 | - | - |
| 20 | complexo | parcela interna nao paga ha dois meses, biro ainda nao reflete | `ANALISE_MANUAL` | - | K9-L | - | - |
| 21 | complexo | inadimplencia interna acima de noventa dias com biro intacto | `REPROVADO` | - | K9 | - | - |
| 22 | complexo | atrasos internos curtos e recorrentes, nenhum caracterizado como perda | `APROVADO_COM_RESSALVAS` | 58 | - | R$ 3.700 | - |
| 23 | complexo | CPF pendente de regularizacao na Receita | `ANALISE_MANUAL` | - | K11 | - | - |
| 24 | complexo | pessoa exposta politicamente com perfil de credito excelente | `ANALISE_MANUAL` | - | K13 | - | - |
| 25 | complexo | CPF suspenso na Receita Federal | `ANALISE_MANUAL` | - | K11 | - | - |
| 26 | complexo | certidao negativa de debitos negada | `APROVADO` | 71 | - | R$ 10.200 | - |
| 27 | complexo | debitos com exigibilidade suspensa: certidao positiva com efeito de negativa | `APROVADO` | 76 | - | R$ 10.200 | - |
| 28 | complexo | cliente nunca passou por validacao cadastral | `APROVADO_COM_RESSALVAS` | 70 | - | R$ 5.600 | C13, C14 |
| 29 | complexo | sem triagem de compliance no cadastro | `APROVADO` | 79 | - | R$ 10.800 | - |
| 30 | complexo | data de nascimento diverge da Receita, biometria nao validada | `APROVADO_COM_RESSALVAS` | 55 | - | R$ 3.500 | - |
| 31 | complexo | nome e nascimento divergem da Receita ao mesmo tempo | `ANALISE_MANUAL` | 42 | - | - | - |
| 32 | complexo | consentimento de open finance revogado | `APROVADO_COM_RESSALVAS` | 64 | - | R$ 4.700 | C7, C8, C9 |
| 33 | complexo | consentimento de open finance expirado | `APROVADO_COM_RESSALVAS` | 57 | - | R$ 3.700 | C7, C8, C9 |
| 34 | complexo | nao-cliente da instituicao, sem historico interno | `APROVADO_COM_RESSALVAS` | 67 | - | R$ 6.500 | C7, C8, C9, C10, C11, C12 |
| 35 | complexo | conta encerrada, relacionamento inativo | `APROVADO_COM_RESSALVAS` | 55 | - | R$ 3.600 | - |
| 36 | complexo | score interno nunca calculado | `APROVADO` | 74 | - | R$ 10.700 | C11 |
| 37 | complexo | cliente novo: sem parcelas internas no periodo | `APROVADO` | 73 | - | R$ 10.600 | C12 |
| 38 | complexo | sem historico de pagamento no biro | `APROVADO_COM_RESSALVAS` | 63 | - | R$ 3.800 | C3 |
| 39 | complexo | autonomo sem renda declarada; biro estima e o fluxo confirma | `APROVADO` | 70 | - | R$ 6.800 | - |
| 40 | complexo | sem renda declarada, sem estimativa e sem vinculo vigente | `ANALISE_MANUAL` | - | M2 | - | - |
| 41 | complexo | biro sem score calculado | `ANALISE_MANUAL` | - | M1 | - | - |
| 42 | complexo | sem comprometimento informado e sem parcelas mensais para calcular | `APROVADO_COM_RESSALVAS` | 67 | - | R$ 3.600 | C2, C5 |
| 43 | complexo | tres criterios exatamente no limite superior da faixa | `APROVADO` | 96 | - | R$ 12.600 | - |
| 44 | complexo | os mesmos criterios um passo abaixo do limite | `APROVADO` | 76 | - | R$ 12.400 | - |
| 45 | complexo | limites superiores das faixas intermediarias | `APROVADO` | 87 | - | R$ 12.600 | - |
| 46 | complexo | um passo abaixo em todas as faixas anteriores | `APROVADO_COM_RESSALVAS` | 67 | - | R$ 6.100 | - |
| 47 | complexo | todos os criterios no piso da faixa que ainda pontua | `APROVADO_COM_RESSALVAS` | 55 | - | R$ 3.200 | - |
| 48 | complexo | todos os criterios um passo abaixo do piso | `REPROVADO` | 31 | - | - | - |
| 49 | complexo | negativacao ativa de valor baixo | `APROVADO_COM_RESSALVAS` | 63 | - | R$ 3.600 | - |
| 50 | complexo | negativacao ativa de valor relevante | `APROVADO_COM_RESSALVAS` | 54 | - | R$ 3.200 | - |
| 51 | complexo | duas negativacoes ativas somando pouco | `ANALISE_MANUAL` | 45 | - | - | - |
| 52 | complexo | negativacao ativa e contestada pelo cliente | `ANALISE_MANUAL` | - | M4 | - | - |
| 53 | complexo | dividas em cobranca logo abaixo do teto da politica | `APROVADO` | 72 | - | R$ 21.000 | - |
| 54 | complexo | as mesmas dividas, agora acima do teto | `REPROVADO` | - | K7 | - | - |
| 55 | complexo | acao civel em andamento, sem penhora | `APROVADO_COM_RESSALVAS` | 68 | - | R$ 4.400 | - |
| 56 | complexo | alerta de fraude de baixa severidade em apuracao | `ANALISE_MANUAL` | - | K2 | - | - |
| 57 | complexo | fluxo de caixa magro limita o credito apesar do bom perfil | `APROVADO` | 79 | - | R$ 4.000 | - |
| 58 | complexo | pre-aprovacao vigente acima do limite calculado pela politica | `APROVADO` | 83 | - | R$ 42.000 | - |
| 59 | complexo | pre-aprovacao existe, mas ja expirou | `APROVADO` | 81 | - | R$ 12.600 | - |
| 60 | complexo | muitos sinais fracos somados, nenhum eliminatorio isolado | `REPROVADO` | 27 | - | - | - |
| 61 | bom | aposentado do INSS com consignado em dia, margem preservada | `APROVADO` | 100 | - | R$ 10.000 | - |
| 62 | bom | servidor com dois vinculos vigentes; vale o de maior salario | `APROVADO` | 100 | - | R$ 22.800 | - |
| 63 | ruim | CPF cancelado na Receita Federal | `REPROVADO` | - | K10 | - | - |
| 64 | ruim | duas negativacoes ativas somando mais de cinco vezes a renda | `REPROVADO` | - | K6 | - | - |
| 65 | complexo | CPF ausente no cadastro, aciona M3 | `ANALISE_MANUAL` | - | M3 | - | - |
| 66 | complexo | alerta de fraude ativo com severidade nula, nunca presumir gravidade | `ANALISE_MANUAL` | - | K2 | - | - |
| 67 | complexo | negativacao e divida com situacao nula, tratamento conservador sem knock-out | `APROVADO` | 72 | - | R$ 9.700 | - |
| 68 | complexo | pensao por morte movimentada apos obito do titular: severidade decide por K10 | `REPROVADO` | - | K10 | - | - |
| 69 | complexo | fraude em apuracao e CPF pendente simultaneos, vence a de menor numero | `ANALISE_MANUAL` | - | K2 | - | - |
| 70 | complexo | triagem de compliance irregular sem PEP e sem lista restritiva | `ANALISE_MANUAL` | - | K13 | - | - |
| 71 | complexo | dois scores de biro com datas diferentes; vale o mais recente | `APROVADO_COM_RESSALVAS` | 67 | - | R$ 4.200 | - |
| 72 | complexo | documento antigo cancelado, validacao mais recente regular nao reprova | `APROVADO` | 76 | - | R$ 10.800 | - |
| 73 | complexo | duas analises de fluxo de caixa em datas distintas; vale a mais recente | `APROVADO_COM_RESSALVAS` | 68 | - | R$ 3.600 | - |
| 74 | complexo | tres pre-aprovacoes ativas do mesmo tipo; vale a de maior valor como piso | `APROVADO` | 95 | - | R$ 35.000 | - |
| 75 | complexo | beneficiario de beneficio social como unica renda, cadastro pobre | `APROVADO_COM_RESSALVAS` | 62 | - | R$ 1.800 | - |
| 76 | complexo | trabalhador informal vivendo de Pix; C2 cai no fallback parcelas / renda | `APROVADO_COM_RESSALVAS` | 64 | - | R$ 3.400 | - |
| 77 | complexo | fluxo com padrao de conta intermediaria: entradas e saidas altas e identicas, sem alerta ainda no biro | `APROVADO` | 74 | - | R$ 2.400 | - |
| 78 | complexo | superendividamento regularizado: negativacoes quitadas e dividas renegociadas fora de cobranca (Lei 14.181/2021) | `APROVADO_COM_RESSALVAS` | 67 | - | R$ 3.000 | - |
| 79 | complexo | mudanca de endereco nao propagada ao biro, financiamento de veiculo e acao civel em curso | `APROVADO` | 76 | - | R$ 9.900 | - |
| 80 | complexo | divergencia de nome por casamento, biometria validada e sem alerta de fraude | `APROVADO` | 73 | - | R$ 11.200 | - |
