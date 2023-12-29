import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';

import ValidateSubtask from '../../../../validation/ValidateSubtask';

import './editsubtasks.sass';

const EditSubtasks = ({ subtasks, setSubtasks }) => {
  const [error, setError] = useState('');

  const addSubtask = () => {
    const titleError = ValidateSubtask.title(subtasks.value);

    if (titleError) {
      setError(titleError);
    } else {
      setSubtasks({
        value: '',
        list: [...subtasks.list, {
          title: subtasks.value,
          done: false,
        }],
      });
      setError('');
    }
  };

  const removeSubtask = (iSubtask) => {
    setSubtasks({
      value: '',
      list: subtasks.list.filter((_, i) => i !== iSubtask),
    });
  };

  return (
    <Form.Group className="EditSubtasks">
      <Form.Label className="Label">
        SUBTASKS
      </Form.Label>

      {subtasks.list.map((subtask, i) => (
        <div key={i} className="ControlWrapper">
          <Form.Control
            className="Input"
            type="text"
            value={subtask.title}
            disabled
          />

          {setSubtasks && (
            <Button
              className="Remove"
              onClick={() => removeSubtask(i)}
              type="button"
            >
              <FontAwesomeIcon className="Icon" icon={faTimes} />
            </Button>
          )}
        </div>
      ))}

      {setSubtasks && (
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
        </div>
      )}
    </Form.Group>
  );
};

EditSubtasks.propTypes = {
  subtasks: PropTypes.exact({
    value: PropTypes.string.isRequired,
    list: PropTypes.arrayOf(
      PropTypes.exact({
        title: PropTypes.string.isRequired,
        done: PropTypes.bool.isRequired,
      }),
    ).isRequired,
  }).isRequired,
  setSubtasks: PropTypes.func,
};

EditSubtasks.defaultProps = { setSubtasks: undefined };

export default EditSubtasks;
