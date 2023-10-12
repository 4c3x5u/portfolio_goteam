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
  description TEXT,
  "order"     INTEGER     NOT NULL
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
    ('board 1'), -- used for DELETE board test
    ('board 2'),
    ('board 3'),
    ('board 4'),
    ('board 5'),
    ('board 6'),
    ('board 7'),
    ('board 8'),
    ('board 9'),
    ('board 10'),
    ('board 11');
-- a record is inserted here with the id of 6 during board POST tests.

INSERT INTO app.user_board(username, boardID, isAdmin)
VALUES
  ('bob123', 1,  TRUE),
  ('bob123', 2,  TRUE),
  ('bob123', 3,  TRUE),
  ('bob123', 4,  FALSE),
  ('bob123', 8,  FALSE),
  ('bob123', 9,  TRUE),
  ('bob123', 11, FALSE);
-- a board is inserted here with the values of ('bob124', 6, TRUE) during board
-- POST tests.

-- insert columns into the second board for testing recursive board deletion
INSERT INTO app."column"(boardID, "order")
VALUES
    (1,  1), -- id: 1 (used for board DELETE tests)
    (1,  2), -- id: 2 (used for board DELETE tests)
    (1,  3), -- id: 3 (used for board DELETE tests)
    (1,  4), -- id: 4 (used for board DELETE tests)
    (2,  1), -- id: 5
    (4,  1), -- id: 6
    (2,  2), -- id: 7
    (6,  1), -- id: 8
    (4,  2), -- id: 9
    (3,  1), -- id: 10
    (3,  2), -- id: 11
    (5,  1), -- id: 12
    (4,  3), -- id: 13
    (3,  3), -- id: 14
    (7,  1), -- id: 15
    (8,  1), -- id: 16
    (9,  1), -- id: 17
    (10, 1), -- id: 18
    (11, 1); -- id: 19

-- insert a task into each column for testing recursive board deletion
INSERT INTO app.task(columnID, title, "order")
VALUES
    (1,  'task 1', 1), -- (used for board DELETE board tests)
    (2,  'task 2', 1),
    (3,  'task 3', 1),
    (4,  'task 4', 1),
    (10, 'task 5', 1),
    (10, 'task 6', 2),
    (8,  'task 7', 1),
    (9,  'task 8', 1),
    (11, 'task 9', 1),
    (15, 'task 10', 1),
    (16, 'task 11', 1),
    (17, 'task 12', 1),
    (18, 'task 13', 1),
    (19, 'task 14', 1);

-- insert a subtask into each task for testing recursive board deletion
INSERT INTO app.subtask(taskID, title, "order", isDone)
VALUES
    (1,  'subtask 1', 1, false),
    (2,  'subtask 2', 1, false),
    (3,  'subtask 3', 1, false),
    (4,  'subtask 4', 1, false),
    (9,  'subtask 5', 1, false),
    (13, 'subtask 6', 1, false),
    (14, 'subtask 7', 1, false);


