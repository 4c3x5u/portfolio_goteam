/* eslint-disable
react/jsx-props-no-spreading,
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React from 'react';
import PropTypes from 'prop-types';
import { Menu, Item, useContextMenu } from 'react-contexify';
import { Draggable } from 'react-beautiful-dnd';

import Subtask from './Subtask/Subtask';
import window from '../../../../../misc/window';

import './task.sass';
import 'react-contexify/dist/ReactContexify.css';

const Task = ({
  id, title, description, order, handleActivate, subtasks,
}) => {
  const MENU_ID = `edit-task-${id}`;

  const { show } = useContextMenu({ id: MENU_ID });

  return (
    <Draggable
      key={id}
      draggableId={`draggable-${id}`}
      index={order}
    >
      {(provided) => (
        <div
          className="Task"
          onContextMenu={show}
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
                {subtasks.sort((subtask) => subtask.order).map((subtask) => (
                  <Subtask
                    id={subtask.id}
                    title={subtask.title}
                    done={subtask.done}
                  />
                ))}
              </ul>
            )}
          </div>

          <Menu className="ContextMenu" id={MENU_ID}>
            <Item
              className="ContextMenuItem"
              onClick={() => handleActivate(window.EDIT_TASK)({
                id,
                title,
                description,
                subtasks,
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
  handleActivate: PropTypes.func.isRequired,
  subtasks: PropTypes.arrayOf({
    id: PropTypes.number.isRequired,
    title: PropTypes.string.isRequired,
    order: PropTypes.number.isRequired,
    done: PropTypes.bool.isRequired,
  }).isRequired,
};

export default Task;
