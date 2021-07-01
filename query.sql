-- CREATE TABLE USERS
INSERT INTO "users"("username", "phone", "password") VALUES ('Lala', '0813', '12345');

-- READ
SELECT * FROM users;

-- UPDATE
UPDATE users SET "username"='Lana' WHERE id=2;

-- UPDATE AND RETURNING UPDATED DATA
UPDATE users SET username='Lana', updated_at=CURRENT_TIMESTAMP WHERE id=2 RETURNING *;

-- DELETE
DELETE FROM users WHERE id=2;

-- CREATE TABLE MESSAGES
CREATE TABLE 
  messages 
  (
    id serial PRIMARY KEY, 
    id_sender int NOT NULL, 
    id_recipient int NOT NULL, 
    content VARCHAR (255) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    -- updated at still not updated 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_sender) REFERENCES users (id),
    FOREIGN KEY (id_recipient) REFERENCES users (id)
  );

-- ADD COLUMN ID_ROOM INTO TABLE MESSAGES 
ALTER TABLE messages
ADD COLUMN id_room INT NOT NULL,
FOREIGN KEY (id_room) REFERENCES rooms (id);

-- CREATE TABLE ROOMS 
CREATE TABLE
  rooms (
    id serial PRIMARY KEY,
    id_user1 int NOT NULL,
    id_user2 int NOT NULL,
    id_last_msg INT NOT NULL,
    last_msg VARCHAR (255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    -- updated at still not updated 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user1) REFERENCES users (id),
    FOREIGN KEY (id_user2) REFERENCES users (id),
    FOREIGN KEY (id_last_msg) REFERENCES messages (id)
  );

CREATE TABLE
  rooms (
    id serial PRIMARY KEY,
    id_recipient int NOT NULL,
    id_last_msg INT NOT NULL,
    last_msg VARCHAR (255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    -- updated at still not updated 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_recipient) REFERENCES users (id),
    FOREIGN KEY (id_last_msg) REFERENCES messages (id)
  );