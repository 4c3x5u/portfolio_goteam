/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React, { useState, useContext } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import AppContext from '../../../AppContext';
import TasksAPI from '../../../api/TasksAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import EditSubtasks from './EditSubtasks/EditSubtasks';
import inputType from '../../../misc/inputType';

import logo from './edittask.svg';
import './edittask.sass';

const EditTask = ({
  id, title, description, subtasks, toggleOff,
}) => {
  const { activeBoard, loadBoard } = useContext(AppContext);
  const [newTitle, setNewTitle] = useState(title);
  const [newDescription, setNewDescription] = useState(description);
  const [newSubtasks, setNewSubtasks] = useState({
    value: '',
    list: subtasks,
  });

  const handleSubmit = (e) => {
    e.preventDefault();

    const data = {
      title: newTitle,
      description: newDescription,
      column: activeBoard.columns[0].id,
      subtasks: newSubtasks.list,
    };

    TasksAPI
      .patch(id, data)
      .then(() => { loadBoard(); toggleOff(); })
      .catch((err) => console.error(err));
  };

  return (
    <div className="EditTask" onClick={toggleOff}>
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
        />

        <FormGroup
          type={inputType.TEXTAREA}
          label="description"
          value={newDescription}
          setValue={setNewDescription}
        />

        <EditSubtasks subtasks={newSubtasks} setSubtasks={setNewSubtasks} />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>
      </Form>
    </div>
  );
};

EditTask.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  subtasks: PropTypes.arrayOf({
    id: PropTypes.number.isRequired,
    title: PropTypes.string.isRequired,
    order: PropTypes.number.isRequired,
    done: PropTypes.bool.isRequired,
  }).isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default EditTask;
