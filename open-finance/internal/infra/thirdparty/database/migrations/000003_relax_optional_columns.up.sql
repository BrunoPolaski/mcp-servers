-- As fixtures do Open Finance simulam dados incompletos vindos das instituicoes
-- (endereco sem rua/CEP/data de mudanca, cadastro sem telefone). O bureau, de
-- onde vem esses mesmos registros, ja trata essas colunas como opcionais.
ALTER TABLE addresses ALTER COLUMN street DROP NOT NULL;
ALTER TABLE addresses ALTER COLUMN zip_code DROP NOT NULL;
ALTER TABLE addresses ALTER COLUMN moved_in_date DROP NOT NULL;
ALTER TABLE personal_informations ALTER COLUMN primary_phone DROP NOT NULL;
