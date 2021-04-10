/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import EditSubtasks from './EditSubtasks/EditSubtasks';
import inputType from '../../../misc/inputType';

import logo from './edittask.svg';
import './edittask.sass';

const EditTask = ({
  id, title, description, subtasks, toggleOff,
}) => {
  const [newTitle, setNewTitle] = useState(title);
  const [newDescription, setNewDescription] = useState(description);
  const [newSubtasks, setNewSubtasks] = useState({
    value: '',
    list: [],
  });

  useEffect(() => (
    // TODO: use task id to find and populate subtasks.
    setNewSubtasks({
      value: '',
      list: subtasks,
    })
  ), []);

  const handleSubmit = () => (
    // TODO: Set title, description, and subtasks as needed, submit to API
    console.log(`Task ID: ${id}`)
  );

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
