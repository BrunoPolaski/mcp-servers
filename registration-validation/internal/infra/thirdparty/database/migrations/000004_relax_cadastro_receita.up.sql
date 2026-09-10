-- A consulta de CPF na Receita Federal devolve nome civil, filiacao, nascimento
-- e situacao cadastral: nao devolve nacionalidade nem telefone. As colunas
-- herdadas do schema do biro precisam aceitar ausencia desses dados.
ALTER TABLE personal_informations ALTER COLUMN nationality DROP NOT NULL;
ALTER TABLE personal_informations ALTER COLUMN primary_phone DROP NOT NULL;
