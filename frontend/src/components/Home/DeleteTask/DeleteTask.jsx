/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React from 'react';
import PropTypes from 'prop-types';
import {
  Button, Form, Row, Col,
} from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import EditSubtasks from '../EditTask/EditSubtasks/EditSubtasks';
import inputType from '../../../misc/inputType';

import logo from './deletetask.svg';
import './deletetask.sass';

const DeleteTask = ({
  id, title, description, subtasks, toggleOff,
}) => {
  // TODO: Use the task ID to delete it here
  const handleSubmit = () => (console.log(`Delete Task ${id}`));

  return (
    <div className="DeleteTask" onClick={toggleOff}>
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
          disabled
        />

        <FormGroup
          type={inputType.TEXTAREA}
          label="description"
          value={description}
          disabled
        />

        <EditSubtasks subtasks={{ list: subtasks }} />

        <Row className="ButtonWrapper">
          <Col className="ButtonCol">
            <Button
              className="Button CancelButton"
              type="submit"
              aria-label="submit"
            >
              CANCEL
            </Button>
          </Col>
          <Col className="ButtonCol">
            <Button
              className="Button DeleteButton"
              type="submit"
              aria-label="submit"
            >
              DELETE
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

DeleteTask.propTypes = {
  id: PropTypes.number.isRequired,
  title: PropTypes.string.isRequired,
  description: PropTypes.string,
  subtasks: PropTypes.arrayOf({
    id: PropTypes.number.isRequired,
    title: PropTypes.string.isRequired,
    order: PropTypes.number.isRequired,
    done: PropTypes.bool.isRequired,
  }).isRequired,
  toggleOff: PropTypes.func.isRequired,
};

DeleteTask.defaultProps = {
  description: null,
};

export default DeleteTask;
