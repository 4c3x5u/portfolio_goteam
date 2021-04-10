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

import './deletemember.sass';
import logo from './deletemember.svg';

const DeleteMember = ({ id, username, toggleOff }) => {
  // TODO: Delete member using the ID here
  const handleSubmit = () => console.log(`Remove member ${id}`);

  return (
    <div className="DeleteMember" onClick={toggleOff}>
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
          label="username"
          value={username}
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

DeleteMember.propTypes = {
  id: PropTypes.number.isRequired,
  username: PropTypes.string.isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default DeleteMember;
