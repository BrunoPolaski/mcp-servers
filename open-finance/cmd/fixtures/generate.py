#!/usr/bin/env python3
"""Gera as fixtures de Open Finance a partir dos perfis definidos abaixo.

Cada perfil é derivado do quadro de crédito do mesmo cliente no birô, de modo
que as duas fontes contem uma história coerente. Rode a partir de
cmd/fixtures/: python3 generate.py
"""

import json
from datetime import datetime, timedelta
from pathlib import Path

OUT = Path(__file__).parent / "fixtures"
TS = "2026-08-31T00:00:00Z"

INSTITUTIONS = {
    "alfa": ("Banco Sintético Alfa", "11222333000181"),
    "beta": ("Banco Sintético Beta", "22333444000172"),
    "gama": ("Banco Sintético Gama", "33444555000163"),
    "delta": ("Banco Sintético Delta", "44555666000154"),
}

# person_id: perfil. "renda" é a renda declarada no birô, usada só para
# documentar a coerência entre as fontes; não vai para nenhuma fixture.
PROFILES = {
    1: {  # Felipe Pereira Santos — score 254, DTI 1,10, 0/7 pagamentos em dia
        "renda": 3690,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 2, "age": 34, "chk": True, "sav": False, "inv": False, "value": None},
        "cash": {"in": 3550.00, "out": 4300.00, "vol": 0.38, "neg_days": 22, "recurring": False},
        "accounts": [("alfa", "checking", 210.00)],
        "recurring": [
            ("expense", "rent", "Aluguel residencial", 1200.00, "Imobiliária Sintética", True),
            ("expense", "utility", "Energia elétrica", 320.00, "Concessionária Sintética", True),
            ("income", "salary", "Salário", 1750.00, "Empregador Sintético", False),
        ],
    },
    2: {  # Henrique Martins Barbosa — score 726, DTI 0,31, 6/6 em dia
        "renda": 7437,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 3, "age": 71, "chk": True, "sav": True, "inv": True, "value": 18500.00},
        "cash": {"in": 7600.00, "out": 6100.00, "vol": 0.08, "neg_days": 0, "recurring": True},
        "accounts": [("beta", "checking", 4200.00), ("alfa", "savings", 9800.00)],
        "recurring": [
            ("income", "salary", "Salário", 6950.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 1900.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Serviço de streaming", 89.00, "Streaming Sintético", True),
        ],
    },
    3: {  # Fernanda Costa Barbosa — score 448, 2 negativações ativas
        "renda": 2477,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 2, "age": 22, "chk": True, "sav": False, "inv": False, "value": None},
        "cash": {"in": 2380.00, "out": 2610.00, "vol": 0.31, "neg_days": 14, "recurring": True},
        "accounts": [("alfa", "checking", 640.00)],
        "recurring": [
            ("income", "salary", "Salário", 2150.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 950.00, "Imobiliária Sintética", True),
            ("expense", "utility", "Água e esgoto", 210.00, "Concessionária Sintética", True),
        ],
    },
    4: {  # Fernanda Rodrigues Pereira — score 349, alerta de fraude
        "renda": 3836,
        "consent": {"bank": "gama", "status": "granted"},
        "profile": {"rel": 3, "age": 15, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 4100.00, "out": 3980.00, "vol": 0.55, "neg_days": 9, "recurring": False},
        "accounts": [("gama", "checking", 880.00), ("beta", "savings", 1500.00)],
        "recurring": [
            ("expense", "rent", "Aluguel residencial", 1100.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de telefonia", 49.00, "Telecom Sintética", True),
        ],
    },
    5: {  # Gabriela Ribeiro Barbosa — score 809 com fluxo apertado: fontes discordam
        "renda": 6915,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 4, "age": 88, "chk": True, "sav": True, "inv": True, "value": 42000.00},
        "cash": {"in": 7100.00, "out": 6980.00, "vol": 0.19, "neg_days": 7, "recurring": True},
        "accounts": [("beta", "checking", 3100.00), ("gama", "savings", 5400.00)],
        "recurring": [
            ("income", "salary", "Salário", 6400.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2400.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento de veículo", 1850.00, "Banco Sintético Beta", True),
            ("expense", "subscription", "Serviço de streaming", 89.00, "Streaming Sintético", True),
        ],
    },
    6: {  # Henrique Almeida Ribeiro — renda indeterminável e consentimento revogado
        "renda": None,
        "consent": {"bank": "alfa", "status": "revoked"},
        "profile": None,
        "cash": None,
        "accounts": [],
        "recurring": [],
    },
    7: {  # Igor Souza Martins — score 931, DTI 0,21
        "renda": 16649,
        "consent": {"bank": "delta", "status": "granted"},
        "profile": {"rel": 4, "age": 132, "chk": True, "sav": True, "inv": True, "value": 310000.00},
        "cash": {"in": 17200.00, "out": 12100.00, "vol": 0.06, "neg_days": 0, "recurring": True},
        "accounts": [("delta", "checking", 12400.00), ("beta", "savings", 48000.00)],
        "recurring": [
            ("income", "salary", "Salário", 15800.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 3800.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de saúde", 129.00, "Operadora Sintética", True),
        ],
    },
    8: {  # Lucas Martins Souza — renda declarada 12.043 contra 4.200 detectados
        "renda": 12043,
        "consent": {"bank": "gama", "status": "granted"},
        "profile": {"rel": 2, "age": 41, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 5900.00, "out": 5750.00, "vol": 0.44, "neg_days": 11, "recurring": True},
        "accounts": [("gama", "checking", 1900.00)],
        "recurring": [
            ("income", "salary", "Salário", 4200.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 1750.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento imobiliário", 2100.00, "Banco Sintético Gama", True),
        ],
    },
    9: {  # Eduardo Barbosa Almeida — score 663, DTI 0,39
        "renda": 10529,
        "consent": {"bank": "beta", "status": "granted"},
        "profile": {"rel": 3, "age": 57, "chk": True, "sav": True, "inv": False, "value": None},
        "cash": {"in": 10800.00, "out": 9500.00, "vol": 0.14, "neg_days": 3, "recurring": True},
        "accounts": [("beta", "checking", 5200.00), ("alfa", "savings", 7300.00)],
        "recurring": [
            ("income", "salary", "Salário", 9600.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2600.00, "Imobiliária Sintética", True),
            ("expense", "financing", "Financiamento de veículo", 1500.00, "Banco Sintético Beta", True),
        ],
    },
    10: {  # Eduardo Ribeiro Ribeiro — score 800, DTI 0,26
        "renda": 7958,
        "consent": {"bank": "alfa", "status": "granted"},
        "profile": {"rel": 3, "age": 96, "chk": True, "sav": True, "inv": True, "value": 61000.00},
        "cash": {"in": 8300.00, "out": 6600.00, "vol": 0.09, "neg_days": 0, "recurring": True},
        "accounts": [("alfa", "checking", 6100.00), ("delta", "savings", 22000.00)],
        "recurring": [
            ("income", "salary", "Salário", 7500.00, "Empregador Sintético", True),
            ("expense", "rent", "Aluguel residencial", 2100.00, "Imobiliária Sintética", True),
            ("expense", "subscription", "Plano de telefonia", 45.00, "Telecom Sintética", True),
        ],
    },
}

# Os três meses cobertos pelos extratos, encerrando em 2026-08-31.
MONTHS = [
    ("2026-06-01T00:00:00Z", "2026-06-30T00:00:00Z"),
    ("2026-07-01T00:00:00Z", "2026-07-31T00:00:00Z"),
    ("2026-08-01T00:00:00Z", "2026-08-31T00:00:00Z"),
]


def base(row_id):
    return {"id": row_id, "created_at": TS, "updated_at": TS, "deleted_at": None}


def build():
    profiles, statements, analyses, recurrences, consents = [], [], [], [], []
    profile_by_person = {}
    next_id = {"profile": 1, "stmt": 1, "cash": 1, "rec": 1, "consent": 1}

    for person_id in sorted(PROFILES):
        p = PROFILES[person_id]

        if p["profile"]:
            pid = next_id["profile"]
            next_id["profile"] += 1
            profile_by_person[person_id] = pid
            profiles.append({
                **base(pid),
                "person_id": person_id,
                "profile_date": "2026-08-31T00:00:00Z",
                "banking_relationships": p["profile"]["rel"],
                "account_age_average": p["profile"]["age"],
                "has_checking_account": p["profile"]["chk"],
                "has_savings_account": p["profile"]["sav"],
                "has_investment_account": p["profile"]["inv"],
                "investments_value": p["profile"]["value"],
            })

        if p["cash"]:
            c = p["cash"]
            net = round(c["in"] - c["out"], 2)
            cid = next_id["cash"]
            next_id["cash"] += 1
            analyses.append({
                **base(cid),
                "person_id": person_id,
                "analysis_date": "2026-08-31T00:00:00Z",
                "period_days": 90,
                "average_monthly_inflow": c["in"],
                "average_monthly_outflow": c["out"],
                "net_cash_flow": net,
                "inflow_volatility": c["vol"],
                "negative_balance_days": c["neg_days"],
                "has_recurring_income": c["recurring"],
            })

        for bank, account_type, opening in p["accounts"]:
            name, cnpj = INSTITUTIONS[bank]
            balance = opening
            for period_start, period_end in MONTHS:
                if account_type == "checking":
                    credits, debits = p["cash"]["in"], p["cash"]["out"]
                    count = 62
                else:
                    # Poupança: movimentação enxuta, com aporte mensal modesto.
                    credits, debits = 400.00, 200.00
                    count = 4
                closing = round(balance + credits - debits, 2)
                sid = next_id["stmt"]
                next_id["stmt"] += 1
                statements.append({
                    **base(sid),
                    "person_id": person_id,
                    "institution": name,
                    "institution_document": cnpj,
                    "account_type": account_type,
                    "period_start": period_start,
                    "period_end": period_end,
                    "opening_balance": balance,
                    "closing_balance": closing,
                    "total_credits": credits,
                    "total_debits": debits,
                    "transaction_count": count,
                    "currency": "BRL",
                })
                balance = closing

        for kind, category, description, amount, counterparty, active in p["recurring"]:
            rid = next_id["rec"]
            next_id["rec"] += 1
            recurrences.append({
                **base(rid),
                "person_id": person_id,
                "transaction_type": kind,
                "category": category,
                "description": description,
                "amount": amount,
                "frequency": "monthly",
                "counterparty": counterparty,
                "first_detected_date": "2025-09-05T00:00:00Z",
                "last_occurrence_date": "2026-08-05T00:00:00Z" if active else "2026-03-05T00:00:00Z",
                "is_active": active,
            })

        c = p["consent"]
        name, _ = INSTITUTIONS[c["bank"]]
        cid = next_id["consent"]
        next_id["consent"] += 1
        revoked = c["status"] == "revoked"
        consents.append({
            **base(cid),
            "person_id": person_id,
            "consent_id": f"urn:openfinance:consent:{person_id:03d}",
            "institution": name,
            "status": c["status"],
            "scope": json.dumps(["ACCOUNTS_READ", "ACCOUNTS_BALANCES_READ", "RESOURCES_READ"]),
            "granted_at": "2026-01-15T00:00:00Z" if revoked else "2026-05-10T00:00:00Z",
            "expires_at": None if revoked else "2027-05-10T00:00:00Z",
            "revoked_at": "2026-04-02T00:00:00Z" if revoked else None,
        })

    return profile_by_person, {
        "bank_account_profiles.json": profiles,
        "bank_statements.json": statements,
        "cash_flow_analyses.json": analyses,
        "recurring_transactions.json": recurrences,
        "data_sharing_consents.json": consents,
    }


def rewrite_persons(profile_by_person):
    """Reescreve persons.json trocando os vínculos do birô pelo perfil bancário."""
    path = OUT / "persons.json"
    rows = json.loads(path.read_text())
    for row in rows:
        row.pop("credit_score_id", None)
        row.pop("financial_profile_id", None)
        row["bank_account_profile_id"] = profile_by_person.get(row["id"])
    path.write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")


def main():
    profile_by_person, files = build()
    for name, rows in files.items():
        (OUT / name).write_text(json.dumps(rows, indent=2, ensure_ascii=False) + "\n")
        print(f"{name}: {len(rows)} registros")
    rewrite_persons(profile_by_person)
    print("persons.json: vínculos reescritos")


if __name__ == "__main__":
    main()
