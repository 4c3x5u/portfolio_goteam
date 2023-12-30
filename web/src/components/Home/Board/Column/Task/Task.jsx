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
import TaskAPI from '../../../../../api/TaskAPI';

import './task.sass';
import 'react-contexify/dist/ReactContexify.css';

const Task = ({
  id,
  title,
  description,
  order,
  assignedUser,
  handleActivate,
  colNo,
  subtasks,
}) => {
  const {
    user, activeBoard, setActiveBoard, members, notify,
  } = useContext(AppContext);

  const MENU_ID = `edit-task-${id}`;
  const { show } = useContextMenu({ id: MENU_ID });

  // TODO: use
  // const assignMember = (username) => {
  //   // Keep an initial state to avoid loadBoard() on API error
  //   const initialActiveBoard = activeBoard;
  //
  //   // Update client state to avoid load time
  //   setActiveBoard({
  //     ...activeBoard,
  //     column: activeBoard.columns.map((column, i) => (
  //       i === colNo ? {
  //         ...column,
  //         tasks: column.tasks.map((task) => (
  //           task.id === id
  //             ? { ...task, user: username }
  //             : task
  //         )),
  //       } : column
  //     )),
  //   });
  //
  //   // Add user to task in database
  //   TaskAPI
  //     .patch(id, { title, description, subtasks, asignee: username })
  //     .catch((err) => {
  //       notify(
  //         'Unable to assign user.',
  //         err?.response?.data?.user || err?.message || 'Server Error.',
  //       );
  //       setActiveBoard(initialActiveBoard);
  //     });
  // };

  const toggleSubtaskDone = (iSubtask) => () => {
    // Keep an initial state to avoid loadBoard() on API error
    const initialActiveBoard = activeBoard;

    let subtasks = subtasks.map((subtask, i) => (
      i === iSubtask
        ? { title: subtask.title, done: !subtask.done }
        : subtask
    ));

    // Update client state to avoid load screen.
    setActiveBoard({
      ...activeBoard,
      columns: activeBoard.columns.map((column, i) => (
        i === colNo ? {
          ...column,
          tasks: column.tasks.map((task) => (
            task.id === id
              ? { ...task, subtasks: subtasks }
              : task
          )),
        } : column
      )),
    });

    // Update subtask in database
    TaskAPI
      .patch(id, { title, description, order, subtasks })
      .catch((err) => {
        notify(
          'Unable to check subtask done.',
          `${err?.message || 'Server Error'}.`,
        );
        setActiveBoard(initialActiveBoard);
      });
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

            {subtasks && subtasks.length > 0 && (
              <ul className="Subtasks">
                {subtasks.map((subtask, i) => (
                  <Subtask
                    key={i}
                    title={subtask.title}
                    done={subtask.done}
                    assignee={assignedUser}
                    toggleDone={toggleSubtaskDone(i)}
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
            {/* TODO: add back once you created assign endpoint */}
            {/* {members.filter((m) => m.isActive && !m.isAdmin).length > 0 && ( */}
            {/*   <Submenu */}
            {/*     className="Submenu" */}
            {/*     label={( */}
            {/*       <div */}
            {/*         style={{ */}
            {/*           textAlign: 'center', */}
            {/*           display: 'flex', */}
            {/*           justifyContent: 'center', */}
            {/*         }} */}
            {/*       > */}
            {/*         ASSIGN */}
            {/*       </div> */}
            {/*     )} */}
            {/*     arrow={<div style={{ display: 'none' }} />} */}
            {/*   > */}
            {/*     {members.map((member) => ( */}
            {/*       // If member is not admin or assigned, display them in list */}
            {/*       !member.isAdmin && member.isActive && ( */}
            {/*         <Item */}
            {/*           key={member.username} */}
            {/*           onClick={() => assignMember(member.username)} */}
            {/*         > */}
            {/*           {member.username.length > 20 */}
            {/*             ? `${member.username.substring(0, 17)}...` */}
            {/*             : member.username} */}
            {/*         </Item> */}
            {/*       ) */}
            {/*     ))} */}
            {/*   </Submenu> */}
            {/* )} */}

            <Item
              className="ContextMenuItem"
              onClick={() => handleActivate(window.EDIT_TASK)({
                id,
                title,
                description,
                subtasks,
                colNo: colNo,
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
                colNo,
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
  id: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  order: PropTypes.number.isRequired,
  // assignee: PropTypes.string,
  colNo: PropTypes.number.isRequired,
  subtasks: PropTypes.arrayOf(
    PropTypes.exact({
      title: PropTypes.string.isRequired,
      done: PropTypes.bool.isRequired,
    }),
  ),
  handleActivate: PropTypes.func.isRequired,
};

export default Task;
