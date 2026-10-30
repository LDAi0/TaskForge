CREATE TABLE public.tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_type TEXT NOT NULL,
    parametres JSONB,
    status TEXT,
    result JSONB DEFAULT NULL,
    error TEXT,
    attempts INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    CONSTRAINT tasks_status_chek CHECK (status IN ('pending','running','finished'))
)