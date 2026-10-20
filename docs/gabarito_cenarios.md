# Gabarito dos cenarios sinteticos

> Gerado por `generate_cenarios.py` aplicando a politica v1.3 de
> `docs/instructions.md`. Regenerado junto com as fixtures.

## Distribuicao

| Faixa | Cenarios | %  |
|---|---|---|
| bom | 5 | 9% |
| complexo | 42 | 80% |
| ruim | 5 | 9% |

| Decisao | Cenarios |
|---|---|
| `APROVADO` | 15 |
| `APROVADO_COM_RESSALVAS` | 13 |
| `ANALISE_MANUAL` | 13 |
| `REPROVADO` | 11 |

## Por cenario

| # | Faixa | O que exercita | Decisao | Pontos | Regra | Criterios sem dado |
|---|---|---|---|---|---|---|
| 1 | bom | assalariado de carreira, quatro fontes concordantes | `APROVADO` | 100 | - | - |
| 2 | bom | aposentado do INSS com consignado em dia, margem preservada | `APROVADO` | 100 | - | - |
| 3 | bom | servidor com dois vinculos vigentes; vale o de maior salario | `APROVADO` | 100 | - | - |
| 4 | ruim | alerta de fraude grave e ativo | `REPROVADO` | - | K1 | - |
| 5 | ruim | score critico, abaixo do piso da politica | `REPROVADO` | - | K4 | - |
| 6 | ruim | pensao por morte movimentada apos obito do titular: severidade decide por K10 | `REPROVADO` | - | K10 | - |
| 7 | complexo | biro otimo e comportamento interno pessimo: fontes discordam | `APROVADO` | 78 | - | - |
| 8 | complexo | vinculo declarado no biro sem contrapartida no eSocial | `APROVADO_COM_RESSALVAS` | 65 | - | - |
| 9 | complexo | fluxo de caixa do open finance desmente a renda declarada | `ANALISE_MANUAL` | - | K8 | - |
| 10 | complexo | inadimplencia interna acima de noventa dias com biro intacto | `REPROVADO` | - | K9 | - |
| 11 | complexo | pessoa exposta politicamente com perfil de credito excelente | `ANALISE_MANUAL` | - | K13 | - |
| 12 | complexo | cliente nunca passou por validacao cadastral | `APROVADO_COM_RESSALVAS` | 63 | - | C13, C14 |
| 13 | complexo | consentimento de open finance revogado | `APROVADO_COM_RESSALVAS` | 59 | - | C7, C8, C9 |
| 14 | complexo | nao-cliente da instituicao, sem historico interno | `APROVADO_COM_RESSALVAS` | 67 | - | C7, C8, C9, C10, C11, C12 |
| 15 | complexo | sem historico de pagamento no biro | `APROVADO_COM_RESSALVAS` | 63 | - | C3 |
| 16 | complexo | sem renda declarada, sem estimativa e sem vinculo vigente | `ANALISE_MANUAL` | - | M2 | - |
| 17 | complexo | biro sem score calculado | `ANALISE_MANUAL` | - | M1 | - |
| 18 | complexo | sem comprometimento informado e sem parcelas mensais para calcular | `APROVADO_COM_RESSALVAS` | 63 | - | C2, C5 |
| 19 | complexo | tres criterios exatamente no limite superior da faixa | `APROVADO` | 92 | - | - |
| 20 | complexo | os mesmos criterios um passo abaixo do limite | `APROVADO` | 76 | - | - |
| 21 | complexo | negativacao ativa e contestada pelo cliente | `ANALISE_MANUAL` | - | M4 | - |
| 22 | complexo | as mesmas dividas, agora acima do teto | `REPROVADO` | - | K7 | - |
| 23 | complexo | CPF ausente no cadastro, aciona M3 | `ANALISE_MANUAL` | - | M3 | - |
| 24 | complexo | fraude em apuracao e CPF pendente simultaneos, vence a de menor numero | `ANALISE_MANUAL` | - | K2 | - |
| 25 | complexo | dois scores de biro com datas diferentes; vale o mais recente | `APROVADO_COM_RESSALVAS` | 65 | - | - |
| 26 | complexo | documento antigo cancelado, validacao mais recente regular nao reprova | `APROVADO` | 76 | - | - |
| 27 | complexo | duas analises de fluxo de caixa em datas distintas; vale a mais recente | `APROVADO_COM_RESSALVAS` | 64 | - | - |
| 28 | complexo | tres pre-aprovacoes ativas: dado de contexto, nao altera a pontuacao | `APROVADO` | 91 | - | - |
| 29 | complexo | trabalhador informal vivendo de Pix; C2 cai no fallback parcelas / renda | `APROVADO_COM_RESSALVAS` | 64 | - | - |
| 30 | complexo | superendividamento regularizado: negativacoes quitadas e dividas renegociadas fora de cobranca (Lei 14.181/2021) | `APROVADO_COM_RESSALVAS` | 67 | - | - |
| 31 | bom | engenheira com vinculo longo, quatro fontes concordantes | `APROVADO` | 100 | - | - |
| 32 | bom | professora estatutaria com consignado em dia e biometria validada | `APROVADO` | 100 | - | - |
| 33 | ruim | score bom, mas duas negativacoes ativas somam mais de cinco rendas (K6) | `REPROVADO` | - | K6 | - |
| 34 | ruim | perfil excelente nas quatro fontes, mas consta em lista restritiva (K12) | `REPROVADO` | - | K12 | - |
| 35 | complexo | pedido de insolvencia civil ativo com biro e registro interno favoraveis (K3) | `REPROVADO` | - | K3 | - |
| 36 | complexo | duas negativacoes ativas e uma de situacao nula: so a nula leva a tres (K5) | `REPROVADO` | - | K5 | - |
| 37 | complexo | alerta de fraude ativo com severidade nula: K2, nunca presumir gravidade | `ANALISE_MANUAL` | - | K2 | - |
| 38 | complexo | divida de situacao nula decide K7: sem ela, a soma fica abaixo do teto | `REPROVADO` | - | K7 | - |
| 39 | complexo | CPF suspenso na Receita como unica restricao (K11 vence) | `ANALISE_MANUAL` | - | K11 | - |
| 40 | complexo | biro limpo e duas parcelas internas recentes com 60 dias de atraso (K9-L vence) | `ANALISE_MANUAL` | - | K9-L | - |
| 41 | complexo | triagem KYC irregular sem PEP e sem lista restritiva (K13, segunda metade) | `ANALISE_MANUAL` | - | K13 | - |
| 42 | complexo | PEP com situacao nula na triagem: avaliar apenas PEP e lista restritiva (K13) | `ANALISE_MANUAL` | - | K13 | - |
| 43 | complexo | perfil de aprovacao com consentimento expirado e utilizacao nula: rebaixamento decide | `APROVADO_COM_RESSALVAS` | 89 | - | C5, C7, C8, C9 |
| 44 | complexo | consentimento ativo sem analise de fluxo de caixa e score interno nulo: rebaixamento decide | `APROVADO_COM_RESSALVAS` | 87 | - | C7, C9, C11 |
| 45 | complexo | consentimento ativo sem receita recorrente e volatilidade nula: C8 zera sem dado ausente | `APROVADO` | 73 | - | - |
| 46 | complexo | nome civil diverge com CPF regular e situacao nula na triagem, sem eliminatoria | `APROVADO` | 78 | - | - |
| 47 | complexo | beneficiaria de programa social, sem triagem de compliance e sem parcelas internas | `APROVADO` | 72 | - | C12 |
| 48 | complexo | negativacao de situacao nula pontuada como ativa e relacionamento encerrado | `APROVADO_COM_RESSALVAS` | 69 | - | - |
| 49 | complexo | acao de cobranca ativa (nao e falencia), vinculo encerrado no eSocial e CND suspensa | `APROVADO` | 72 | - | - |
| 50 | complexo | pontuacao exatamente no limiar de aprovacao (70), com C1 e C2 na fronteira | `APROVADO` | 70 | - | - |
| 51 | complexo | sem eliminatoria, mas pontuacao baixa leva a analise manual pelo scorecard | `ANALISE_MANUAL` | 46 | - | - |
| 52 | complexo | score baixo sem cair no piso, fontes ruins sem eliminatoria: reprova pelo scorecard | `REPROVADO` | 24 | - | - |
