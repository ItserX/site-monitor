CREATE TABLE IF NOT EXISTS sites (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    active BOOLEAN NOT NULL DEFAULT true
);

INSERT INTO sites (id, url, active) VALUES

('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://yandex.ru', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'https://mail.ru', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'https://vk.com', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'https://avito.ru', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'https://wildberries.ru', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'https://google.com', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'https://github.com', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', 'https://stackoveqweqweqweeerflow.com', true),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21', 'https://stackoeqewqeqeverflow.com', true)
ON CONFLICT (url) DO NOTHING;









