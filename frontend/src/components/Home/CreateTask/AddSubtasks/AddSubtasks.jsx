import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';

import './addsubtasks.sass';
import ValidateSubtask from '../../../../validation/ValidateSubtask';

const AddSubtasks = ({ subtasks, setSubtasks }) => {
  const [error, setError] = useState('');

  const addSubtask = () => {
    const titleError = ValidateSubtask.title(subtasks.value);

    if (titleError) {
      setError(titleError);
    } else {
      setSubtasks({
        value: '',
        list: [...subtasks.list, subtasks.value],
      });
      setError('');
    }
  };

  const removeSubtask = (subtaskIndex) => {
    setSubtasks({
      value: '',
      list: subtasks.list.filter((_, i) => i !== subtaskIndex),
    });
  };

  return (
    <div className="AddSubtasks">
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
            onClick={() => removeSubtask(index)}
            type="button"
          >
            <FontAwesomeIcon className="Icon" icon={faTimes} />
          </Button>
        </div>
      ))}

      <Form.Group className="ControlWrapper">
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
          onClick={addSubtask}
          type="button"
        >
          <FontAwesomeIcon className="Icon" icon={faPlus} />
        </Button>

        {error && (
          <span className="Error text-danger">
            {error}
          </span>
        )}
      </Form.Group>
    </div>
  );
};

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
