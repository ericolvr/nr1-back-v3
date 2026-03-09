-- Migration: Add activity_attachments table
-- Description: Adiciona tabela para armazenar anexos/evidências de atividades com metadados

CREATE TABLE IF NOT EXISTS activity_attachments (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES action_plan_activities(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Índices para performance
    CONSTRAINT fk_activity_attachments_activity FOREIGN KEY (activity_id) REFERENCES action_plan_activities(id) ON DELETE CASCADE
);

CREATE INDEX idx_activity_attachments_activity_id ON activity_attachments(activity_id);
CREATE INDEX idx_activity_attachments_created_at ON activity_attachments(created_at);

-- Comentários
COMMENT ON TABLE activity_attachments IS 'Armazena anexos/evidências de atividades de planos de ação';
COMMENT ON COLUMN activity_attachments.activity_id IS 'ID da atividade relacionada';
COMMENT ON COLUMN activity_attachments.file_name IS 'Nome original do arquivo';
COMMENT ON COLUMN activity_attachments.file_url IS 'URL do arquivo armazenado';
COMMENT ON COLUMN activity_attachments.file_size IS 'Tamanho do arquivo em bytes';
COMMENT ON COLUMN activity_attachments.mime_type IS 'Tipo MIME do arquivo (image/jpeg, application/pdf, etc)';
