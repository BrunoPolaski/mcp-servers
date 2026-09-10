#!/usr/bin/env python3
"""Gera o bloco cadastral de cada servidor MCP a partir da visao do banco.

Em sistemas reais as quatro fontes nao guardam o mesmo cadastro: o birô monta o
dele com o que os credores informam, e por isso atrasa; o banco tem o cadastro
mais atual, porque e ele quem faz o KYC; a outra instituicao, via Open Finance,
compartilha os dados cadastrais que ela mesma mantem, sem RG, titulo ou CTPS; e
a Receita Federal devolve apenas nome civil, filiacao, nascimento e situacao do
CPF - nao devolve endereco, telefone nem e-mail.

E dessa diferenca que vive o servidor de validacao cadastral: sem divergencia
entre as bases, `name_matches` e `birth_date_matches` nunca teriam o que apontar.
As divergencias plantadas (DIVERGENCIAS) batem com o que o servidor de validacao
cadastral afirma sobre as mesmas pessoas.

O cadastro canonico e o do internal-registry, escrito por generate_cenarios.py.
As divergencias nao sao declaradas aqui: sao lidas do que o servidor de validacao
cadastral afirma sobre cada pessoa (name_matches, birth_date_matches), de modo
que as duas pontas nao tem como discordar. Rode este script depois daquele,
sempre a partir de mcp-servers/: python3 generate_cadastro.py
"""

import json
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).parent
BANCO = "internal-registry"  # cadastro canonico: o mais atual das quatro bases


def fixtures(service, name):
    return ROOT / service / "cmd" / "fixtures" / "fixtures" / f"{name}.json"


def read(service, name):
    path = fixtures(service, name)
    return json.loads(path.read_text()) if path.exists() else []


def write(service, name, rows):
    fixtures(service, name).write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")
    print(f"{service}/{name}.json: {len(rows)} registros")


# Campos que a Receita Federal devolve na consulta de CPF. O resto do cadastro
# nao existe nessa fonte.
RECEITA_FIELDS = {"id", "created_at", "updated_at", "deleted_at",
                  "full_name", "mother_name", "birth_date", "document"}

# Colunas que so existem no schema do banco.
SO_BANCO = ("document_validated", "biometric_validated", "receita_federal_status")

# Campos que o Open Finance nao compartilha (a API de dados cadastrais devolve
# nome, filiacao, nascimento, estado civil, documento, enderecos e contatos).
SEM_OPEN_FINANCE = ("rg", "rg_issuer", "rg_issue_date", "voter_id", "work_card")

def divergencias():
    """Le a validacao cadastral e monta a visao divergente da Receita.

    Quando ela diz que o nome nao confere, o nome civil na Receita tem um
    sobrenome que as bases privadas nao registraram; quando diz que o nascimento
    nao confere, as bases privadas trocaram dia e mes. Sao as duas divergencias
    que existem no mundo real e que essas flags servem para apontar.
    """
    base_pessoas = {row["id"]: row for row in read(BANCO, "personal_informations")}
    plantadas = {}
    for validacao in read("registration-validation", "document_validations"):
        pessoa = base_pessoas.get(validacao["person_id"])
        if pessoa is None:
            continue
        divergente = {}
        if not validacao["name_matches"]:
            divergente["full_name"] = f"{pessoa['full_name']} Neto"
        if not validacao["birth_date_matches"]:
            nascimento = datetime.strptime(pessoa["birth_date"][:10], "%Y-%m-%d")
            if nascimento.day > 12 or nascimento.day == nascimento.month:
                raise SystemExit(
                    f"pessoa {pessoa['id']}: nascida em {nascimento:%d/%m}, trocar dia e mes nao "
                    "produz divergencia; marque-a em alguem nascido ate o dia 12 e com dia != mes")
            divergente["birth_date"] = nascimento.replace(
                day=nascimento.month, month=nascimento.day).strftime("%Y-%m-%dT00:00:00Z")
        if divergente:
            plantadas[pessoa["id"]] = divergente
    return plantadas


# Telefone que o birô ainda tem em cadastro, anterior ao que o cliente informou
# ao banco.
TELEFONE_ANTIGO = {4: "81988001122", 6: "81997654321"}

