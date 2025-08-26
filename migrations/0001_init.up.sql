CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       username VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL
);

CREATE TABLE todo_list (
                           id SERIAL PRIMARY KEY,
                           title VARCHAR(150) NOT NULL,
                           description TEXT
);

CREATE TABLE todo_item (
                           id SERIAL PRIMARY KEY,
                           title VARCHAR(150) NOT NULL,
                           description TEXT,
                           done BOOLEAN DEFAULT FALSE
);

CREATE TABLE users_lists (
                             id SERIAL PRIMARY KEY,
                             user_id INT NOT NULL,
                             list_id INT NOT NULL,
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                             FOREIGN KEY (list_id) REFERENCES todo_list(id) ON DELETE CASCADE
);

CREATE TABLE list_items (
                            id SERIAL PRIMARY KEY,
                            list_id INT NOT NULL,
                            item_id INT NOT NULL,
                            FOREIGN KEY (list_id) REFERENCES todo_list(id) ON DELETE CASCADE,
                            FOREIGN KEY (item_id) REFERENCES todo_item(id) ON DELETE CASCADE
);
