/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import BoardsAPI from '../../../api/BoardsAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import ValidateBoard from '../../../validation/ValidateBoard';
import inputType from '../../../misc/inputType';

import logo from './editboard.svg';
import './editboard.sass';

const EditBoard = ({ id, name, toggleOff }) => {
  const { loadBoard, notify } = useContext(AppContext);
  const [newName, setNewName] = useState(name);
  const [nameError, setNameError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientNameError = ValidateBoard.name(newName);

    if (clientNameError) {
      setNameError(clientNameError);
    } else {
      BoardsAPI
        .patch(id, { name: newName })
        .then(() => {
          toggleOff();
          loadBoard();
        })
        .catch((err) => {
          const serverNameError = err?.response?.data?.name;
          if (serverNameError) {
            setNameError(serverNameError);
          } else if (err?.message) {
            notify(err.message || 'Server Error', 'Edit board failure.');
          }
        });
    }
  };

  return (
    <div className="EditBoard" onClick={toggleOff}>
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
          label="name"
          value={newName}
          setValue={setNewName}
          error={nameError}
        />

        <Row className="ButtonWrapper">
          <Col className="ButtonCol">
            <Button
              className="Button CancelButton"
              type="button"
              aria-label="cancel"
              onClick={toggleOff}
            >
              CANCEL
            </Button>
          </Col>

          <Col className="ButtonCol">
            <Button
              className="Button GoButton"
              type="submit"
              aria-label="submit"
            >
              GO!
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

EditBoard.propTypes = {
  id: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default EditBoard;
