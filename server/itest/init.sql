-- noinspection SqlNoDataSourceInspectionForFile

-- Initialise the database for server use.

CREATE SCHEMA app;

CREATE TABLE app.team (
    id         INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    inviteCode UUID    NOT NULL
);

CREATE TABLE app."user" (
    username VARCHAR(15) PRIMARY KEY,
    password BYTEA       NOT NULL,
    teamID   INTEGER     NOT NULL    REFERENCES app.team,
    isAdmin  BOOL        NOT NULL
);

CREATE TABLE app.board (
    id     INTEGER     PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name   VARCHAR(35) NOT NULL,
    teamID INTEGER     NOT NULL    REFERENCES app.team
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

INSERT INTO app.team(inviteCode)
VALUES 
    ('afeadc4a-68b0-4c33-9e83-4648d20ff26a'),
    ('66ca0ddf-5f62-4713-bcc9-36cb0954eb7b'),
    ('74c80ae5-64f3-4298-a8ff-48f8f920c7d4'),
    ('3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc');


INSERT INTO app."user"(username, password, teamID, isAdmin) 
VALUES 
    ('team1Admin', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 1, true),
    ('team1Member', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 1, false),
    ('team2Admin', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 2, true),
    ('team2Member', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 2, false),
    ('team3Admin', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 3, true),
    ('team4Admin', '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO', 4, true);

INSERT INTO app.board(name, teamID)
VALUES
    ('Team 1 Board 1',  1), -- used to verify max boards case for POST board
    ('Team 1 Board 2',  1), -- used to verify max boards case for POST board
                            -- also used in GET board tests
    ('Team 1 Board 3',  1), -- used to verify max boards case for POST board
    ('Team 3 Board 1',  3); -- gets deleted during DELETE board tests

INSERT INTO app."column"(boardID, "order")
VALUES
    (4, 1), -- gets deleted during DELETE board tests
    (4, 2), -- gets deleted during DELETE board tests
    (4, 3), -- gets deleted during DELETE board tests
    (4, 4), -- gets deleted during DELETE board tests
    (1, 1), -- used as source for PATCH column tests
    (1, 2), -- used as target for PATCH column tests
    (1, 3),
    (2, 1), -- used in GET board tests
    (2, 2), -- used in GET board tests
    (2, 3), -- used in GET board tests
    (2, 4); -- used in GET board tests

INSERT INTO app.task(columnID, title, description, "order")
VALUES
    (1,  'task 1',  null,   1), -- gets deleted during DELETE board tests
    (2,  'task 2',  null,   1), -- gets deleted during DELETE board tests
    (3,  'task 3',  null,   1), -- gets deleted during DELETE board tests
    (4,  'task 4',  null,   1), -- gets deleted during DELETE board tests
    (5,  'task 5',  null,   1), -- gets moved from column 5 to 6 during PATCH column test
    (7,  'task 6',  null,   1),
    (7,  'task 7',  null,   1),
    (7,  'task 8',  null,   1),
    (7,  'task 9',  null,   1), -- gets deleted during DELETE task tests
    (8,  'task 10', 'desc', 1), -- used in GET board tests
    (10, 'task 11', null,   1); -- used in GET board tests

INSERT INTO app.subtask(taskID, title, "order", isDone)
VALUES
    (1,  'subtask 1', 1, false), -- gets deleted during DELETE board tests
    (2,  'subtask 2', 1, false), -- gets deleted during DELETE board tests
    (3,  'subtask 3', 1, false), -- gets deleted during DELETE board tests
    (4,  'subtask 4', 1, false), -- gets deleted during DELETE board tests
    (8,  'subtask 5', 1, false), -- gets deleted during DELETE subtask tests
    (9,  'subtask 6', 1, false), -- gets deleted during DELETE task tests
    (9,  'subtask 7', 1, false), -- gets deleted during DELETE task tests
    (10, 'subtask 8', 1, false), -- used in GET board tests
    (10, 'subtask 9', 2, true); -- used in GET board tests
