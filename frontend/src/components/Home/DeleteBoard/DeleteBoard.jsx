import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import AppContext from '../../../AppContext';
import BoardsAPI from '../../../api/BoardsAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';

import logo from './deleteboard.svg';
import './deleteboard.sass';

const DeleteBoard = ({ id, name, toggleOff }) => {
  const {
    activeBoard, boards, setBoards, loadBoard, notify,
  } = useContext(AppContext);

  const handleSubmit = (e) => {
    e.preventDefault();

    // Keep an initial state to avoid loadBoard() on API error
    const initialBoards = boards;

    // Update client state to avoid load time
    const newBoards = boards.filter((board) => board.id !== id);
    setBoards(newBoards);

    // Delete board in database
    BoardsAPI
      .delete(id)
      .then(() => {
        toggleOff();
        sessionStorage.removeItem('board-id');
        return (activeBoard.id === id && loadBoard(newBoards[0].id));
      })
      .catch((err) => {
        notify(
          'Unable to delete board.',
          `${err.message || 'Server Error'}.`,
        );
        setBoards(initialBoards);
      });
  };

  return (
    <div className="DeleteBoard">
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
          disabled
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

DeleteBoard.propTypes = {
  id: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default DeleteBoard;
