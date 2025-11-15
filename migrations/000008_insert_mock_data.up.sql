INSERT INTO teams (team_name) VALUES
    ('backend'),
    ('frontend'),
    ('devops');

INSERT INTO users (user_id, username, team_id, is_active) VALUES
    ('user1', 'Alice', (SELECT id FROM teams WHERE team_name = 'backend'), true),
    ('user2', 'Bob', (SELECT id FROM teams WHERE team_name = 'backend'), true),
    ('user3', 'Charlie', (SELECT id FROM teams WHERE team_name = 'backend'), true),
    ('user4', 'David', (SELECT id FROM teams WHERE team_name = 'backend'), false),
    ('user5', 'Eve', (SELECT id FROM teams WHERE team_name = 'frontend'), true),
    ('user6', 'Frank', (SELECT id FROM teams WHERE team_name = 'frontend'), true),
    ('user7', 'Grace', (SELECT id FROM teams WHERE team_name = 'frontend'), true),
    ('user8', 'Henry', (SELECT id FROM teams WHERE team_name = 'devops'), true),
    ('user9', 'Ivy', (SELECT id FROM teams WHERE team_name = 'devops'), true);

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) VALUES
    ('pr-1001', 'Add authentication', 'user1', 'OPEN', now() - interval '2 days', NULL),
    ('pr-1002', 'Fix login bug', 'user2', 'OPEN', now() - interval '1 day', NULL),
    ('pr-1003', 'Update API docs', 'user3', 'OPEN', now() - interval '3 hours', NULL),
    ('pr-1004', 'Refactor database', 'user1', 'MERGED', now() - interval '5 days', now() - interval '4 days'),
    ('pr-1005', 'Add tests', 'user2', 'MERGED', now() - interval '7 days', now() - interval '6 days'),
    ('pr-2001', 'Update UI components', 'user5', 'OPEN', now() - interval '1 day', NULL),
    ('pr-2002', 'Fix CSS bug', 'user6', 'OPEN', now() - interval '5 hours', NULL);

INSERT INTO pull_request_reviewers (pr_id, user_id, assigned_at) VALUES
    ('pr-1001', 'user2', now() - interval '2 days'),
    ('pr-1001', 'user3', now() - interval '2 days'),
    ('pr-1002', 'user1', now() - interval '1 day'),
    ('pr-1002', 'user3', now() - interval '1 day'),
    ('pr-1003', 'user1', now() - interval '3 hours'),
    ('pr-1003', 'user2', now() - interval '3 hours'),
    ('pr-1004', 'user2', now() - interval '5 days'),
    ('pr-1005', 'user1', now() - interval '7 days'),
    ('pr-1005', 'user3', now() - interval '7 days'),
    ('pr-2001', 'user6', now() - interval '1 day'),
    ('pr-2001', 'user7', now() - interval '1 day'),
    ('pr-2002', 'user5', now() - interval '5 hours');