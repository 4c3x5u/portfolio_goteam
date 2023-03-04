-- Initialise the database for server use.

CREATE SCHEMA app;

CREATE TABLE app."user" (
  id       VARCHAR(15) PRIMARY KEY,
  password BYTEA       NOT NULL
);

CREATE TABLE app.board (
  id   INTEGER     PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name VARCHAR(30) NOT NULL
);

CREATE TABLE app.user_board (
  id      INTEGER     PRIMARY KEY GENERATED ALWAYS      AS IDENTITY,
  userID  VARCHAR(15) NOT NULL    REFERENCES app."user",
  boardID INTEGER     NOT NULL    REFERENCES app.board,
  isAdmin BOOLEAN     NOT NULL
);

CREATE TABLE app."column" (
  id      INTEGER  PRIMARY KEY GENERATED ALWAYS     AS IDENTITY,
  tableID INTEGER  NOT NULL    REFERENCES app.board,
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

INSERT INTO app."user"(id, password) 
VALUES 
  ('bob123', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO'),
  ('bob124', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO');

INSERT INTO app.board(name) 
VALUES ('Board #1'), ('Board #2'), ('Board #3');

INSERT INTO app.user_board(userID, boardID, isAdmin) 
VALUES ('bob123', 1, TRUE), ('bob123', 2, TRUE), ('bob123', 3, TRUE);

