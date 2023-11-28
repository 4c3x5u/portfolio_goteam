import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import BoardAPI from '../../../api/BoardAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import ValidateBoard from '../../../validation/ValidateBoard';
import inputType from '../../../misc/inputType';

import logo from './editboard.svg';
import './editboard.sass';

const EditBoard = ({ id, name, toggleOff }) => {
  const { boards, setBoards, notify } = useContext(AppContext);
  const [newName, setNewName] = useState(name);
  const [nameError, setNameError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientNameError = ValidateBoard.name(newName);

    if (clientNameError) {
      setNameError(clientNameError);
    } else {
      // Keep an initial state to avoid loadBoard() on API error
      const initialBoards = boards;

      // Update client state to avoid load time
      setBoards(boards.map((board) => (
        board.id === id
          ? { ...board, name: newName }
          : board
      )));

      // Edit board in database
      BoardAPI
        .patch(id, { name: newName })
        .then(toggleOff)
        .catch((err) => {
          const serverNameError = err?.response?.data?.name;
          const editBoardError = err.response.data.message;
          if (serverNameError) {
            setNameError(serverNameError);
          } else if (editBoardError) {
            notify(
              'Unable to edit board.',
              `${editBoardError || 'Server Error'}.`,
            );
          }
          setBoards(initialBoards);
        });
    }
  };

  return (
    <div className="EditBoard">
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
              SUBMIT
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
