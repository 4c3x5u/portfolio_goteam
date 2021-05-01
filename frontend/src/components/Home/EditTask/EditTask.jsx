import React, { useState, useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Form, Row, Col,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import TasksAPI from '../../../api/TasksAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import EditSubtasks from './EditSubtasks/EditSubtasks';
import inputType from '../../../misc/inputType';
import ValidateTask from '../../../validation/ValidateTask';

import logo from './edittask.svg';
import './edittask.sass';

const EditTask = ({
  id, title, description, subtasks, columnId, toggleOff,
}) => {
  const {
    activeBoard, setActiveBoard, loadBoard, notify,
  } = useContext(AppContext);
  const [newTitle, setNewTitle] = useState(title);
  const [newDescription, setNewDescription] = useState(description);
  const [newSubtasks, setNewSubtasks] = useState({
    value: '',
    list: subtasks,
  });
  const [titleError, setTitleError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientTitleError = ValidateTask.title(newTitle);

    if (clientTitleError) {
      setTitleError(clientTitleError);
    } else {
      // Keep an initial state to avoid loadBoard() on API error
      const initialActiveBoard = activeBoard;

      // Update client state to avoid load time
      setActiveBoard({
        ...activeBoard,
        columns: activeBoard.columns.map((column) => (
          column.id === columnId ? {
            ...column,
            tasks: column.tasks.map((task) => (
              task.id === id ? {
                ...task,
                title: newTitle,
                description: newDescription,
                subtasks: newSubtasks.list.map((subtask, i) => ({
                  id: subtask.id || -100 + i,
                  title: subtask.title,
                  order: subtask.order || -100 + i,
                  done: !!subtask.done,
                })),
              } : task
            )),
          } : column
        )),
      });

      // Update task in database
      TasksAPI
        .patch(id, {
          title: newTitle,
          description: newDescription,
          column: activeBoard.columns[0].id,
          subtasks: newSubtasks.list,
        })
        .then(() => {
          // Load board to retrieve the "actual" subtask IDs
          loadBoard();
          toggleOff();
        })
        .catch((err) => {
          const serverTitleError = err?.response?.data?.title || '';
          if (serverTitleError) {
            setTitleError(serverTitleError);
          } else {
            notify(
              'Unable to edit task.',
              `${err?.message || 'Server Error'}.`,
            );
          }
          setActiveBoard(initialActiveBoard);
        });
    }
  };

  return (
    <div className="EditTask">
      <Form
        className="Form"
        onSubmit={handleSubmit}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup
          type={inputType.TEXT}
          label="title"
          value={newTitle}
          setValue={setNewTitle}
          error={titleError}
        />

        <FormGroup
          type={inputType.TEXTAREA}
          label="description"
          value={newDescription}
          setValue={setNewDescription}
        />

        <EditSubtasks
          subtasks={newSubtasks}
          setSubtasks={setNewSubtasks}
        />

        <Row className="ButtonWrapper">
          <Col className="ButtonCol">
            <Button
              className="Button CancelButton"
              type="button"
              aria-label="cancel"
              onClick={toggleOff}
            >
              CANCEL
            </Button>
          </Col>

          <Col className="ButtonCol">
            <Button
              className="Button GoButton"
              type="submit"
              aria-label="submit"
            >
              SUBMIT
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

EditTask.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  subtasks: PropTypes.arrayOf(
    PropTypes.exact({
      id: PropTypes.number.isRequired,
      title: PropTypes.string.isRequired,
      order: PropTypes.number.isRequired,
      done: PropTypes.bool.isRequired,
    }),
  ).isRequired,
  columnId: PropTypes.number.isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default EditTask;
