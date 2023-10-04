-- Initialise the database for server use.

CREATE SCHEMA app;

CREATE TABLE app."user" (
  username VARCHAR(15) PRIMARY KEY,
  password BYTEA       NOT NULL
);

CREATE TABLE app.board (
  id   INTEGER     PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name VARCHAR(35) NOT NULL
);

CREATE TABLE app.user_board (
  id       INTEGER     PRIMARY KEY GENERATED ALWAYS      AS IDENTITY,
  username VARCHAR(15) NOT NULL    REFERENCES app."user",
  boardID  INTEGER     NOT NULL    REFERENCES app.board,
  isAdmin  BOOLEAN     NOT NULL
);

CREATE TABLE app."column" (
  id      INTEGER  PRIMARY KEY GENERATED ALWAYS     AS IDENTITY,
  boardID INTEGER  NOT NULL    REFERENCES app.board,
  "order" SMALLINT NOT NULL
);

CREATE TABLE app.task (
  id          INTEGER     PRIMARY KEY GENERATED ALWAYS        AS IDENTITY,
  columnID    INTEGER     NOT NULL    REFERENCES app."column",
  title       VARCHAR(50) NOT NULL,
  description TEXT
);

CREATE TABLE app.subtask (
  id     INTEGER     PRIMARY KEY GENERATED ALWAYS    AS IDENTITY,
  taskID INTEGER     NOT NULL    REFERENCES app.task,
  title  VARCHAR(50) NOT NULL,
  isDone BOOLEAN     NOT NULL
);

INSERT INTO app."user"(username, password) 
VALUES 
  ('bob123', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO'),
  ('bob124', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO');

INSERT INTO app.board(name) 
VALUES
    ('Board #1'), -- used for DELETE tests
    ('Board #2'),
    ('Board #3'),
    ('Board #4');

INSERT INTO app.user_board(username, boardID, isAdmin)
VALUES
  ('bob123', 1, TRUE),
  ('bob123', 2, TRUE),
  ('bob123', 3, TRUE),
  ('bob123', 4, FALSE);

-- insert columns into the second board for testing recursive board deletion
INSERT INTO app."column"(boardID, "order")
VALUES (1, 1), (1, 2), (1, 3), (1, 4);

-- insert a task into each column for testing recursive board deletion
INSERT INTO app.task(columnID, title)
VALUES (1, 'task A'), (2, 'task B'), (3, 'task C'), (4, 'task D');

-- insert a subtask into each task for testing recursive board deletion
INSERT INTO app.subtask(taskID, title, isDone)
VALUES
    (1, 'subtask A', false),
    (2, 'subtask B', false),
    (3, 'subtask C', false),
    (4, 'subtask D', false);

