import React from 'react';
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';

import './addsubtasks.sass';

const AddSubtasks = ({ subtasks, setSubtasks }) => (
  <Form.Group className="FormGroup">
    <Form.Label className="Label">
      SUBTASKS
    </Form.Label>

    {subtasks.list.map((st) => (
      <Form.Control
        key={st}
        className="Input"
        type="text"
        value={st}
        onChange={() => console.log('NOT ALLOWED')}
      />
    ))}

    <Form.Control
      className="Input"
      type="text"
      value={subtasks.value}
      onChange={(e) => setSubtasks({
        ...subtasks,
        value: e.target.value,
      })}
    />
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
