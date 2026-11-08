CREATE TABLE public.tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_type TEXT NOT NULL,
    parameters JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    result JSONB DEFAULT NULL,
    error TEXT DEFAULT NULL,
    attempts INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ DEFAULT NULL,
    finished_at TIMESTAMPTZ DEFAULT NULL,
    lease_until TIMESTAMPTZ DEFAULT NULL,
    CONSTRAINT tasks_status_chek CHECK (status IN ('pending','running','completed','failed')),
    CONSTRAINT attempts_not_below_zero CHECK (attempts >= 0)
)