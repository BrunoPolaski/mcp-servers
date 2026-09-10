#!/usr/bin/env python3
"""Gera os 30 cenarios de avaliacao (ids 11 a 40) nas quatro fontes.

Sao o conjunto definitivo de avaliacao e os unicos clientes das fixtures - os
dez preliminares, usados na validacao inicial so com o biro, foram removidos.
Estao divididos em tres faixas de complexidade:

  simples (1-10)         historico regular, renda comprovada, sem restricoes
  intermediarios (11-20) score medio, atrasos ja regularizados, renda variavel,
                         consultas recentes em varias instituicoes
  complexos (21-30)      score baixo com justificativa plausivel (21-25) ou
                         score alto com sinais de alerta (26-30)

Cada cenario e declarado uma unica vez em CENARIOS e derivado dali para as
quatro fontes, de modo que elas contem a mesma historia: quem tem score baixo
por inadimplencia antiga aparece com a negativacao quitada no biro e com fluxo
de caixa positivo no open finance; quem tem alerta de fraude aparece com
divergencia cadastral na Receita; quem nao tem vinculo no biro tambem nao tem
vinculo no eSocial.

As lacunas de dados (`falta`) sao deliberadas: existem para exercitar as regras
de dado ausente da politica de credito parametrizada.

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
# Os 30 cenarios. Cada chave e um knob; o resto e derivado.
#
#   score        None = biro sem score (falha a pre-condicao de dados minimos)
#   renda        renda mensal declarada; None = nao declarada
#   estimada     renda estimada pelo biro; None = nao estimada
#   dti          comprometimento de renda
#   utilizacao   utilizacao do limite rotativo
#   pontual      proporcao de parcelas pagas em dia nos ultimos 12 meses
#   negativacoes (valor, status, contestada) - status: active | paid
#   dividas      (valor, em_cobranca)
#   processos    tipos de acao judicial
#   fraude       (severidade, status) ou None
#   consultas    consultas de credito nos ultimos 90 dias
#   vinculo      (empregador, tipo, inicio, salario, verificado) ou None
#   fluxo        (entrada, saida, volatilidade, dias_negativo, renda_recorrente)
#   banco        (meses de relacionamento, segmento, limite pre-aprovado)
#   receita      (situacao do CPF, CND, PEP, sancoes)
#   divergencias campos do cadastro que nao conferem com a Receita Federal
#   falta        campos deliberadamente ausentes
#   nota         o que o cenario exercita
# ---------------------------------------------------------------------------
def cenario(faixa, score, renda, dti, utilizacao, pontual, nota, *, estimada=None,
            negativacoes=(), dividas=(), processos=(), fraude=None, consultas=1,
            vinculo=None, fluxo=None, banco=(24, "varejo", 0.0), receita=("regular", "regular", False, False),
            divergencias=(), mudou=False, falta=()):
    return dict(faixa=faixa, score=score, renda=renda, estimada=estimada, dti=dti,
                utilizacao=utilizacao, pontual=pontual, negativacoes=list(negativacoes),
                dividas=list(dividas), processos=list(processos), fraude=fraude,
                consultas=consultas, vinculo=vinculo, fluxo=fluxo, banco=banco,
                receita=receita, divergencias=set(divergencias), mudou=mudou,
                falta=set(falta), nota=nota)


CLT = "CLT"
PJ = "PJ"
AUT = "autonomous"

CENARIOS = {
    # ---------------- simples (1-10): baixo risco, aprovacao quase automatica -
    1: cenario("simples", 812, 6800, 0.18, 0.12, 1.00, "assalariado estavel, sem restricoes",
                vinculo=("Industria XYZ", CLT, "2017-03-01", 6800, True),
                fluxo=(6900, 5100, 0.05, 0, True), banco=(96, "varejo", 25000)),
    2: cenario("simples", 845, 9200, 0.15, 0.09, 1.00, "renda alta, relacionamento longo",
                vinculo=("Tech Solutions Ltda", CLT, "2015-08-01", 9200, True),
                fluxo=(9400, 6800, 0.04, 0, True), banco=(132, "alta renda", 60000)),
    3: cenario("simples", 764, 4300, 0.22, 0.19, 0.98,
               "servidor publico; mudou de endereco e o birô nao foi informado",
               mudou=True,
                vinculo=("Prefeitura Sintetica", "estatutario", "2019-02-01", 4300, True),
                fluxo=(4350, 3300, 0.03, 0, True), banco=(72, "varejo", 15000)),
    4: cenario("simples", 795, 5600, 0.20, 0.14, 0.97, "CLT com investimentos",
                vinculo=("Comercio ABC", CLT, "2018-06-01", 5600, True),
                fluxo=(5700, 4100, 0.06, 0, True), banco=(84, "varejo", 20000)),
    5: cenario("simples", 878, 12500, 0.11, 0.07, 1.00, "alta renda, sem passivo relevante",
                vinculo=("Consultoria Delta", PJ, "2016-01-01", 12500, True),
                fluxo=(12800, 8200, 0.07, 0, True), banco=(120, "alta renda", 90000)),
    6: cenario("simples", 731, 3900, 0.26, 0.24, 0.96, "renda menor, mas pontual",
                vinculo=("Comercio ABC", CLT, "2020-04-01", 3900, True),
                fluxo=(3950, 3050, 0.05, 0, True), banco=(60, "varejo", 9000)),
    7: cenario("simples", 803, 7400, 0.17, 0.11, 1.00, "dois vinculos, o anterior encerrado sem pendencia",
                vinculo=("Industria XYZ", CLT, "2021-09-01", 7400, True),
                fluxo=(7500, 5300, 0.05, 0, True), banco=(48, "varejo", 28000)),
    8: cenario("simples", 758, 4800, 0.24, 0.21, 0.95, "cliente novo no banco, historico externo limpo",
                vinculo=("Comercio ABC", CLT, "2019-11-01", 4800, True),
                fluxo=(4850, 3700, 0.06, 0, True), banco=(14, "varejo", 7000)),
    9: cenario("simples", 826, 8100, 0.16, 0.10, 1.00, "MEI consolidado, faturamento estavel",
                vinculo=("MEI Servicos", AUT, "2017-05-01", 8100, True),
                fluxo=(8200, 6000, 0.09, 0, True), banco=(90, "varejo", 32000)),
    10: cenario("simples", 771, 5200, 0.23, 0.18, 0.98, "sem telefone no cadastro do biro",
                vinculo=("Comercio ABC", CLT, "2018-10-01", 5200, True),
                fluxo=(5250, 4000, 0.05, 0, True), banco=(66, "varejo", 16000),
                falta=("telefone",)),

    # ------------- intermediarios (11-20): exigem analise contextual ---------
    11: cenario("intermediario", 648, 4100, 0.34, 0.41, 0.92, "dois atrasos curtos ja regularizados",
                negativacoes=[(1850.0, "paid", False)], consultas=3,
                vinculo=("Comercio ABC", CLT, "2021-02-01", 4100, True),
                fluxo=(4150, 3600, 0.14, 2, True), banco=(38, "varejo", 6000)),
    12: cenario("intermediario", 604, None, 0.38, 0.47, 0.88, "autonomo sem renda declarada; biro estima",
                estimada=3800, consultas=4,
                vinculo=("Consultoria Delta", AUT, "2020-07-01", None, False),
                fluxo=(3750, 3400, 0.28, 5, True), banco=(30, "varejo", 3500),
                falta=("renda_declarada",)),
    13: cenario("intermediario", 671, 5400, 0.31, 0.36, 0.94, "renda variavel de comissao",
                consultas=3, vinculo=("Comercio ABC", CLT, "2019-05-01", 5400, True),
                fluxo=(5500, 4700, 0.31, 3, True), banco=(52, "varejo", 8000)),
    14: cenario("intermediario", 559, 3200, 0.45, 0.58, 0.83, "endividamento perto do limite da politica",
                dividas=[(4200.0, False)], consultas=5,
                vinculo=("Industria XYZ", CLT, "2022-03-01", 3200, True),
                fluxo=(3250, 3050, 0.19, 6, True), banco=(26, "varejo", 2500)),
    15: cenario("intermediario", 622, 4600, 0.36, 0.44, 0.90, "cinco consultas em 90 dias",
                consultas=7, vinculo=("Comercio ABC", CLT, "2020-11-01", 4600, True),
                fluxo=(4650, 3900, 0.16, 1, True), banco=(44, "varejo", 5500)),
    16: cenario("intermediario", 587, 2900, 0.42, 0.52, 0.86, "negativacao quitada ha um ano",
                negativacoes=[(2300.0, "paid", False)], consultas=4,
                vinculo=("Comercio ABC", CLT, "2021-08-01", 2900, True),
                fluxo=(2950, 2600, 0.21, 4, True), banco=(33, "varejo", 2000)),
    17: cenario("intermediario", 663, 6100, 0.29, 0.33, 0.93, "PJ com sazonalidade forte",
                consultas=2, vinculo=("Consultoria Delta", PJ, "2018-09-01", 6100, True),
                fluxo=(6300, 5400, 0.37, 7, True), banco=(58, "varejo", 11000)),
    18: cenario("intermediario", 531, 3600, 0.48, 0.63, 0.79,
                "vinculo nao confirmado no eSocial e nascimento divergente da Receita",
                divergencias=("nascimento",),
                dividas=[(3100.0, False)], consultas=5,
                vinculo=("Comercio ABC", CLT, "2023-01-01", 3600, False),
                fluxo=(3600, 3350, 0.23, 8, True), banco=(20, "varejo", 1500),
                receita=("regular", "irregular", False, False)),
    19: cenario("intermediario", 695, 7200, 0.28, 0.30, 0.95, "processo civil em andamento, sem penhora",
                processos=["civil"], consultas=3,
                vinculo=("Industria XYZ", CLT, "2017-12-01", 7200, True),
                fluxo=(7300, 5900, 0.12, 0, True), banco=(70, "varejo", 14000)),
    20: cenario("intermediario", 576, 3400, 0.44, 0.55, 0.85, "sem e-mail e sem renda estimada no biro",
                consultas=4, vinculo=("Comercio ABC", CLT, "2021-06-01", 3400, True),
                fluxo=(3450, 3000, 0.18, 3, True), banco=(28, "varejo", 2200),
                falta=("email", "renda_estimada")),

    # ------- complexos A (21-25): score baixo com justificativa plausivel ----
    21: cenario("complexo", 342, 4500, 0.29, 0.35, 0.92,
                "inadimplencia de 2022 quitada; fluxo atual positivo no open finance",
                negativacoes=[(6800.0, "paid", False)], consultas=2,
                vinculo=("Industria XYZ", CLT, "2023-04-01", 4500, True),
                fluxo=(4600, 3500, 0.09, 0, True), banco=(41, "varejo", 0.0)),
    22: cenario("complexo", 388, 5200, 0.33, 0.40, 0.90,
                "recuperacao judicial ativa da empresa do titular: regra eliminatoria",
                negativacoes=[(4100.0, "paid", False)], processos=["recuperacao_judicial"], consultas=3,
                vinculo=("Comercio ABC", CLT, "2024-02-01", 5200, True),
                fluxo=(5300, 4100, 0.11, 0, True), banco=(63, "varejo", 0.0)),
    23: cenario("complexo", 296, 3800, 0.37, 0.49, 0.87,
                "score abaixo de 300: reprovacao por regra eliminatoria", mudou=True,
                negativacoes=[(5200.0, "paid", False)], dividas=[(2800.0, False)], consultas=4,
                vinculo=("Comercio ABC", CLT, "2023-09-01", 3800, True),
                fluxo=(3850, 3200, 0.15, 1, True), banco=(35, "varejo", 0.0)),
    24: cenario("complexo", 412, 6300, 0.31, 0.38, 0.91,
                "negativacao ativa e contestada ao mesmo tempo: falha a pre-condicao",
                negativacoes=[(3900.0, "active", True)], consultas=3,
                vinculo=("Consultoria Delta", PJ, "2021-10-01", 6300, True),
                fluxo=(6400, 4900, 0.13, 0, True), banco=(55, "varejo", 0.0)),
    25: cenario("complexo", 357, None, 0.36, 0.44, 0.89,
                "sem renda determinavel: sem declaracao, sem estimativa e sem vinculo vigente",
                negativacoes=[(2600.0, "paid", False)], consultas=2, vinculo=None,
                fluxo=(2400, 2300, 0.34, 9, False), banco=(22, "varejo", 0.0),
                receita=("regular", "irregular", False, False),
                falta=("renda_declarada", "renda_estimada", "vinculo")),

    # ------- complexos B (26-30): score alto com sinais de alerta ------------
    26: cenario("complexo", 806, 7600, 0.47, 0.71, 0.94,
                "score alto, mas dividas em cobranca acima de dez vezes a renda",
                dividas=[(48000.0, True), (34000.0, True)], consultas=6,
                vinculo=("Industria XYZ", CLT, "2019-03-01", 7600, True),
                fluxo=(7700, 7400, 0.22, 4, True), banco=(77, "alta renda", 12000)),
    27: cenario("complexo", 772, 8400, 0.29, 0.26, 0.96,
                "score alto, tres negativacoes ativas e aumento subito de despesas",
                negativacoes=[(16000.0, "active", False), (15000.0, "active", False),
                              (13000.0, "active", False)], consultas=8, vinculo=("Tech Solutions Ltda", CLT, "2018-01-01", 8400, True),
                fluxo=(8500, 9700, 0.41, 11, True), banco=(88, "alta renda", 30000)),
    28: cenario("complexo", 758, 6900, 0.34, 0.33, 0.95,
                "alerta de fraude grave ativo, com divergencia de nome: reprovacao",
                fraude=("high", "active"), consultas=5, divergencias=("nome",),
                vinculo=("Comercio ABC", CLT, "2020-05-01", 6900, True),
                fluxo=(7000, 5600, 0.14, 1, True), banco=(59, "varejo", 18000),
                receita=("pending", "regular", False, False)),
    29: cenario("complexo", 791, 9600, 0.30, 0.28, 0.97,
                "PEP com contrato publico e alerta em apuracao: encaminhamento manual",
                fraude=("medium", "investigating"), consultas=3, vinculo=("Consultoria Delta", PJ, "2017-07-01", 9600, True),
                fluxo=(9700, 7100, 0.10, 0, True), banco=(104, "alta renda", 45000),
                receita=("regular", "regular", True, False)),
    30: cenario("complexo", None, 5800, 0.32, 0.31, 0.93,
                "biro sem score calculado: falha a pre-condicao de dados minimos",
                consultas=2, vinculo=("Comercio ABC", CLT, "2022-08-01", 5800, True),
                fluxo=(5900, 4600, 0.12, 0, True), banco=(47, "varejo", 8000),
                falta=("score",)),
}

NOMES = ["Amanda Ferreira Lima", "Bruno Carvalho Nunes", "Carla Menezes Rocha", "Diego Tavares Pinto",
         "Elaine Moraes Cardoso", "Fabio Antunes Vieira", "Giovana Peixoto Braga", "Helio Marques Fontes",
         "Isadora Campos Teles", "Joao Batista Siqueira", "Karina Duarte Prado", "Leandro Bastos Farias",
         "Mariana Freitas Lopes", "Nelson Aguiar Bittencourt", "Olivia Rezende Sampaio", "Paulo Cesar Andrade",
         "Queila Monteiro Serra", "Rafael Guedes Pontes", "Simone Vasques Coelho", "Tiago Meireles Bandeira",
         "Ursula Barreto Galvao", "Vinicius Aragao Quintela", "Wanda Cordeiro Lisboa", "Xavier Toledo Amancio",
         "Yara Bonfim Salgado", "Zeno Villela Krause", "Alice Padilha Moura", "Benedito Falcao Ximenes",
         "Clarice Uchoa Bezerra", "Danilo Espindola Rangel"]
MAES = ["Marta Ferreira", "Sonia Carvalho", "Regina Menezes", "Laura Tavares", "Cristina Moraes",
        "Vera Antunes", "Beatriz Peixoto", "Alice Marques", "Neusa Campos", "Rita Batista",
        "Eliane Duarte", "Silvia Bastos", "Angela Freitas", "Marlene Aguiar", "Cecilia Rezende",
        "Dulce Andrade", "Ivone Monteiro", "Nadia Guedes", "Tereza Vasques", "Aparecida Meireles",
        "Rosane Barreto", "Ilda Aragao", "Selma Cordeiro", "Miriam Toledo", "Julia Bonfim",
        "Norma Villela", "Solange Padilha", "Adelia Falcao", "Zuleica Uchoa", "Berenice Espindola"]


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


def renda_considerada(c):
    """Precedencia da politica: declarada, estimada, salario de vinculo vigente."""
    if c["renda"]:
        return c["renda"]
    if c["estimada"]:
        return c["estimada"]
    if c["vinculo"] and c["vinculo"][3]:
        return c["vinculo"][3]
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
        nascimento = datetime(1996 - offset % 30, 1 + (offset * 5) % 12, 1 + (offset * 7) % 27)

        pessoas.append({
            **base(pi_id + offset),
            "full_name": nome,
            "mother_name": MAES[offset],
            "birth_date": dia(nascimento),
            "gender": "female" if offset % 2 else "male",
            "nationality": "Brazilian",
            "marital_status": ["single", "married", "divorced"][offset % 3],
            "document": cpf(person_id),
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
            "biometric_validated": c["fraude"] is None and offset % 5 != 0,
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

        if c["score"] is not None:
            score_por_pessoa[person_id] = add("credit_scores", {
                "person_id": person_id, "score": c["score"], "score_date": dia("2026-08-15"),
                "score_model": "serasa_v3",
                "score_reason": c["nota"],
                "payment_history": int(c["pontual"] * 40), "credit_usage": int((1 - c["utilizacao"]) * 25),
                "credit_age": 10 + person_id % 15, "credit_mix": 8 + person_id % 10,
                "recent_inquiries": max(0, 15 - c["consultas"] * 2),
                "risk_level": ("very_low" if c["score"] >= 800 else "low" if c["score"] >= 700
                               else "medium" if c["score"] >= 550 else "high" if c["score"] >= 400
                               else "very_high"),
                "default_probability": round(max(0.01, min(0.95, (1000 - c["score"]) / 1000)), 3),
            })

        perfil_por_pessoa[person_id] = add("financial_profiles", {
            "person_id": person_id, "profile_date": dia("2026-08-15"),
            "declared_monthly_income": c["renda"],
            "estimated_monthly_income": c["estimada"] or (
                None if "renda_estimada" in c["falta"] else (round(c["renda"] * 0.97) if c["renda"] else None)),
            "income_source": "salary" if c["vinculo"] and c["vinculo"][1] == CLT else "mixed",
            "total_assets": None, "real_estate_value": None, "vehicles_value": None,
            "total_liabilities": round((renda or 0) * c["dti"] * 12, 2),
            "total_monthly_payments": round((renda or 0) * c["dti"], 2),
            "debt_to_income_ratio": c["dti"],
            "available_credit": None,
            "credit_utilization": c["utilizacao"],
        })

        if c["vinculo"]:
            empregador, tipo, inicio, salario, verificado = c["vinculo"]
            add("employment_records", {
                "person_id": person_id, "employer_name": empregador,
                "employer_document": f"{person_id * 1234567890123 % 90000000000000 + 10000000000000}",
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
        for mes in range(12):
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
            add("debts", {
                "person_id": person_id, "debt_type": "loan", "creditor": banco[0],
                "creditor_document": banco[1], "original_amount": valor,
                "current_amount": round(valor * 1.08, 2), "interest_rate": 4.5, "fees": None,
                "origin_date": dia("2025-11-10"), "due_date": dia("2026-06-10"),
                "status": "overdue" if em_cobranca else "open",
                "in_collection": em_cobranca,
                "collection_date": dia("2026-07-01") if em_cobranca else None,
                "collection_agency": "Cobranca Sintetica" if em_cobranca else None,
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
            "assessment_type": "credit", "risk_score": c["score"] or 0,
            "risk_level": "low" if (c["score"] or 0) >= 700 else "medium" if (c["score"] or 0) >= 500 else "high",
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
        entrada, saida, volatilidade, dias_negativos, recorrente = c["fluxo"]
        banco = list(BANCOS.values())[person_id % 4]
        meses = c["banco"][0]

        perfil_por_pessoa[person_id] = add("bank_account_profiles", {
            "person_id": person_id, "profile_date": dia("2026-08-20"),
            "banking_relationships": 1 + person_id % 3,
            "account_age_average": meses,
            "has_checking_account": True,
            "has_savings_account": entrada > saida,
            "has_investment_account": c["faixa"] == "simples" and entrada - saida > 1500,
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

        add("cash_flow_analyses", {
            "person_id": person_id, "analysis_date": dia("2026-08-20"), "period_days": 90,
            "average_monthly_inflow": round(entrada, 2), "average_monthly_outflow": round(saida, 2),
            "net_cash_flow": round(entrada - saida, 2), "inflow_volatility": volatilidade,
            "negative_balance_days": dias_negativos, "has_recurring_income": recorrente,
        })

        recorrentes = [("income", "salary", "Salario", round(entrada * 0.85, 2), recorrente),
                       ("expense", "rent", "Aluguel residencial", round(saida * 0.35, 2), True),
                       ("expense", "utility", "Energia eletrica", round(saida * 0.08, 2), True)]
        if c["dti"] >= 0.4:
            recorrentes.append(("expense", "loan", "Parcela de emprestimo", round(saida * 0.22, 2), True))
        for tipo, categoria, descricao, valor, ativo in recorrentes:
            add("recurring_transactions", {
                "person_id": person_id, "transaction_type": tipo, "category": categoria,
                "description": descricao, "amount": valor, "frequency": "monthly",
                "counterparty": banco[0] if tipo == "expense" else (c["vinculo"][0] if c["vinculo"] else None),
                "first_detected_date": dia(HOJE - timedelta(days=270)),
                "last_occurrence_date": dia(HOJE - timedelta(days=12)),
                "is_active": ativo,
            })

        add("data_sharing_consents", {
            "person_id": person_id, "consent_id": f"urn:sintetico:consent:{person_id:04d}",
            "institution": banco[0], "status": "granted",
            "scope": json.dumps(["ACCOUNTS_READ", "ACCOUNTS_BALANCES_READ", "RESOURCES_READ"]),
            "granted_at": dia(HOJE - timedelta(days=45)),
            "expires_at": dia(HOJE + timedelta(days=135)), "revoked_at": None,
        })

    total = 0
    for tabela, (prev, _) in tabelas.items():
        total += grava("open-finance", tabela, prev, novos[tabela])
    print(f"open finance: {total} registros")
    return perfil_por_pessoa


# --------------------------------------------------------------------------- cadastro interno
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

        relacionamento_por_pessoa[person_id] = add("customer_relationships", {
            "person_id": person_id,
            "customer_since": dia(HOJE - timedelta(days=30 * meses)),
            "relationship_months": meses, "segment": segmento,
            "branch": f"{1000 + person_id}", "is_active": True,
            "churn_risk": "low" if meses >= 48 else "medium" if meses >= 18 else "high",
            "internal_score": min(1000, (c["score"] or 500) + (10 if meses >= 48 else 0)),
        })

        produto = add("contracted_products", {
            "person_id": person_id, "product_type": "checking_account",
            "product_name": "Conta Corrente", "contract_number": f"CC-{person_id:05d}",
            "contracted_date": dia(HOJE - timedelta(days=30 * meses)), "status": "active",
            "balance": round((renda or 2000) * 0.3, 2), "monthly_value": 29.9,
        })
        if c["dti"] >= 0.3:
            add("contracted_products", {
                "person_id": person_id, "product_type": "loan",
                "product_name": "Emprestimo Pessoal", "contract_number": f"LN-{person_id:05d}",
                "contracted_date": dia(HOJE - timedelta(days=400)), "status": "active",
                "balance": round((renda or 2000) * c["dti"] * 8, 2),
                "monthly_value": round((renda or 2000) * c["dti"] * 0.6, 2),
            })

        # os atrasos internos acompanham a pontualidade observada no birô
        atrasos = round((1 - c["pontual"]) * 12)
        for mes in range(12):
            vencimento = HOJE - timedelta(days=30 * (12 - mes))
            atrasado = mes < atrasos
            add("internal_payment_records", {
                "person_id": person_id, "contracted_product_id": produto,
                "reference_month": dia(vencimento.replace(day=1)), "due_date": dia(vencimento),
                "payment_date": dia(vencimento, 9 if atrasado else 0),
                "amount_due": 29.9, "amount_paid": 29.9,
                "status": "late" if atrasado else "paid", "days_late": 9 if atrasado else 0,
            })

        if limite:
            add("pre_approved_limits", {
                "person_id": person_id, "product_type": "personal_loan",
                "approved_amount": float(limite), "interest_rate": 3.2,
                "calculated_date": dia("2026-08-20"), "valid_until": dia(HOJE + timedelta(days=60)),
                "policy_version": "v1.2", "is_active": True,
            })

        if c["renda"]:
            add("income_declarations", {
                "person_id": person_id, "declaration_date": dia("2026-02-15"),
                "income_type": "salary" if c["vinculo"] and c["vinculo"][1] == CLT else "self_employed",
                "monthly_amount": float(c["renda"]), "yearly_amount": float(c["renda"] * 12),
                "source": "folha de pagamento" if c["vinculo"] else "declaracao do cliente",
                "verified": bool(c["vinculo"] and c["vinculo"][4]),
                "verified_by": "analista" if c["vinculo"] and c["vinculo"][4] else None,
                "proof_file_id": None,
            })

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

        add("document_validations", {
            "person_id": person_id, "validation_date": dia("2026-08-10"),
            "document_number": documento, "document_type": "cpf",
            "receita_federal_status": situacao,
            "is_valid": situacao == "regular",
            # o alerta de fraude do birô e a divergencia cadastral sao o mesmo evento
            "name_matches": "nome" not in c["divergencias"],
            "birth_date_matches": "nascimento" not in c["divergencias"],
            "biometric_validated": c["fraude"] is None and person_id % 5 != 0,
            "source": "receita_federal",
            "raw_response": json.dumps({"situacao": situacao, "consulta": dia("2026-08-10")},
                                       ensure_ascii=False),
        })

        pendencias = []
        if cnd != "regular":
            pendencias = ["divida_ativa_uniao", "cnd_negada"]
        add("fiscal_regularities", {
            "person_id": person_id, "check_date": dia("2026-08-10"),
            "has_debts": cnd != "regular", "cnd_status": cnd,
            "cnd_number": f"CND-2026-{person_id:07d}" if cnd == "regular" else None,
            "cnd_issue_date": dia("2026-08-10") if cnd == "regular" else None,
            "cnd_valid_until": dia("2026-11-08") if cnd == "regular" else None,
            "pending_issues": json.dumps(pendencias, ensure_ascii=False) if pendencias else None,
        })

        if c["vinculo"]:
            empregador, tipo, inicio, _, verificado = c["vinculo"]
            add("employment_link_validations", {
                "person_id": person_id, "validation_date": dia("2026-08-10"),
                "employer_name": empregador,
                "employer_document": f"{person_id * 1234567890123 % 90000000000000 + 10000000000000}",
                "employment_type": tipo, "status": "active", "start_date": dia(inicio),
                "end_date": None, "source": "eSocial", "verified": verificado,
            })

        status = "irregular" if situacao != "regular" or sancoes else "attention" if pep else "regular"
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

    faixas = {}
    for c in CENARIOS.values():
        faixas[c["faixa"]] = faixas.get(c["faixa"], 0) + 1
    print("cenarios por faixa:", faixas)
    print("rode agora: python3 generate_cadastro.py")


if __name__ == "__main__":
    main()
