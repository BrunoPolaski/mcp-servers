#!/usr/bin/env python3
"""Gera as fixtures de validacao cadastral a partir dos perfis definidos abaixo.

Os dez clientes sao os mesmos do biro, do registro interno e do open finance
(mesmos CPFs e mesmos ids), e cada perfil aqui e derivado do quadro de credito
do cliente no biro, de modo que as fontes contem uma historia coerente: quem
tem negativacao e processo aparece com CND irregular, quem tem alerta de fraude
aparece com divergencia cadastral, e quem nao tem vinculo no biro tambem nao
tem vinculo no eSocial. Rode a partir de cmd/fixtures/: python3 generate.py
"""

import json
from datetime import datetime, timedelta
from pathlib import Path

OUT = Path(__file__).parent / "fixtures"
TS = "2026-08-31T00:00:00Z"

# person_id: perfil. "biro" resume o que sustenta cada decisao e nao vai para
# nenhuma fixture; serve so para documentar a coerencia entre as fontes.
PROFILES = {
    1: {  # Felipe Pereira Santos - score 254, 1 processo, CLT ativo
        "biro": "score 254, 0 negativacoes ativas, 1 processo",
        "documents": [
            ("cpf", "83609091800", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", True, "irregular", None,
                   ["irpf_2024_pendente", "divida_ativa_uniao"]),
        "jobs": [("Comercio ABC", "44312559206332", "CLT", "active", "2022-01-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    2: {  # Henrique Martins Barbosa - score 726, historico limpo, MEI
        "biro": "score 726, sem negativacoes, sem processos",
        "documents": [
            ("cpf", "15393772025", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", False, "regular", "CND-2026-0000002", []),
        "jobs": [("MEI Servicos", "06341759988101", "autonomous", "active", "2022-12-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    3: {  # Fernanda Costa Barbosa - score 448, 2 negativacoes ativas, 1 processo
        "biro": "score 448, 2 negativacoes ativas, 1 processo",
        "documents": [
            ("cpf", "45516693934", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", True, "irregular", None,
                   ["parcelamento_simples_nacional", "cnd_negada"]),
        "jobs": [("Consultoria Delta", "75113003027845", "autonomous", "active", "2021-05-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    4: {  # Fernanda Rodrigues Pereira - score 349 e alerta de fraude no biro
        "biro": "score 349, 1 alerta de fraude",
        "documents": [
            # Divergencia de nome e biometria reprovada: e o mesmo evento que
            # gerou o alerta de fraude no biro.
            ("cpf", "26959892961", "2026-02-10", "pending", False, False, True, False),
        ],
        "fiscal": ("2026-08-01", True, "irregular", None, ["irpf_2024_pendente"]),
        "jobs": [("Comercio ABC", "05496200523462", "autonomous", "active", "2023-10-01", None, False)],
        "compliance": [("kyc_full", "2026-08-01", "review", False, None, False, None,
                        {"motivo": "divergencia_cadastral", "origem": "biro_fraud_alert"})],
    },
    5: {  # Gabriela Ribeiro Barbosa - score 809, socia PJ, 1 alerta de fraude
        "biro": "score 809, 1 negativacao, 1 alerta de fraude, vinculo PJ",
        "documents": [
            ("cpf", "35509139404", "2026-02-10", "regular", True, True, True, False),
            ("cnpj", "69603521830167", "2026-02-12", "regular", True, True, False, False),
        ],
        "fiscal": ("2026-08-01", False, "regular", "CND-2026-0000005", []),
        "jobs": [("Comercio ABC", "69603521830167", "PJ", "active", "2023-06-01", None, True)],
        "compliance": [
            ("kyc_full", "2026-08-01", "attention", True,
             "Socia de empresa com contrato vigente com orgao publico", False, None, None),
            ("pep_screening", "2026-08-01", "attention", True,
             "Enquadrada como PEP por relacionamento societario", False, None, None),
        ],
    },
    6: {  # Henrique Almeida Ribeiro - score 294, sem vinculo empregaticio
        "biro": "score 294, 2 negativacoes ativas, 2 processos, sem vinculo",
        "documents": [
            # Nascimento diverge das bases privadas (dia e mes trocados), por
            # isso birth_date_matches=false.
            ("cpf", "07432970084", "2026-02-10", "suspended", False, True, False, False),
        ],
        "fiscal": ("2026-08-01", True, "irregular", None,
                   ["divida_ativa_uniao", "cnd_negada", "processos_judiciais"]),
        "jobs": [],  # nenhum vinculo encontrado no eSocial, como no biro
        "compliance": [
            ("kyc_full", "2026-08-01", "irregular", False, None, True,
             "Homonimo em lista restritiva; requer analise manual",
             {"motivo": "cpf_suspenso_receita"}),
            ("sanctions_screening", "2026-08-01", "irregular", False, None, True,
             "Homonimo em lista restritiva; requer analise manual", None),
        ],
    },
    7: {  # Igor Souza Martins - score 931, historico limpo
        "biro": "score 931, sem negativacoes, sem processos",
        "documents": [
            ("cpf", "40161087990", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", False, "regular", "CND-2026-0000007", []),
        "jobs": [("Industria XYZ", "85131505577925", "autonomous", "active", "2019-11-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    8: {  # Lucas Martins Souza - score 461, 1 negativacao, 1 processo
        "biro": "score 461, 1 negativacao ativa, 1 processo",
        "documents": [
            ("cpf", "18979021232", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", True, "irregular", None, ["divida_ativa_uniao"]),
        "jobs": [("Comercio ABC", "28652418388745", "CLT", "active", "2022-08-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    9: {  # Eduardo Barbosa Almeida - score 663, dois vinculos no biro
        "biro": "score 663, 1 negativacao, vinculo anterior encerrado em 2022",
        "documents": [
            ("cpf", "08458651459", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", False, "regular", "CND-2026-0000009", []),
        "jobs": [
            ("Tech Solutions Ltda", "55905068740903", "autonomous", "terminated",
             "2021-01-01", "2022-09-01", True),
            ("Comercio ABC", "92131479763109", "CLT", "active", "2023-06-01", None, True),
        ],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
    10: {  # Eduardo Ribeiro Ribeiro - score 800, CLT desde 2018
        "biro": "score 800, sem negativacoes, CLT desde 2018",
        "documents": [
            ("cpf", "76009108535", "2026-02-10", "regular", True, True, True, True),
        ],
        "fiscal": ("2026-08-01", False, "regular", "CND-2026-0000010", []),
        "jobs": [("Comercio ABC", "86793706428521", "CLT", "active", "2018-11-01", None, True)],
        "compliance": [("kyc_full", "2026-08-01", "regular", False, None, False, None, None)],
    },
}


def day(date, plus=0):
    """'2026-08-01' -> '2026-08-01T00:00:00Z', opcionalmente somando dias."""
    if date is None:
        return None
    return (datetime.strptime(date, "%Y-%m-%d") + timedelta(days=plus)).strftime("%Y-%m-%dT00:00:00Z")


def base(row_id):
    return {"id": row_id, "created_at": TS, "updated_at": TS, "deleted_at": None}


def build():
    documents, fiscais, jobs, checks = [], [], [], []

    for person_id, profile in PROFILES.items():
        for doc_type, number, date, status, valid, name_ok, birth_ok, bio in profile["documents"]:
            documents.append({
                **base(len(documents) + 1),
                "person_id": person_id,
                "validation_date": day(date),
                "document_number": number,
                "document_type": doc_type,
                "receita_federal_status": status,
                "is_valid": valid,
                "name_matches": name_ok,
                "birth_date_matches": birth_ok,
                "biometric_validated": bio,
                "source": "receita_federal",
                "raw_response": json.dumps(
                    {"situacao": status, "consulta": day(date)}, ensure_ascii=False),
            })

        date, has_debts, cnd_status, cnd_number, pending = profile["fiscal"]
        fiscais.append({
            **base(len(fiscais) + 1),
            "person_id": person_id,
            "check_date": day(date),
            "has_debts": has_debts,
            "cnd_status": cnd_status,
            "cnd_number": cnd_number,
            # Sem CND emitida nao ha data de emissao nem validade.
            "cnd_issue_date": day(date) if cnd_number else None,
            "cnd_valid_until": day(date, 90) if cnd_number else None,
            "pending_issues": json.dumps(pending, ensure_ascii=False) if pending else None,
        })

        for employer, document, kind, status, start, end, verified in profile["jobs"]:
            jobs.append({
                **base(len(jobs) + 1),
                "person_id": person_id,
                "validation_date": day("2026-08-01"),
                "employer_name": employer,
                "employer_document": document,
                "employment_type": kind,
                "status": status,
                "start_date": day(start),
                "end_date": day(end),
                "source": "eSocial",
                "verified": verified,
            })

        for kind, date, status, is_pep, pep, sanctions, sanction_details, details in profile["compliance"]:
            checks.append({
                **base(len(checks) + 1),
                "person_id": person_id,
                "check_type": kind,
                "check_date": day(date),
                "status": status,
                "details": json.dumps(details, ensure_ascii=False) if details else None,
                "is_pep": is_pep,
                "pep_details": pep,
                "on_sanctions_list": sanctions,
                "sanctions_details": sanction_details,
                "valid_until": day(date, 180),
            })

    return {
        "document_validations.json": documents,
        "fiscal_regularities.json": fiscais,
        "employment_link_validations.json": jobs,
        "compliance_checks.json": checks,
    }


def rewrite_persons():
    """Limpa de persons.json as colunas que so existem nas outras fontes."""
    path = OUT / "persons.json"
    rows = json.loads(path.read_text())
    for row in rows:
        # Colunas de outras fontes que nao existem no schema deste servico.
        for field in ("customer_relationship_id", "credit_score_id",
                      "financial_profile_id", "consent_status", "consent_granted_at"):
            row.pop(field, None)
    path.write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")


def main():
    for name, rows in build().items():
        (OUT / name).write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")
        print(f"{name}: {len(rows)} registros")
    rewrite_persons()
    print("persons.json: colunas de outras fontes removidas")


if __name__ == "__main__":
    main()