# Quem tem mais de um endereco no cadastro do banco mudou de casa; o birô so
# conhece o anterior, e ainda o trata como atual.


def lacunas_cadastro_biro(person_id):
    """Campos ausentes no cadastro do birô.

    A regra e deterministica e existe para exercitar o tratamento de dado
    ausente da politica de credito.
    """
    faltando = []
    if person_id % 4 == 0:
        faltando.append("primary_phone")
    if person_id % 5 == 0:
        faltando.append("email")
    if person_id % 7 == 0:
        faltando += ["rg", "rg_issuer", "rg_issue_date"]
    if person_id % 6 == 0:
        faltando.append("voter_id")
    return faltando


def lacunas_endereco_biro(address_id):
    faltando = []
    if address_id % 9 == 0:
        faltando.append("zip_code")
    if address_id % 8 == 0:
        faltando.append("neighborhood")
    if address_id % 11 == 0:
        faltando.append("moved_in_date")
    return faltando


def cadastro_biro(base):
    """Cadastro do birô: defasado, e com lacunas de preenchimento."""
    rows = []
    for row in base:
        novo = dict(row)
        for field in SO_BANCO:
            novo.pop(field, None)
        for field in lacunas_cadastro_biro(row["id"]):
            novo[field] = None
        if row["id"] in TELEFONE_ANTIGO:
            novo["primary_phone"] = TELEFONE_ANTIGO[row["id"]]
        rows.append(novo)
    return rows


def cadastro_open_finance(base):
    """Dados cadastrais compartilhados pela outra instituicao."""
    rows = []
    for row in base:
        novo = dict(row)
        for field in SEM_OPEN_FINANCE:
            novo[field] = None
        for field in SO_BANCO:
            novo.pop(field, None)
        rows.append(novo)
    return rows


def cadastro_receita(base, plantadas):
    """Consulta de CPF: nome civil, filiacao, nascimento e nada mais."""
    rows = []
    for row in base:
        novo = {k: (v if k in RECEITA_FIELDS else None) for k, v in row.items()}
        for field in SO_BANCO:
            novo.pop(field, None)
        for field in ("email_verified", "phone_verified"):
            novo[field] = False
        novo.update(plantadas.get(row["id"], {}))
        rows.append(novo)
    return rows


def enderecos(base_enderecos, base_vinculos, visao):
    """Endereco por fonte. O birô nao foi informado das mudancas."""
    por_pessoa = {}
    for vinculo in base_vinculos:
        por_pessoa.setdefault(vinculo["personal_information_id"], []).append(vinculo["address_id"])
    # de quem mudou, o birô so conhece o endereco mais antigo
    desconhecidos = {ids[-1] for ids in por_pessoa.values() if len(ids) > 1}

    rows = []
    for row in base_enderecos:
        if visao == "bureau" and row["id"] in desconhecidos:
            continue
        novo = dict(row)
        if visao == "bureau":
            if not novo["is_current"]:
                novo.update(is_current=True, moved_out_date=None)
            for field in lacunas_endereco_biro(row["id"]):
                novo[field] = None
        rows.append(novo)

    vinculos = [v for v in base_vinculos
                if not (visao == "bureau" and v["address_id"] in desconhecidos)]
    return rows, vinculos


def main():
    base = read(BANCO, "personal_informations")
    base_enderecos = read(BANCO, "addresses")
    base_vinculos = read(BANCO, "person_addresses")
    plantadas = divergencias()

    write("bureau", "personal_informations", cadastro_biro(base))
    write("open-finance", "personal_informations", cadastro_open_finance(base))
    write("registration-validation", "personal_informations", cadastro_receita(base, plantadas))

    for visao in ("bureau", "open-finance"):
        rows, vinculos = enderecos(base_enderecos, base_vinculos, visao)
        write(visao, "addresses", rows)
        write(visao, "person_addresses", vinculos)

    # A Receita nao devolve endereco, nem documentos digitalizados.
    for name in ("addresses", "person_addresses", "person_documents", "files"):
        write("registration-validation", name, [])


if __name__ == "__main__":
    main()
