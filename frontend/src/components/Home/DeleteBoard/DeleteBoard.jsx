/* eslint-disable
jsx-a11y/click-events-have-key-events,
jsx-a11y/no-static-element-interactions */

import React from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';

import logo from './deleteboard.svg';
import './deleteboard.sass';

const DeleteBoard = ({ id, name, toggleOff }) => {
  // TODO: Delete board using the ID here
  const handleSubmit = () => console.log(`deleteBoard(${id})`);

  return (
    <div className="DeleteBoard" onClick={toggleOff}>
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
