import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../AppContext';
import TeamControlsMenuItem from './Item/TeamControlsMenuItem';

import './teamcontrolsmenu.sass';

const TeamControlsMenu = ({ handleCreate, handleDelete }) => {
  const { members } = useContext(AppContext);

  return (
    <div className="TeamControlsMenu">
      {members.map((member) => (
        <TeamControlsMenuItem
          username={member.username}
          isAdmin={member.isAdmin}
          isActive={member.isActive}
          handleDelete={handleDelete}
        />
      ))}
      <button className="CreateButton" type="button" onClick={handleCreate}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

TeamControlsMenu.propTypes = {
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
};

export default TeamControlsMenu;
