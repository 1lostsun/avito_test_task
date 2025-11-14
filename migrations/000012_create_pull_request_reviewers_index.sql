-- Индекс для поиска PR по ревьюверу
CREATE INDEX idx_pr_reviewers_user_id ON pull_request_reviewers(user_id);