/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import { Form, Button } from 'react-bootstrap';

import AppContext from '../../../AppContext';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';
import { postBoard } from '../../../misc/api';

import logo from './createboard.svg';
import './createboard.sass';

const CreateBoard = ({ toggleOff }) => {
  const { currentUser, loadBoard } = useContext(AppContext);
  const [name, setName] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    postBoard({ name, team_id: currentUser.teamId })
      .then((res) => {
        toggleOff();
        loadBoard(res.data.id);
      })
      .catch((err) => {
        // TODO: Toast
        console.error(err);
      });
  };

  return (
    <div className="CreateBoard" onClick={toggleOff}>
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
          label="title"
          value={name}
          setValue={setName}
        />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>
      </Form>
    </div>
  );
};

CreateBoard.propTypes = {
  toggleOff: PropTypes.func.isRequired,
};

export default CreateBoard;
