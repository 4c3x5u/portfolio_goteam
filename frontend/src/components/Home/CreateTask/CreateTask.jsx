import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import AddSubtasks from './AddSubtasks/AddSubtasks';
import { inputType } from '../../../misc/enums';

import logo from '../../../assets/createtask.svg';
import './createtask.sass';

const CreateTask = ({ toggleOff }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [subtasks, setSubtasks] = useState({ value: '', list: [] });

  const handleSubmit = () => console.log('TODO: IMPLEMENT');

  return (
    <button className="CreateTask" onClick={toggleOff} type="button">
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
        />

        <FormGroup
          type={inputType.TEXTAREA}
          label="description"
          value={description}
          setValue={setDescription}
        />

        <AddSubtasks subtasks={subtasks} setSubtasks={setSubtasks} />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>
      </Form>
    </button>
  );
};

CreateTask.propTypes = { toggleOff: PropTypes.func.isRequired };

export default CreateTask;
