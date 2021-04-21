import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Button, Col, Form, Row,
} from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';
import AppContext from '../../../AppContext';
import UsersAPI from '../../../api/UsersAPI';

import logo from './deletemember.svg';
import './deletemember.sass';

const DeleteMember = ({ username, toggleOff }) => {
  const { loadBoard, notify } = useContext(AppContext);

  const handleSubmit = (e) => {
    e.preventDefault();

    UsersAPI
      .delete(username)
      .then(() => {
        toggleOff();
        loadBoard();
      })
      .catch((err) => {
        notify(
          'Unable to delete member.',
          `${err?.message || 'Server Error'}.`,
        );
      });
  };

  return (
    <div className="DeleteMember">
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
  username: PropTypes.string.isRequired,
  toggleOff: PropTypes.func.isRequired,
};

export default DeleteMember;
