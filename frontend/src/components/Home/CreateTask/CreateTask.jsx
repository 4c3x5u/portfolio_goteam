/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import AppContext from '../../../AppContext';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import AddSubtasks from './AddSubtasks/AddSubtasks';
import inputType from '../../../misc/inputType';
import { postTask } from '../../../misc/api';

import logo from './createtask.svg';
import './createtask.sass';

const CreateTask = ({ toggleOff }) => {
  const { activeBoard, loadActiveBoard } = useContext(AppContext);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [subtasks, setSubtasks] = useState({ value: '', list: [] });

  const handleSubmit = (e) => {
    e.preventDefault();

    // TODO: See if the ordering will work with this. If not, pass the order
    //       manually
    postTask({
      title,
      description,
      column: activeBoard.columns[0].id,
      subtasks: subtasks.list,
    }).then((res) => {
      loadActiveBoard();

      console.log(res.msg);
      toggleOff();
    }).catch((err) => {
      console.error(err);
    });
  };

  return (
    <div className="CreateTask" onClick={toggleOff}>
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
    </div>
  );
};

CreateTask.propTypes = { toggleOff: PropTypes.func.isRequired };

export default CreateTask;
