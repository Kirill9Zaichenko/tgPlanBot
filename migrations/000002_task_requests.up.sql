CREATE TABLE IF NOT EXISTS task_requests (
                                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                                             task_id INTEGER NOT NULL,
                                             sender_user_id INTEGER NOT NULL,
                                             receiver_user_id INTEGER NOT NULL,
                                             status TEXT NOT NULL,
                                             comment TEXT NOT NULL DEFAULT '',
                                             decided_at DATETIME NULL,
                                             created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                             FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_task_requests_task_id ON task_requests(task_id);
CREATE INDEX IF NOT EXISTS idx_task_requests_receiver_user_id ON task_requests(receiver_user_id);
CREATE INDEX IF NOT EXISTS idx_task_requests_status ON task_requests(status);