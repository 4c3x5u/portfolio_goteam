import React from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';

import './addsubtasks.sass';

const AddSubtasks = ({ subtasks, setSubtasks }) => (
  <Form.Group className="AddSubtasks">
    <Form.Label className="Label">
      SUBTASKS
    </Form.Label>

    {subtasks.list.map((subtask, index) => (
      <div className="ControlWrapper">
        <Form.Control
          key={subtask}
          className="Input"
          type="text"
          value={subtask}
          onChange={() => console.log('NOT ALLOWED')}
        />

        <Button
          className="Remove"
          onClick={() => {
            const subtaskList = subtasks.list;
            subtasks.list.splice(index - 1, 1);
            setSubtasks({ value: '', list: subtaskList });
          }}
          type="button"
        >
          <FontAwesomeIcon className="Icon" icon={faTimes} />
        </Button>
      </div>
    ))}

    <div className="ControlWrapper">
      <Form.Control
        className="Input"
        type="text"
        value={subtasks.value}
        onChange={(e) => setSubtasks({
          ...subtasks,
          value: e.target.value,
        })}
      />

      <Button
        className="Add"
        onClick={() => setSubtasks({
          value: '',
          list: [...subtasks.list, subtasks.value],
        })}
        type="button"
      >
        <FontAwesomeIcon className="Icon" icon={faPlus} />
      </Button>
    </div>
  </Form.Group>
);

AddSubtasks.propTypes = {
  subtasks: PropTypes.objectOf({
    value: PropTypes.string.isRequired,
    list: PropTypes.arrayOf(
      PropTypes.string.isRequired,
    ).isRequired,
  }).isRequired,
  setSubtasks: PropTypes.func.isRequired,
};

export default AddSubtasks;
