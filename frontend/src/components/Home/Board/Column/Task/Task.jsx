/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Menu, Item, useContextMenu } from 'react-contexify';

import Subtask from './Subtask/Subtask';
import window from '../../../../../misc/window';

import './task.sass';
import 'react-contexify/dist/ReactContexify.css';

const Task = ({
  id, title, description, handleActivate,
}) => {
  const [subtasks, setSubtasks] = useState([]);

  const MENU_ID = `edit-task-${id}`;

  const { show } = useContextMenu({ id: MENU_ID });

  // TODO: API call based on task ID
  console.log(`Task ID: ${id}`);
  useEffect(() => {
    setSubtasks([
      {
        id: 0,
        title: 'Some subtask',
        order: 0,
        done: false,
      },
      {
        id: 1,
        title: 'Some subtask',
        order: 1,
        done: false,
      },
    ]);
  }, []);

  return (
    <div className="Task" onContextMenu={show}>
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
  );
};

Task.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  handleActivate: PropTypes.func.isRequired,
};

export default Task;
