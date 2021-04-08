import React from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';

import './editsubtasks.sass';

const EditSubtasks = ({ subtasks, setSubtasks }) => (
  <Form.Group className="EditSubtasks">
    <Form.Label className="Label">
      SUBTASKS
    </Form.Label>

    {subtasks.list.map((subtask) => (
      <div className="ControlWrapper">
        <Form.Control
          key={subtask.id}
          className="Input"
          type="text"
          value={subtask.title}
          disabled
        />

        {setSubtasks && (
          <Button
            className="Remove"
            onClick={() => {
              setSubtasks({
                value: '',
                list: subtasks.list.filter((st) => st.id !== subtask.id),
              });
            }}
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
          onClick={() => setSubtasks({
            value: '',
            list: [...subtasks.list, {
              id: null,
              title: subtasks.value,
              order: subtasks.list.length,
              done: false,
            }],
          })}
          type="button"
        >
          <FontAwesomeIcon className="Icon" icon={faPlus} />
        </Button>
      </div>
    )}
  </Form.Group>
);

EditSubtasks.propTypes = {
  subtasks: PropTypes.objectOf({
    value: PropTypes.string.isRequired,
    list: PropTypes.arrayOf(
      PropTypes.string.isRequired,
    ).isRequired,
  }).isRequired,
  setSubtasks: PropTypes.func,
};

EditSubtasks.defaultProps = { setSubtasks: null };

export default EditSubtasks;
