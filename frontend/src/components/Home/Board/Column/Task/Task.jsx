/* eslint-disable
react/jsx-props-no-spreading,
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Menu, Submenu, Item, useContextMenu,
} from 'react-contexify';
import { Draggable } from 'react-beautiful-dnd';
import _ from 'lodash/core';

import AppContext from '../../../../../AppContext';
import Subtask from './Subtask/Subtask';
import window from '../../../../../misc/window';
import TasksAPI from '../../../../../api/TasksAPI';

import './task.sass';
import 'react-contexify/dist/ReactContexify.css';

const Task = ({
  id,
  title,
  description,
  order,
  assignedUser,
  handleActivate,
  columnId,
  subtasks,
}) => {
  const {
    activeBoard, setActiveBoard, user, members, loadBoard, notify,
  } = useContext(AppContext);

  const MENU_ID = `edit-task-${id}`;
  const { show } = useContextMenu({ id: MENU_ID });

  const assignMember = (username) => {
    // Update client state to avoid load time
    setActiveBoard({
      ...activeBoard,
      columns: activeBoard.columns.map((column) => (
        column.id === columnId ? {
          ...column,
          tasks: column.tasks.map((task) => (
            task.id === id
              ? { ...task, user: username }
              : task
          )),
        } : column
      )),
    });

    // Add user to task in database
    TasksAPI
      .patch(id, { user: username })
      .catch((err) => {
        notify(
          'Unable to assign user.',
          err?.response?.data?.user || err?.message || 'Server Error.',
        );
      })
      .finally(() => loadBoard());
  };

  return (
    <Draggable
      draggableId={`draggable-${id}`}
      index={order}
      isDragDisabled={!user.isAdmin && user.username !== assignedUser}
    >
      {(provided) => (
        <div
          className="Task"
          onContextMenu={user.isAdmin && show}
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
        >
          <div className="TaskBody">
            <h1 className="Title">
              {title}
            </h1>

            <p className="Description">
              {description}
            </p>

            {subtasks.length > 0 && (
              <ul className="Subtasks">
                {_.sortBy(subtasks, (subtask) => subtask.order)
                  .map((subtask) => (
                    <Subtask
                      id={subtask.id}
                      key={subtask.id}
                      title={subtask.title}
                      done={subtask.done}
                      assignedUser={assignedUser}
                    />
                  ))}
              </ul>
            )}

            {assignedUser && (
              <div
                className={
                  `AssignedUser ${assignedUser === user.username && 'Me'}`
                }
              >
                @
                {assignedUser.length > 20
                  ? `${assignedUser.substring(0, 17)}...`
                  : assignedUser}
              </div>
            )}
          </div>

          <Menu className="ContextMenu" id={MENU_ID}>
            {members.length >= 2 && (
              <Submenu
                className="Submenu"
                label={(
                  <div
                    style={{
                      textAlign: 'center',
                      display: 'flex',
                      justifyContent: 'center',
                    }}
                  >
                    ASSIGN
                  </div>
                )}
                arrow={<div style={{ display: 'none' }} />}
              >
                {members.map((member) => (
                  // If member is not admin or assigned, display them in list
                  !member.isAdmin && member.isActive && (
                    <Item
                      key={member.username}
                      onClick={() => assignMember(member.username)}
                    >
                      {member.username.length > 20
                        ? `${member.username.substring(0, 17)}...`
                        : member.username}
                    </Item>
                  )
                ))}
              </Submenu>
            )}

            <Item
              className="ContextMenuItem"
              onClick={() => handleActivate(window.EDIT_TASK)({
                id,
                title,
                description,
                subtasks,
                columnId,
                toggleOff: handleActivate(window.NONE),
              })}
            >
              EDIT
            </Item>

            <Item
              className="ContextMenuItem"
              onClick={() => handleActivate(window.DELETE_TASK)({
                id,
                title,
                description,
                subtasks,
                columnId,
                toggleOff: handleActivate(window.NONE),
              })}
            >
              DELETE
            </Item>
          </Menu>
        </div>
      )}
    </Draggable>
  );
};

Task.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  order: PropTypes.number.isRequired,
  assignedUser: PropTypes.string.isRequired,
  columnId: PropTypes.number.isRequired,
  subtasks: PropTypes.arrayOf(
    PropTypes.exact({
      id: PropTypes.number.isRequired,
      title: PropTypes.string.isRequired,
      order: PropTypes.number.isRequired,
      done: PropTypes.bool.isRequired,
    }),
  ).isRequired,
  handleActivate: PropTypes.func.isRequired,
};

export default Task;
