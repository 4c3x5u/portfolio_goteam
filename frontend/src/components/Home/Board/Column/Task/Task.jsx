/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */
import {
  Menu,
  Item,
  // Separator,
  // Submenu,
  useContextMenu,
} from 'react-contexify';

import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import Subtask from './Subtask/Subtask';
import EditTask from './EditTask/EditTask';

import './task.sass';
import 'react-contexify/dist/ReactContexify.css';

const Task = ({ id, title, description }) => {
  const [subtasks, setSubtasks] = useState([]);
  const [showEdit, setShowEdit] = useState(false);

  const MENU_ID = `edit-task-${id}`;

  const { show } = useContextMenu({ id: MENU_ID });

  // const handleItemClick = ({ event, props, triggerEvent, data }) => (
  //   console.log('Menu item clicked')
  // );

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

      <Menu id={MENU_ID}>
        <Item onClick={() => setShowEdit(true)}>
          Edit Task
        </Item>
      </Menu>

      {
        showEdit && (
          <EditTask
            id={id}
            title={title}
            description={description}
            subtasks={subtasks}
            toggleOff={() => setShowEdit(false)}
          />
        )
      }
    </div>
  );
};

Task.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
};

export default Task;
