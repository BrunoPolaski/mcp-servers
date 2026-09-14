#!/usr/bin/env python3
"""Gera os 60 cenarios de avaliacao nas quatro fontes.

Distribuicao deliberada, calcada na realidade de uma carteira: casos limpos e
casos obviamente reprovaveis sao raros; a maior parte exige leitura conjunta
das quatro fontes, que frequentemente discordam entre si.

  bom (1-6)        10%  historico alinhado nas quatro fontes, aprovacao direta
  ruim (7-12)      10%  regra eliminatoria evidente, reprovacao direta
  complexo (13-60) 80%  sinais conflitantes, dados faltantes, valores em cima
                        das fronteiras de faixa e de decisao

Cada cenario e declarado uma unica vez em CENARIOS e derivado dali para as
quatro fontes, de modo que elas contem a mesma historia - inclusive quando a
historia e contraditoria de proposito: score de biro alto com comportamento
interno pessimo, CPF regular na Receita com divergencia de nome, vinculo
declarado no biro sem contrapartida no eSocial.

As lacunas de dados (`falta`, `consentimento`, `relacionamento`, `cadastral`)
sao deliberadas: existem para exercitar as regras de dado ausente da politica.

Roda a partir de mcp-servers/: python3 generate_cenarios.py
Depois rode generate_cadastro.py, que propaga o cadastro para as demais fontes.
"""

import hashlib
import json
from datetime import datetime, timedelta
from pathlib import Path

ROOT = Path(__file__).parent
HOJE = datetime(2026, 8, 31)
TS = "2026-08-31T00:00:00Z"
PRIMEIRO_ID = 1

BANCOS = {
    "alfa": ("Banco Sintético Alfa", "11222333000181"),
    "beta": ("Banco Sintético Beta", "22333444000172"),
    "gama": ("Banco Sintético Gama", "33444555000163"),
    "delta": ("Banco Sintético Delta", "44555666000154"),
}
CIDADES = [("Sao Paulo", "SP"), ("Curitiba", "PR"), ("Recife", "PE"), ("Salvador", "BA"),
           ("Belo Horizonte", "MG"), ("Porto Alegre", "RS"), ("Fortaleza", "CE"), ("Goiania", "GO")]
RUAS = ["Rua das Acacias", "Avenida Central", "Rua Sete de Setembro", "Rua Sao Jorge",
        "Avenida Brasil", "Rua dos Ipes", "Rua Marechal Deodoro", "Avenida Getulio Vargas"]

# ---------------------------------------------------------------------------
# Cada chave e um knob; o resto e derivado.
#
#   score        None = biro sem score (falha a pre-condicao de dados minimos);
#                int = um so registro; lista de (score, data) = varios registros
#                em credit_scores (usa-se o de data mais recente, §5 C1)
#   renda        renda mensal declarada; None = nao declarada
#   estimada     renda estimada pelo biro; None = nao estimada
#   dti          comprometimento de renda
#   utilizacao   utilizacao do limite rotativo (fracao de 0 a 1)
#   pontual      proporcao de parcelas pagas em dia no biro, ultimos 12 meses
#   negativacoes (valor, status, contestada) - status: active | paid | None
#                (None = situacao nula, tratada como ativa - §8)
#   dividas      (valor, em_cobranca) - em_cobranca: True | False | None
#                (None = situacao nula, tratada como nao quitada - §8)
#   processos    tipos de acao judicial
#   fraude       (severidade, status) ou None; severidade pode ser None
#                (severidade nula = tratar por K2, nunca presumir gravidade)
#   consultas    consultas de credito nos ultimos 90 dias
#   vinculo      (empregador, tipo, inicio, salario, verificado) ou None;
#                lista de tuplas = varios vinculos em employment_records e em
#                employment_links (usa-se sempre o vigente de maior salario)
#   fluxo        (entrada, saida, volatilidade, dias_negativo, renda_recorrente)
#                ou lista de (data, entrada, saida, volatilidade, dias_negativo,
#                renda_recorrente) = varias analises em cash_flow_analyses
#                (usa-se sempre a de data mais recente, C7/C8/C9/K8)
#   banco        (meses de relacionamento, segmento, limite pre-aprovado);
#                limite pre-aprovado pode ser um valor ou uma lista de valores
#                (varias linhas em pre_approved_limits; usa-se a de maior valor
#                aprovado ativo, §7)
#   receita      (situacao do CPF, CND, PEP, sancoes)
#                situacao: regular | pending | suspended | canceled | deceased
#                CND:      regular | irregular | suspended
#   divergencias campos do cadastro que nao conferem com a Receita Federal
#   falta        campos deliberadamente ausentes; "cpf" = documento nulo em
#                personal_informations (aciona M3); "dti" = indice e parcelas
#                nulos (C2 incalculavel, dado ausente); "dti_informado" = so o
#                indice nulo, parcelas continuam informadas (C2 cai no
#                fallback parcelas/renda, nao e dado ausente)
#   triagem      "regular" | "irregular" - triagem de compliance irregular sem
#                PEP e sem lista de sancoes (segunda metade de K13)
#   produtos     tipos extras de produto contratado, alem de conta corrente e
#                emprestimo: "consignado" | "credit_card" | "vehicle_financing"
#   origem_receita  natureza da receita recorrente detectada no open finance:
#                "salario" | "inss" | "beneficio_social" | "pix"; None = salario
#   validacao_anterior  (situacao_antiga, data_antiga) = uma validacao de
#                documento anterior e divergente da atual em document_validations
#                (K10/C13 continuam usando a mais recente)
#
#   -- knobs que desacoplam as fontes umas das outras --
#   consentimento  granted | expired | revoked | None (sem registro de consentimento)
#   relacionamento True | False (nao-cliente) | "inativo" (conta encerrada)
#   score_interno  None = deriva do score do biro; int = diverge de proposito
#   interno        (proporcao em dia, pior atraso em dias) das parcelas internas;
#                  None = espelha `pontual`. Atraso > 30 dias vira `missed`,
#                  que e o que alimenta as regras de inadimplencia interna.
#   esocial        "auto" | "encerrado" | "ausente"
#   cadastral      "completo" | "sem_compliance" | "ausente"
#   biometria      biometria validada na Receita
#   pre_aprovado_valido  False = pre-aprovacao existe mas ja expirou
# ---------------------------------------------------------------------------
def cenario(faixa, score, renda, dti, utilizacao, pontual, nota, *, estimada=None,
            negativacoes=(), dividas=(), processos=(), fraude=None, consultas=1,
            vinculo=None, fluxo=None, banco=(24, "varejo", 0.0),
            receita=("regular", "regular", False, False), divergencias=(), mudou=False,
            falta=(), consentimento="granted", relacionamento=True, score_interno=None,
            interno=None, esocial="auto", cadastral="completo", biometria=True,
            pre_aprovado_valido=True, triagem="regular", produtos=(), origem_receita=None,
            validacao_anterior=None):
    return dict(faixa=faixa, score=score, renda=renda, estimada=estimada, dti=dti,
                utilizacao=utilizacao, pontual=pontual, negativacoes=list(negativacoes),
                dividas=list(dividas), processos=list(processos), fraude=fraude,
                consultas=consultas, vinculo=vinculo, fluxo=fluxo, banco=banco,
                receita=receita, divergencias=set(divergencias), mudou=mudou,
                falta=set(falta), consentimento=consentimento,
                relacionamento=relacionamento, score_interno=score_interno,
                interno=interno, esocial=esocial, cadastral=cadastral,
                biometria=biometria, pre_aprovado_valido=pre_aprovado_valido,
                triagem=triagem, produtos=tuple(produtos), origem_receita=origem_receita,
                validacao_anterior=validacao_anterior, nota=nota)


CLT = "CLT"
PJ = "PJ"
AUT = "autonomous"
EST = "estatutario"
APO = "aposentado"
PEN = "pensionista"

# produtos contratados alem da conta corrente e do emprestimo pessoal padrao;
# chave = valor aceito pelo knob `produtos`, valor = (nome, prefixo do contrato)
PRODUTOS_EXTRA = {
    "consignado": ("Emprestimo Consignado", "CONSIG"),
    "credit_card": ("Cartao de Credito", "CC-PROD"),
    "vehicle_financing": ("Financiamento de Veiculo", "FIN"),
}

