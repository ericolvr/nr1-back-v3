-- Atualizar template para usar variáveis
UPDATE action_plan_templates 
SET 
  title_template = 'Reduzir {category} - Departamento: {department_name}',
  description_template = 'Detectado risco {risk_level} na categoria {category}. Score médio: {average_score} em {question_count} perguntas. Ação necessária!',
  updated_at = NOW()
WHERE id = 2;

-- Verificar atualização
SELECT id, category, title_template, description_template, min_risk_level, auto_create, active 
FROM action_plan_templates 
WHERE id = 2;
