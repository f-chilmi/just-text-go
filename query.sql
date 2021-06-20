-- CREATE
INSERT INTO "users"("username", "phone", "password") VALUES ('Lala', '0813', '12345');

-- READ
SELECT * FROM users;

-- UPDATE
UPDATE users SET "username"='Lana' WHERE id=2;

-- UPDATE AND RETURNING UPDATED DATA
UPDATE users SET username='Lana', updated_at=CURRENT_TIMESTAMP WHERE id=2 RETURNING *;

-- DELETE
DELETE FROM users WHERE id=2;