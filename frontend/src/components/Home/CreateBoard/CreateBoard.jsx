/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import { Form, Button } from 'react-bootstrap';

import AppContext from '../../../AppContext';
import BoardsAPI from '../../../api/BoardsAPI';
import FormGroup from '../../_shared/FormGroup/FormGroup';
import inputType from '../../../misc/inputType';
import ValidateBoard from '../../../validation/ValidateBoard';

import logo from './createboard.svg';
import './createboard.sass';

const CreateBoard = ({ toggleOff }) => {
  const { user, loadBoard } = useContext(AppContext);
  const [name, setName] = useState('');
  const [nameError, setNameError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientNameError = ValidateBoard.name(name);

    BoardsAPI
      .post({ name, team_id: user.teamId })
      .then((res) => { toggleOff(); loadBoard(res.data.id); })
      .catch((err) => {
        setNameError(clientNameError || err.response.data.name || '');
        // TODO: Handle other errors that may arise (Toast)
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
          label="name"
          value={name}
          setValue={setName}
          error={nameError}
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
