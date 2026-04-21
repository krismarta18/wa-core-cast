-- Add settings table for dynamic branding and config
CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed initial settings
INSERT INTO settings (key, value, description) VALUES 
('app_name', 'WACAST', 'The name of the application displayed in the dashboard'),
('app_logo', '', 'URL or Base64 of the application logo'),
('footer_text', 'WACAST Core v1.0.0', 'Text displayed in the footer of the sidebar')
ON CONFLICT (key) DO NOTHING;
