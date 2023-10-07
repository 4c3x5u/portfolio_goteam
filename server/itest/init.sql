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
  "order"     INTEGER     NOT NULL,
  description TEXT
);

CREATE TABLE app.subtask (
  id      INTEGER     PRIMARY KEY GENERATED ALWAYS    AS IDENTITY,
  taskID  INTEGER     NOT NULL    REFERENCES app.task,
  title   VARCHAR(50) NOT NULL,
  "order" INTEGER     NOT NULL,
  isDone  BOOLEAN     NOT NULL
);

INSERT INTO app."user"(username, password) 
VALUES 
  ('bob123', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO'),
  ('bob124', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO');

INSERT INTO app.board(name) 
VALUES
    ('Board #1'), -- id: 1 (used for board DELETE tests)
    ('Board #2'), -- id: 2
    ('Board #3'), -- id: 3
    ('Board #4'), -- id: 4
    ('Board #5'); -- id: 5
-- a record is inserted here with the id of 6 during board POST tests.

INSERT INTO app.user_board(username, boardID, isAdmin)
VALUES
  ('bob123', 1, TRUE),
  ('bob123', 2, TRUE),
  ('bob123', 3, TRUE),
  ('bob123', 4, FALSE);
-- a board is inserted here with the values of ('bob124', 6, TRUE) during board
-- POST tests.

-- insert columns into the second board for testing recursive board deletion
INSERT INTO app."column"(boardID, "order")
VALUES
    (1, 1), -- id: 1 (used for board DELETE tests)
    (1, 2), -- id: 2 (used for board DELETE tests)
    (1, 3), -- id: 3 (used for board DELETE tests)
    (1, 4), -- id: 4 (used for board DELETE tests)
    (2, 1), -- id: 5
    (4, 1), -- id: 6
    (2, 2); -- id: 7

-- insert a task into each column for testing recursive board deletion
INSERT INTO app.task(columnID, title, "order")
VALUES
    (1, 'task A', 1), -- id: 1 (used for board DELETE tests)
    (2, 'task B', 1), -- id: 2
    (3, 'task C', 1), -- id: 3
    (4, 'task D', 1), -- id: 4
    (5, 'task D', 1); -- id: 4

-- insert a subtask into each task for testing recursive board deletion
INSERT INTO app.subtask(taskID, title, "order", isDone)
VALUES
    (1, 'subtask A', 1, false),
    (2, 'subtask B', 1, false),
    (3, 'subtask C', 1, false),
    (4, 'subtask D', 1, false);

