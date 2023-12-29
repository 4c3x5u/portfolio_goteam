import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import TaskAPI from '../../../api/TaskAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import AddSubtasks from './AddSubtasks/AddSubtasks';
import ValidateTask from '../../../validation/ValidateTask';
import inputType from '../../../misc/inputType';

import logo from './createtask.svg';
import './createtask.sass';

const CreateTask = ({ toggleOff }) => {
  const {
    activeBoard, setActiveBoard, loadBoard, notify,
  } = useContext(AppContext);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [subtasks, setSubtasks] = useState({ value: '', list: [] });
  const [titleError, setTitleError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientTitleError = ValidateTask.title(title);

    if (clientTitleError) {
      setTitleError(clientTitleError);
    } else {
      // Keep an initial state to avoid loadBoard() on API error
      const initialActiveBoard = activeBoard;

      const subts = subtasks.list.map((subtask) => ({ title: subtask }))

      // Update client state to avoid load screen
      setActiveBoard({
        ...activeBoard,
        columns: activeBoard.columns.map((column) => (
          column.order === 1 ? {
            ...column,
            tasks: [
              ...column.tasks,
              {
                teamID: "",
                boardID: "",
                id: "",
                title,
                description,
                colNo: column.order,
                order: -1,
                user: '',
                subtasks: subts,
              },
            ],
          } : column
        )),
      });

      // Create task in the database
      TaskAPI
        .post({
          boardID: activeBoard.id,
          title,
          description,
          colNo: 0,
          subtasks: subts,
        })
        .then(() => {
          // Load board to retrieve the "actual" ID of the created task
          loadBoard();
          toggleOff();
        })
        .catch((err) => {
          const serverTitleError = err?.response?.data?.title || '';
          const createTaskError = err?.response?.data?.error;
          if (serverTitleError) {
            setTitleError(serverTitleError);
          } else {
            notify(
              'Unable to create task.',
              `${createTaskError || 'Server Error'}.`,
            );
          }
          setActiveBoard(initialActiveBoard);
        });
    }
  };

  return (
    <div className="CreateTask">
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
          value={title}
          setValue={setTitle}
          error={titleError}
        />

        <FormGroup
          type={inputType.TEXTAREA}
          label="description"
          value={description}
          setValue={setDescription}
        />

        <AddSubtasks
          subtasks={subtasks}
          setSubtasks={setSubtasks}
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
              CREATE
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

CreateTask.propTypes = { toggleOff: PropTypes.func.isRequired };

export default CreateTask;
