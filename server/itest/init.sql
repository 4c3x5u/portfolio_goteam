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

INSERT INTO app.team(inviteCode)
VALUES 
    ('afeadc4a-68b0-4c33-9e83-4648d20ff26a'),
    ('66ca0ddf-5f62-4713-bcc9-36cb0954eb7b'),
    ('74c80ae5-64f3-4298-a8ff-48f8f920c7d4');

INSERT INTO app."user"(username, password, teamID, isAdmin) 
VALUES 
    (
        'team1Admin',
        '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO',
        1,
        true
    ),
    (
        'team1Member',
        '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO',
        1,
        false
    ),
    (
        'team2Admin',
        '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO',
        2,
        true
    ),
    (
        'team3Admin',
        '$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO',
        3,
        false
    );

-- DO NOT insert a board for team 2 because POST board assertions rely on total
-- board count
INSERT INTO app.board(name, teamID)
VALUES
    ('Team 1 Board 1',  1),
    ('Team 1 Board 2',  1),
    ('Team 1 Board 3',  1),
    ('Team 3 Board 1',  3); -- gets deleted during DELETE board tests

-- insert columns into the second board for testing recursive board deletion
INSERT INTO app."column"(boardID, "order")
VALUES
    (3,  1), -- gets deleted during DELETE board tests
    (3,  2), -- gets deleted during DELETE board tests
    (3,  3), -- gets deleted during DELETE board tests
    (3,  4); -- gets deleted during DELETE board tests
--     (2,  1), -- id: 5
--     (4,  1), -- id: 6
--     (2,  2), -- id: 7
--     (6,  1), -- id: 8
--     (4,  2), -- id: 9
--     (3,  1), -- id: 10
--     (3,  2), -- id: 11
--     (5,  1), -- id: 12
--     (4,  3), -- id: 13
--     (3,  3), -- id: 14
--     (7,  1), -- id: 15
--     (8,  1), -- id: 16
--     (9,  1), -- id: 17
--     (10, 1), -- id: 18
--     (11, 1), -- id: 19
--     (12, 1); -- id: 20

-- insert a task into each column for testing recursive board deletion
INSERT INTO app.task(columnID, title, "order")
VALUES
    (1,  'task 1', 1), -- gets deleted during DELETE board tests
    (2,  'task 2', 1), -- gets deleted during DELETE board tests
    (3,  'task 3', 1), -- gets deleted during DELETE board tests
    (4,  'task 4', 1); -- gets deleted during DELETE board tests
--     (10, 'task 5', 1),
--     (10, 'task 6', 2),
--     (8,  'task 7', 1),
--     (9,  'task 8', 1),
--     (11, 'task 9', 1),
--     (15, 'task 10', 1),
--     (16, 'task 11', 1),
--     (17, 'task 12', 1),
--     (18, 'task 13', 1),
--     (19, 'task 14', 1),
--     (20, 'task 15', 1);

-- insert a subtask into each task for testing recursive board deletion
-- INSERT INTO app.subtask(taskID, title, "order", isDone)
-- VALUES
--     (1,  'subtask 1', 1, false),
--     (2,  'subtask 2', 1, false),
--     (3,  'subtask 3', 1, false),
--     (4,  'subtask 4', 1, false),
--     (9,  'subtask 5', 1, false),
--     (13, 'subtask 6', 1, false),
--     (14, 'subtask 7', 1, false),
--     (15, 'subtask 8', 1, false);
