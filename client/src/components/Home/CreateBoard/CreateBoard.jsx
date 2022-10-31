import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Form, Button, Row, Col,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import BoardsAPI from '../../../api/BoardsAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';
import ValidateBoard from '../../../validation/ValidateBoard';

import logo from './createboard.svg';
import './createboard.sass';

const CreateBoard = ({ toggleOff }) => {
  const {
    boards, setBoards, user, loadBoard, notify,
  } = useContext(AppContext);
  const [name, setName] = useState('');
  const [nameError, setNameError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientNameError = ValidateBoard.name(name);

    if (clientNameError) {
      setNameError(clientNameError);
    } else {
      // Keep an initial state to avoid loadBoard() on API error
      const initialBoards = boards;

      // Update client state to avoid load time
      setBoards([...boards, { id: Number.MAX_SAFE_INTEGER, name }]);

      // Create board in database
      BoardsAPI
        .post({ name, team_id: user.teamId })
        .then(() => {
          loadBoard();
          toggleOff();
        })
        .catch((err) => {
          const serverNameError = err?.response?.data?.name;
          if (serverNameError) {
            setNameError(serverNameError);
          } else if (err?.message) {
            notify(
              'Unable to create board.',
              `${err.message || 'Server Error'}.`,
            );
          }
          setBoards(initialBoards);
        });
    }
  };

  return (
    <div className="CreateBoard">
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
          value={name}
          setValue={setName}
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
              CREATE
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

CreateBoard.propTypes = {
  toggleOff: PropTypes.func.isRequired,
};

export default CreateBoard;
