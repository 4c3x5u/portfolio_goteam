import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import TasksAPI from '../../../api/TasksAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import AddSubtasks from './AddSubtasks/AddSubtasks';
import ValidateTask from '../../../validation/ValidateTask';
import inputType from '../../../misc/inputType';

import logo from './createtask.svg';
import './createtask.sass';

const CreateTask = ({ toggleOff }) => {
  const { activeBoard, loadBoard, notify } = useContext(AppContext);
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
      TasksAPI
        .post({
          title,
          description,
          column: activeBoard.columns[0].id,
          subtasks: subtasks.list,
        })
        .then(() => {
          loadBoard();
          toggleOff();
        })
        .catch((err) => {
          const serverTitleError = err?.response?.data?.title || '';

          if (serverTitleError) {
            setTitleError(serverTitleError);
          } else {
            notify(
              'Unable to create task.',
              `${err?.message || 'Server Error'}.`,
            );
          }
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
              GO!
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

CreateTask.propTypes = { toggleOff: PropTypes.func.isRequired };

export default CreateTask;