CENARIOS = {
    # ================= bons (1-6): raros, alinhados nas quatro fontes ========
    1: cenario("bom", 842, 9800, 0.16, 0.10, 1.00,
               "assalariado de carreira, quatro fontes concordantes",
               vinculo=("Industria XYZ", CLT, "2014-03-01", 9800, True),
               fluxo=(9900, 6900, 0.04, 0, True), banco=(108, "alta renda", 40000),
               interno=(1.00, 0)),
    2: cenario("bom", 806, 6400, 0.19, 0.13, 0.98, "CLT estavel, relacionamento de sete anos",
               vinculo=("Comercio ABC", CLT, "2017-06-01", 6400, True),
               fluxo=(6500, 4800, 0.05, 0, True), banco=(84, "varejo", 22000),
               interno=(0.98, 0)),
    3: cenario("bom", 878, 15200, 0.12, 0.07, 1.00, "alta renda, sem passivo relevante",
               vinculo=("Consultoria Delta", PJ, "2013-01-01", 15200, True),
               fluxo=(15500, 9800, 0.06, 0, True), banco=(144, "alta renda", 90000),
               interno=(1.00, 0)),
    4: cenario("bom", 771, 5100, 0.24, 0.18, 0.97, "servidor publico, renda comprovada",
               vinculo=("Prefeitura Sintetica", EST, "2018-02-01", 5100, True),
               fluxo=(5150, 3900, 0.03, 0, True), banco=(72, "varejo", 14000),
               interno=(0.97, 0)),
    5: cenario("bom", 819, 7700, 0.17, 0.11, 1.00, "MEI consolidado, faturamento estavel",
               vinculo=("MEI Servicos", AUT, "2016-05-01", 7700, True),
               fluxo=(7800, 5600, 0.08, 0, True), banco=(96, "varejo", 30000),
               interno=(1.00, 0)),
    6: cenario("bom", 795, 5900, 0.21, 0.15, 0.98, "CLT com poupanca, sem restricoes",
               vinculo=("Tech Solutions Ltda", CLT, "2019-04-01", 5900, True),
               fluxo=(6000, 4400, 0.05, 0, True), banco=(66, "varejo", 18000),
               interno=(0.98, 0)),

    # ================= ruins (7-12): raros, eliminatoria evidente ============
    7: cenario("ruim", 512, 4200, 0.39, 0.48, 0.84, "alerta de fraude grave e ativo",
               fraude=("critical", "active"), divergencias=("nome",), consultas=9,
               vinculo=("Comercio ABC", CLT, "2024-01-01", 4200, False),
               fluxo=(4250, 4100, 0.29, 7, True), banco=(19, "varejo", 0.0),
               receita=("pending", "irregular", False, False), biometria=False,
               interno=(0.78, 40)),
    8: cenario("ruim", 241, 3100, 0.52, 0.77, 0.66, "score critico, abaixo do piso da politica",
               negativacoes=[(7400.0, "active", False)], dividas=[(9200.0, True)], consultas=8,
               vinculo=("Comercio ABC", CLT, "2024-06-01", 3100, False),
               fluxo=(3150, 3300, 0.33, 14, True), banco=(11, "varejo", 0.0),
               interno=(0.62, 55)),
    9: cenario("ruim", 388, 3600, 0.49, 0.71, 0.71, "tres negativacoes ativas",
               negativacoes=[(2900.0, "active", False), (4300.0, "active", False),
                             (1800.0, "active", False)], consultas=7,
               vinculo=("Comercio ABC", CLT, "2023-03-01", 3600, True),
               fluxo=(3650, 3550, 0.26, 9, True), banco=(29, "varejo", 0.0),
               interno=(0.68, 25)),
    10: cenario("ruim", 704, 5400, 0.27, 0.22, 0.96, "CPF cancelado por obito na Receita Federal",
                vinculo=("Industria XYZ", CLT, "2016-08-01", 5400, True),
                fluxo=(5450, 4200, 0.06, 0, True), banco=(81, "varejo", 12000),
                receita=("deceased", "regular", False, False), biometria=False,
                interno=(0.96, 0)),
    11: cenario("ruim", 736, 8300, 0.25, 0.19, 0.97, "homonimo confirmado em lista restritiva",
                vinculo=("Consultoria Delta", PJ, "2018-11-01", 8300, True),
                fluxo=(8400, 6100, 0.09, 0, True), banco=(74, "alta renda", 26000),
                receita=("regular", "regular", False, True), interno=(0.97, 0)),
    12: cenario("ruim", 419, 4700, 0.46, 0.64, 0.77,
                "falencia ativa somada a dividas em cobranca acima do teto",
                processos=["bankruptcy"], dividas=[(38000.0, True), (21000.0, True)],
                negativacoes=[(5100.0, "active", False)], consultas=6,
                vinculo=("Comercio ABC", CLT, "2022-09-01", 4700, False),
                fluxo=(4750, 4900, 0.31, 12, True), banco=(37, "varejo", 0.0),
                receita=("regular", "irregular", False, False), interno=(0.70, 48)),

    # ========== complexos (13-60): a maioria, sinais conflitantes ============
    # -- fontes que discordam entre si --
    13: cenario("complexo", 812, 7200, 0.44, 0.62, 0.9,
                "biro otimo e comportamento interno pessimo: fontes discordam",
                vinculo=("Industria XYZ", CLT, "2017-10-01", 7200, True),
                fluxo=(7300, 5500, 0.07, 0, True), banco=(63, "varejo", 15000),
                score_interno=372, interno=(0.58, 22)),
    14: cenario("complexo", 468, 4400, 0.41, 0.57, 0.81,
                "biro ruim e relacionamento interno impecavel de oito anos",
                negativacoes=[(3100.0, "paid", False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2016-02-01", 4400, True),
                fluxo=(4500, 3600, 0.11, 1, True), banco=(102, "varejo", 9000),
                score_interno=788, interno=(1.00, 0)),
    15: cenario("complexo", 612, 5600, 0.42, 0.58, 0.86,
                "vinculo declarado no biro sem contrapartida no eSocial",
                consultas=3, vinculo=("Comercio ABC", CLT, "2021-05-01", 5600, True),
                fluxo=(5650, 4600, 0.13, 2, True), banco=(47, "varejo", 7000),
                esocial="ausente", interno=(0.93, 0)),
    16: cenario("complexo", 604, 6100, 0.45, 0.61, 0.84,
                "eSocial mostra o vinculo encerrado; o biro ainda o da como vigente",
                consultas=2, vinculo=("Industria XYZ", CLT, "2019-07-01", 6100, True),
                fluxo=(6150, 4900, 0.12, 1, True), banco=(55, "varejo", 10000),
                esocial="encerrado", interno=(0.95, 0)),
    17: cenario("complexo", 656, 6800, 0.39, 0.53, 0.88,
                "nome diverge da Receita sem alerta de fraude correspondente no biro",
                divergencias=("nome",), consultas=2,
                vinculo=("Comercio ABC", CLT, "2018-03-01", 6800, True),
                fluxo=(6900, 5300, 0.08, 0, True), banco=(69, "varejo", 16000),
                biometria=False, interno=(0.96, 0)),
    18: cenario("complexo", 634, 4900, 0.36, 0.43, 0.91,
                "renda declarada muito acima do que o open finance observa",
                consultas=4, vinculo=("Consultoria Delta", AUT, "2020-09-01", 4900, False),
                fluxo=(2100, 1950, 0.35, 6, True), banco=(41, "varejo", 3000),
                interno=(0.91, 0)),
    19: cenario("complexo", 596, 3800, 0.43, 0.55, 0.87,
                "fluxo de caixa do open finance desmente a renda declarada",
                negativacoes=[(1900.0, "paid", False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2022-01-01", 3800, True),
                fluxo=(1700, 1820, 0.44, 13, True), banco=(32, "varejo", 2000),
                interno=(0.87, 0)),

    # -- inadimplencia interna, o sinal mais forte do registro interno --
    20: cenario("complexo", 741, 6300, 0.28, 0.26, 0.96,
                "parcela interna nao paga ha dois meses, biro ainda nao reflete",
                consultas=3, vinculo=("Industria XYZ", CLT, "2019-11-01", 6300, True),
                fluxo=(6350, 5000, 0.10, 1, True), banco=(58, "varejo", 11000),
                interno=(0.75, 62)),
    21: cenario("complexo", 767, 7100, 0.25, 0.21, 0.97,
                "inadimplencia interna acima de noventa dias com biro intacto",
                consultas=2, vinculo=("Tech Solutions Ltda", CLT, "2018-05-01", 7100, True),
                fluxo=(7200, 5600, 0.09, 0, True), banco=(76, "varejo", 19000),
                interno=(0.70, 128)),
    22: cenario("complexo", 588, 4600, 0.46, 0.64, 0.83,
                "atrasos internos curtos e recorrentes, nenhum caracterizado como perda",
                consultas=4, vinculo=("Comercio ABC", CLT, "2021-03-01", 4600, True),
                fluxo=(4650, 3900, 0.15, 3, True), banco=(44, "varejo", 5000),
                interno=(0.67, 18)),

    # -- situacao cadastral e compliance --
    23: cenario("complexo", 784, 7900, 0.23, 0.18, 0.98, "CPF pendente de regularizacao na Receita",
                consultas=2, vinculo=("Consultoria Delta", PJ, "2017-04-01", 7900, True),
                fluxo=(8000, 6000, 0.08, 0, True), banco=(79, "alta renda", 24000),
                receita=("pending", "regular", False, False), interno=(0.98, 0)),
    24: cenario("complexo", 828, 11200, 0.18, 0.12, 1.00,
                "pessoa exposta politicamente com perfil de credito excelente",
                vinculo=("Prefeitura Sintetica", EST, "2015-02-01", 11200, True),
                fluxo=(11400, 7900, 0.05, 0, True), banco=(118, "alta renda", 55000),
                receita=("regular", "regular", True, False), interno=(1.00, 0)),
    25: cenario("complexo", 692, 5300, 0.31, 0.34, 0.94, "CPF suspenso na Receita Federal",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-08-01", 5300, True),
                fluxo=(5350, 4300, 0.12, 2, True), banco=(50, "varejo", 6000),
                receita=("suspended", "regular", False, False), interno=(0.94, 0)),
    26: cenario("complexo", 617, 5800, 0.41, 0.56, 0.87, "certidao negativa de debitos negada",
                consultas=3, vinculo=("Comercio ABC", CLT, "2019-09-01", 5800, True),
                fluxo=(5900, 4700, 0.11, 1, True), banco=(53, "varejo", 9000),
                receita=("regular", "irregular", False, False), interno=(0.95, 0)),
    27: cenario("complexo", 634, 5500, 0.38, 0.49, 0.89,
                "debitos com exigibilidade suspensa: certidao positiva com efeito de negativa",
                consultas=3, vinculo=("Industria XYZ", CLT, "2020-02-01", 5500, True),
                fluxo=(5600, 4500, 0.13, 2, True), banco=(49, "varejo", 7500),
                receita=("regular", "suspended", False, False), interno=(0.94, 0)),
    28: cenario("complexo", 601, 6600, 0.43, 0.57, 0.86, "cliente nunca passou por validacao cadastral",
                consultas=2, vinculo=("Comercio ABC", CLT, "2021-01-01", 6600, True),
                fluxo=(6700, 5200, 0.10, 1, True), banco=(45, "varejo", 12000),
                cadastral="ausente", interno=(0.96, 0)),
    29: cenario("complexo", 623, 6000, 0.4, 0.54, 0.88, "sem triagem de compliance no cadastro",
                consultas=2, vinculo=("Comercio ABC", CLT, "2020-06-01", 6000, True),
                fluxo=(6100, 4800, 0.11, 1, True), banco=(51, "varejo", 10000),
                cadastral="sem_compliance", interno=(0.95, 0)),
    30: cenario("complexo", 521, 5200, 0.55, 0.81, 0.74,
                "data de nascimento diverge da Receita, biometria nao validada",
                divergencias=("nascimento",), consultas=3,
                vinculo=("Comercio ABC", CLT, "2021-07-01", 5200, True),
                fluxo=(5250, 4300, 0.14, 3, True), banco=(42, "varejo", 5500),
                biometria=False, interno=(0.93, 0)),
    31: cenario("complexo", 486, 4800, 0.58, 0.86, 0.69,
                "nome e nascimento divergem da Receita ao mesmo tempo",
                divergencias=("nome", "nascimento"), consultas=4,
                vinculo=("Comercio ABC", CLT, "2022-04-01", 4800, False),
                fluxo=(4850, 4100, 0.17, 4, True), banco=(35, "varejo", 3500),
                biometria=False, interno=(0.92, 0)),

    # -- ausencia de fontes inteiras --
    32: cenario("complexo", 596, 5700, 0.44, 0.59, 0.85, "consentimento de open finance revogado",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-10-01", 5700, True),
                fluxo=(5750, 4600, 0.12, 2, True), banco=(48, "varejo", 8000),
                consentimento="revoked", interno=(0.95, 0)),
    33: cenario("complexo", 571, 4700, 0.47, 0.63, 0.82, "consentimento de open finance expirado",
                consultas=4, vinculo=("Comercio ABC", CLT, "2021-11-01", 4700, True),
                fluxo=(4750, 3950, 0.16, 3, True), banco=(39, "varejo", 4500),
                consentimento="expired", interno=(0.92, 0)),
    34: cenario("complexo", 648, 6900, 0.37, 0.47, 0.9, "nao-cliente da instituicao, sem historico interno",
                consultas=2, vinculo=("Tech Solutions Ltda", CLT, "2018-08-01", 6900, True),
                fluxo=(7000, 5400, 0.09, 0, True), banco=(0, "varejo", 0.0),
                relacionamento=False, consentimento=None, interno=(0.0, 0)),
    35: cenario("complexo", 538, 5000, 0.52, 0.76, 0.77, "conta encerrada, relacionamento inativo",
                consultas=3, vinculo=("Comercio ABC", CLT, "2019-12-01", 5000, True),
                fluxo=(5100, 4200, 0.13, 2, True), banco=(56, "varejo", 0.0),
                relacionamento="inativo", interno=(0.93, 0)),
    36: cenario("complexo", 607, 6200, 0.42, 0.55, 0.87, "score interno nunca calculado",
                consultas=2, vinculo=("Industria XYZ", CLT, "2020-03-01", 6200, True),
                fluxo=(6300, 5000, 0.10, 1, True), banco=(52, "varejo", 9500),
                score_interno=0, interno=(0.96, 0)),
    37: cenario("complexo", 619, 5900, 0.4, 0.52, 0.88, "cliente novo: sem parcelas internas no periodo",
                consultas=3, vinculo=("Comercio ABC", CLT, "2021-04-01", 5900, True),
                fluxo=(6000, 4800, 0.11, 1, True), banco=(4, "varejo", 0.0),
                interno=("vazio", 0)),
    38: cenario("complexo", 592, 4500, 0.43, 0.58, 0.86, "sem historico de pagamento no biro",
                consultas=4, vinculo=("Comercio ABC", CLT, "2022-06-01", 4500, True),
                fluxo=(4550, 3800, 0.18, 4, True), banco=(27, "varejo", 3000),
                falta=("historico_pagamento",), interno=(0.90, 0)),
    39: cenario("complexo", 621, None, 0.44, 0.56, 0.87,
                "autonomo sem renda declarada; biro estima e o fluxo confirma",
                estimada=4100, consultas=4,
                vinculo=("Consultoria Delta", AUT, "2020-05-01", None, False),
                fluxo=(4200, 3500, 0.24, 4, True), banco=(36, "varejo", 3500),
                falta=("renda_declarada",), interno=(0.89, 0)),
    40: cenario("complexo", 588, None, 0.44, 0.59, 0.85,
                "sem renda declarada, sem estimativa e sem vinculo vigente",
                consultas=5, vinculo=None,
                fluxo=(2300, 2250, 0.37, 10, False), banco=(24, "varejo", 0.0),
                falta=("renda_declarada", "renda_estimada", "vinculo"),
                esocial="ausente", interno=(0.85, 0)),
    41: cenario("complexo", None, 5600, 0.31, 0.33, 0.94, "biro sem score calculado",
                consultas=3, vinculo=("Comercio ABC", CLT, "2021-09-01", 5600, True),
                fluxo=(5700, 4600, 0.12, 2, True), banco=(46, "varejo", 7000),
                falta=("score",), interno=(0.94, 0)),
    42: cenario("complexo", 604, 4900, 0.34, 0.4, 0.86,
                "sem comprometimento informado e sem parcelas mensais para calcular",
                consultas=4, vinculo=("Comercio ABC", CLT, "2021-12-01", 4900, True),
                fluxo=(4950, 4100, 0.15, 3, True), banco=(38, "varejo", 4000),
                falta=("dti", "utilizacao"), interno=(0.91, 0)),

    # -- valores em cima das fronteiras de faixa --
    43: cenario("complexo", 700, 6000, 0.30, 0.30, 0.95, "tres criterios exatamente no limite superior da faixa",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-01-01", 6000, True),
                fluxo=(6600, 5400, 0.12, 0, True), banco=(60, "varejo", 8000),
                score_interno=800, interno=(0.95, 0)),
    44: cenario("complexo", 699, 6000, 0.31, 0.31, 0.949, "os mesmos criterios um passo abaixo do limite",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-01-01", 6000, True),
                fluxo=(6550, 5400, 0.12, 1, True), banco=(59, "varejo", 8000),
                score_interno=799, interno=(0.949, 0)),
    45: cenario("complexo", 800, 7000, 0.40, 0.899, 0.85, "limites superiores das faixas intermediarias",
                consultas=4, vinculo=("Industria XYZ", CLT, "2019-01-01", 7000, True),
                fluxo=(7700, 6300, 0.14, 5, True), banco=(24, "varejo", 10000),
                interno=(0.85, 0)),
    46: cenario("complexo", 799, 7000, 0.41, 0.90, 0.849, "um passo abaixo em todas as faixas anteriores",
                consultas=4, vinculo=("Industria XYZ", CLT, "2019-01-01", 7000, True),
                fluxo=(7650, 6300, 0.14, 6, True), banco=(23, "varejo", 10000),
                interno=(0.849, 0)),
    47: cenario("complexo", 601, 4300, 0.50, 0.60, 0.70, "todos os criterios no piso da faixa que ainda pontua",
                negativacoes=[(980.0, "active", False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2022-02-01", 4300, False),
                fluxo=(4300, 4300, 0.22, 15, True), banco=(6, "varejo", 1500),
                interno=(0.70, 0)),
    48: cenario("complexo", 599, 4300, 0.51, 0.61, 0.699, "todos os criterios um passo abaixo do piso",
                negativacoes=[(1020.0, "active", False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2022-02-01", 4300, False),
                fluxo=(4250, 4300, 0.22, 16, True), banco=(5, "varejo", 1200),
                interno=(0.699, 0)),

    # -- passivo: negativacoes, dividas e processos --
    49: cenario("complexo", 592, 4400, 0.44, 0.59, 0.85, "negativacao ativa de valor baixo",
                negativacoes=[(870.0, "active", False)], consultas=4,
                vinculo=("Comercio ABC", CLT, "2021-10-01", 4400, True),
                fluxo=(4450, 3800, 0.17, 4, True), banco=(40, "varejo", 3000),
                interno=(0.90, 0)),
    50: cenario("complexo", 571, 4100, 0.47, 0.63, 0.83, "negativacao ativa de valor relevante",
                negativacoes=[(6200.0, "active", False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2022-05-01", 4100, True),
                fluxo=(4150, 3700, 0.20, 6, True), banco=(31, "varejo", 2000),
                interno=(0.88, 0)),
    51: cenario("complexo", 524, 3700, 0.54, 0.77, 0.75, "duas negativacoes ativas somando pouco",
                negativacoes=[(1400.0, "active", False), (2100.0, "active", False)], consultas=6,
                vinculo=("Comercio ABC", CLT, "2023-02-01", 3700, True),
                fluxo=(3750, 3500, 0.23, 8, True), banco=(25, "varejo", 1500),
                interno=(0.84, 0)),
    52: cenario("complexo", 662, 5100, 0.35, 0.39, 0.92, "negativacao ativa e contestada pelo cliente",
                negativacoes=[(4600.0, "active", True)], consultas=4,
                vinculo=("Industria XYZ", CLT, "2021-06-01", 5100, True),
                fluxo=(5200, 4300, 0.14, 2, True), banco=(43, "varejo", 5000),
                interno=(0.92, 0)),
    53: cenario("complexo", 682, 8600, 0.48, 0.72, 0.9, "dividas em cobranca logo abaixo do teto da politica",
                dividas=[(41000.0, True), (38000.0, True)], consultas=6,
                vinculo=("Tech Solutions Ltda", CLT, "2017-08-01", 8600, True),
                fluxo=(8700, 8300, 0.19, 3, True), banco=(83, "alta renda", 21000),
                interno=(0.95, 0)),
    54: cenario("complexo", 788, 8600, 0.46, 0.66, 0.94, "as mesmas dividas, agora acima do teto",
                dividas=[(52000.0, True), (44000.0, True)], consultas=7,
                vinculo=("Tech Solutions Ltda", CLT, "2017-08-01", 8600, True),
                fluxo=(8700, 8500, 0.21, 5, True), banco=(83, "alta renda", 21000),
                interno=(0.94, 0)),
    55: cenario("complexo", 588, 5400, 0.45, 0.61, 0.85, "acao civel em andamento, sem penhora",
                processos=["civil"], consultas=3,
                vinculo=("Industria XYZ", CLT, "2018-12-01", 5400, True),
                fluxo=(5500, 4500, 0.12, 1, True), banco=(64, "varejo", 8500),
                interno=(0.93, 0)),
    56: cenario("complexo", 706, 5800, 0.30, 0.29, 0.95, "alerta de fraude de baixa severidade em apuracao",
                fraude=("low", "investigating"), consultas=4,
                vinculo=("Comercio ABC", CLT, "2020-04-01", 5800, True),
                fluxo=(5900, 4700, 0.11, 1, True), banco=(54, "varejo", 9000),
                biometria=False, interno=(0.95, 0)),

    # -- limite de credito: teto de fluxo e piso de pre-aprovacao --
    57: cenario("complexo", 664, 6500, 0.36, 0.44, 0.91,
                "fluxo de caixa magro limita o credito apesar do bom perfil",
                consultas=2, vinculo=("Comercio ABC", CLT, "2019-05-01", 6500, True),
                fluxo=(6560, 6320, 0.10, 0, True), banco=(68, "varejo", 4000),
                interno=(0.97, 0)),
    58: cenario("complexo", 691, 7300, 0.33, 0.38, 0.93,
                "pre-aprovacao vigente acima do limite calculado pela politica",
                consultas=2, vinculo=("Industria XYZ", CLT, "2018-07-01", 7300, True),
                fluxo=(7500, 5100, 0.08, 0, True), banco=(87, "alta renda", 42000),
                interno=(0.98, 0)),
    59: cenario("complexo", 686, 6400, 0.34, 0.4, 0.92, "pre-aprovacao existe, mas ja expirou",
                consultas=3, vinculo=("Comercio ABC", CLT, "2019-10-01", 6400, True),
                fluxo=(6500, 5100, 0.10, 1, True), banco=(61, "varejo", 35000),
                pre_aprovado_valido=False, interno=(0.96, 0)),
    60: cenario("complexo", 498, 3300, 0.56, 0.79, 0.71,
                "muitos sinais fracos somados, nenhum eliminatorio isolado",
                negativacoes=[(1600.0, "active", False)], dividas=[(5200.0, False)],
                consultas=7, vinculo=("Comercio ABC", CLT, "2023-05-01", 3300, False),
                fluxo=(3350, 3400, 0.28, 11, True), banco=(15, "varejo", 0.0),
                receita=("regular", "irregular", False, False), biometria=False,
                divergencias=("nascimento",), esocial="encerrado", interno=(0.72, 20)),

    # ================= bons (61-62): raros, alinhados nas quatro fontes ======
    61: cenario("bom", 831, 4200, 0.20, 0.14, 0.98,
                "aposentado do INSS com consignado em dia, margem preservada",
                vinculo=("INSS", APO, "2011-04-01", 4200, True), produtos=("consignado",),
                origem_receita="inss", fluxo=(4300, 3200, 0.05, 0, True),
                banco=(120, "varejo", 8000), interno=(0.98, 0)),
    62: cenario("bom", 810, None, 0.20, 0.12, 0.97,
                "servidor com dois vinculos vigentes; vale o de maior salario",
                vinculo=[("Prefeitura Sintetica", EST, "2015-01-01", 9500, True),
                         ("Consultoria Secundaria", PJ, "2022-01-01", 2200, True)],
                fluxo=(9600, 6900, 0.05, 0, True), banco=(84, "alta renda", 12000),
                interno=(0.97, 0)),

    # ================= ruins (63-64): raros, eliminatoria evidente ===========
    63: cenario("ruim", 712, 5600, 0.26, 0.20, 0.96, "CPF cancelado na Receita Federal",
                vinculo=("Industria XYZ", CLT, "2016-05-01", 5600, True),
                fluxo=(5650, 4300, 0.06, 0, True), banco=(78, "varejo", 13000),
                receita=("canceled", "regular", False, False), biometria=False,
                interno=(0.96, 0)),
    64: cenario("ruim", 450, 3000, 0.45, 0.55, 0.75,
                "duas negativacoes ativas somando mais de cinco vezes a renda",
                negativacoes=[(8000.0, "active", False), (8200.0, "active", False)],
                consultas=6, vinculo=("Comercio ABC", CLT, "2021-01-01", 3000, True),
                fluxo=(3050, 3100, 0.30, 10, True), banco=(30, "varejo", 0.0),
                interno=(0.85, 0)),

    # ================= complexos (65-80): 16 cenarios, tres blocos ===========
    # -- regras da politica sem cobertura ate aqui --
    65: cenario("complexo", 650, 5000, 0.35, 0.40, 0.90, "CPF ausente no cadastro, aciona M3",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-01-01", 5000, True),
                fluxo=(5050, 4200, 0.15, 2, True), banco=(40, "varejo", 6000),
                falta=("cpf",), interno=(0.90, 0)),
    66: cenario("complexo", 620, 5200, 0.38, 0.42, 0.88,
                "alerta de fraude ativo com severidade nula, nunca presumir gravidade",
                fraude=(None, "active"), consultas=3,
                vinculo=("Comercio ABC", CLT, "2020-05-01", 5200, True),
                fluxo=(5250, 4300, 0.13, 2, True), banco=(45, "varejo", 7000),
                interno=(0.90, 0)),
    67: cenario("complexo", 610, 5000, 0.35, 0.40, 0.85,
                "negativacao e divida com situacao nula, tratamento conservador sem knock-out",
                negativacoes=[(500.0, None, False)], dividas=[(2000.0, None)], consultas=3,
                vinculo=("Comercio ABC", CLT, "2019-08-01", 5000, True),
                fluxo=(5050, 4200, 0.14, 2, True), banco=(35, "varejo", 5000),
                interno=(0.85, 0)),
    68: cenario("complexo", 680, 4200, 0.30, 0.35, 0.80,
                "pensao por morte movimentada apos obito do titular: severidade decide por K10",
                vinculo=("INSS - Pensao por Morte", PEN, "2010-03-01", 4200, True),
                fluxo=(4250, 3600, 0.10, 2, True), banco=(60, "varejo", 5000),
                receita=("deceased", "regular", False, False), biometria=False,
                interno=(0.80, 45)),
    69: cenario("complexo", 600, 4800, 0.35, 0.38, 0.85,
                "fraude em apuracao e CPF pendente simultaneos, vence a de menor numero",
                fraude=("medium", "investigating"), consultas=3,
                vinculo=("Comercio ABC", CLT, "2020-01-01", 4800, True),
                fluxo=(4850, 4000, 0.14, 2, True), banco=(35, "varejo", 5000),
                receita=("pending", "regular", False, False), interno=(0.85, 0)),
    70: cenario("complexo", 630, 5100, 0.33, 0.36, 0.89,
                "triagem de compliance irregular sem PEP e sem lista restritiva",
                consultas=2, vinculo=("Comercio ABC", CLT, "2019-10-01", 5100, True),
                fluxo=(5150, 4200, 0.12, 1, True), banco=(50, "varejo", 8000),
                triagem="irregular", interno=(0.90, 0)),

    # -- regras de desempate entre varios registros --
    71: cenario("complexo", [(820, "2025-01-10"), (560, "2026-08-01")], 4600, 0.38, 0.45, 0.83,
                "dois scores de biro com datas diferentes; vale o mais recente",
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-02-01", 4600, True),
                fluxo=(4650, 3900, 0.15, 3, True), banco=(30, "varejo", 4000),
                interno=(0.85, 0)),
    72: cenario("complexo", 680, 5300, 0.32, 0.34, 0.90,
                "documento antigo cancelado, validacao mais recente regular nao reprova",
                validacao_anterior=("canceled", "2025-03-15"), consultas=2,
                vinculo=("Comercio ABC", CLT, "2019-03-01", 5300, True),
                fluxo=(5350, 4400, 0.11, 1, True), banco=(55, "varejo", 9000),
                interno=(0.92, 0)),
    73: cenario("complexo", 590, 5000, 0.35, 0.40, 0.87,
                "duas analises de fluxo de caixa em datas distintas; vale a mais recente",
                fluxo=[("2025-10-01", 6000, 4000, 0.08, 0, True),
                       ("2026-08-25", 5000, 4700, 0.30, 12, True)],
                consultas=3, vinculo=("Comercio ABC", CLT, "2020-04-01", 5000, True),
                banco=(35, "varejo", 5000), interno=(0.88, 0)),
    74: cenario("complexo", 760, 6000, 0.25, 0.20, 0.95,
                "tres pre-aprovacoes ativas do mesmo tipo; vale a de maior valor como piso",
                consultas=2, vinculo=("Comercio ABC", CLT, "2018-01-01", 6000, True),
                fluxo=(6100, 4800, 0.08, 0, True), banco=(70, "varejo", [5000, 12000, 35000]),
                interno=(0.95, 0)),

    # -- realismo brasileiro ausente --
    75: cenario("complexo", 560, None, 0.30, 0.50, 0.80,
                "beneficiario de beneficio social como unica renda, cadastro pobre",
                estimada=1800, vinculo=None, origem_receita="beneficio_social",
                fluxo=(1850, 1600, 0.10, 1, True), banco=(10, "varejo", 0.0),
                falta=("renda_declarada", "email", "telefone"), esocial="ausente",
                consultas=2, interno=(0.80, 0)),
    76: cenario("complexo", 560, None, 0.28, 0.55, 0.72,
                "trabalhador informal vivendo de Pix; C2 cai no fallback parcelas / renda",
                estimada=3200, vinculo=None, origem_receita="pix", esocial="ausente",
                fluxo=(4200, 3100, 0.42, 5, True), banco=(8, "varejo", 0.0),
                falta=("renda_declarada", "dti_informado"), consultas=4, interno=(0.75, 0)),
    77: cenario("complexo", 640, 2200, 0.30, 0.35, 0.90,
                "fluxo com padrao de conta intermediaria: entradas e saidas altas e "
                "identicas, sem alerta ainda no biro",
                consultas=9, vinculo=("Comercio ABC", CLT, "2021-01-01", 2200, True),
                fluxo=(18000, 17800, 0.55, 18, True), banco=(18, "varejo", 0.0),
                interno=(0.90, 0)),
    78: cenario("complexo", 660, 4500, 0.55, 0.40, 0.92,
                "superendividamento regularizado: negativacoes quitadas e dividas "
                "renegociadas fora de cobranca (Lei 14.181/2021)",
                negativacoes=[(3200.0, "paid", False), (1800.0, "paid", False)],
                dividas=[(15000.0, False)], produtos=("consignado", "credit_card"),
                consultas=3, vinculo=("Comercio ABC", CLT, "2018-01-01", 4500, True),
                fluxo=(4600, 4300, 0.15, 3, True), banco=(50, "varejo", 3000),
                interno=(0.97, 0)),
    79: cenario("complexo", 610, 5200, 0.36, 0.42, 0.86,
                "mudanca de endereco nao propagada ao biro, financiamento de veiculo "
                "e acao civel em curso",
                mudou=True, produtos=("vehicle_financing",), processos=["civil"],
                consultas=3, vinculo=("Comercio ABC", CLT, "2019-06-01", 5200, True),
                fluxo=(5250, 4300, 0.14, 2, True), banco=(45, "varejo", 6000),
                interno=(0.88, 0)),
    80: cenario("complexo", 690, 5600, 0.33, 0.37, 0.91,
                "divergencia de nome por casamento, biometria validada e sem alerta de fraude",
                divergencias=("nome",), biometria=True, consultas=2,
                vinculo=("Comercio ABC", CLT, "2019-02-01", 5600, True),
                fluxo=(5650, 4600, 0.12, 1, True), banco=(48, "varejo", 7000),
                interno=(0.92, 0)),
}

NOMES = ["Amanda Ferreira Lima", "Bruno Carvalho Nunes", "Carla Menezes Rocha", "Diego Tavares Pinto",
         "Elaine Moraes Cardoso", "Fabio Antunes Vieira", "Giovana Peixoto Braga", "Helio Marques Fontes",
         "Isadora Campos Teles", "Joao Batista Siqueira", "Karina Duarte Prado", "Leandro Bastos Farias",
         "Mariana Freitas Lopes", "Nelson Aguiar Bittencourt", "Olivia Rezende Sampaio", "Paulo Cesar Andrade",
         "Queila Monteiro Serra", "Rafael Guedes Pontes", "Simone Vasques Coelho", "Tiago Meireles Bandeira",
         "Ursula Barreto Galvao", "Vinicius Aragao Quintela", "Wanda Cordeiro Lisboa", "Xavier Toledo Amancio",
         "Yara Bonfim Salgado", "Zeno Villela Krause", "Alice Padilha Moura", "Benedito Falcao Ximenes",
         "Clarice Uchoa Bezerra", "Danilo Espindola Rangel", "Eduarda Nogueira Vilela", "Fernando Quirino Sa",
         "Gabriela Assuncao Prates", "Heitor Zanetti Macedo", "Iara Valadares Pimenta", "Juliano Correia Estrela",
         "Katia Lemos Wanderley", "Lucas Sarmento Feitosa", "Marcela Bulhoes Tenorio", "Nicolas Peçanha Adorno",
         "Otavia Simoes Caldeira", "Pedro Henrique Grangeiro", "Roberta Xavier Munhoz", "Samuel Trindade Beltrao",
         "Tatiana Oliveira Rebouças", "Ulisses Prado Camargo", "Viviane Nobrega Cotrim", "Wesley Amorim Godoi",
         "Ximena Ferraz Delgado", "Yuri Bacelar Pontual", "Adriana Lustosa Verissimo", "Bernardo Tinoco Alvim",
         "Camila Reis Fontoura", "Douglas Vergara Nepomuceno", "Erika Bandeira Quaresma", "Flavio Mesquita Salles",
         "Gisele Antunes Barroca", "Humberto Leao Passarinho", "Ingrid Colaco Vasconcelos", "Jorge Belmonte Tavora",
         "Kleber Monteiro Vasconcellos", "Luana Ribeiro Castanheira", "Mauricio Salgueiro Andrade",
         "Nadia Figueiredo Bezerra", "Otavio Cavalcanti Serpa", "Priscila Bandeira Moutinho",
         "Rogerio Almeida Quaresma", "Sabrina Cortez Malheiros", "Thiago Loureiro Sandoval",
         "Valeria Pimentel Guimaraes", "Wagner Sepulveda Teixeira", "Yasmin Carrilho Fagundes",
         "Alexandre Bittar Nascimento", "Bianca Sequeira Marinho", "Cassio Werneck Azambuja",
         "Debora Chaves Montenegro", "Emerson Dalmaso Ribas", "Flavia Junqueira Ottoni",
         "Gustavo Meneghel Serrano", "Helena Bicalho Frota"]
MAES = ["Marta Ferreira", "Sonia Carvalho", "Regina Menezes", "Laura Tavares", "Cristina Moraes",
        "Vera Antunes", "Beatriz Peixoto", "Alice Marques", "Neusa Campos", "Rita Batista",
        "Eliane Duarte", "Silvia Bastos", "Angela Freitas", "Marlene Aguiar", "Cecilia Rezende",
        "Dulce Andrade", "Ivone Monteiro", "Nadia Guedes", "Tereza Vasques", "Aparecida Meireles",
        "Rosane Barreto", "Ilda Aragao", "Selma Cordeiro", "Miriam Toledo", "Julia Bonfim",
        "Norma Villela", "Solange Padilha", "Adelia Falcao", "Zuleica Uchoa", "Berenice Espindola",
        "Lourdes Nogueira", "Iracema Quirino", "Diva Assuncao", "Elza Zanetti", "Custodia Valadares",
        "Wilma Correia", "Genoveva Lemos", "Antonieta Sarmento", "Perpetua Bulhoes", "Otilia Peçanha",
        "Jandira Simoes", "Efigenia Grangeiro", "Consuelo Xavier", "Amelia Trindade", "Josefa Oliveira",
        "Leonor Prado", "Guiomar Nobrega", "Benedita Amorim", "Raquel Ferraz", "Hilda Bacelar",
        "Delfina Lustosa", "Ondina Tinoco", "Clotilde Reis", "Sebastiana Vergara", "Anesia Bandeira",
        "Zilda Mesquita", "Lindaura Antunes", "Etelvina Leao", "Nilza Colaco", "Cordelia Belmonte",
        "Terezinha Bicalho", "Marisa Cavalcanti", "Odete Salgueiro", "Firmina Figueiredo",
        "Clarinda Serpa", "Joventina Valente", "Almerinda Quaresma", "Petronilha Malheiros",
        "Aurea Loureiro", "Hortencia Pimentel", "Isaura Sepulveda", "Filomena Carrilho",
        "Argentina Bittar", "Romilda Sequeira", "Eunice Werneck", "Adalgisa Chaves",
        "Creuza Dalmaso", "Belmira Junqueira", "Wanderlea Meneghel", "Corina Bicalho"]

# --------------------------------------------------------------------------- helpers
def cpf(seed):
    """CPF sintetico com digitos verificadores validos, deterministico por seed."""
    digest = hashlib.sha256(f"tcc-cenario-{seed}".encode()).hexdigest()
    base = [int(d) for d in digest if d.isdigit()][:9]
    for _ in range(2):
        peso = len(base) + 1
        total = sum(d * (peso - i) for i, d in enumerate(base))
        resto = total % 11
        base.append(0 if resto < 2 else 11 - resto)
    return "".join(str(d) for d in base)


def dia(date, plus=0):
    if date is None:
        return None
    if isinstance(date, str):
        date = datetime.strptime(date, "%Y-%m-%d")
    return (date + timedelta(days=plus)).strftime("%Y-%m-%dT00:00:00Z")


def base(row_id):
    return {"id": row_id, "created_at": TS, "updated_at": TS, "deleted_at": None}


# -- normalizacao dos knobs que passaram a aceitar mais de um registro -------
# Cada knob continua aceitando o formato escalar de sempre (um valor, uma
# tupla, um par (valor, data)); as funcoes abaixo o tornam lista de um so
# elemento, de modo que o resto do gerador so precise iterar. Isso preserva
# byte a byte o comportamento anterior quando o cenario nao usa lista.
def _lista(valor):
    """Normaliza um escalar/None/0/lista em lista. Falsy vira lista vazia."""
    if not valor:
        return []
    return valor if isinstance(valor, list) else [valor]


def _score_lista(score):
    """Normaliza o knob score: int vira uma entrada unica na data padrao."""
    if score is None:
        return []
    return score if isinstance(score, list) else [(score, "2026-08-15")]


def _score_atual(score):
    """Score de data mais recente (§5 C1); None se nao houver nenhum."""
    lista = _score_lista(score)
    return max(lista, key=lambda sd: sd[1])[0] if lista else None


def _vinc_lista(vinculo):
    """Normaliza o knob vinculo: uma tupla vira lista de um elemento."""
    if not vinculo:
        return []
    return vinculo if isinstance(vinculo, list) else [vinculo]


def _vinc_principal(vinculo):
    """Vinculo vigente de maior salario (renda considerada, C6, C15)."""
    lista = _vinc_lista(vinculo)
    return max(lista, key=lambda v: v[3]) if lista else None


def _fluxo_lista(fluxo):
    """Normaliza o knob fluxo: tupla unica vira lista de uma entrada datada."""
    if isinstance(fluxo, list):
        return fluxo
    entrada, saida, volatilidade, dias_negativo, recorrente = fluxo
    return [("2026-08-20", entrada, saida, volatilidade, dias_negativo, recorrente)]


def _fluxo_atual(fluxo):
    """Analise de fluxo de caixa de data mais recente (C7, C8, C9, K8)."""
    _, entrada, saida, volatilidade, dias_negativo, recorrente = max(
        _fluxo_lista(fluxo), key=lambda f: f[0])
    return entrada, saida, volatilidade, dias_negativo, recorrente


def renda_considerada(c):
    """Precedencia da politica: declarada, estimada, salario de vinculo vigente."""
    if c["renda"]:
        return c["renda"]
    if c["estimada"]:
        return c["estimada"]
    principal = _vinc_principal(c["vinculo"])
    if principal and principal[3]:
        return principal[3]
    return None


def fixtures(service, name):
    return ROOT / service / "cmd" / "fixtures" / "fixtures" / f"{name}.json"


def carrega(service, name):
    path = fixtures(service, name)
    return json.loads(path.read_text()) if path.exists() else []


def proximo_id(rows):
    return max((r.get("id", 0) for r in rows), default=0) + 1


def grava(service, name, preliminares, novos):
    """Mantem os registros do conjunto preliminar e substitui os dos cenarios."""
    fixtures(service, name).write_text(
        json.dumps(preliminares + novos, indent=2, ensure_ascii=False) + "\n")
    return len(novos)


def separa(service, name, campo="person_id"):
    """Divide as fixtures entre o conjunto preliminar (<= 10) e os cenarios."""
    rows = carrega(service, name)
    preliminares = [r for r in rows if r.get(campo, 0) < PRIMEIRO_ID]
    return preliminares, proximo_id(preliminares)


# --------------------------------------------------------------------------- cadastro
def endereco(address_id, person_id, offset, anterior=False):
    cidade, uf = CIDADES[offset % len(CIDADES)]
    entrada = datetime(2019 + offset % 6, 1 + offset % 12, 1 + offset % 27)
    return {
        **base(address_id),
        "zip_code": f"{10000 + person_id * 37:05d}-{person_id * 3 % 1000:03d}",
        "state": uf,
        "city": cidade,
        "neighborhood": ["Centro", "Jardins", "Boa Vista", "Aldeota"][offset % 4],
        "street": RUAS[(offset + (3 if anterior else 0)) % len(RUAS)],
        "number": str(100 + person_id * 13 + (7 if anterior else 0)),
        "complement": None if offset % 3 else f"Apto {10 + offset}",
        "reference_point": None,
        "address_type": "residential",
        "latitude": None,
        "longitude": None,
        "validated_by_post": offset % 4 != 0,
        "risk_score": 10 + (person_id * 7) % 80,
        "is_current": not anterior,
        "is_correspondence": not anterior,
        "moved_in_date": dia(entrada.replace(year=entrada.year - 3) if anterior else entrada),
        "moved_out_date": dia(entrada) if anterior else None,
        "verification_status": ["verified", "unverified", "disputed"][offset % 3],
    }


def cadastro():
    """Cadastro canonico (visao do banco); generate_cadastro.py deriva as demais."""
    pi_prev, pi_id = separa("internal-registry", "personal_informations", "id")
    end_prev, end_id = separa("internal-registry", "addresses", "id")
    vin_prev, vin_id = separa("internal-registry", "person_addresses", "personal_information_id")

    pessoas, enderecos, vinculos = [], [], []
    for offset, (person_id, c) in enumerate(sorted(CENARIOS.items())):
        nome = NOMES[offset]
        cidade, uf = CIDADES[offset % len(CIDADES)]
        # dia sempre <= 12 e diferente do mes: e a condicao para que
        # generate_cadastro.py consiga plantar a divergencia de data de
        # nascimento (que troca dia por mes) em qualquer pessoa
        mes = 1 + (offset * 5) % 12
        dia_nasc = 1 + (offset * 7) % 12
        if dia_nasc == mes:
            dia_nasc = 1 + dia_nasc % 12
        nascimento = datetime(1996 - offset % 30, mes, dia_nasc)

        pessoas.append({
            **base(pi_id + offset),
            "full_name": nome,
            "mother_name": MAES[offset],
            "birth_date": dia(nascimento),
            "gender": "female" if offset % 2 else "male",
            "nationality": "Brazilian",
            "marital_status": ["single", "married", "divorced"][offset % 3],
            "document": None if "cpf" in c["falta"] else cpf(person_id),
            "rg": f"{20 + offset}.{100 + offset}.{200 + offset}-{offset % 10}",
            "rg_issuer": f"SSP-{uf}",
            "rg_issue_date": dia(nascimento.replace(year=nascimento.year + 18)),
            "voter_id": f"{700000000000 + person_id * 137}",
            "work_card": f"CTPS-{person_id:05d}",
            "primary_phone": f"{11 + offset % 80}9{80000000 + person_id * 7919:08d}"[:11],
            "secondary_phone": None,
            "email": None if "email" in c["falta"] else f"{nome.split()[0].lower()}.{person_id}@example.com",
            "alternative_email": None,
            "profile_photo_id": None,
            "document_validated": c["receita"][0] == "regular",
            "email_verified": "email" not in c["falta"] and offset % 3 != 0,
            "phone_verified": "telefone" not in c["falta"] and offset % 4 != 0,
            "biometric_validated": c["biometria"],
            "receita_federal_status": c["receita"][0],
        })

        if c["mudou"]:
            # endereco anterior: o birô nao foi informado da mudanca
            enderecos.append(endereco(end_id + len(enderecos), person_id, offset, anterior=True))
            vinculos.append({
                **base(vin_id + len(vinculos)),
                "personal_information_id": pi_id + offset,
                "address_id": end_id + len(enderecos) - 1,
            })

        enderecos.append(endereco(end_id + len(enderecos), person_id, offset))

        vinculos.append({
            **base(vin_id + len(vinculos)),
            "personal_information_id": pi_id + offset,
            "address_id": end_id + len(enderecos) - 1,
        })

    grava("internal-registry", "personal_informations", pi_prev, pessoas)
    grava("internal-registry", "addresses", end_prev, enderecos)
    grava("internal-registry", "person_addresses", vin_prev, vinculos)
    print(f"cadastro canonico: {len(pessoas)} pessoas, {len(enderecos)} enderecos")
    return {person_id: pi_id + offset for offset, person_id in enumerate(sorted(CENARIOS))}


def pessoas_por_fonte(pi_por_pessoa, chaves):
    """persons.json de cada fonte, com o FK de dominio que ela usa."""
    for service, campo, valor in chaves:
        prev, next_id = separa(service, "persons", "id")
        novos = []
        for offset, person_id in enumerate(sorted(CENARIOS)):
            row = {
                **base(next_id + offset),
                "personal_information_id": pi_por_pessoa[person_id],
                "last_verified_at": dia("2026-08-20"),
            }
            if campo:
                row[campo] = valor(person_id, offset)
            novos.append(row)
        grava(service, "persons", prev, novos)


# --------------------------------------------------------------------------- birô
def biro():
    tabelas = {t: separa("bureau", t) for t in
               ["credit_scores", "financial_profiles", "employment_records", "credit_accounts",
                "payment_histories", "debts", "negative_records", "legal_records",
                "fraud_alerts", "credit_inquiries", "risk_assessments"]}
    novos = {t: [] for t in tabelas}

    def add(tabela, row):
        prev, next_id = tabelas[tabela]
        novos[tabela].append({**base(next_id + len(novos[tabela])), **row})
        return next_id + len(novos[tabela]) - 1

    score_por_pessoa = {}
    perfil_por_pessoa = {}
    for person_id, c in sorted(CENARIOS.items()):
        renda = renda_considerada(c)
        vinculo_principal = _vinc_principal(c["vinculo"])

        # havendo mais de um score (lista de (score, data)), todos sao emitidos;
        # score_atual (usado no resto desta funcao) e o de data mais recente
        ultima_data_score = None
        for score_valor, score_data in _score_lista(c["score"]):
            score_id = add("credit_scores", {
                "person_id": person_id, "score": score_valor, "score_date": dia(score_data),
                "score_model": "serasa_v3",
                "score_reason": c["nota"],
                "payment_history": int(c["pontual"] * 40), "credit_usage": int((1 - c["utilizacao"]) * 25),
                "credit_age": 10 + person_id % 15, "credit_mix": 8 + person_id % 10,
                "recent_inquiries": max(0, 15 - c["consultas"] * 2),
                "risk_level": ("very_low" if score_valor >= 800 else "low" if score_valor >= 700
                               else "medium" if score_valor >= 550 else "high" if score_valor >= 400
                               else "very_high"),
                "default_probability": round(max(0.01, min(0.95, (1000 - score_valor) / 1000)), 3),
            })
            if ultima_data_score is None or score_data > ultima_data_score:
                score_por_pessoa[person_id] = score_id
                ultima_data_score = score_data
        score_atual = _score_atual(c["score"])

        perfil_por_pessoa[person_id] = add("financial_profiles", {
            "person_id": person_id, "profile_date": dia("2026-08-15"),
            "declared_monthly_income": c["renda"],
            "estimated_monthly_income": c["estimada"] or (
                None if "renda_estimada" in c["falta"] else (round(c["renda"] * 0.97) if c["renda"] else None)),
            "income_source": ("salary" if vinculo_principal and vinculo_principal[1] == CLT
                              else "benefit" if vinculo_principal and vinculo_principal[1] in (APO, PEN)
                              else "mixed"),
            "total_assets": None, "real_estate_value": None, "vehicles_value": None,
            "total_liabilities": None if "dti" in c["falta"] else round((renda or 0) * c["dti"] * 12, 2),
            "total_monthly_payments": None if "dti" in c["falta"] else round((renda or 0) * c["dti"], 2),
            # "dti_informado": so o indice sai nulo; as parcelas continuam
            # informadas, o que habilita o fallback de C2 (parcelas / renda)
            "debt_to_income_ratio": (None if "dti" in c["falta"] or "dti_informado" in c["falta"]
                                     else c["dti"]),
            "available_credit": None,
            "credit_utilization": None if "utilizacao" in c["falta"] else c["utilizacao"],
        })

        # havendo mais de um vinculo, um registro por entrada; a renda
        # considerada e o C6 usam sempre o vigente de maior salario
        for i, (empregador, tipo, inicio, salario, verificado) in enumerate(_vinc_lista(c["vinculo"])):
            add("employment_records", {
                "person_id": person_id, "employer_name": empregador,
                "employer_document": f"{(person_id * 1234567890123 + i) % 90000000000000 + 10000000000000}",
                "job_title": None, "employment_type": tipo, "salary": salario,
                "start_date": dia(inicio), "end_date": None, "is_current": True,
                "verification_status": "verified" if verificado else "unverified",
                "data_source": "eSocial",
            })

        # uma conta rotativa e, quando ha comprometimento relevante, um parcelado
        banco = list(BANCOS.values())[person_id % 4]
        limite = round(max(1000, (renda or 2000) * 1.5), 2)
        conta = add("credit_accounts", {
            "person_id": person_id, "account_type": "credit_card", "creditor": banco[0],
            "creditor_document": banco[1], "account_number": f"ACC{person_id:08d}",
            "opened_date": dia(datetime(2020 + person_id % 5, 1 + person_id % 12, 10)),
            "closed_date": None, "status": "active",
            "credit_limit": limite, "current_balance": round(limite * c["utilizacao"], 2),
            "available_credit": round(limite * (1 - c["utilizacao"]), 2),
            "original_amount": None, "remaining_amount": None, "interest_rate": 12.9,
            "monthly_payment": round(limite * c["utilizacao"] * 0.15, 2),
            "payment_due_day": 1 + person_id % 27, "number_of_payments": None,
            "remaining_payments": None,
            "payment_status": "current" if c["pontual"] >= 0.95 else "late" if c["pontual"] < 0.8 else "current",
            "days_late": 0 if c["pontual"] >= 0.9 else 15,
            "highest_days_late": int((1 - c["pontual"]) * 90),
            "times_late_30_days": int((1 - c["pontual"]) * 8),
            "times_late_60_days": int((1 - c["pontual"]) * 3),
            "times_late_90_days": int((1 - c["pontual"]) * 1),
            "last_reported_date": dia("2026-08-15"),
        })

        # doze meses de historico, com os atrasos distribuidos conforme a pontualidade
        atrasos = round((1 - c["pontual"]) * 12)
        for mes in range(0 if "historico_pagamento" in c["falta"] else 12):
            vencimento = HOJE - timedelta(days=30 * (12 - mes))
            atrasado = mes < atrasos
            add("payment_histories", {
                "person_id": person_id, "credit_account_id": conta, "debt_id": None,
                "payment_date": dia(vencimento, 12 if atrasado else 0),
                "due_date": dia(vencimento),
                "amount": round(limite * c["utilizacao"] * 0.15, 2),
                "amount_due": round(limite * c["utilizacao"] * 0.15, 2),
                "status": "late" if atrasado else "on_time",
                "days_late": 12 if atrasado else 0,
            })

        for valor, em_cobranca in c["dividas"]:
            # situacao nula (§8): a divida esta observavelmente em cobranca (por
            # isso in_collection fica true), mas nao da para saber se foi
            # quitada - so o status (que e o campo de situacao) fica nulo.
            # in_collection e nao-ponteiro no Go (bureau/entities/debt.go),
            # entao nunca pode ser emitido null - so status, que e *string.
            em_cobranca_efetivo = em_cobranca is not False
            add("debts", {
                "person_id": person_id, "debt_type": "loan", "creditor": banco[0],
                "creditor_document": banco[1], "original_amount": valor,
                "current_amount": round(valor * 1.08, 2), "interest_rate": 4.5, "fees": None,
                "origin_date": dia("2025-11-10"), "due_date": dia("2026-06-10"),
                "status": None if em_cobranca is None else "overdue" if em_cobranca else "open",
                "in_collection": em_cobranca_efetivo,
                "collection_date": dia("2026-07-01") if em_cobranca_efetivo else None,
                "collection_agency": "Cobranca Sintetica" if em_cobranca_efetivo else None,
                "settlement_amount": None, "settlement_date": None,
            })

        for valor, status, contestada in c["negativacoes"]:
            add("negative_records", {
                "person_id": person_id, "record_type": "spc", "creditor": banco[0],
                "creditor_document": banco[1], "amount": valor,
                "inclusion_date": dia("2022-05-18" if status == "paid" else "2026-03-12"),
                "contract_number": f"CT-{person_id:05d}", "status": status,
                "removal_date": dia("2023-08-01") if status == "paid" else None,
                "removal_reason": "payment" if status == "paid" else None,
                "process_number": None, "notary": None,
                "is_disputed": contestada,
                "dispute_date": dia("2026-04-02") if contestada else None,
                "dispute_reason": "Cliente alega quitacao" if contestada else None,
            })

        for tipo in c["processos"]:
            add("legal_records", {
                "person_id": person_id, "record_type": tipo,
                "process_number": f"{person_id:07d}-12.2025.8.26.0100",
                "court": "TJSP", "filing_date": dia("2025-04-22"), "status": "active",
                "amount": round((renda or 3000) * 2, 2),
                "description": "Acao de cobranca em andamento", "resolution": None,
                "resolution_date": None,
            })

        if c["fraude"]:
            severidade, status = c["fraude"]
            add("fraud_alerts", {
                "person_id": person_id, "alert_type": "identity", "severity": severidade,
                "description": "Divergencia cadastral apontada na validacao de documento",
                "detected_date": dia("2026-07-28"), "status": status,
                "resolved_date": None, "resolved_by": None, "notes": None,
            })

        for i in range(c["consultas"]):
            outro = list(BANCOS.values())[(person_id + i) % 4]
            add("credit_inquiries", {
                "person_id": person_id, "inquiry_date": dia(HOJE - timedelta(days=12 * (i + 1))),
                "inquiry_type": "credit_application", "creditor": outro[0],
                "creditor_document": outro[1], "purpose": "personal_loan",
                "amount": round((renda or 3000) * 3, 2),
                "result": "approved" if i % 3 == 0 else "denied",
            })

        fatores = []
        if any(s == "active" for _, s, _ in c["negativacoes"]):
            fatores.append("active_delinquency")
        if any(s == "paid" for _, s, _ in c["negativacoes"]):
            fatores.append("settled_delinquency")
        if c["utilizacao"] >= 0.6:
            fatores.append("high_utilization")
        if c["dti"] >= 0.4:
            fatores.append("high_debt_to_income")
        if c["consultas"] >= 5:
            fatores.append("multiple_recent_inquiries")
        if not fatores:
            fatores = ["clean_history", "stable_income"]
        add("risk_assessments", {
            "person_id": person_id, "assessment_date": dia("2026-08-20"),
            "assessment_type": "credit", "risk_score": score_atual or 0,
            "risk_level": "low" if (score_atual or 0) >= 700 else "medium" if (score_atual or 0) >= 500 else "high",
            "risk_factors": json.dumps(fatores, ensure_ascii=False),
            "recommendation": c["nota"], "model_version": "v3",
        })

    # metadados das fontes consultadas: liga cada pessoa a duas fontes existentes
    fontes = carrega("bureau", "data_sources")
    vinculos_prev = [v for v in carrega("bureau", "person_data_sources")
                     if v["person_id"] < PRIMEIRO_ID]
    vinculos = [{"person_id": person_id, "data_source_id": fontes[(person_id + i) % len(fontes)]["id"]}
                for person_id in sorted(CENARIOS) for i in range(2)] if fontes else []
    grava("bureau", "person_data_sources", vinculos_prev, vinculos)

    total = len(vinculos)
    for tabela, (prev, _) in tabelas.items():
        total += grava("bureau", tabela, prev, novos[tabela])
    print(f"birô: {total} registros")
    return score_por_pessoa, perfil_por_pessoa


# --------------------------------------------------------------------------- open finance
def open_finance():
    tabelas = {t: separa("open-finance", t) for t in
               ["bank_account_profiles", "bank_statements", "cash_flow_analyses",
                "recurring_transactions", "data_sharing_consents"]}
    novos = {t: [] for t in tabelas}

    def add(tabela, row):
        prev, next_id = tabelas[tabela]
        novos[tabela].append({**base(next_id + len(novos[tabela])), **row})
        return next_id + len(novos[tabela]) - 1

    perfil_por_pessoa = {}
    for person_id, c in sorted(CENARIOS.items()):
        # os demais dados de open finance (perfil, extratos, receitas recorrentes)
        # refletem sempre a analise de fluxo mais recente
        entrada, saida, volatilidade, dias_negativos, recorrente = _fluxo_atual(c["fluxo"])
        banco = list(BANCOS.values())[person_id % 4]
        meses = c["banco"][0]

        perfil_por_pessoa[person_id] = add("bank_account_profiles", {
            "person_id": person_id, "profile_date": dia("2026-08-20"),
            "banking_relationships": 1 + person_id % 3,
            "account_age_average": meses,
            "has_checking_account": True,
            "has_savings_account": entrada > saida,
            "has_investment_account": c["faixa"] == "bom" and entrada - saida > 1500,
            "investments_value": round((entrada - saida) * 12, 2) if entrada - saida > 1500 else None,
        })

        # tres extratos mensais, cobrindo os 90 dias previstos na metodologia
        saldo = round(max(50.0, (entrada - saida) * 2), 2)
        for mes in range(3):
            inicio = HOJE - timedelta(days=30 * (3 - mes))
            fim = inicio + timedelta(days=29)
            fechamento = round(saldo + (entrada - saida), 2)
            add("bank_statements", {
                "person_id": person_id, "institution": banco[0], "institution_document": banco[1],
                "account_type": "checking", "period_start": dia(inicio), "period_end": dia(fim),
                "opening_balance": saldo, "closing_balance": fechamento,
                "total_credits": round(entrada, 2), "total_debits": round(saida, 2),
                "transaction_count": 30 + person_id % 40, "currency": "BRL",
            })
            saldo = fechamento

        # havendo mais de uma analise (datas distintas), uma linha por entrada;
        # C7 e C9 usam sempre a mais recente (ja isolada acima em entrada/saida/etc)
        for data_analise, entrada_i, saida_i, volatilidade_i, dias_negativo_i, recorrente_i in \
                _fluxo_lista(c["fluxo"]):
            add("cash_flow_analyses", {
                "person_id": person_id, "analysis_date": dia(data_analise), "period_days": 90,
                "average_monthly_inflow": round(entrada_i, 2), "average_monthly_outflow": round(saida_i, 2),
                "net_cash_flow": round(entrada_i - saida_i, 2), "inflow_volatility": volatilidade_i,
                "negative_balance_days": dias_negativo_i, "has_recurring_income": recorrente_i,
            })

        categoria_receita, descricao_receita = {
            "salario": ("salary", "Salario"),
            "inss": ("benefit", "Beneficio INSS"),
            "beneficio_social": ("benefit", "Beneficio Social"),
            "pix": ("transfer", "Recebimentos via Pix"),
        }.get(c["origem_receita"], ("salary", "Salario"))
        recorrentes = [("income", categoria_receita, descricao_receita, round(entrada * 0.85, 2), recorrente),
                       ("expense", "rent", "Aluguel residencial", round(saida * 0.35, 2), True),
                       ("expense", "utility", "Energia eletrica", round(saida * 0.08, 2), True)]
        if c["dti"] >= 0.4:
            recorrentes.append(("expense", "loan", "Parcela de emprestimo", round(saida * 0.22, 2), True))
        vinculo_principal = _vinc_principal(c["vinculo"])
        for tipo, categoria, descricao, valor, ativo in recorrentes:
            add("recurring_transactions", {
                "person_id": person_id, "transaction_type": tipo, "category": categoria,
                "description": descricao, "amount": valor, "frequency": "monthly",
                "counterparty": banco[0] if tipo == "expense" else (
                    vinculo_principal[0] if vinculo_principal else None),
                "first_detected_date": dia(HOJE - timedelta(days=270)),
                "last_occurrence_date": dia(HOJE - timedelta(days=12)),
                "is_active": ativo,
            })

        consentimento = c["consentimento"]
        if consentimento is not None:
            expirado = consentimento == "expired"
            revogado = consentimento == "revoked"
            add("data_sharing_consents", {
                "person_id": person_id, "consent_id": f"urn:sintetico:consent:{person_id:04d}",
                "institution": banco[0],
                "status": "revoked" if revogado else "expired" if expirado else "granted",
                "scope": json.dumps(["ACCOUNTS_READ", "ACCOUNTS_BALANCES_READ", "RESOURCES_READ"]),
                "granted_at": dia(HOJE - timedelta(days=45 if not expirado else 400)),
                "expires_at": dia(HOJE + timedelta(days=135)) if not expirado
                              else dia(HOJE - timedelta(days=35)),
                "revoked_at": dia(HOJE - timedelta(days=20)) if revogado else None,
            })

    total = 0
    for tabela, (prev, _) in tabelas.items():
        total += grava("open-finance", tabela, prev, novos[tabela])
    print(f"open finance: {total} registros")
    return perfil_por_pessoa


# --------------------------------------------------------------------------- cadastro interno
def declaracao_de_renda(person_id, c):
    principal = _vinc_principal(c["vinculo"])
    return {
        "person_id": person_id, "declaration_date": dia("2026-02-15"),
        "income_type": ("salary" if principal and principal[1] == CLT
                        else "benefit" if principal and principal[1] in (APO, PEN)
                        else "self_employed"),
        "monthly_amount": float(c["renda"]), "yearly_amount": float(c["renda"] * 12),
        "source": "folha de pagamento" if principal else "declaracao do cliente",
        "verified": bool(principal and principal[4]),
        "verified_by": "analista" if principal and principal[4] else None,
        "proof_file_id": None,
    }


def cadastro_interno():
    tabelas = {t: separa("internal-registry", t) for t in
               ["customer_relationships", "contracted_products", "internal_payment_records",
                "pre_approved_limits", "income_declarations"]}
    novos = {t: [] for t in tabelas}

    def add(tabela, row):
        prev, next_id = tabelas[tabela]
        novos[tabela].append({**base(next_id + len(novos[tabela])), **row})
        return next_id + len(novos[tabela]) - 1

    relacionamento_por_pessoa = {}
    for person_id, c in sorted(CENARIOS.items()):
        meses, segmento, limite = c["banco"]
        renda = renda_considerada(c)

        # nao-cliente: nenhuma linha no registro interno
        if c["relacionamento"] is False:
            if c["renda"]:
                add("income_declarations", declaracao_de_renda(person_id, c))
            continue

        ativo = c["relacionamento"] != "inativo"
        interno = c["score_interno"]
        relacionamento_por_pessoa[person_id] = add("customer_relationships", {
            "person_id": person_id,
            "customer_since": dia(HOJE - timedelta(days=30 * meses)),
            "relationship_months": meses, "segment": segmento,
            "branch": f"{1000 + person_id}", "is_active": ativo,
            "churn_risk": "low" if meses >= 48 else "medium" if meses >= 18 else "high",
            # score_interno None deriva do birô; 0 significa "nunca calculado"
            "internal_score": (None if interno == 0 else interno if interno is not None
                               else min(1000, (_score_atual(c["score"]) or 500) + (10 if meses >= 48 else 0))),
        })

        # parcelas internas (abaixo) se vinculam sempre a conta corrente
        produto = add("contracted_products", {
            "person_id": person_id, "product_type": "checking_account",
            "product_name": "Conta Corrente", "contract_number": f"CC-{person_id:05d}",
            "contracted_date": dia(HOJE - timedelta(days=30 * meses)),
            "status": "active" if ativo else "closed",
            "balance": round((renda or 2000) * 0.3, 2), "monthly_value": 29.9,
        })
        if c["dti"] >= 0.3:
            add("contracted_products", {
                "person_id": person_id, "product_type": "loan",
                "product_name": "Emprestimo Pessoal", "contract_number": f"LN-{person_id:05d}",
                "contracted_date": dia(HOJE - timedelta(days=400)),
                "status": "active" if ativo else "closed",
                "balance": round((renda or 2000) * c["dti"] * 8, 2),
                "monthly_value": round((renda or 2000) * c["dti"] * 0.6, 2),
            })
        for tipo_produto in c["produtos"]:
            nome, prefixo = PRODUTOS_EXTRA[tipo_produto]
            add("contracted_products", {
                "person_id": person_id, "product_type": tipo_produto,
                "product_name": nome, "contract_number": f"{prefixo}-{person_id:05d}",
                "contracted_date": dia(HOJE - timedelta(days=30 * meses)),
                "status": "active" if ativo else "closed",
                "balance": round((renda or 2000) * 0.5, 2), "monthly_value": round((renda or 2000) * 0.05, 2),
            })

        # o comportamento interno e independente do birô: e o proprio knob que
        # decide quantas parcelas sairam do prazo e com que gravidade
        pontual_interno, atraso = c["interno"] if c["interno"] else (c["pontual"], 0)
        if pontual_interno != "vazio":
            fora_do_prazo = 12 - round(pontual_interno * 12)
            # ate 30 dias e atraso; acima disso a instituicao trata como perda
            situacao = "missed" if atraso > 30 else "late"
            for mes in range(12):
                vencimento = HOJE - timedelta(days=30 * (12 - mes))
                # os casos fora do prazo ficam nos meses mais recentes
                atrasado = mes >= 12 - fora_do_prazo
                add("internal_payment_records", {
                    "person_id": person_id, "contracted_product_id": produto,
                    "reference_month": dia(vencimento.replace(day=1)), "due_date": dia(vencimento),
                    "payment_date": None if atrasado and situacao == "missed"
                                    else dia(vencimento, atraso if atrasado else 0),
                    "amount_due": 29.9,
                    "amount_paid": 0.0 if atrasado and situacao == "missed" else 29.9,
                    "status": situacao if atrasado else "on_time",
                    "days_late": atraso if atrasado else 0,
                })

        # havendo mais de um valor pre-aprovado, uma linha por entrada; §7 usa
        # sempre a de maior valor aprovado ativo
        valido = c["pre_aprovado_valido"]
        for valor_pre in _lista(limite):
            add("pre_approved_limits", {
                "person_id": person_id, "product_type": "personal_loan",
                "approved_amount": float(valor_pre), "interest_rate": 3.2,
                "calculated_date": dia("2026-08-20"),
                "valid_until": dia(HOJE + timedelta(days=60)) if valido
                               else dia(HOJE - timedelta(days=25)),
                "policy_version": "v1.3", "is_active": valido,
            })

        if c["renda"]:
            add("income_declarations", declaracao_de_renda(person_id, c))

    total = 0
    for tabela, (prev, _) in tabelas.items():
        total += grava("internal-registry", tabela, prev, novos[tabela])
    print(f"cadastro interno: {total} registros")
    return relacionamento_por_pessoa


# --------------------------------------------------------------------------- validação cadastral
def validacao_cadastral():
    tabelas = {t: separa("registration-validation", t) for t in
               ["document_validations", "fiscal_regularities", "employment_link_validations",
                "compliance_checks"]}
    novos = {t: [] for t in tabelas}

    def add(tabela, row):
        prev, next_id = tabelas[tabela]
        novos[tabela].append({**base(next_id + len(novos[tabela])), **row})

    for person_id, c in sorted(CENARIOS.items()):
        situacao, cnd, pep, sancoes = c["receita"]
        documento = cpf(person_id)

        # "ausente": o cliente nunca passou por validacao cadastral
        if c["cadastral"] == "ausente":
            continue

        if c["validacao_anterior"]:
            # validacao antiga e divergente da atual: K10 e C13 usam sempre a
            # mais recente (a que segue abaixo, sempre datada de 2026-08-10)
            situacao_antiga, data_antiga = c["validacao_anterior"]
            add("document_validations", {
                "person_id": person_id, "validation_date": dia(data_antiga),
                "document_number": documento, "document_type": "cpf",
                "receita_federal_status": situacao_antiga,
                "is_valid": situacao_antiga not in ("canceled", "deceased"),
                "name_matches": "nome" not in c["divergencias"],
                "birth_date_matches": "nascimento" not in c["divergencias"],
                "biometric_validated": c["biometria"],
                "source": "receita_federal",
                "raw_response": json.dumps({"situacao": situacao_antiga, "consulta": dia(data_antiga)},
                                           ensure_ascii=False),
            })

        add("document_validations", {
            "person_id": person_id, "validation_date": dia("2026-08-10"),
            "document_number": documento, "document_type": "cpf",
            "receita_federal_status": situacao,
            # CPF pendente ou suspenso continua sendo um documento valido;
            # so cancelamento e obito invalidam o documento em si
            "is_valid": situacao not in ("canceled", "deceased"),
            "name_matches": "nome" not in c["divergencias"],
            "birth_date_matches": "nascimento" not in c["divergencias"],
            "biometric_validated": c["biometria"],
            "source": "receita_federal",
            "raw_response": json.dumps({"situacao": situacao, "consulta": dia("2026-08-10")},
                                       ensure_ascii=False),
        })

        pendencias = []
        if cnd == "irregular":
            pendencias = ["divida_ativa_uniao", "cnd_negada"]
        elif cnd == "suspended":
            pendencias = ["parcelamento_em_curso"]
        add("fiscal_regularities", {
            "person_id": person_id, "check_date": dia("2026-08-10"),
            "has_debts": cnd != "regular", "cnd_status": cnd,
            "cnd_number": f"CND-2026-{person_id:07d}" if cnd == "regular" else None,
            "cnd_issue_date": dia("2026-08-10") if cnd == "regular" else None,
            "cnd_valid_until": dia("2026-11-08") if cnd == "regular" else None,
            "pending_issues": json.dumps(pendencias, ensure_ascii=False) if pendencias else None,
        })

        # o eSocial e fonte independente: pode nao ter o(s) vinculo(s) que o
        # birô declara; havendo mais de um vinculo, um registro por entrada
        if c["esocial"] != "ausente":
            encerrado = c["esocial"] == "encerrado"
            for i, (empregador, tipo, inicio, _, verificado) in enumerate(_vinc_lista(c["vinculo"])):
                add("employment_link_validations", {
                    "person_id": person_id, "validation_date": dia("2026-08-10"),
                    "employer_name": empregador,
                    "employer_document": f"{(person_id * 1234567890123 + i) % 90000000000000 + 10000000000000}",
                    "employment_type": tipo,
                    "status": "terminated" if encerrado else "active",
                    "start_date": dia(inicio),
                    "end_date": dia(HOJE - timedelta(days=120)) if encerrado else None,
                    "source": "eSocial", "verified": verificado,
                })

        if c["cadastral"] == "sem_compliance":
            continue

        impeditivo = sancoes or situacao in ("canceled", "deceased")
        status = ("irregular" if impeditivo or situacao in ("pending", "suspended")
                  or c["triagem"] == "irregular"
                  else "attention" if pep else "regular")
        add("compliance_checks", {
            "person_id": person_id, "check_type": "kyc_full", "check_date": dia("2026-08-10"),
            "status": status,
            "details": json.dumps({"motivo": c["nota"]}, ensure_ascii=False) if status != "regular" else None,
            "is_pep": pep,
            "pep_details": "Contrato vigente com orgao publico" if pep else None,
            "on_sanctions_list": sancoes,
            "sanctions_details": "Homonimo em lista restritiva" if sancoes else None,
            "valid_until": dia("2027-02-06"),
        })

    total = 0
    for tabela, (prev, _) in tabelas.items():
        total += grava("registration-validation", tabela, prev, novos[tabela])
    print(f"validacao cadastral: {total} registros")


# --------------------------------------------------------------------------- gabarito
# Implementacao da politica de credito v1.3 (docs/instructions.md) aplicada aos
# knobs dos cenarios. Serve de gabarito para conferir a resposta do agente e,
# durante a geracao, para provar que a distribuicao e a pretendida.
def _faixa(valor, faixas, fora=0):
    for limite, pontos in faixas:
        if valor >= limite:
            return pontos
    return fora


def avalia(c):
    ausentes, criterios = [], {}
    renda = renda_considerada(c)
    score_atual = _score_atual(c["score"])
    vinculo_principal = _vinc_principal(c["vinculo"])
    fluxo_atual = _fluxo_atual(c["fluxo"])
    # situacao nula em negativacao e tratada como ativa (conservador, §8)
    negativas_ativas = [v for v, st, _ in c["negativacoes"] if st == "active" or st is None]
    consentido = c["consentimento"] == "granted"
    situacao, cnd, pep, sancoes = c["receita"]
    pontual_interno, atraso_interno = c["interno"] if c["interno"] else (c["pontual"], 0)
    sem_interno = c["relacionamento"] is False or pontual_interno == "vazio"

    def marca(nome, pontos, observado, ausente=False):
        criterios[nome] = {"pontos": pontos, "observado": observado, "ausente": ausente}
        if ausente:
            ausentes.append(nome)
        return pontos

    # ---- §3 dados minimos -------------------------------------------------
    if score_atual is None:
        return {"decisao": "ANALISE_MANUAL", "regra": "M1", "pontuacao": None,
                "criterios": {}, "ausentes": [], "limite": None}
    if renda is None:
        return {"decisao": "ANALISE_MANUAL", "regra": "M2", "pontuacao": None,
                "criterios": {}, "ausentes": [], "limite": None}
    if "cpf" in c["falta"]:
        return {"decisao": "ANALISE_MANUAL", "regra": "M3", "pontuacao": None,
                "criterios": {}, "ausentes": [], "limite": None}
    if any(st == "active" and contestada for _, st, contestada in c["negativacoes"]):
        return {"decisao": "ANALISE_MANUAL", "regra": "M4", "pontuacao": None,
                "criterios": {}, "ausentes": [], "limite": None}

    # ---- §4 eliminatorias, reprovacao antes de analise manual --------------
    receita_recorrente = round(fluxo_atual[0] * 0.85, 2) if fluxo_atual[4] else 0.0
    # situacao nula em divida e tratada como nao quitada (conservador, §8)
    divida_cobranca = sum(round(v * 1.08, 2) for v, cobranca in c["dividas"] if cobranca is not False)

    reprova, manual = [], []
    if c["fraude"] and c["fraude"][1] == "active" and c["fraude"][0] in ("high", "critical"):
        reprova.append("K1")
    elif c["fraude"] and c["fraude"][1] in ("active", "investigating"):
        manual.append("K2")
    if "bankruptcy" in c["processos"]:
        reprova.append("K3")
    if score_atual < 300:
        reprova.append("K4")
    if len(negativas_ativas) >= 3:
        reprova.append("K5")
    if sum(negativas_ativas) > 5 * renda:
        reprova.append("K6")
    if divida_cobranca > 10 * renda:
        reprova.append("K7")
    if consentido and receita_recorrente > 0 and renda > 2 * receita_recorrente:
        manual.append("K8")
    if not sem_interno and pontual_interno < 1.0 and atraso_interno > 90:
        reprova.append("K9")
    elif not sem_interno and pontual_interno < 1.0 and 30 < atraso_interno <= 90:
        manual.append("K9-L")
    if c["cadastral"] != "ausente":
        if situacao in ("canceled", "deceased"):
            reprova.append("K10")
        elif situacao in ("pending", "suspended"):
            manual.append("K11")
        if c["cadastral"] != "sem_compliance":
            if sancoes:
                reprova.append("K12")
            elif pep:
                manual.append("K13")
            elif c["triagem"] == "irregular":
                manual.append("K13")

    if reprova or manual:
        vencedora = reprova[0] if reprova else manual[0]
        return {"decisao": "REPROVADO" if reprova else "ANALISE_MANUAL", "regra": vencedora,
                "pontuacao": None, "criterios": {}, "ausentes": [], "limite": None}

    # ---- §5 scorecard ------------------------------------------------------
    total = 0
    total += marca("C1", _faixa(score_atual, [(800, 22), (700, 19), (600, 14), (500, 10), (400, 5), (300, 2)]),
                   f"score {score_atual}")

    if "dti" in c["falta"]:
        dti_efetivo = 0.50
        total += marca("C2", 0, "comprometimento indeterminavel", ausente=True)
    elif "dti_informado" in c["falta"]:
        # indice nulo, mas parcela mensal informada: cai no fallback da
        # politica (parcelas / renda considerada); nao e dado ausente
        parcelas = round((renda or 0) * c["dti"], 2)
        dti_efetivo = round(parcelas / renda, 4)
        total += marca("C2", 11 if dti_efetivo <= 0.30 else 8 if dti_efetivo <= 0.40
                       else 4 if dti_efetivo <= 0.50 else 0,
                       f"comprometimento {dti_efetivo:.2f} (calculado por parcelas / renda)")
    else:
        dti_efetivo = c["dti"]
        total += marca("C2", 11 if c["dti"] <= 0.30 else 8 if c["dti"] <= 0.40 else 4 if c["dti"] <= 0.50 else 0,
                       f"comprometimento {c['dti']:.2f}")

    if "historico_pagamento" in c["falta"]:
        total += marca("C3", 4, "sem pagamentos no periodo", ausente=True)
    else:
        total += marca("C3", _faixa(c["pontual"], [(0.95, 11), (0.85, 8), (0.70, 4)]),
                       f"{c['pontual']:.1%} em dia")

    total += marca("C4", 9 if not negativas_ativas
                   else (5 if sum(negativas_ativas) <= 1000 else 2) if len(negativas_ativas) == 1
                   else 1,
                   f"{len(negativas_ativas)} negativacao(oes) ativa(s) somando R$ {sum(negativas_ativas):,.2f}")

    if "utilizacao" in c["falta"]:
        total += marca("C5", 1, "utilizacao nao informada", ausente=True)
    else:
        total += marca("C5", 3 if c["utilizacao"] < 0.30 else 2 if c["utilizacao"] < 0.60
                       else 1 if c["utilizacao"] < 0.90 else 0, f"utilizacao {c['utilizacao']:.2f}")

    total += marca("C6", 2 if vinculo_principal and vinculo_principal[4] else 1 if vinculo_principal else 0,
                   "vinculo verificado" if vinculo_principal and vinculo_principal[4]
                   else "vinculo nao verificado" if vinculo_principal else "sem vinculo atual")

    if not consentido:
        total += marca("C7", 3, "sem consentimento ativo", ausente=True)
        total += marca("C8", 2, "sem consentimento ativo", ausente=True)
        total += marca("C9", 1, "sem consentimento ativo", ausente=True)
        fluxo_liquido = None
    else:
        entrada, saida, _, dias_negativo, _ = fluxo_atual
        fluxo_liquido = entrada - saida
        razao = fluxo_liquido / renda
        total += marca("C7", 7 if razao >= 0.20 else 5 if razao >= 0.10 else 3 if razao >= 0 else 0,
                       f"fluxo liquido R$ {fluxo_liquido:,.2f} ({razao:.2f} da renda)")
        total += marca("C8", 4 if receita_recorrente >= 0.8 * renda else 2 if receita_recorrente > 0 else 0,
                       f"receita recorrente R$ {receita_recorrente:,.2f}")
        total += marca("C9", _faixa(-dias_negativo, [(0, 4), (-5, 2), (-15, 1)]),
                       f"{dias_negativo} dia(s) com saldo negativo")

    meses = c["banco"][0]
    if c["relacionamento"] is False:
        total += marca("C10", 2, "nao-cliente da instituicao", ausente=True)
        total += marca("C11", 2, "nao-cliente da instituicao", ausente=True)
    else:
        inativo = c["relacionamento"] == "inativo"
        total += marca("C10", 0 if inativo else _faixa(meses, [(60, 5), (24, 3), (6, 2)]),
                       "relacionamento encerrado" if inativo else f"{meses} meses de relacionamento")
        interno = c["score_interno"]
        efetivo = (None if interno == 0 else interno if interno is not None
                   else min(1000, score_atual + (10 if meses >= 48 else 0)))
        if efetivo is None:
            total += marca("C11", 2, "score interno nunca calculado", ausente=True)
        else:
            total += marca("C11", _faixa(efetivo, [(800, 5), (600, 3), (400, 2)]), f"score interno {efetivo}")

    if sem_interno:
        total += marca("C12", 2, "sem parcelas internas no periodo", ausente=True)
    else:
        total += marca("C12", _faixa(pontual_interno, [(0.95, 5), (0.85, 4), (0.70, 2)]),
                       f"{pontual_interno:.1%} das parcelas internas em dia")

    if c["cadastral"] == "ausente":
        total += marca("C13", 2, "documento nunca validado", ausente=True)
        total += marca("C14", 2, "regularidade fiscal nunca verificada", ausente=True)
    else:
        divergencias = len(c["divergencias"] & {"nome", "nascimento"})
        total += marca("C13", 0 if divergencias >= 2 else 2 if divergencias == 1
                       else 5 if c["biometria"] else 4,
                       f"{divergencias} divergencia(s), biometria "
                       + ("validada" if c["biometria"] else "nao validada"))
        total += marca("C14", 4 if cnd == "regular" else 2 if cnd == "suspended" else 0,
                       f"CND {cnd}")

    if not vinculo_principal or c["esocial"] == "ausente":
        total += marca("C15", 0, "sem registro no eSocial")
    elif c["esocial"] == "encerrado":
        total += marca("C15", 1, "vinculo encerrado no eSocial")
    else:
        total += marca("C15", 3 if vinculo_principal[4] else 2,
                       "vinculo ativo confirmado" if vinculo_principal[4] else "vinculo ativo nao confirmado")

    # ---- §6 decisao --------------------------------------------------------
    decisao = ("APROVADO" if total >= 70 else "APROVADO_COM_RESSALVAS" if total >= 50
               else "ANALISE_MANUAL" if total >= 35 else "REPROVADO")
    degradou = decisao == "APROVADO" and len(ausentes) >= 2
    if degradou:
        decisao = "APROVADO_COM_RESSALVAS"

    # ---- §7 limite ---------------------------------------------------------
    limite = None
    if decisao in ("APROVADO", "APROVADO_COM_RESSALVAS"):
        fator = 3.0 if decisao == "APROVADO" else 1.5
        limite = renda * fator * (1 - dti_efetivo)
        if fluxo_liquido and fluxo_liquido > 0:
            limite = min(limite, fluxo_liquido * 12)
        # havendo mais de uma pre-aprovacao ativa, usa-se a de maior valor
        pre_lista = _lista(c["banco"][2])
        pre = max(pre_lista) if pre_lista else None
        if decisao == "APROVADO" and pre and c["pre_aprovado_valido"]:
            limite = max(limite, float(pre))
        limite = int(limite // 100) * 100

    return {"decisao": decisao, "regra": None, "pontuacao": total, "criterios": criterios,
            "ausentes": ausentes, "limite": limite, "degradou": degradou}


def gabarito():
    linhas = ["# Gabarito dos cenarios sinteticos",
              "",
              "> Gerado por `generate_cenarios.py` aplicando a politica v1.3 de",
              "> `docs/instructions.md`. Regenerado junto com as fixtures.",
              ""]
    contagem, faixas = {}, {}
    detalhe = []
    for person_id, c in sorted(CENARIOS.items()):
        r = avalia(c)
        contagem[r["decisao"]] = contagem.get(r["decisao"], 0) + 1
        faixas[c["faixa"]] = faixas.get(c["faixa"], 0) + 1
        pontos = "-" if r["pontuacao"] is None else str(r["pontuacao"])
        limite = "-" if r["limite"] is None else f"R$ {r['limite']:,.0f}".replace(",", ".")
        detalhe.append((person_id, c, r, pontos, limite))

    linhas += ["## Distribuicao", "",
               "| Faixa | Cenarios | %  |", "|---|---|---|"]
    for faixa in ("bom", "complexo", "ruim"):
        n = faixas.get(faixa, 0)
        linhas.append(f"| {faixa} | {n} | {n * 100 // len(CENARIOS)}% |")
    linhas += ["", "| Decisao | Cenarios |", "|---|---|"]
    for decisao, n in sorted(contagem.items(), key=lambda kv: -kv[1]):
        linhas.append(f"| `{decisao}` | {n} |")

    linhas += ["", "## Por cenario", "",
               "| # | Faixa | O que exercita | Decisao | Pontos | Regra | Limite | Criterios sem dado |",
               "|---|---|---|---|---|---|---|---|"]
    for person_id, c, r, pontos, limite in detalhe:
        ausentes = ", ".join(r["ausentes"]) or "-"
        linhas.append(f"| {person_id} | {c['faixa']} | {c['nota']} | `{r['decisao']}` | "
                      f"{pontos} | {r['regra'] or '-'} | {limite} | {ausentes} |")
    linhas.append("")
    (ROOT / "docs" / "gabarito_cenarios.md").write_text("\n".join(linhas))
    return contagem, faixas


def main():
    pi_por_pessoa = cadastro()
    scores, perfis = biro()
    contas = open_finance()
    relacionamentos = cadastro_interno()
    validacao_cadastral()

    pessoas_por_fonte(pi_por_pessoa, [
        ("bureau", "credit_score_id", lambda p, _: scores.get(p)),
        ("open-finance", "bank_account_profile_id", lambda p, _: contas.get(p)),
        ("internal-registry", "customer_relationship_id", lambda p, _: relacionamentos.get(p)),
        ("registration-validation", None, None),
    ])
    # o birô guarda tambem o perfil financeiro na pessoa
    prev, next_id = separa("bureau", "persons", "id")
    rows = carrega("bureau", "persons")
    for offset, person_id in enumerate(sorted(CENARIOS)):
        rows[len(prev) + offset]["financial_profile_id"] = perfis.get(person_id)
    fixtures("bureau", "persons").write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")

    contagem, faixas = gabarito()
    print("cenarios por faixa:", faixas)
    print("decisao esperada:", contagem)
    print("gabarito: docs/gabarito_cenarios.md")
    print("rode agora: python3 generate_cadastro.py")


if __name__ == "__main__":
    main()
